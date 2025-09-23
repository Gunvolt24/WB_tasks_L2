package chunk_test

import (
	"bufio"
	"bytes"
	"os"
	"sort"
	"testing"

	"github.com/Gunvolt24/wb_l2/L2.10/internal/chunk"
	"github.com/Gunvolt24/wb_l2/L2.10/internal/comparator"
	"github.com/Gunvolt24/wb_l2/L2.10/internal/options"
)

func TestSplitSort_SplitsSortsAndDedups(t *testing.T) {
	t.Helper()
	tmp := t.TempDir()

	// Вход: намеренно неотсортированный + дубликаты.
	input := "b\nc\na\na\nz\ny\nx\n"
	maxLines := 2 // заставим сделать много чанков

	less, eq := comparator.CreateComparator(options.Options{}) // лексикографика по всей строке
	paths, err := chunk.SplitSort(bytes.NewBufferString(input), maxLines, less, eq, tmp, true)
	if err != nil {
		t.Fatalf("SplitSort error: %v", err)
	}
	// Удалим созданные файлы после теста.
	t.Cleanup(func() {
		for _, p := range paths {
			_ = os.Remove(p)
		}
	})

	if len(paths) < 3 {
		t.Fatalf("expected >= 3 chunk files, got %d", len(paths))
	}

	// Проверим каждый чанк: строки отсортированы и без дублей.
	for _, p := range paths {
		lines, err := readLines(p)
		if err != nil {
			t.Fatalf("read %s: %v", p, err)
		}
		if !isSorted(lines, less) {
			t.Fatalf("chunk %s is not sorted: %v", p, lines)
		}
		if hasAdjacentDup(lines) {
			t.Fatalf("chunk %s still has duplicates: %v", p, lines)
		}
	}
}

func TestSplitSort_EmptyInput_CreatesEmptyFile(t *testing.T) {
	tmp := t.TempDir()
	less, eq := comparator.CreateComparator(options.Options{})
	paths, err := chunk.SplitSort(bytes.NewBuffer(nil), 10, less, eq, tmp, true)
	if err != nil {
		t.Fatalf("SplitSort error: %v", err)
	}
	t.Cleanup(func() {
		for _, p := range paths {
			_ = os.Remove(p)
		}
	})
	if len(paths) != 1 {
		t.Fatalf("expected 1 empty file, got %d", len(paths))
	}
	info, err := os.Stat(paths[0])
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Size() != 0 {
		t.Fatalf("empty chunk file must have size=0, got %d", info.Size())
	}
}

// --- функции-помощники ---

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		out = append(out, sc.Text())
	}
	return out, sc.Err()
}

func isSorted(lines []string, less func(a, b string) bool) bool {
	return sort.SliceIsSorted(lines, func(i, j int) bool { return !less(lines[j], lines[i]) })
}

func hasAdjacentDup(lines []string) bool {
	for i := 1; i < len(lines); i++ {
		if lines[i] == lines[i-1] {
			return true
		}
	}
	return false
}
