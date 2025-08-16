package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestExpandCombinedShortFlags(t *testing.T) {
	cases := []struct {
		in  []string
		out []string
	}{
		{[]string{"-nr"}, []string{"-n", "-r"}},
		{[]string{"-k2"}, []string{"-k", "2"}},
		{[]string{"-t,"}, []string{"-t", ","}},
		{[]string{"-n", "-r"}, []string{"-n", "-r"}},
		{[]string{"--long", "file.txt"}, []string{"--long", "file.txt"}},
		{[]string{"--", "-nr", "x"}, []string{"--", "-nr", "x"}},
	}
	for _, c := range cases {
		got := expandCombinedShortFlags(c.in)
		if !reflect.DeepEqual(got, c.out) {
			t.Fatalf("expand mismatch: in=%v got=%v want=%v", c.in, got, c.out)
		}
	}
}

func TestSortLinesLexicographic(t *testing.T) {
	lines := []string{"banana", "apple", "cherry"}
	opt := Options{}
	out := SortLines(lines, opt)
	expected := []string{"apple", "banana", "cherry"}
	if strings.Join(out, "\n") != strings.Join(expected, "\n") {
		t.Fatalf("lexicographic sort mismatch: got %q want %q", out, expected)
	}
}

func TestSortLinesNumeric_ByColumn(t *testing.T) {
	lines := []string{"item1 10", "item2 2", "item3 100"}
	opt := Options{KeyColumn: 2, Numeric: true, Delimiter: " "}
	out := SortLines(lines, opt)
	expected := []string{"item2 2", "item1 10", "item3 100"}
	if strings.Join(out, "\n") != strings.Join(expected, "\n") {
		t.Fatalf("numeric column sort mismatch: got %q want %q", out, expected)
	}
}

func TestSortLinesReverse(t *testing.T) {
	lines := []string{"a", "b", "c"}
	opt := Options{Reverse: true}
	out := SortLines(lines, opt)
	expected := []string{"c", "b", "a"}
	if strings.Join(out, "\n") != strings.Join(expected, "\n") {
		t.Fatalf("reverse sort mismatch: got %q want %q", out, expected)
	}
}

func TestSortLinesUniqueStable(t *testing.T) {
	lines := []string{"b", "a", "a", "b", "c"}
	opt := Options{Unique: true}
	out := SortLines(lines, opt)
	expected := []string{"a", "b", "c"}
	if strings.Join(out, "\n") != strings.Join(expected, "\n") {
		t.Fatalf("unique sort mismatch: got %q want %q", out, expected)
	}
}

func TestSortLinesMonth(t *testing.T) {
	lines := []string{"Apr", "Jan", "Dec"}
	opt := Options{Month: true}
	out := SortLines(lines, opt)
	expected := []string{"Jan", "Apr", "Dec"}
	if strings.Join(out, "\n") != strings.Join(expected, "\n") {
		t.Fatalf("month sort mismatch: got %q want %q", out, expected)
	}
}

func TestSortLinesIgnoreBlankTrailing(t *testing.T) {
	lines := []string{"a", "a   ", "b"}
	opt := Options{IgnoreTrailBlanks: true, Unique: true}
	out := SortLines(lines, opt)
	expected := []string{"a", "b"}
	if strings.Join(out, "\n") != strings.Join(expected, "\n") {
		t.Fatalf("ignore blanks unique mismatch: got %q want %q", out, expected)
	}
}

func TestSortLinesHumanNumericByColumn(t *testing.T) {
	lines := []string{"file1 10K", "file2 2M", "file3 500"}
	opt := Options{KeyColumn: 2, HumanNumeric: true, Delimiter: " "}
	out := SortLines(lines, opt)
	expected := []string{"file3 500", "file1 10K", "file2 2M"}
	if strings.Join(out, "\n") != strings.Join(expected, "\n") {
		t.Fatalf("human numeric sort mismatch: got %q want %q", out, expected)
	}
}

func TestExtractColumnValue(t *testing.T) {
	line := "a,bb,3"
	if got := extractColumnValue(line, ",", 1); got != "a" {
		t.Fatalf("col1 mismatch: got %q", got)
	}
	if got := extractColumnValue(line, ",", 2); got != "bb" {
		t.Fatalf("col2 mismatch: got %q", got)
	}
	if got := extractColumnValue(line, ",", 3); got != "3" {
		t.Fatalf("col3 mismatch: got %q", got)
	}
	if got := extractColumnValue(line, ",", 4); got != "" {
		t.Fatalf("col4 should be empty, got %q", got)
	}
}

func TestIsSortedReaderCRLFSorted(t *testing.T) {
	data := "apple\r\nbanana\r\ncherry\r\n"
	r := strings.NewReader(data)
	ok, idx, err := IsSortedReader(r, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok || idx != 0 {
		t.Fatalf("expected sorted ok, got ok=%v idx=%d", ok, idx)
	}
}

func TestIsSortedReaderUnsorted(t *testing.T) {
	data := "a\n c\n b\n"
	r := strings.NewReader(data)
	ok, idx, err := IsSortedReader(r, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok || idx == 0 {
		t.Fatalf("expected unsorted, got ok=%v idx=%d", ok, idx)
	}
}
