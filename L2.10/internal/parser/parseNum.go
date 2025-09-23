package parser

import (
	"strconv"
	"strings"
)

// CmpNumSort - функция сравнения по числовому представлению
func CmpNumSort(a, b string) int {
	aFloat, aErr := strconv.ParseFloat(strings.TrimSpace(a), 64)
	bFloat, bErr := strconv.ParseFloat(strings.TrimSpace(b), 64)

	switch {
	case aErr != nil && bErr != nil:
		return strings.Compare(a, b) // если преобразование не удалось, сравниваем по строковому представлению
	case aErr != nil:
		return -1
	case bErr != nil:
		return 1
	case aFloat < bFloat:
		return -1
	case aFloat > bFloat:
		return 1
	default:
		return 0
	}
}
