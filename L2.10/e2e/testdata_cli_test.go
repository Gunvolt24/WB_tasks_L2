package e2e_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir := filepath.Dir(curFile())
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found; run tests from repo root")
		}
		dir = parent
	}
}

func curFile() string {
	_, file, _, _ := runtime.Caller(0)
	return file
}

func buildBinary(t *testing.T) string {
	t.Helper()
	root := findRepoRoot(t)
	out := filepath.Join(t.TempDir(), "sort")
	if runtime.GOOS == "windows" {
		out += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", out, "./cmd")
	cmd.Dir = root
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build failed: %v\n%s", err, b)
	}
	return out
}

func run(bin string, args ...string) (stdout, stderr string, code int) {
	cmd := exec.Command(bin, args...)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	code = 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else {
			code = -1
		}
	}
	return out.String(), errb.String(), code
}

func testdata(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(findRepoRoot(t), "testdata", name)
}

func assertOrder(t *testing.T, got string, want []string) {
	t.Helper()
	g := strings.Split(strings.TrimSpace(got), "\n")
	if len(g) != len(want) {
		t.Fatalf("len mismatch: got %d, want %d\nGOT:\n%s\nWANT:\n%s",
			len(g), len(want), strings.Join(g, "\n"), strings.Join(want, "\n"))
	}
	for i := range want {
		if g[i] != want[i] {
			t.Fatalf("at %d: %q != %q\nGOT:\n%s\nWANT:\n%s",
				i, g[i], want[i], strings.Join(g, "\n"), strings.Join(want, "\n"))
		}
	}
}

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

func Test_MonthsSizes_K2_M_And_Human_K3_h_r(t *testing.T) {
	bin := buildBinary(t)
	path := testdata(t, "months_sizes.txt")

	// По месяцам (возрастание) по 2-му столбцу
	stdout, stderr, code := run(bin, "-k", "2", "-M", path)
	if code != 0 || stderr != "" {
		t.Fatalf("sort -k2 -M failed: code=%d, stderr=%q", code, stderr)
	}
	assertOrder(t, stdout, []string{
		"id2\tJan\t512",
		"id4\tJan\t1K",
		"id1\tFeb\t1K",
		"id3\tMar\t2K",
		"id5\tApr\t900M",
	})

	// Человекочитаемая сортировка по числовому значению с учетом суффиксов (убывание) по 3-ему столбцу
	stdout, stderr, code = run(bin, "-k", "3", "-h", "-r", path)
	if code != 0 || stderr != "" {
		t.Fatalf("sort -k3 -h -r failed: code=%d, stderr=%q", code, stderr)
	}
	assertOrder(t, stdout, []string{
		"id5\tApr\t900M",
		"id3\tMar\t2K",
		"id1\tFeb\t1K",
		"id4\tJan\t1K",
		"id2\tJan\t512",
	})

	// Проверка -c по месяцам
	stdout2, _, _ := run(bin, "-k", "2", "-M", path)
	tmp := filepath.Join(t.TempDir(), "out.txt")
	if err := os.WriteFile(tmp, []byte(stdout2), 0o600); err != nil {
		t.Fatalf("write out: %v", err)
	}
	_, stderr, code = run(bin, "-k", "2", "-M", "-c", tmp)
	if code != 0 || stderr != "" {
		t.Fatalf("check -c failed: code=%d, stderr=%q", code, stderr)
	}
}

func Test_Numbers_n_and_r(t *testing.T) {
	bin := buildBinary(t)
	path := testdata(t, "numbers.txt") // один столбец

	stdout, stderr, code := run(bin, "-n", "-r", path)
	if code != 0 || stderr != "" {
		t.Fatalf("sort -n -r failed: code=%d, stderr=%q", code, stderr)
	}
	assertOrder(t, stdout, []string{
		"100",
		"10",
		"3.14",
		"2",
		"0",
		"-5",
	})

	// И проверка -c для возрастания
	stdout, _, _ = run(bin, "-n", path)
	tmp := filepath.Join(t.TempDir(), "asc.txt")
	_ = os.WriteFile(tmp, []byte(stdout), 0o600)
	_, stderr, code = run(bin, "-n", "-c", tmp)
	if code != 0 || stderr != "" {
		t.Fatalf("check -c for numbers failed: code=%d, stderr=%q", code, stderr)
	}
}

func Test_DupsByCol2_u(t *testing.T) {
	bin := buildBinary(t)
	path := testdata(t, "dups_by_col2.txt") // uN<TAB>a/b<TAB>size

	stdout, stderr, code := run(bin, "-k", "2", "-u", path)
	if code != 0 || stderr != "" {
		t.Fatalf("sort -k2 -u failed: code=%d, stderr=%q", code, stderr)
	}
	// Ожидаем по одному для 'a' и 'b' (стабильно остаётся первый)
	assertOrder(t, stdout, []string{
		"u1\ta\t10K",
		"u3\tb\t10K",
	})
}

func Test_TrailingBlanks_b_u(t *testing.T) {
	bin := buildBinary(t)
	path := testdata(t, "trailing_blanks.txt") // x<TAB>bar[spaces/tabs]

	stdout, stderr, code := run(bin, "-k", "2", "-b", "-u", path)
	if code != 0 || stderr != "" {
		t.Fatalf("sort -k2 -b -u failed: code=%d, stderr=%q", code, stderr)
	}
	if strings.TrimSpace(stdout) != "x\tbar" {
		t.Fatalf("expected single 'x<TAB>bar', got %q", stdout)
	}
}

func Test_MissingCols_K3(t *testing.T) {
	bin := buildBinary(t)
	path := testdata(t, "missing_cols.txt") // без 3-его столбца

	stdout, stderr, code := run(bin, "-k", "3", path)
	if code != 0 || stderr != "" {
		t.Fatalf("sort -k3 failed: code=%d, stderr=%q", code, stderr)
	}
	// Поскольку у первых двух строк ключ "", стабильная сортировка сохранит их исходный порядок.
	assertOrder(t, stdout, []string{
		"only-one",
		"two\tcols",
		"three\tcolumns\there",
	})

	// Проверка -c
	tmp := filepath.Join(t.TempDir(), "mc.txt")
	_ = os.WriteFile(tmp, []byte(stdout), 0o600)
	_, stderr, code = run(bin, "-k", "3", "-c", tmp)
	if code != 0 || stderr != "" {
		t.Fatalf("check -c for missing_cols failed: code=%d, stderr=%q", code, stderr)
	}
}

func Test_MixedNonNumeric_K2_n_SelfCheckAndMultiset(t *testing.T) {
	bin := buildBinary(t)
	path := testdata(t, "mixed_non_numeric.txt") // a<TAB>foo; b<TAB>10; c<TAB>x20; d<TAB>-3

	// Сортируем -n по 2-ему столбцу и проверяем, что -c проходит по тому же правилу.
	stdout, stderr, code := run(bin, "-k", "2", "-n", path)
	if code != 0 || stderr != "" {
		t.Fatalf("sort -k2 -n failed: code=%d, stderr=%q", code, stderr)
	}
	tmp := filepath.Join(t.TempDir(), "mix.txt")
	_ = os.WriteFile(tmp, []byte(stdout), 0o600)
	_, stderr, code = run(bin, "-k", "2", "-n", "-c", tmp)
	if code != 0 || stderr != "" {
		t.Fatalf("self-check -c failed: code=%d, stderr=%q", code, stderr)
	}

	// Проверка, что мультимножество строк не изменилось.
	origBytes, _ := os.ReadFile(path)
	orig := strings.Split(strings.TrimSpace(string(origBytes)), "\n")
	got := strings.Split(strings.TrimSpace(stdout), "\n")
	if !sameMultiset(orig, got) {
		t.Fatalf("multiset mismatch for mixed_non_numeric.txt\ngor: %v\nwant: %v", got, orig)
	}
}
