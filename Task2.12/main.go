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

// parseFields разбирает строку с номерами полей и диапазонами.
// Возвращает слайс с номерами полей (индексация с 0).
// Например, "1,3-5" -> [0, 2, 3, 4]
func parseFields(fieldsStr string) ([]int, error) {
	result := make([]int, 0)
	parts := strings.Split(fieldsStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}

			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid number in range: %s", rangeParts[0])
			}
			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid number in range: %s", rangeParts[1])
			}

			if start > end {
				return nil, fmt.Errorf("invalid range: start > end (%d > %d)", start, end)
			}

			for i := start; i <= end; i++ {
				result = append(result, i-1)
			}
		} else {
			fieldNum, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid field number: %s", part)
			}
			result = append(result, fieldNum-1)
		}
	}
	return result, nil
}

func main() {
	var (
		fieldsStr = flag.String("f", "", "Fields to print (comma-separated, can include ranges)")
		delimiter = flag.String("d", "\t", "Field delimiter")
		separated = flag.Bool("s", false, "Only lines containing delimiter")
	)
	flag.Parse()
	var fields []int
	if *fieldsStr != "" {
		var err error
		fields, err = parseFields(*fieldsStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing fields: %v\n", err)
			os.Exit(1)
		}
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}

		line = strings.TrimSuffix(line, "\n")
		if *separated && !strings.Contains(line, *delimiter) {
			// if *separated && !strings.Contains(line, "\t") {

			continue
		}
		line = "field1\tfield2\tfield3\tfield4"

		parts := strings.Split(line, *delimiter)
		fmt.Printf("Parts %v\n", parts)
		// parts := strings.Split(line, "\t")
		output := make([]string, 0)

		for _, field := range fields {
			if field >= 0 && field < len(parts) {
				output = append(output, parts[field])
			}
		}

		fmt.Println(strings.Join(output, *delimiter))
	}
}
