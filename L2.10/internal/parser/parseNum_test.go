package parser

import "testing"

func TestCmpNumSort(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"2", "10", -1},
		{"-5", "-3", -1},
		{"3.14", "3.140", 0},
		{"  7 ", "7", 0},
		{"x", "10", -1},
		{"10", "x", 1},
		{"foo", "bar", stringsCompare("foo", "bar")},
	}
	for i, tt := range tests {
		if got := CmpNumSort(tt.a, tt.b); got != tt.want {
			t.Fatalf("case %d: CmpNumSort(%q,%q) = %d, want %d", i, tt.a, tt.b, got, tt.want)
		}
	}
}
