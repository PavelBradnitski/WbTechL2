package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

// Options описывает параметры сортировки.
type Options struct {
	// KeyColumn — номер колонки (1-based) для сортировки; 0 — сортировать по всей строке.
	KeyColumn int
	// Numeric — сортировка по числовому значению.
	Numeric bool
	// Reverse — обратный порядок сортировки.
	Reverse bool
	// Unique — выводить только уникальные строки.
	Unique bool
	// Month — сравнение по названию месяца (Jan..Dec).
	Month bool
	// IgnoreTrailBlanks — обрезать хвостовые пробелы перед сравнением.
	IgnoreTrailBlanks bool
	// CheckOnly — только проверить отсортирован ли ввод; сообщить о первом нарушении.
	CheckIfSorted bool
	// HumanNumeric — сравнение чисел с суффиксами (например, 1K, 10M).
	HumanNumeric bool
	// Delimiter — разделитель колонок (по умолчанию TAB).
	Delimiter string
}

type record struct {
	line  string
	key   key
	index int
}

type key struct {
	raw      string
	monthVal int
	numVal   float64
	isNum    bool
}

var (
	flagKeyColumn         = pflag.IntP("key", "k", 0, "sort by column N (1-based), tab-separated by default")
	flagNumeric           = pflag.BoolP("numeric", "n", false, "compare by numeric value")
	flagReverse           = pflag.BoolP("reverse", "r", false, "reverse order")
	flagUnique            = pflag.BoolP("unique", "u", false, "output only unique lines")
	flagMonth             = pflag.BoolP("month", "M", false, "compare by month name (Jan..Dec)")
	flagIgnoreTrailBlanks = pflag.BoolP("ignore-blanks", "b", false, "ignore trailing blanks")
	flagCheckIfSorted     = pflag.BoolP("check", "c", false, "check whether input is sorted")
	flagHumanNumbers      = pflag.BoolP("human-numeric", "h", false, "compare numbers with suffixes (K,M,G, etc.)")
	flagDelimiter         = pflag.StringP("delimiter", "t", "\t", "input column delimiter (default TAB)")
)

func main() {
	// Поддержка объединённых коротких флагов (-nr) и присоединённых значений (-k2, -t,)
	// достигается предварительным расширением os.Args перед разбором
	expanded := expandCombinedShortFlags(os.Args[1:])
	if err := pflag.CommandLine.Parse(expanded); err != nil {
		log.Fatalf("parse flags: %v", err)
	}

	var input *os.File
	args := pflag.Args()
	if len(args) > 1 {
		log.Fatalf("too many arguments: expected at most 1 file, got %d", len(args))
	}
	if len(args) == 1 {
		f, err := os.Open(args[0])
		if err != nil {
			log.Fatalf("open file: %v", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatalf("close file: %v", err)
			}
		}()
		input = f
	} else {
		input = os.Stdin
	}

	options := Options{
		KeyColumn:         *flagKeyColumn,
		Numeric:           *flagNumeric,
		Reverse:           *flagReverse,
		Unique:            *flagUnique,
		Month:             *flagMonth,
		IgnoreTrailBlanks: *flagIgnoreTrailBlanks,
		CheckIfSorted:     *flagCheckIfSorted,
		HumanNumeric:      *flagHumanNumbers,
		Delimiter:         *flagDelimiter,
	}

	if *flagCheckIfSorted {
		// Потоковая проверка без загрузки всего ввода в память
		ok, idx, err := IsSortedReader(input, options)
		if err != nil {
			log.Fatalf("ошибка чтения ввода: %v", err)
		}
		if !ok {
			_, err = fmt.Fprintf(os.Stderr, "не отсортировано: нарушение на строке %d\n", idx)
			if err != nil {
				log.Fatalf("ошибка вывода: %v", err)
			}
			os.Exit(1)
		}
		return
	}

	// Читаем весь ввод в память только для режима сортировки
	scanner := bufio.NewScanner(input)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)
	lines := make([]string, 0, 1024)
	for scanner.Scan() {
		// Нормализация CRLF (Windows) к LF для корректных сравнений
		line := strings.TrimRight(scanner.Text(), "\r")
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("read input: %v", err)
	}

	sorted := SortLines(lines, options)
	writer := bufio.NewWriter(os.Stdout)
	defer func() {
		if err := writer.Flush(); err != nil {
			log.Fatalf("flush writer: %v", err)
		}
	}()
	for _, l := range sorted {
		_, err := fmt.Fprintln(writer, l)
		if err != nil {
			log.Fatalf("write line: %v", err)
		}
	}
}

// expandCombinedShortFlags расширяет объединённые короткие булевы флаги (например, "-nr" → "-n" "-r")
// и обрабатывает присоединённые значения для флагов с аргументами (например, "-k2" → "-k" "2", "-t," → "-t" ",").
// Длинные флаги ("--...") и разделитель "--" сохраняются. Это обеспечивает GNU-подобный UX
// при использовании стандартного механизма разбора.
func expandCombinedShortFlags(args []string) []string {
	if len(args) == 0 {
		return args
	}
	valueFlags := map[byte]bool{
		'k': true, // column index
		't': true, // delimiter
	}
	booleanFlags := map[byte]bool{
		'n': true,
		'r': true,
		'u': true,
		'M': true,
		'b': true,
		'c': true,
		'h': true,
	}

	out := make([]string, 0, len(args)*2)
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" { // прекращаем разбор флагов; передаем остальные аргументы как есть
			out = append(out, a)
			out = append(out, args[i+1:]...)
			break
		}
		if len(a) < 2 || a[0] != '-' {
			out = append(out, a)
			continue
		}
		// Длинный флаг или отрицательное число "-2" — оставляем без изменений
		if strings.HasPrefix(a, "--") || (len(a) == 2 && a[1] >= '0' && a[1] <= '9') {
			out = append(out, a)
			continue
		}
		// Формы с '=' типа -k=2 или -t=, оставляем как есть (pflag сам обработает)
		if strings.Contains(a, "=") {
			out = append(out, a)
			continue
		}
		if len(a) == 2 { // одиночный короткий флаг, например -n
			out = append(out, a)
			continue
		}
		// проходим по символам после '-'
		j := 1
		for j < len(a) {
			ch := a[j]
			if valueFlags[ch] {
				out = append(out, "-"+string(ch))
				val := a[j+1:]
				if val != "" {
					out = append(out, val)
				}
				// значение занимает остаток строки
				break
			}
			if booleanFlags[ch] {
				out = append(out, "-"+string(ch))
				j++
				continue
			}
			// Неизвестный флаг: для надёжности оставляем исходный аргумент без изменений
			out = append(out, a)
			break
		}
	}
	return out
}

// SortLines возвращает новый слайс с отсортированными строками согласно опциям.
func SortLines(lines []string, opt Options) []string {
	if len(lines) == 0 {
		return nil
	}

	// Подготовить записи с предвычисленными ключами для эффективного сравнения.
	items := make([]record, 0, len(lines))
	for idx, l := range lines {
		key := extractKey(l, opt)
		items = append(items, record{line: l, key: key, index: idx})
	}

	sort.SliceStable(items, func(i, j int) bool {
		ci := items[i]
		cj := items[j]
		cmp := compareKeys(ci.key, cj.key, opt)
		if opt.Reverse {
			cmp = -cmp
		}
		if cmp == 0 {
			// При равенстве ключей сохраняем исходный порядок по индексу.
			return ci.index < cj.index
		}
		return cmp < 0
	})

	out := make([]string, 0, len(items))
	var lastKey *key
	for _, it := range items {
		if opt.Unique {
			if lastKey == nil {
				k := it.key
				lastKey = &k
				out = append(out, it.line)
				continue
			}
			if compareKeys(*lastKey, it.key, opt) == 0 {
				// Одинаковы согласно выбранным опциям сравнения — пропускаем дубликаты.
				continue
			}
			k := it.key
			lastKey = &k
		}
		out = append(out, it.line)
	}
	return out
}

func extractKey(line string, opt Options) key {
	val := line
	if opt.KeyColumn > 0 {
		val = extractColumnValue(line, opt.Delimiter, opt.KeyColumn)
	}
	v := val
	if opt.IgnoreTrailBlanks {
		v = strings.TrimRight(v, " \t")
	}
	k := key{raw: v}
	if opt.Month {
		k.monthVal = monthIndex(v)
	}
	if opt.HumanNumeric {
		if nv, ok := parseHumanNumber(v); ok {
			k.isNum = true
			k.numVal = nv
		}
	} else if opt.Numeric {
		if nv, ok := parseFloat(v); ok {
			k.isNum = true
			k.numVal = nv
		}
	}
	return k
}

// extractColumnValue возвращает значение N-й (1-based) колонки из строки,
// используя указанный разделитель. Если колонки нет — возвращает пустую строку.
func extractColumnValue(line, delimiter string, columnIndex int) string {
	if columnIndex <= 0 {
		return line
	}
	target := columnIndex - 1
	if delimiter == "" {
		return ""
	}
	if len(delimiter) == 1 {
		// Оптимизация для одиночного байта-разделителя
		d := delimiter[0]
		current := 0
		start := 0
		for i := 0; i < len(line); i++ {
			if line[i] == d {
				if current == target {
					return line[start:i]
				}
				current++
				start = i + 1
			}
		}
		if current == target {
			return line[start:]
		}
		return ""
	}
	// Общий случай для строкового разделителя
	current := 0
	start := 0
	for {
		if current == target {
			pos := strings.Index(line[start:], delimiter)
			if pos == -1 {
				return line[start:]
			}
			return line[start : start+pos]
		}
		pos := strings.Index(line[start:], delimiter)
		if pos == -1 {
			return ""
		}
		start = start + pos + len(delimiter)
		current++
	}
}

// IsSortedReader выполняет потоковую проверку отсортированности ввода без загрузки
// всего содержимого в память. Возвращает ok, 1-based индекс строки нарушения и ошибку чтения (если была).
func IsSortedReader(r io.Reader, opt Options) (bool, int, error) {
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)
	var (
		lineIndex int
		prev      key
		havePrev  bool
	)
	for scanner.Scan() {
		lineIndex++
		line := strings.TrimRight(scanner.Text(), "\r")
		cur := extractKey(line, opt)
		if !havePrev {
			prev = cur
			havePrev = true
			continue
		}
		cmp := compareKeys(prev, cur, opt)
		if opt.Reverse {
			cmp = -cmp
		}
		if cmp > 0 || (opt.Unique && cmp == 0) {
			return false, lineIndex, nil
		}
		prev = cur
	}
	if err := scanner.Err(); err != nil {
		return false, 0, err
	}
	return true, 0, nil
}

func compareKeys(a, b key, opt Options) int {
	if opt.Month {
		if a.monthVal != b.monthVal {
			return a.monthVal - b.monthVal
		}
	}
	if a.isNum || b.isNum {
		if !a.isNum {
			return -1
		}
		if !b.isNum {
			return 1
		}
		if a.numVal < b.numVal {
			return -1
		}
		if a.numVal > b.numVal {
			return 1
		}
		return 0
	}
	if a.raw < b.raw {
		return -1
	}
	if a.raw > b.raw {
		return 1
	}
	return 0
}

func parseFloat(s string) (float64, bool) {
	if s == "" {
		return 0, false
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

var monthNames = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
	"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

func monthIndex(s string) int {
	if s == "" {
		return 0
	}
	lower := strings.ToLower(strings.TrimSpace(s))
	if len(lower) >= 3 {
		lower = lower[:3]
	}
	if v, ok := monthNames[lower]; ok {
		return v
	}
	return 0
}

func parseHumanNumber(s string) (float64, bool) {
	if s == "" {
		return 0, false
	}
	s = strings.TrimSpace(s)
	// Опциональные суффиксы: B/K/M/G/T (без учёта регистра).
	// Примеры: 10K, 10KB, 1.5M, 2G, 3T
	base := s
	mult := 1.0

	// Затем проверяем буквенный множитель и отрезаем его
	if n := len(base); n > 0 {
		switch base[n-1] {
		case 'B', 'b':
			base = base[:n-1]
		case 'K', 'k':
			mult = 1024
			base = base[:n-1]
		case 'M', 'm':
			mult = 1024 * 1024
			base = base[:n-1]
		case 'G', 'g':
			mult = 1024 * 1024 * 1024
			base = base[:n-1]
		case 'T', 't':
			mult = 1024 * 1024 * 1024 * 1024
			base = base[:n-1]
		}
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(base), 64)
	if err != nil {
		return 0, false
	}
	val := f * mult
	if math.IsInf(val, 0) || math.IsNaN(val) {
		return 0, false
	}
	return val, true
}
