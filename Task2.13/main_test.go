package main

import (
	"testing"
)

func TestSelectFields(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		options  cutOptions
		expected string
	}{
		{
			name:     "без опций – вернуть строку целиком",
			line:     "a b c d",
			options:  cutOptions{fields: "", delimiter: " ", separated: false},
			expected: "a b c d",
		},
		{
			name:     "разделитель таб, выбрать первое поле",
			line:     "foo\tbar\tbaz",
			options:  cutOptions{fields: "1", delimiter: "\t", separated: false},
			expected: "foo",
		},
		{
			name:     "разделитель запятая, выбрать второе поле",
			line:     "1,2,3,4",
			options:  cutOptions{fields: "2", delimiter: ",", separated: false},
			expected: "2",
		},
		{
			name:     "несколько полей",
			line:     "a,b,c,d",
			options:  cutOptions{fields: "2,4", delimiter: ",", separated: false},
			expected: "b,d",
		},
		{
			name:     "невалидный индекс (слишком большой)",
			line:     "x,y",
			options:  cutOptions{fields: "3", delimiter: ",", separated: false},
			expected: "",
		},
		{
			name:     "некорректные индексы в списке",
			line:     "p q r",
			options:  cutOptions{fields: "0,-1,abc,2", delimiter: " ", separated: false},
			expected: "q",
		},
		{
			name:     "разделитель таб, объединение через пробел",
			line:     "one\ttwo\tthree",
			options:  cutOptions{fields: "1,3", delimiter: "\t", separated: false},
			expected: "one three",
		},
		{
			name:     "separated=true, строка без разделителя → отбрасывается",
			line:     "abcdef",
			options:  cutOptions{fields: "1", delimiter: ",", separated: true},
			expected: "",
		},
		{
			name:     "separated=true, строка с разделителем → работает как обычно",
			line:     "foo,bar,baz",
			options:  cutOptions{fields: "2", delimiter: ",", separated: true},
			expected: "bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldsIdx := parseFieldsList(tt.options.fields)
			result := selectFields(tt.line, tt.options, fieldsIdx)
			if result != tt.expected {
				t.Errorf("ожидалось %q, получили %q", tt.expected, result)
			}
		})
	}
}
