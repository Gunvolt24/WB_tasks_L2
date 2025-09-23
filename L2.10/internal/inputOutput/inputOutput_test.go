package inputoutput_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	inputoutput "github.com/Gunvolt24/wb_l2/L2.10/internal/inputOutput"
)

func TestCopyFileToStdout_Success(t *testing.T) {
	tmp := t.TempDir()
	want := "line1\nline2\n"
	path := filepath.Join(tmp, "data.txt")
	if err := os.WriteFile(path, []byte(want), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	// Перехватываем stdout.
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w

	// Вызов
	callErr := inputoutput.CopyFileToStdout(path)

	// Закрытие и сбор вывода
	_ = w.Close()
	var got bytes.Buffer
	_, _ = io.Copy(&got, r)
	_ = r.Close()
	os.Stdout = oldStdout

	if callErr != nil {
		t.Fatalf("CopyFileToStdout error: %v", callErr)
	}
	if got.String() != want {
		t.Fatalf("stdout mismatch:\n got: %q\nwant: %q", got.String(), want)
	}
}

func TestCopyFileToStdout_FileNotFound(t *testing.T) {
	err := inputoutput.CopyFileToStdout("no/such/file.txt")
	if err == nil {
		t.Fatalf("expected error for missing file, got nil")
	}
}
