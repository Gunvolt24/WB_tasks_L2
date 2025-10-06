package runcmdctx

import (
	"context"
	"fmt"
	"strings"
)

// runAndOrCmdCtx - выполняет команды, разделенные по токенам "&&" и "||"
func runAndOrCmdCtx(ctx context.Context, line string) error {
	tokens := strings.Fields(line) // разбиваем строку на токены по пробелам
	var segments [][]string        // сегменты команд между && и ||
	var ops []string               // операторы && и ||
	current := []string{}          // текущий сегмент команды

	for _, token := range tokens {
		if token == "&&" || token == "||" {
			if len(current) == 0 {
				return fmt.Errorf("syntax error near unexpected token `%s`", token)
			}
			segments = append(segments, current)
			ops = append(ops, token)
			current = []string{}
			continue
		}
		current = append(current, token)
	}
	if len(current) > 0 {
		segments = append(segments, current)
	}

	var lastErr error
	for i, segment := range segments {
		cmdStr := strings.Join(segment, " ")
		lastErr = RunCmdCtx(ctx, cmdStr) // допускаем пайпы/редиректы/builtin внутри сегмента

		if i < len(ops) {
			if ops[i] == "&&" && lastErr != nil {
				break // если команда завершилась с ошибкой и оператор "&&" - прерываем выполнение следующих команд
			}
			if ops[i] == "||" && lastErr == nil {
				break // если команда завершилась успешно и оператор "||" - прерываем выполнение следующих команд
			}
		}
	}
	return lastErr
}
