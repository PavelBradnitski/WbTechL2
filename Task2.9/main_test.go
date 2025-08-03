package main

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"", "", false},
		{"45", "", true},
		{"qwe\\4\\5", "qwe45", false},
		{"qwe\\45", "qwe44444", false},
		{"3abc", "", true},
		{"qwe\\", "", true},

		{"\\3abc", "3abc", false},
		{"z\\9y2", "z9yy", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := UnpackString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
			if got != tt.expected {
				t.Errorf("expected=%q, got=%q", tt.expected, got)
			}
		})
	}
}
