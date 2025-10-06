package e2e

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func buildBin(t *testing.T) string {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	root := filepath.Dir(filepath.Dir(thisFile)) // .../L2.15

	// убедимся, что каталог bin существует
	binDir := filepath.Join(root, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("mkdir bin: %v", err)
	}

	// добавить .exe на Windows
	name := "minish_e2e"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	bin := filepath.Join(binDir, name)

	cmd := exec.Command("go", "build", "-o", bin, "./cmd")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build error: %v\n%s", err, out)
	}
	return bin
}

func runShellWithInput(t *testing.T, bin string, workdir string, script string, timeout time.Duration) (stdout, stderr string, code int) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin)
	cmd.Dir = workdir

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	cmd.Stdin = strings.NewReader(script)

	err := cmd.Run()

	stdout = outBuf.String()
	stderr = errBuf.String()
	code = 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else if ctx.Err() == context.DeadlineExceeded {
			t.Fatalf("shell timed out; stdout:\n%s\nstderr:\n%s", stdout, stderr)
		} else {
			t.Fatalf("shell run error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout, stderr)
		}
	}
	return
}

func Test_EchoAndRedirections(t *testing.T) {
	bin := buildBin(t)
	tmp := t.TempDir()

	script := "" +
		"echo one > a.txt\n" +
		"echo two >> a.txt\n" +
		"pwd > p.txt\n"

	_, _, _ = runShellWithInput(t, bin, tmp, script, 5*time.Second)

	// Проверяем содержимое файлов
	aBytes, err := os.ReadFile(filepath.Join(tmp, "a.txt"))
	if err != nil {
		t.Fatalf("read a.txt: %v", err)
	}
	gotA := strings.ReplaceAll(string(aBytes), "\r\n", "\n")
	wantA := "one\n" + "two\n"
	if gotA != wantA {
		t.Fatalf("a.txt mismatch:\n--- got ---\n%q\n--- want ---\n%q", gotA, wantA)
	}

	// pwd > p.txt должен записать каталог tmp
	pBytes, err := os.ReadFile(filepath.Join(tmp, "p.txt"))
	if err != nil {
		t.Fatalf("read p.txt: %v", err)
	}
	gotP := strings.TrimSpace(strings.ReplaceAll(string(pBytes), "\r\n", "\n"))
	if gotP != tmp {
		t.Fatalf("pwd mismatch: got %q want %q", gotP, tmp)
	}
}

func Test_CdAndPwd(t *testing.T) {
	bin := buildBin(t)
	tmp := t.TempDir()
	sub := filepath.Join(tmp, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}

	script := "" +
		"pwd > p1.txt\n" +
		"cd sub\n" +
		"pwd > p2.txt\n"

	_, _, _ = runShellWithInput(t, bin, tmp, script, 5*time.Second)

	// p1.txt в корне tmp
	p1Path := filepath.Join(tmp, "p1.txt")
	p1Bytes, err := os.ReadFile(p1Path)
	if err != nil {
		t.Fatalf("read %s: %v", p1Path, err)
	}
	got1 := strings.TrimSpace(strings.ReplaceAll(string(p1Bytes), "\r\n", "\n"))
	if got1 != tmp {
		t.Fatalf("p1 mismatch: got %q want %q", got1, tmp)
	}

	// p2.txt уже в каталоге sub (после cd sub)
	p2Path := filepath.Join(sub, "p2.txt")
	p2Bytes, err := os.ReadFile(p2Path)
	if err != nil {
		t.Fatalf("read %s: %v", p2Path, err)
	}
	got2 := strings.TrimSpace(strings.ReplaceAll(string(p2Bytes), "\r\n", "\n"))
	if got2 != sub {
		t.Fatalf("p2 mismatch: got %q want %q", got2, sub)
	}
}

func Test_AndOr(t *testing.T) {
	bin := buildBin(t)
	tmp := t.TempDir()

	var script string
	if runtime.GOOS == "windows" {
		script = "" +
			"cmd /c exit 1 && echo bad > and.txt\n" +
			"cmd /c exit 0 || echo bad2 > or.txt\n" +
			"echo ok > keep.txt\n"
	} else {
		script = "" +
			"false && echo bad > and.txt\n" +
			"true || echo bad2 > or.txt\n" +
			"echo ok > keep.txt\n"
	}

	_, _, _ = runShellWithInput(t, bin, tmp, script, 5*time.Second)

	if _, err := os.Stat(filepath.Join(tmp, "and.txt")); !os.IsNotExist(err) {
		t.Fatalf("and.txt must NOT exist")
	}
	if _, err := os.Stat(filepath.Join(tmp, "or.txt")); !os.IsNotExist(err) {
		t.Fatalf("or.txt must NOT exist")
	}
	if _, err := os.Stat(filepath.Join(tmp, "keep.txt")); err != nil {
		t.Fatalf("keep.txt must exist: %v", err)
	}
}

func Test_Pipeline(t *testing.T) {
	bin := buildBin(t)
	tmp := t.TempDir()

	var script string
	if runtime.GOOS == "windows" {
		// только внешние команды
		script = "cmd /c echo foo | findstr foo > pipe.txt\n"
	} else {
		script = "/bin/echo foo | grep foo > pipe.txt\n"
	}

	_, _, _ = runShellWithInput(t, bin, tmp, script, 5*time.Second)

	data, err := os.ReadFile(filepath.Join(tmp, "pipe.txt"))
	if err != nil {
		t.Fatalf("read pipe.txt: %v", err)
	}
	got := strings.TrimSpace(strings.ReplaceAll(string(data), "\r\n", "\n"))
	if got != "foo" {
		t.Fatalf("pipeline output mismatch: got %q want %q", got, "foo")
	}
}
