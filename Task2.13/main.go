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

		selectedLine := selectFields(line, options)
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

// selectFields выбирает указанные поля (колонки) из строки
func selectFields(line string, options cutOptions) string {
	if options.fields == "" {
		return line
	}

	if options.separated && !strings.Contains(line, options.delimiter) {
		return ""
	}

	var fields []string
	if options.delimiter == "\t" {
		fields = strings.Fields(line)
	} else {
		fields = strings.Split(line, options.delimiter)
	}

	fieldsIdx := parseFieldsList(options.fields)
	selectedFields := make([]string, 0, len(fieldsIdx))

	for _, idx := range fieldsIdx {
		if idx > 0 && idx <= len(fields) {
			selectedFields = append(selectedFields, fields[idx-1])
		}
	}

	var delimiter string
	if options.delimiter == "\t" {
		delimiter = " "
	} else {
		delimiter = options.delimiter
	}
	return strings.Join(selectedFields, delimiter)
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
