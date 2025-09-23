package comparator_test

import (
	"testing"

	"github.com/Gunvolt24/wb_l2/L2.10/internal/comparator"
	"github.com/Gunvolt24/wb_l2/L2.10/internal/options"
)

func TestCreateComparator_Lexicographic_ByK2(t *testing.T) {
	// Две колонки, разделитель — пробел.
	less, eq := comparator.CreateComparator(options.Options{
		Column: 2, Splitter: " ",
	})
	a := "id1 aa"
	b := "id2 bb"
	if !less(a, b) {
		t.Fatalf("expected aa < bb")
	}
	if eq(a, b) {
		t.Fatalf("expected keys not equal")
	}
}

func TestCreateComparator_Months_K2(t *testing.T) {
	// Проверяем -M по второй колонке.
	less, _ := comparator.CreateComparator(options.Options{
		Column: 2, MonthNames: true, Splitter: " ",
	})
	if !less("id2 Jan 512", "id1 Feb 1K") {
		t.Fatalf("expected Jan < Feb")
	}
	if !less("id1 Feb 1K", "id3 Mar 2K") {
		t.Fatalf("expected Feb < Mar")
	}
}

func TestCreateComparator_Human_K3_Reverse(t *testing.T) {
	// -h по третьей колонке, обратный порядок.
	less, _ := comparator.CreateComparator(options.Options{
		Column: 3, HumanReadable: true, Reverse: true, Splitter: " ",
	})
	// В reverse «больше» должно идти раньше.
	if !less("idX X 2G", "idY Y 900M") {
		t.Fatalf("reverse human: 2G should precede 900M")
	}
}

func TestCreateComparator_Numeric_K1(t *testing.T) {
	// Числовая сортировка по первой колонке.
	less, _ := comparator.CreateComparator(options.Options{
		Column: 1, NumSort: true, Splitter: " ",
	})
	if !less("2 b", "10 a") {
		t.Fatalf("numeric: 2 < 10")
	}
	if less("10 b", "2 a") {
		t.Fatalf("numeric: 10 !< 2")
	}
}

func TestCreateComparator_TrailingBlanks_B_Equal(t *testing.T) {
	// -b: игнор хвостовых пробелов в ключе, проверим equal (для -u).
	_, eq := comparator.CreateComparator(options.Options{
		Column: 2, TrailBlanks: true, Splitter: " ",
	})
	if !eq("x bar   ", "x bar") {
		t.Fatalf("trailing blanks must be ignored for key equality")
	}
}

func TestCreateComparator_MissingColumn_IsEmptyKey(t *testing.T) {
	// Нет 3-й колонки -> ключ = "", должен быть «меньше» непустых.
	less, _ := comparator.CreateComparator(options.Options{
		Column: 3, Splitter: " ",
	})
	a := "only-two cols"
	b := "three cols here"
	if !less(a, b) {
		t.Fatalf("missing key must sort before present key")
	}
}
