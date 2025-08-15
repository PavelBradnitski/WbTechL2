package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	anagrams := FindAnagrams(words)
	fmt.Println(anagrams)
}

// FindAnagrams находит все множества анаграмм в заданном массиве строк.
func FindAnagrams(words []string) map[string][]string {
	anagramGroups := make(map[string][]string)
	wordToSignature := make(map[string]string)

	for _, word := range words {
		lowerWord := strings.ToLower(word)
		// сортируем символы для сравнения строк
		chars := strings.Split(lowerWord, "")
		sort.Strings(chars)
		signature := strings.Join(chars, "")

		if _, ok := anagramGroups[signature]; !ok {
			anagramGroups[signature] = []string{}
		}

		anagramGroups[signature] = append(anagramGroups[signature], lowerWord)
		if _, ok := wordToSignature[signature]; !ok {
			wordToSignature[signature] = lowerWord
		}

	}

	result := make(map[string][]string)
	for signature, group := range anagramGroups {
		if len(group) > 1 {
			sort.Strings(group)
			result[wordToSignature[signature]] = group
		}
	}

	return result
}
