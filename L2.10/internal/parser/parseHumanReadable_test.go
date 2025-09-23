package parser

import "testing"

func TestCmpHuman_Order(t *testing.T) {
	tests := []struct {
		a, b string
		want int // -1 if a<b, 0 if a==b, 1 if a>b
	}{
		{"512", "1K", -1},
		{"1K", "900M", -1},
		{"900M", "2G", -1},
		{"3.5G", "2G", 1},
		{"1K", "1024", 0},    // 1K == 1024
		{"  2m", "2048K", 0}, // регистр и пробелы
		{"foo", "bar", sign(stringsCompare("foo", "bar"))}, // оба нераспознаны → строковая политика
	}
	for i, tt := range tests {
		got := CmpHuman(tt.a, tt.b)
		if got != tt.want {
			t.Fatalf("case %d: CmpHuman(%q,%q) = %d, want %d", i, tt.a, tt.b, got, tt.want)
		}
	}
}

// stringsCompare и sign - функции-помощники, чтобы не тянуть strings.Compare прямо в таблицу.
func stringsCompare(a, b string) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
func sign(x int) int {
	switch {
	case x < 0:
		return -1
	case x > 0:
		return 1
	default:
		return 0
	}
}
