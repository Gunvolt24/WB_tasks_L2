package inputoutput

import (
	"io"
	"os"
)

// CopyFileToStdout - копирует содержимое файла в stdout
func CopyFileToStdout(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(os.Stdout, file)
	return err
}
