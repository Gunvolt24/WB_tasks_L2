package parser

import "fmt"

// ParseRedirection - парсит аргументы команды, выделяя редиректы.
// Возвращает "чистые" аргументы (без >, >>, <), входной файл (если есть), выходной файл (если есть),
// режим добавления (true - >>, false - >), и ошибку (если синтаксис неверен).
func ParseRedirection(args []string) (pure []string, inFile, outFile string, appendMode bool, err error) {
	pure = make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case ">":
			if i+1 >= len(args) {
				return nil, "", "", false, fmt.Errorf("syntax error: > FILE expected")
			}
			outFile = args[i+1]
			appendMode = false
			i++
		case ">>":
			if i+1 >= len(args) {
				return nil, "", "", false, fmt.Errorf("syntax error: >> FILE expected")
			}
			outFile = args[i+1]
			appendMode = true
			i++
		case "<":
			if i+1 >= len(args) {
				return nil, "", "", false, fmt.Errorf("syntax error: < FILE expected")
			}
			inFile = args[i+1]
			i++
		default:
			pure = append(pure, arg)
		}
	}

	return

}
