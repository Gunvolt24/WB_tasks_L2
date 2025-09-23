package comparator

import (
	"strings"

	"github.com/Gunvolt24/wb_l2/L2.10/internal/options"
	"github.com/Gunvolt24/wb_l2/L2.10/internal/parser"
)

// CreateComparator - создает функцию сравнения (less) и функцию проверки равенства (equal) по ключу.
func CreateComparator(o options.Options) (less func(a, b string) bool, equal func(a, b string) bool) {
	extract := func(s string) string {
		if o.TrailBlanks {
			s = trimTrailing(s)
		}
		if o.Column <= 0 {
			return s
		}
		return column(s, o.Splitter, o.Column)
	}

	compare := func(a, b string) int {
		keyA, keyB := extract(a), extract(b)
		switch {
		case o.HumanReadable:
			return parser.CmpHuman(keyA, keyB)
		case o.NumSort:
			return parser.CmpNumSort(keyA, keyB)
		case o.MonthNames:
			return parser.CmpMonth(keyA, keyB)
		default:
			switch {
			case keyA < keyB:
				return -1
			case keyA > keyB:
				return 1
			default:
				return 0
			}
		}
	}

	less = func(a, b string) bool {
		c := compare(a, b)
		if o.Reverse {
			return c > 0
		}
		return c < 0
	}

	equal = func(a, b string) bool {
		return compare(a, b) == 0
	}

	return less, equal
}

// column - возвращает N-ый столбец строки s, разделённой d.
// При отсутствии N-ого столбца, возвращает пустую строку.
func column(s, d string, idx1 int) string {
	if d == "" {
		d = "\t"
	}

	parts := strings.Split(s, d)
	i := idx1 - 1
	if i < 0 || i >= len(parts) {
		return ""
	}
	return parts[i]
}

// trimTrailing - удаляет только пробельные символы в конце строки (ведущие пробелы остаются).
func trimTrailing(s string) string {
	end := len(s)
	for end > 0 {
		switch s[end-1] {
		case ' ', '\t':
			end--
		default:
			return s[:end]
		}
	}
	return s[:end]
}
