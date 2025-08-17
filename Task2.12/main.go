package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
	Задание:

	Реализовать утилиту фильтрации по аналогии с консольной утилитой
	(man grep — смотрим описание и основные параметры).

	Реализовать поддержку утилитой следующих ключей:
	-A - "after": печатать +N строк после совпадения;
	-B - "before": печатать +N строк до совпадения;
	-C - "context": (A+B) печатать ±N строк вокруг совпадения;
	-c - "count": количество строк;
	-i - "ignore-case": игнорировать регистр;
	-v - "invert": вместо совпадения, исключать;
	-F - "fixed": точное совпадение со строкой, не паттерн;
	-n - "line num": напечатать номер строки.
*/

const outputFileName = "grep_result.txt"

// SearchParams - структура для хранения параметров поиска
type SearchParams struct {
	afterLines   int
	beforeLines  int
	contextLines int
	countOnly    bool
	ignoreCase   bool
	invertMatch  bool
	fixedString  bool
	lineNumber   bool
}

func main() {
	// Флаги для запуска утилиты
	afterLines := flag.Int("A", 0, "печатать +N строк после совпадения")
	beforeLines := flag.Int("B", 0, "печатать +N строк до совпадения")
	contextLines := flag.Int("C", 0, "печатать ±N строк вокруг совпадения")
	countOnly := flag.Bool("c", false, "количество строк")
	ignoreCase := flag.Bool("i", false, "игнорировать регистр")
	invertMatch := flag.Bool("v", false, "вместо совпадения, исключать")
	fixedString := flag.Bool("F", false, "точное совпадение со строкой, не паттерн")
	lineNumber := flag.Bool("n", false, "печатать номер строки")

	flag.Parse()

	args := flag.Args()

	if len(args) < 2 {
		fmt.Println("Необходимо указать паттерн и имя файла.")
		os.Exit(1)
	}

	pattern := args[0]
	fileName := args[1]

	// Создаем структуру с параметрами поиска
	params := &SearchParams{
		afterLines:   *afterLines,
		beforeLines:  *beforeLines,
		contextLines: *contextLines,
		countOnly:    *countOnly,
		ignoreCase:   *ignoreCase,
		invertMatch:  *invertMatch,
		fixedString:  *fixedString,
		lineNumber:   *lineNumber,
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Считываем строки из файла
	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		os.Exit(1)
	}

	// Выполняем поиск и получаем результат
	result := findMatchingLines(lines, pattern, params)

	// Открытие файла для записи
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	// Пишем результаты в файл
	if params.countOnly {
		_, err = outputFile.WriteString(strconv.Itoa(len(result)))
	} else {
		for idx, line := range result {
			_, err = outputFile.WriteString(line)
			if idx != len(result)-1 {
				_, err = outputFile.WriteString("\n")
			}
		}
	}
	if err != nil {
		fmt.Println("Ошибка записи в файл:", err)
		os.Exit(1)
	}
}

// contains проверяет, содержится ли элемент в срезе целых чисел.
func contains(slice []int, element int) bool {
	for _, value := range slice {
		if value == element {
			return true
		}
	}
	return false
}

// processLine определяет, соответствует ли строка паттерну с учётом флагов
func processLine(line string, pattern string, params *SearchParams) bool {
	if params.fixedString {
		if params.ignoreCase {
			return strings.Contains(strings.ToLower(line), strings.ToLower(pattern)) != params.invertMatch
		}
		return strings.Contains(line, pattern) != params.invertMatch
	}

	// Обработка паттерна с '.' (замена на произвольный символ)
	var n []int
	for idx, value := range pattern {
		if value == '.' {
			n = append(n, idx)
		}
	}

	newPattern := strings.ReplaceAll(pattern, ".", "")
	newLines := strings.Split(line, " ")

	for _, newLine := range newLines {
		var builder strings.Builder
		for i, char := range newLine {
			// contains проверяет, содержится ли элемент в срезе целых чисел.
			if !contains(n, i) {
				builder.WriteRune(char)
			}
		}

		if params.ignoreCase {
			if strings.Contains(strings.ToLower(builder.String()), strings.ToLower(newPattern)) != params.invertMatch {
				return true
			}
		} else {
			if strings.Contains(builder.String(), newPattern) != params.invertMatch {
				return true
			}
		}
	}

	return false
}

// findMatchingLines ищет строки, соответствующие паттерну, и возвращает их с контекстом
func findMatchingLines(lines []string, pattern string, params *SearchParams) []string {
	var result []string
	var matchingLineNumbers []int

	for num, line := range lines {
		if processLine(line, pattern, params) {
			matchingLineNumbers = append(matchingLineNumbers, num+1)
		}
	}

	for _, lineNumber := range matchingLineNumbers {
		// Определяем начало и конец контекста
		contextStart := lineNumber - params.beforeLines - 1
		contextEnd := lineNumber + params.afterLines - 1

		// Корректируем контекст если есть флаг -C
		if params.contextLines > 0 {
			contextStart = lineNumber - params.contextLines - 1
			contextEnd = lineNumber + params.contextLines - 1
		}

		// Обрезаем контекст чтобы не вылезти за границы массива строк
		if contextStart < 0 {
			contextStart = 0
		}
		if contextEnd >= len(lines) {
			contextEnd = len(lines) - 1
		}

		for i := contextStart; i <= contextEnd; i++ {
			outputLine := lines[i]
			if params.lineNumber {
				outputLine = strconv.Itoa(i+1) + ":" + outputLine
			}
			result = append(result, outputLine)
		}
	}
	return result
}
