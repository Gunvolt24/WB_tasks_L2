package main

import (
	"bytes"
	"strings"
	"testing"
)

func runGrep(t *testing.T, input string, c config) (out string, matches int, err error) {
	t.Helper()
	var buf bytes.Buffer
	m, e := grep(strings.NewReader(input), &buf, c)
	return buf.String(), m, e
}

func TestFlagF_FixedSubstring(t *testing.T) {
	in := "a.c\naXc\nabc\n"
	c := config{fixed: true, pattern: "a.c"}
	out, matches, err := runGrep(t, in, c)
	if err != nil {
		t.Fatal(err)
	}
	want := "a.c\n"
	if out != want || matches != 1 {
		t.Fatalf("got out=%q matches=%d; want out=%q matches=1", out, matches, want)
	}
}

func TestRegex_Default(t *testing.T) {
	in := "a.c\naXc\nabc\n"
	c := config{fixed: false, pattern: "a.c"}
	out, matches, err := runGrep(t, in, c)
	if err != nil {
		t.Fatal(err)
	}
	wantLines := []string{"a.c", "aXc", "abc", ""}
	if out != strings.Join(wantLines, "\n") || matches != 3 {
		t.Fatalf("unexpected out=%q matches=%d", out, matches)
	}
}

func TestIgnoreCase_i(t *testing.T) {
	in := "foo\nMATCH here\nbar\n"
	c := config{ignore: true, pattern: "match"}
	out, matches, err := runGrep(t, in, c)
	if err != nil {
		t.Fatal(err)
	}
	want := "MATCH here\n"
	if out != want || matches != 1 {
		t.Fatalf("out=%q matches=%d", out, matches)
	}
}

func TestInvert_v(t *testing.T) {
	in := "ok\nerror now\nfine\n"
	c := config{invert: true, fixed: true, pattern: "error"}
	out, matches, err := runGrep(t, in, c)
	if err != nil {
		t.Fatal(err)
	}
	want := "ok\nfine\n"
	if out != want || matches != 2 {
		t.Fatalf("out=%q matches=%d", out, matches)
	}
}

func TestNumber_n(t *testing.T) {
	in := "alpha\nbeta\nmatch line\ngamma\n"
	c := config{number: true, fixed: true, pattern: "match"}
	out, matches, err := runGrep(t, in, c)
	if err != nil {
		t.Fatal(err)
	}
	want := "3:match line\n"
	if out != want || matches != 1 {
		t.Fatalf("out=%q matches=%d", out, matches)
	}
}

func TestCount_c(t *testing.T) {
	in := "x\nmatch\ny\nmatch\nz\n"
	c := config{count: true, fixed: true, pattern: "match"}
	out, matches, err := runGrep(t, in, c)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(out) != "2" || matches != 2 {
		t.Fatalf("out=%q matches=%d", out, matches)
	}
}

func TestAfter_A(t *testing.T) {
	in := "l1\nmatch\nl2\nl3\nl4\n"
	c := config{after: 2, fixed: true, pattern: "match"}
	out, matches, err := runGrep(t, in, c)
	if err != nil {
		t.Fatal(err)
	}
	want := "match\nl2\nl3\n"
	if out != want || matches != 1 {
		t.Fatalf("out=%q matches=%d", out, matches)
	}
}

func TestBefore_B(t *testing.T) {
	in := "p1\np2\np3\nmatch\np4\n"
	c := config{before: 2, fixed: true, pattern: "match"}
	out, matches, err := runGrep(t, in, c)
	if err != nil {
		t.Fatal(err)
	}
	want := "p2\np3\nmatch\n"
	if out != want || matches != 1 {
		t.Fatalf("out=%q matches=%d", out, matches)
	}
}

func TestAround_C(t *testing.T) {
	in := "l1\nm1\nl3\nm2\nl5\n"
	c := config{before: 1, after: 1, fixed: true, pattern: "m"}
	out, matches, err := runGrep(t, in, c)
	if err != nil {
		t.Fatal(err)
	}
	want := "l1\nm1\nl3\nm2\nl5\n"
	if out != want || matches != 2 {
		t.Fatalf("out=%q matches=%d", out, matches)
	}
}

func TestOverlappingAfter_NoDup(t *testing.T) {
	in := "m1\nl2\nm2\nl4\n"
	c := config{after: 2, fixed: true, pattern: "m"}
	out, matches, err := runGrep(t, in, c)
	if err != nil {
		t.Fatal(err)
	}
	want := "m1\nl2\nm2\nl4\n"
	if out != want || matches != 2 {
		t.Fatalf("out=%q matches=%d", out, matches)
	}
}

func TestInvalidRegex_Error(t *testing.T) {
	in := "anything\n"
	c := config{fixed: false, pattern: "("}
	_, _, err := runGrep(t, in, c)
	if err == nil {
		t.Fatal("expected error on invalid regex, got nil")
	}
}
