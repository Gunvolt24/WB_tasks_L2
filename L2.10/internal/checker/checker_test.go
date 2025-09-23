package checker_test

import (
	"bytes"
	"testing"

	"github.com/Gunvolt24/wb_l2/L2.10/internal/checker"
	"github.com/Gunvolt24/wb_l2/L2.10/internal/comparator"
	"github.com/Gunvolt24/wb_l2/L2.10/internal/options"
)

func TestCheckSorted_OK_Lexicographic(t *testing.T) {
	less, _ := comparator.CreateComparator(options.Options{}) // по всей строке, лексикографика
	ok, i, j, err := checker.CheckSorted(bytes.NewBufferString("a\nb\nc\n"), less)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected sorted, got not sorted at %d,%d", i, j)
	}
	if i != -1 || j != -1 {
		t.Fatalf("expected i=j=-1 on ok=true, got i=%d j=%d", i, j)
	}
}

func TestCheckSorted_Fail_Lexicographic(t *testing.T) {
	less, _ := comparator.CreateComparator(options.Options{})
	ok, i, j, err := checker.CheckSorted(bytes.NewBufferString("a\nc\nb\n"), less)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("expected not sorted, got ok=true")
	}
	// нарушение между 2-й (c) и 3-й (b) строками → i=1, j=2 (0-based)
	if i != 1 || j != 2 {
		t.Fatalf("expected i=1 j=2, got i=%d j=%d", i, j)
	}
}

func TestCheckSorted_EmptyInput(t *testing.T) {
	less, _ := comparator.CreateComparator(options.Options{})
	ok, i, j, err := checker.CheckSorted(bytes.NewBuffer(nil), less)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok || i != -1 || j != -1 {
		t.Fatalf("expected ok with i=j=-1, got ok=%v i=%d j=%d", ok, i, j)
	}
}

func TestCheckSorted_Months_K2_WithSplitter(t *testing.T) {
	// Проверяем режим -k 2 -M с пробельным разделителем.
	less, _ := comparator.CreateComparator(options.Options{
		Column: 2, MonthNames: true, Splitter: " ",
	})
	// Отсортированный по месяцу поток
	ok, _, _, err := checker.CheckSorted(bytes.NewBufferString(
		"id2 Jan 512\nid4 Jan 1K\nid1 Feb 1K\nid3 Mar 2K\nid5 Apr 900M\n",
	), less)
	if err != nil || !ok {
		t.Fatalf("expected sorted months, err=%v ok=%v", err, ok)
	}

	// Неотсортированный: Jan, Mar, Feb -> нарушение между Mar (2) и Feb (3)
	ok, i, j, err := checker.CheckSorted(bytes.NewBufferString(
		"id2 Jan 512\nid3 Mar 2K\nid1 Feb 1K\n",
	), less)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok || !(i == 1 && j == 2) {
		t.Fatalf("expected fail at lines 2/3 (0-based 1/2), got ok=%v i=%d j=%d", ok, i, j)
	}
}
