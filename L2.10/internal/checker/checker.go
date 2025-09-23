package checker

import (
	"bufio"
	"io"
)

// CheckSorted - проверяет отсортированность потока на лету.
// Возвращает true, если поток отсортирован, иначе false.
func CheckSorted(r io.Reader, less func(a, b string) bool) (ok bool, i, j int, err error) {
	scanner := bufio.NewScanner(r)
	buffer := make([]byte, 64*1024)
	scanner.Buffer(buffer, 4*1024*1024) // до 4 МБ на одну строку

	line := 0
	
	var (
		prev string
		has bool
	)

	for scanner.Scan() {
		cur := scanner.Text()
		if has && less(cur, prev) {
			return false, line - 1, line, nil
		}
		prev = cur
		has = true
		line++
	}

	if err := scanner.Err(); err != nil {
		return false, -1, -1, err
	}

	return true, -1, -1, nil
}
	
