package parser

import "testing"

func TestCmpMonth_BasicAndCase(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"Jan", "Feb", -1},
		{" feb ", "JAN", 1},
		{"Mar", "Mar", 0},
		{"Xxx", "Apr", -1},
		{"Apr", "Xxx", 1},
		{"X", "Y", stringsCompare("X", "Y")},
	}
	for i, tt := range tests {
		if got := CmpMonth(tt.a, tt.b); got != tt.want {
			t.Fatalf("case %d: CmpMonth(%q,%q) = %d, want %d", i, tt.a, tt.b, got, tt.want)
		}
	}
}
