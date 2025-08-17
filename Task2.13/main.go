package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// cutOptions структура для хранения параметров командной строки
type cutOptions struct {
	fields    string
	delimiter string
	separated bool // если true, то выводить только строки с разделителем
}

func main() {
	options := parseCommandLineFlags()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		if options.separated && !strings.Contains(line, options.delimiter) {
			continue
		}
		fieldsIdx := parseFieldsList(options.fields)
		selectedLine := selectFields(line, options, fieldsIdx)
		fmt.Println(selectedLine)
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "ошибка чтения: %v\n", err)
		os.Exit(1)
	}
}

func parseCommandLineFlags() cutOptions {
	flagF := flag.String("f", "", "выбрать поля (колонки)")
	flagD := flag.String("d", "\t", "использовать другой разделитель")
	flagS := flag.Bool("s", false, "только строки с разделителем")
	flag.Parse()

	if *flagF == "" && !*flagS {
		fmt.Println("необходимо указать хотя бы один из флагов: -f или -s")
		os.Exit(1)
	}

	options := cutOptions{
		fields:    *flagF,
		delimiter: *flagD,
		separated: *flagS,
	}

	return options
}
func selectFields(line string, options cutOptions, fieldsIdx []int) string {
	// если поля не заданы – возвращаем строку целиком
	if len(fieldsIdx) == 0 {
		return line
	}

	// если включён separated, но в строке нет разделителя → пропускаем
	if options.separated && !strings.Contains(line, options.delimiter) {
		return ""
	}

	// специальный случай для табуляции
	if options.delimiter == "\t" {
		fields := strings.Fields(line)
		var b strings.Builder
		for i, idx := range fieldsIdx {
			if idx > 0 && idx <= len(fields) {
				if i > 0 {
					b.WriteByte(' ') // cut заменяет \t на пробел
				}
				b.WriteString(fields[idx-1])
			}
		}
		return b.String()
	}

	// оптимизация для 1-символьного разделителя (например, , ;)
	if len(options.delimiter) == 1 {
		return fastSelect(line, options.delimiter[0], fieldsIdx)
	}

	// общий случай (много-символьный разделитель)
	fields := strings.Split(line, options.delimiter)
	var b strings.Builder
	for i, idx := range fieldsIdx {
		if idx > 0 && idx <= len(fields) {
			if i > 0 {
				b.WriteString(options.delimiter)
			}
			b.WriteString(fields[idx-1])
		}
	}
	return b.String()
}

// оптимизированный вариант для 1-символьного разделителя
func fastSelect(line string, delim byte, fieldsIdx []int) string {
	var b strings.Builder
	fieldStart := 0
	fieldNum := 1
	idxPos := 0

	for i := 0; i <= len(line); i++ {
		if i == len(line) || line[i] == delim {
			if idxPos < len(fieldsIdx) && fieldsIdx[idxPos] == fieldNum {
				if b.Len() > 0 {
					b.WriteByte(delim)
				}
				b.WriteString(line[fieldStart:i])
				idxPos++
			}
			fieldStart = i + 1
			fieldNum++
		}
	}
	return b.String()
}

// parseFieldsList разбирает строку с номерами полей, разделенными запятыми, и возвращает слайс индексов
func parseFieldsList(fieldsList string) []int {
	fieldIndexes := make([]int, 0)
	fields := strings.Split(fieldsList, ",")

	for _, field := range fields {
		idx, err := strconv.Atoi(field)
		if err == nil && idx > 0 {
			fieldIndexes = append(fieldIndexes, idx)
		}
	}

	return fieldIndexes
}
