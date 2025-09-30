package cut

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	checkinterval "github.com/Gunvolt24/wb_l2/L2.13/internal/checkInterval"
	"github.com/Gunvolt24/wb_l2/L2.13/internal/config"
	"github.com/Gunvolt24/wb_l2/L2.13/internal/fields"
)

// Cut - функция для обрезки полей
func Cut(r io.Reader, w io.Writer, cfg config.Config) error {
	// если флаг "-d" не задан, то устанавливаем табуляцию
	if cfg.Delimiter == "" {
		cfg.Delimiter = "\t"
	}

	// определяем диапазоны
	intervals, err := fields.DefineFields(cfg.Fields)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(r)      // создаем сканнер
	buffer := make([]byte, 64*1024)     // устанавливаем буфер
	scanner.Buffer(buffer, 4*1024*1024) // до 4 МБ на одну строку

	// проходим по каждой строке
	for scanner.Scan() {
		line := scanner.Text()

		// если строка не содержит флаг "-d", то выводим ее
		if !strings.Contains(line, cfg.Delimiter) {
			if cfg.Separator {
				continue
			}

			fmt.Fprintln(w, line)
			continue
		}

		// разбиваем строку на поля
		fields := strings.Split(line, cfg.Delimiter)
		// формируем новую строку
		out := make([]string, 0, len(fields))

		// проходим по каждому полю и проверяем, находится ли он в диапазоне
		for i := 1; i <= len(fields); i++ {
			if checkinterval.IsInInterval(i, intervals) {
				out = append(out, fields[i-1])
			}
		}

		// выводим строку
		fmt.Fprintln(w, strings.Join(out, cfg.Delimiter))
	}

	return scanner.Err()
}
