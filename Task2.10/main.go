package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-k N] filename\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	kFlag := flag.Int("k", 0, "column number (1-based) to sort by (default: whole line). delimiter is tab)")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	filename := flag.Arg(0)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file %s: %v", filename, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %v", err)
	}
	lines = sortLines(lines, *kFlag)

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
}

// sortLines сортирует массив строк lines по ключу -k (номер колонки, 1-based).
// Разделитель колонок - табуляция.
// Возвращает отсортированный срез строк.
func sortLines(lines []string, k int) []string {
	keys := make([]string, len(lines))

	for i, line := range lines {
		fields := strings.Split(line, "\t")
		if k <= len(fields) {
			keys[i] = fields[k-1]
		} else {
			keys[i] = ""
		}
	}

	idx := make([]int, len(lines))
	for i := range idx {
		idx[i] = i
	}
	sort.SliceStable(idx, func(i, j int) bool {
		return keys[idx[i]] < keys[idx[j]]
	})

	sorted := make([]string, len(lines))
	for i, id := range idx {
		sorted[i] = lines[id]
	}
	return sorted
}
