package unpackstr_test

import (
	"errors"
	"strings"
	"testing"

	unpackStr "github.com/Gunvolt24/wb_l2/L2.9/unpackstr"
)

func TestUnpackString(t *testing.T) {
	tests := []struct {
		input, output string
	}{
		{input: `a4bc2d5e`, output: `aaaabccddddde`},
		{input: `abcd`, output: `abcd`},
		{input: `45`, output: "error=строка состоит только из цифр"},
		{input: ``, output: ``},
		{input: `qwe\4\5`, output: `qwe45`},
		{input: `qwe\45`, output: `qwe44444`},
	}

	for _, test := range tests {
		res, err := unpackStr.UnpackString(test.input)

		// Убираем префикс "error=" и проверяем, что ошибка не nil и соответствует ожидаемой
		if wantErr, ok := strings.CutPrefix(test.output, "error="); ok {
			if err == nil {
				t.Fatalf("expected error %q got nil", wantErr)
			}
			if !errors.Is(err, unpackStr.ErrAllDigits) {
				t.Errorf("input=%q | expected error=%q | got error=%q", test.input, wantErr, err.Error())
			}
			continue
		}

		if err != nil {
			t.Errorf("input=%q | error=%v\n", test.input, err)
		} else if res != test.output {
			t.Errorf("input=%q | output=%q | expected=%q\n", test.input, res, test.output)
		}
	}

}
