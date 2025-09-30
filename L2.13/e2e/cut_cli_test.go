// L2.13/e2e/cut_cli_test.go
package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func buildBin(t *testing.T) string {
	t.Helper()

	// Находим каталог L2.13
	_, thisFile, _, _ := runtime.Caller(0)
	l213 := filepath.Dir(thisFile) // .../L2.13/e2e
	l213 = filepath.Dir(l213)      // .../L2.13

	// Куда собирать бинарь — во временную директорию теста
	outDir := t.TempDir()
	name := "cut_e2e"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	bin := filepath.Join(outDir, name)

	// Сборка
	cmd := exec.Command("go", "build", "-o", bin, "./cmd")
	cmd.Dir = l213
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	// На всякий случай проверим, что файл реально есть
	if _, err := os.Stat(bin); err != nil {
		t.Fatalf("built binary not found: %v", err)
	}
	return bin
}

func TestCLI_Tab_Default(t *testing.T) {
	bin := buildBin(t)

	input := "a\tb\tc\nd\te\tf\nplain\nx\ty\n"
	tmp := t.TempDir()
	inPath := filepath.Join(tmp, "in.tsv")
	if err := os.WriteFile(inPath, []byte(input), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bin, "-f", "1,3-5", inPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run failed: %v\n%s", err, out)
	}
	want := "a\tc\nd\tf\nplain\nx\n"
	if string(out) != want {
		t.Fatalf("got:\n%q\nwant:\n%q", out, want)
	}
}

func TestCLI_SkipWithoutDelimiter(t *testing.T) {
	bin := buildBin(t)

	input := "a\tb\tc\nplain\nx\ty\n"
	tmp := t.TempDir()
	inPath := filepath.Join(tmp, "in.tsv")
	if err := os.WriteFile(inPath, []byte(input), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bin, "-s", "-f", "1,3-5", inPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run failed: %v\n%s", err, out)
	}
	want := "a\tc\nx\n"
	if string(out) != want {
		t.Fatalf("got:\n%q\nwant:\n%q", out, want)
	}
}
