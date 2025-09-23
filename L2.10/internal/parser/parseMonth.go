package parser

import "strings"

var months = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
	"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

// CmpMonth - функция сравнения по названию месяца
func CmpMonth(a, b string) int {
	aInt := months[strings.ToLower(strings.TrimSpace(a))]
	bInt := months[strings.ToLower(strings.TrimSpace(b))]

	switch {
	case aInt == 0 && bInt == 0:
		return strings.Compare(a, b)
	case aInt == 0:
		return -1
	case bInt == 0:
		return 1
	case aInt < bInt:
		return -1
	case aInt > bInt:
		return 1
	default:
		return 0
	}
}
