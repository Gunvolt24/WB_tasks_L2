package fields_test

import (
	"reflect"
	"testing"

	"github.com/Gunvolt24/wb_l2/L2.13/internal/dto"
	"github.com/Gunvolt24/wb_l2/L2.13/internal/fields"
)

func TestDefineFields_OK(t *testing.T) {
	got, err := fields.DefineFields("1, 3-5, -2, 7-")
	if err != nil {
		t.Fatalf("DefineFields error: %v", err)
	}
	// Открытые границы мы не сравниваем до бесконечности; проверим первые интервалы и форму.
	if len(got) != 4 {
		t.Fatalf("len=%d, want 4", len(got))
	}
	want0 := dto.Fields{Start: 1, End: 1}
	want1 := dto.Fields{Start: 3, End: 5}
	want2 := dto.Fields{Start: 1, End: 2} // "-2"
	if !reflect.DeepEqual(got[0], want0) || !reflect.DeepEqual(got[1], want1) || !reflect.DeepEqual(got[2], want2) {
		t.Fatalf("prefix mismatch: got=%v", got[:3])
	}
	// 4-й — открытый вправо "7-" (End должен быть >= 7)
	if got[3].Start != 7 || got[3].End < 7 {
		t.Fatalf("bad open-right: %+v", got[3])
	}
}

func TestDefineFields_Errors(t *testing.T) {
	cases := []string{
		"",    // обязателен -f
		"0",   // номера с 1
		"-0",  // правая граница >=1
		"3-1", // from > to
		"abc", // не число
		"1-a", // правая граница не число
	}
	for _, spec := range cases {
		if _, err := fields.DefineFields(spec); err == nil {
			t.Fatalf("want error for %q", spec)
		}
	}
}
