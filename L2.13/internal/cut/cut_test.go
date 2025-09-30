package cut_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Gunvolt24/wb_l2/L2.13/internal/config"
	"github.com/Gunvolt24/wb_l2/L2.13/internal/cut"
)

func TestCut_Tab_Default_NoS(t *testing.T) {
	in := "a\tb\tc\nd\te\tf\nno_delimiter\nx\ty\n"
	cfg := config.Config{Fields: "1,3-5", Delimiter: "\t", Separator: false}

	var out bytes.Buffer
	if err := cut.Cut(strings.NewReader(in), &out, cfg); err != nil {
		t.Fatalf("cut: %v", err)
	}
	want := "a\tc\nd\tf\nno_delimiter\nx\n"
	if out.String() != want {
		t.Fatalf("got:\n%q\nwant:\n%q", out.String(), want)
	}
}

func TestCut_Tab_SkipWithoutDelimiter(t *testing.T) {
	in := "a\tb\tc\nplain\nx\ty\n"
	cfg := config.Config{Fields: "1,3-5", Delimiter: "\t", Separator: true}

	var out bytes.Buffer
	if err := cut.Cut(strings.NewReader(in), &out, cfg); err != nil {
		t.Fatalf("cut: %v", err)
	}
	want := "a\tc\nx\n"
	if out.String() != want {
		t.Fatalf("got:\n%q\nwant:\n%q", out.String(), want)
	}
}

func TestCut_OpenRanges(t *testing.T) {
	in := "id\tname\temail\np\tq\tr\ts\n"
	// "-2" -> первые два поля; "3-" -> от третьего до конца
	cfg1 := config.Config{Fields: "-2", Delimiter: "\t"}
	cfg2 := config.Config{Fields: "3-", Delimiter: "\t"}

	var out1, out2 bytes.Buffer
	if err := cut.Cut(strings.NewReader(in), &out1, cfg1); err != nil {
		t.Fatal(err)
	}
	if err := cut.Cut(strings.NewReader(in), &out2, cfg2); err != nil {
		t.Fatal(err)
	}
	want1 := "id\tname\np\tq\n"
	want2 := "email\nr\ts\n"
	if out1.String() != want1 {
		t.Fatalf("open-left got:\n%q\nwant:\n%q", out1.String(), want1)
	}
	if out2.String() != want2 {
		t.Fatalf("open-right got:\n%q\nwant:\n%q", out2.String(), want2)
	}
}

func TestCut_CustomDelimiterCSV(t *testing.T) {
	in := "a,b,c\n1,2,3\nno_delim_line\n"
	cfg := config.Config{Fields: "-2", Delimiter: ",", Separator: false}

	var out bytes.Buffer
	if err := cut.Cut(strings.NewReader(in), &out, cfg); err != nil {
		t.Fatal(err)
	}
	want := "a,b\n1,2\nno_delim_line\n"
	if out.String() != want {
		t.Fatalf("got:\n%q\nwant:\n%q", out.String(), want)
	}
}
