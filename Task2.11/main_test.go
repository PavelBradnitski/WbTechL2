package main

import (
	"reflect"
	"testing"
)

func TestFindAnagrams(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected map[string][]string
	}{
		{
			name:     "Empty input",
			input:    []string{},
			expected: map[string][]string{},
		},
		{
			name:     "No anagrams",
			input:    []string{"стол", "стул", "лампа"},
			expected: map[string][]string{},
		},
		{
			name:  "Simple anagrams",
			input: []string{"пятак", "пятка", "тяпка"},
			expected: map[string][]string{
				"пятак": {"пятак", "пятка", "тяпка"},
			},
		},
		{
			name:  "Multiple anagram groups",
			input: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"},
			expected: map[string][]string{
				"листок": {"листок", "слиток", "столик"},
				"пятак":  {"пятак", "пятка", "тяпка"},
			},
		},
		{
			name:  "Mixed case anagrams",
			input: []string{"ПяТАк", "ПятКа", "тяпКа"},
			expected: map[string][]string{
				"пятак": {"пятак", "пятка", "тяпка"},
			},
		},
		{
			name:  "Anagrams with same word",
			input: []string{"стол", "стол", "кот", "ток", "стол"},
			expected: map[string][]string{
				"кот":  {"кот", "ток"},
				"стол": {"стол", "стол", "стол"},
			},
		},
		{
			name:  "Single letter words",
			input: []string{"a", "b", "a"},
			expected: map[string][]string{
				"a": {"a", "a"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := FindAnagrams(tc.input)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test Case: %s\nExpected: %v\nActual: %v", tc.name, tc.expected, actual)
			}
		})
	}
}
