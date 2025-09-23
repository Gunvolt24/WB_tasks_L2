package merger_test

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Gunvolt24/wb_l2/L2.10/internal/comparator"
	"github.com/Gunvolt24/wb_l2/L2.10/internal/merger"
	"github.com/Gunvolt24/wb_l2/L2.10/internal/options"
)

func TestKWayMerge_DedupAcrossChunks(t *testing.T) {
	tmp := t.TempDir()
	p1 := writeFile(t, tmp, "c1.txt", []string{"a", "c", "e"})
	p2 := writeFile(t, tmp, "c2.txt", []string{"b", "c", "d"})

	less, eq := comparator.CreateComparator(options.Options{})
	var out strings.Builder
	if err := merger.KWayMerge([]string{p1, p2}, less, eq, &out, true); err != nil {
		t.Fatalf("KWayMerge error: %v", err)
	}

	got := strings.Split(strings.TrimSpace(out.String()), "\n")
	want := []string{"a", "b", "c", "d", "e"}
	if len(got) != len(want) {
		t.Fatalf("len mismatch: %v vs %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("at %d: %q != %q", i, got[i], want[i])
		}
	}
}

func TestKWayMerge_HumanSizes_ByK3(t *testing.T) {
	tmp := t.TempDir()
	// Оба файла индивидуально отсортированы по human-возрастанию.
	p1 := writeFile(t, tmp, "h1.txt", []string{
		"id2 Jan 512",
		"id4 Jan 1K",
	})
	p2 := writeFile(t, tmp, "h2.txt", []string{
		"id1 Feb 1K",
		"id3 Mar 2K",
		"id5 Apr 900M",
	})

	less, eq := comparator.CreateComparator(options.Options{
		Column: 3, HumanReadable: true, Splitter: " ",
	})

	var out strings.Builder
	if err := merger.KWayMerge([]string{p1, p2}, less, eq, &out, false); err != nil {
		t.Fatalf("KWayMerge error: %v", err)
	}

	got := strings.Split(strings.TrimSpace(out.String()), "\n")
	want := []string{
		"id2 Jan 512",
		"id4 Jan 1K",
		"id1 Feb 1K",
		"id3 Mar 2K",
		"id5 Apr 900M",
	}

	// 1) Проверяем, что got отсортирован по тому же less (нестрого неубывание).
	for i := 1; i < len(got); i++ {
		if less(got[i], got[i-1]) {
			t.Fatalf("not sorted at %d: %q < %q", i, got[i], got[i-1])
		}
	}

	// 2) Проверяем мультимножество (без учёта порядка).
	if !sameMultiset(got, want) {
		t.Fatalf("multiset mismatch:\n got:  %v\n want: %v", got, want)
	}
}

// --- функции-помощники ---

func sameMultiset(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	m := map[string]int{}
	for _, s := range a {
		m[s]++
	}
	for _, s := range b {
		m[s]--
		if m[s] < 0 {
			return false
		}
	}
	for _, v := range m {
		if v != 0 {
			return false
		}
	}
	return true
}

func writeFile(t *testing.T, dir, name string, lines []string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	f, err := os.Create(p)
	if err != nil {
		t.Fatalf("create %s: %v", p, err)
	}
	w := bufio.NewWriter(f)
	for _, s := range lines {
		_, _ = w.WriteString(s)
		_ = w.WriteByte('\n')
	}
	_ = w.Flush()
	_ = f.Close()
	return p
}
