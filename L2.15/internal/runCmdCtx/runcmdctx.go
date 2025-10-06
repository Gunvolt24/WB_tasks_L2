package runcmdctx

import (
	"context"
	"os"
	"strings"

	"github.com/Gunvolt24/wb_l2/L2.15/internal/runCmdCtx/pipelines"
)

// RunCmdCtx - выполняет подстановку переменных окружения и выбор способа выполнения команды
// (простая команда, с пайпами или с && и ||)
// в качестве контекста используется контекст с отменой по сигналу прерывания (Ctrl+C)
func RunCmdCtx(ctx context.Context, cmdStr string) error {
	cmdStr = strings.TrimSpace(cmdStr)
	if cmdStr == "" {
		return nil
	}

	// подставляем переменные окружения
	cmdStr = os.ExpandEnv(cmdStr)

	// если есть пайпы - вызываем RunPipeCmdCtx
	if strings.Contains(cmdStr, "&&") || strings.Contains(cmdStr, "||") {
		return runAndOrCmdCtx(ctx, cmdStr)
	}

	// пайпы внутри команды
	if strings.Contains(cmdStr, "|") {
		return pipelines.RunPipelineCtx(ctx, cmdStr)
	}

	return runSimpleCmdCtx(ctx, cmdStr)
}
