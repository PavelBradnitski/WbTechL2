package main

import (
	"reflect"
	"strconv"
	"testing"
)

func TestFindMatchingLinesFlags(t *testing.T) {
	lines := []string{
		"Hello world",
		"Go is awesome",
		"This is a test",
		"Another line",
		"HELLO again",
	}

	tests := []struct {
		name     string
		pattern  string
		params   *SearchParams
		expected []string
	}{
		{
			name:    "-A flag (after lines)",
			pattern: "Go",
			params: &SearchParams{
				fixedString: true,
				afterLines:  1,
			},
			expected: []string{"Go is awesome", "This is a test"},
		},
		{
			name:    "-B flag (before lines)",
			pattern: "Go",
			params: &SearchParams{
				fixedString: true,
				beforeLines: 1,
			},
			expected: []string{"Hello world", "Go is awesome"},
		},
		{
			name:    "-C flag (context lines override -A/-B)",
			pattern: "Go",
			params: &SearchParams{
				fixedString:  true,
				beforeLines:  0,
				afterLines:   0,
				contextLines: 2,
			},
			expected: []string{"Hello world", "Go is awesome", "This is a test", "Another line"},
		},
		{
			name:    "-c flag (count only, case-sensitive)",
			pattern: "Hello",
			params: &SearchParams{
				fixedString: true,
				countOnly:   true,
			},
			expected: []string{strconv.Itoa(1)}, // только "Hello world"
		},
		{
			name:    "-c flag (count only, ignore case)",
			pattern: "Hello",
			params: &SearchParams{
				fixedString: true,
				countOnly:   true,
				ignoreCase:  true,
			},
			expected: []string{strconv.Itoa(2)}, // "Hello world" и "HELLO again"
		},
		{
			name:    "-i flag (ignore case)",
			pattern: "hello",
			params: &SearchParams{
				fixedString: true,
				ignoreCase:  true,
			},
			expected: []string{"Hello world", "HELLO again"},
		},
		{
			name:    "-v flag (invert match, case-sensitive)",
			pattern: "Hello",
			params: &SearchParams{
				fixedString: true,
				invertMatch: true,
			},
			expected: []string{"Go is awesome", "This is a test", "Another line", "HELLO again"},
		},
		{
			name:    "-v flag (invert match, ignore case)",
			pattern: "hello",
			params: &SearchParams{
				fixedString: true,
				invertMatch: true,
				ignoreCase:  true,
			},
			expected: []string{"Go is awesome", "This is a test", "Another line"},
		},
		{
			name:    "-F flag (fixed string match)",
			pattern: "Go is awesome",
			params: &SearchParams{
				fixedString: true,
			},
			expected: []string{"Go is awesome"},
		},
		{
			name:    "-n flag (line numbers)",
			pattern: "test",
			params: &SearchParams{
				fixedString: true,
				lineNumber:  true,
			},
			expected: []string{"3:This is a test"},
		},
		{
			name:    "Combined flags: -i -n -C",
			pattern: "hello",
			params: &SearchParams{
				fixedString:  true,
				ignoreCase:   true,
				lineNumber:   true,
				contextLines: 1,
			},
			expected: []string{"1:Hello world", "2:Go is awesome", "4:Another line", "5:HELLO again"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findMatchingLines(lines, tt.pattern, tt.params)

			if tt.params.countOnly {
				count := len(result)
				result = []string{strconv.Itoa(count)}
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}
