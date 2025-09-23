package parser

import (
	"strconv"
	"strings"
)

// CmpHuman сравнивает человекочитаемые размеры (1K, 2M, 3.5G, 512 и т.п.). Базовая единица - 1024
func CmpHuman(a, b string) int {
	aVal, aOk := parseHuman(strings.TrimSpace(a))
	bVal, bOk := parseHuman(strings.TrimSpace(b))

	switch {
	case !aOk && !bOk:
		return strings.Compare(a, b)
	case !aOk:
		return -1
	case !bOk:
		return 1
	case aVal < bVal:
		return -1
	case aVal > bVal:
		return 1
	default:
		return 0
	}
}

// parseHuman парсит 10K, 3.5M, 2G, 512 (без суффикса) в число в байтах и флаг успеха
func parseHuman(s string) (float64, bool) {
	if s == "" {
		return 0, false
	}
	last := s[len(s)-1]
	mult := 1.0
	switch last {
	case 'K', 'k':
		mult = 1 << 10 // килобайт
		s = s[:len(s)-1]
	case 'M', 'm':
		mult = 1 << 20 // мегабайт
		s = s[:len(s)-1]
	case 'G', 'g':
		mult = 1 << 30 // гигабайт
		s = s[:len(s)-1]
	case 'T', 't': // терабайт
		mult = 1 << 40
		s = s[:len(s)-1]
	case 'P', 'p': // петабайт
		mult = 1 << 50
		s = s[:len(s)-1]
	case 'E', 'e': // эксабайт
		mult = 1 << 60
		s = s[:len(s)-1]
	default:
		// без суффикса - просто число
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, false
	}
	return v * mult, true
}
