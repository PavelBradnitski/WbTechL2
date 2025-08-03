package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func main() {
	fmt.Println(UnpackString("a4bc2d5e"))
	fmt.Println(UnpackString("abcd"))
	fmt.Println(UnpackString("45"))
	fmt.Println(UnpackString(""))
	fmt.Println(UnpackString("qwe\\4\\5"))
	fmt.Println(UnpackString("qwe\\45"))
}

// UnpackString распаковывает строку с повторяющимися символами и поддержкой экранирования.
func UnpackString(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	var sb strings.Builder
	var prev rune
	escaped := false
	started := false

	for _, r := range str {
		switch {
		case escaped:
			sb.WriteRune(r)
			prev = r
			escaped = false
			started = true

		case r == '\\':
			escaped = true

		case unicode.IsDigit(r):
			if !started {
				return "", errors.New("invalid string: starts with digit")
			}
			count := int(r - '0')
			if count > 0 {
				sb.WriteString(strings.Repeat(string(prev), count-1))
			}

		default:
			sb.WriteRune(r)
			prev = r
			started = true
		}
	}

	if escaped {
		return "", errors.New("invalid string: ends with unfinished escape")
	}

	return sb.String(), nil
}
