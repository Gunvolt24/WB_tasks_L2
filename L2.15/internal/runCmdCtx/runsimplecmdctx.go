package runcmdctx

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	externalwithfiles "github.com/Gunvolt24/wb_l2/L2.15/internal/runCmdCtx/externalWithFiles"
	"github.com/Gunvolt24/wb_l2/L2.15/internal/runCmdCtx/helpers"
	"github.com/Gunvolt24/wb_l2/L2.15/internal/runCmdCtx/parser"
)

// runSimpleCmdCtx — одна команда без пайпов, с редиректами и builtin-командами.
func runSimpleCmdCtx(ctx context.Context, cmdStr string) error {
	cmdStr = strings.TrimSpace(cmdStr)
	if cmdStr == "" {
		return nil
	}
	args := strings.Fields(cmdStr)
	if len(args) == 0 {
		return nil
	}

	// 1) Сначала разбираем редиректы: получаем "чистые" аргументы и файлы
	pure, inFile, outFile, appendMode, err := parser.ParseRedirection(args)
	if err != nil {
		return err
	}

	// 2) Если после удаления редиректов команда отсутствует (например, "> out.txt")
	if len(pure) == 0 {
		if outFile != "" {
			flags := os.O_CREATE | os.O_WRONLY
			if appendMode {
				flags |= os.O_APPEND
			} else {
				flags |= os.O_TRUNC
			}
			f, err := os.OpenFile(outFile, flags, 0o644)
			if err != nil {
				return err
			}
			_ = f.Close()
		}
		return nil
	}

	// 3) Дальше работаем только с "чистыми" аргументами (без >, >>, <)
	switch pure[0] {
	case "exit":
		os.Exit(0)
		return nil

	case "cd":
		if len(pure) < 2 {
			homeDir, _ := os.UserHomeDir()
			cdPath(homeDir, nil)
		} else {
			cdPath(pure[1], nil)
		}
		return nil

	case "pwd":
		// поддержка редиректов STDOUT/STDIN для builtin
		return helpers.WithRedirForBuiltin(inFile, outFile, appendMode, func() { pwd() })

	case "echo":
		return helpers.WithRedirForBuiltin(inFile, outFile, appendMode, func() {
			if len(pure) > 1 {
				echo(pure[1:])
			} else {
				echo(nil)
			}
		})

	case "kill":
		if len(pure) < 2 {
			return fmt.Errorf("kill: not enough arguments")
		}
		return helpers.WithRedirForBuiltin(inFile, outFile, appendMode, func() { kill(pure[1]) })

	case "ps":
		// твой ps вызывает внешнюю утилиту, но редирект stdout тоже поддержим
		return helpers.WithRedirForBuiltin(inFile, outFile, appendMode, func() { ps() })

	default:
		// 4) Внешняя команда — редиректы уже распарсены
		return externalwithfiles.RunExternalWithFiles(ctx, pure[0], pure[1:], inFile, outFile, appendMode)
	}
}

// === builtin команды ===

func cdPath(path string, err error) {
	if path == "" {
		if homeDir, err1 := os.UserHomeDir(); err1 == nil {
			path = homeDir
		}
	}
	err = os.Chdir(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cd:", err)
	}
}

func pwd() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "pwd:", err)
		return
	}
	fmt.Println(dir)
}

func echo(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func kill(pidStr string) {
	pidN, err := strconv.Atoi(pidStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "kill:", err)
		return
	}
	proc, err := os.FindProcess(pidN)
	if err != nil {
		fmt.Fprintln(os.Stderr, "kill:", err)
		return
	}
	if err := proc.Kill(); err != nil {
		fmt.Fprintln(os.Stderr, "kill:", err)
	}
}

func ps() {
	name := "ps"
	args := []string{}
	if runtime.GOOS == "windows" {
		name = "tasklist"
	}
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "ps:", err)
	}
}
