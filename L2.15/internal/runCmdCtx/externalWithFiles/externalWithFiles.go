package externalwithfiles

import (
	"context"
	"os"
	"os/exec"
)

// RunExternalWithFiles - выполняет внешнюю команду с возможностью указания файлов для ввода и вывода
func RunExternalWithFiles(ctx context.Context, name string, args []string, inFile, outFile string, appendMode bool) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	var (
		in, out *os.File
		err     error
	)

	if inFile != "" {
		in, err = os.Open(inFile)
		if err != nil {
			return err
		}
		defer in.Close()
		cmd.Stdin = in
	}

	if outFile != "" {
		flag := os.O_CREATE | os.O_WRONLY
		if appendMode {
			flag |= os.O_APPEND
		} else {
			flag |= os.O_TRUNC
		}
		out, err = os.OpenFile(outFile, flag, 0o644)
		if err != nil {
			return err
		}
		defer out.Close()
		cmd.Stdout = out
	}

	return cmd.Run()
}
