package pipelines

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/Gunvolt24/wb_l2/L2.15/internal/runCmdCtx/parser"
)

// RunPipelineCtx - выполняет пайплайны только с внешними командами.
// Поддерживает редиректы: '<' у первой команды, '>'/ '>>' у последней.
func RunPipelineCtx(ctx context.Context, line string) error {
	parts := strings.Split(line, "|")
	n := len(parts)
	if n == 0 {
		return nil
	}

	cmds := make([]*exec.Cmd, n)
	var prevStdout io.ReadCloser

	// файлы для закрытия
	var inFiles []*os.File
	var outFile *os.File

	for i := 0; i < n; i++ {
		segment := strings.TrimSpace(parts[i])
		if segment == "" {
			return fmt.Errorf("empty command in pipeline")
		}

		args := strings.Fields(segment)
		pure, inF, outF, appendMode, err := parser.ParseRedirection(args)
		if err != nil {
			return err
		}
		if len(pure) == 0 {
			return fmt.Errorf("empty command in pipeline segment %d", i)
		}
		// ограничим: входной редирект - только у первой команды, выходной - только у последней
		if i != 0 && inF != "" {
			return fmt.Errorf("input redirect (<) allowed only on first pipeline stage")
		}
		if i != n-1 && outF != "" {
			return fmt.Errorf("output redirect (>/>>) allowed only on last pipeline stage")
		}

		cmd := exec.CommandContext(ctx, pure[0], pure[1:]...)
		cmd.Stderr = os.Stderr

		// stdin
		if i == 0 {
			if inF != "" {
				f, err := os.Open(inF)
				if err != nil {
					return err
				}
				inFiles = append(inFiles, f)
				cmd.Stdin = f
			} else {
				cmd.Stdin = os.Stdin
			}
		} else {
			cmd.Stdin = prevStdout
		}

		// stdout
		if i < n-1 {
			stdoutPipe, err := cmd.StdoutPipe()
			if err != nil {
				// закрыть уже открытые файлы
				for _, f := range inFiles {
					_ = f.Close()
				}
				return err
			}
			prevStdout = stdoutPipe
		} else {
			// последняя команда
			if outF != "" {
				flags := os.O_CREATE | os.O_WRONLY
				if appendMode {
					flags |= os.O_APPEND
				} else {
					flags |= os.O_TRUNC
				}
				f, err := os.OpenFile(outF, flags, 0o644)
				if err != nil {
					for _, f2 := range inFiles {
						_ = f2.Close()
					}
					return err
				}
				outFile = f
				cmd.Stdout = outFile
			} else {
				cmd.Stdout = os.Stdout
			}
		}

		cmds[i] = cmd
	}

	// стартуем всё
	for _, c := range cmds {
		if err := c.Start(); err != nil {
			if outFile != nil {
				_ = outFile.Close()
			}
			for _, f := range inFiles {
				_ = f.Close()
			}
			return err
		}
	}

	// ждём всё
	var firstErr error
	for _, c := range cmds {
		if err := c.Wait(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	if outFile != nil {
		_ = outFile.Close()
	}
	for _, f := range inFiles {
		_ = f.Close()
	}
	return firstErr
}
