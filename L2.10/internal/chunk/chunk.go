package chunk

import (
	"bufio"
	"io"
	"os"
	"sort"
)

// SplitSort - разбивает поток на чанки (до MaxLines),
// сортирует каждый чанк и сохраняет во временные файлы.
// Если Unique=true, то удаляет дублирующиеся строки по ключу (equal).
func SplitSort(
	r io.Reader,
	maxLines int,
	less func(a, b string) bool,
	equal func(a, b string) bool,
	tempDir string,
	unique bool,
) ([]string, error) {
	scanner := bufio.NewScanner(r)
	buffer := make([]byte, 64*1024)
	scanner.Buffer(buffer, 4*1024*1024) // до 4 МБ на одну строку

	lines := make([]string, 0, maxLines)

	var files []string

	flush := func() error {
		if len(lines) == 0 {
			return nil
		}
		sort.SliceStable(lines, func(i, j int) bool {
			return less(lines[i], lines[j])
		})
		if unique && equal != nil {
			lines = uniqueInPlace(lines, equal)
		}
		file, err := os.CreateTemp(tempDir, "sort_chunk_*.txt")
		if err != nil {
			return err
		}
		w := bufio.NewWriter(file)
		for _, line := range lines {
			if _, err := w.WriteString(line); err != nil {
				_ = file.Close()
				return err
			}
			if err := w.WriteByte('\n'); err != nil {
				_ = file.Close()
				return err
			}
		}
		if err := w.Flush(); err != nil {
			_ = file.Close()
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}

		files = append(files, file.Name())
		lines = lines[:0]
		return nil
	}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) == maxLines {
			if err := flush(); err != nil {
				return nil, err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := flush(); err != nil {
		return nil, err
	}

	// Если вход был пуст, то создаем пустой файл для унификации обработки
	if len(files) == 0 {
		file, err := os.CreateTemp(tempDir, "sort_empty_*.txt")
		if err != nil {
			return nil, err
		}
		_ = file.Close()
		files = append(files, file.Name())
	}
	return files, nil
}

// uniqueInPlace удаляет дублирующиеся строки в уже отсортированном срезе.
func uniqueInPlace(a []string, equal func(a, b string) bool) []string {
	if len(a) == 0 {
		return a
	}
	w := 1
	for i := 1; i < len(a); i++ {
		if !equal(a[i], a[w-1]) {
			a[w] = a[i]
			w++
		}
	}
	return a[:w]
}
