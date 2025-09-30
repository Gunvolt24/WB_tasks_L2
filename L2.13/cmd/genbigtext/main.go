package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	var (
		n      = flag.Int("n", 100, "количество строк")
		out    = flag.String("o", "testdata/big_cut.tsv", "путь к выходному файлу")
		del    = flag.String("d", "\\t", "разделитель полей (например, \"\\t\" или \",\")")
		seed   = flag.Int64("seed", time.Now().UnixNano(), "seed для генератора случайных чисел")
		plainP = flag.Float64("plain", 0.15, "доля строк без разделителя (plain)")
		twoP   = flag.Float64("two", 0.20, "доля строк с 2 полями (id, name)")
		fourP  = flag.Float64("four", 0.25, "доля строк с 4 полями (id, name, email, ext)")
		dupP   = flag.Float64("dup", 0.05, "доля дубликатов (повтор предыдущей строки)")
	)
	flag.Parse()

	// Небольшая декодировка популярных escape-последовательностей в -d
	delimiter := decodeDelim(*del)

	// Создаём каталог назначения (если он указан в пути)
	_ = os.MkdirAll(filepath.Dir(*out), 0o755)

	rnd := rand.New(rand.NewSource(*seed))

	f, err := os.Create(*out)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	var prev string
	for i := 1; i <= *n; i++ {
		// Иногда делаем точный дубль предыдущей строки (для проверки -s, стабильности и т.д.)
		if i > 1 && rnd.Float64() < *dupP {
			if prev != "" {
				_, _ = w.WriteString(prev + "\n")
				continue
			}
		}

		// Строка без разделителя?
		if rnd.Float64() < *plainP {
			row := "plain_line_without_delimiter"
			_, _ = w.WriteString(row + "\n")
			prev = row
			continue
		}

		// Выбираем количество полей: 2, 3 (по умолчанию), или 4.
		// Приоритетный выбор: сначала four, затем two, иначе three.
		fieldsCount := 3
		r := rnd.Float64()
		switch {
		case r < *fourP:
			fieldsCount = 4
		case r < *fourP+*twoP:
			fieldsCount = 2
		default:
			fieldsCount = 3
		}

		// Формируем поля
		id := fmt.Sprintf("id%03d", i)
		name := fmt.Sprintf("name%03d", i)
		email := fmt.Sprintf("email%03d@example.com", i)
		ext := fmt.Sprintf("ext%03d", i)

		var parts []string
		switch fieldsCount {
		case 2:
			parts = []string{id, name}
		case 3:
			parts = []string{id, name, email}
		default: // 4
			parts = []string{id, name, email, ext}
		}

		row := strings.Join(parts, delimiter)
		_, _ = w.WriteString(row + "\n")
		prev = row
	}

	// На всякий случай «дожмём» буфер (defer уже делает Flush, но пусть будет)
	_ = w.Flush()
}

// decodeDelim преобразует строки вида "\t" в настоящий символ табуляции.
func decodeDelim(s string) string {
	switch s {
	case `\t`:
		return "\t"
	case `\n`:
		return "\n"
	case `\r`:
		return "\r"
	default:
		return s
	}
}
