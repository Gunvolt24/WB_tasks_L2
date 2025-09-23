package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// утилита для генерации большого текстового файла
func main() {
	var (
		n     = flag.Int("n", 1_000, "количество строк")
		out   = flag.String("o", "testdata/big.txt", "путь к выходному файлу")
		seed  = flag.Int64("seed", time.Now().UnixNano(), "seed для PRNG")
		missP = flag.Float64("miss", 0.01, "доля строк с пропущенными колонками")
		dupP  = flag.Float64("dup", 0.02, "доля дубликатов (повтор предыдущей строки)")
	)
	flag.Parse()
	_ = os.MkdirAll("testdata", 0o755)

	rnd := rand.New(rand.NewSource(*seed))
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	hSuf := []string{"", "K", "M", "G", "T"}
	pad := func(s string) string {
		// случайные хвостовые пробелы для проверки -b
		return s + strings.Repeat(" ", rnd.Intn(4))
	}

	f, err := os.Create(*out)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	var prev string
	for i := 0; i < *n; i++ {
		// иногда делаем точный дубль предыдущей строки (для -u, merge-стыков и т.п.)
		if i > 0 && rnd.Float64() < *dupP {
			_, _ = w.WriteString(prev + "\n")
			continue
		}

		id := fmt.Sprintf("id%07d", i)
		mon := months[rnd.Intn(len(months))]
		// human-size: либо целое, либо с .5; с разными суффиксами
		base := rnd.Intn(1800) + 1 // 1..1800
		val := fmt.Sprintf("%d", base)
		if rnd.Intn(10) == 0 { // иногда дробные
			val = fmt.Sprintf("%d.5", base)
		}
		size := val + hSuf[rnd.Intn(len(hSuf))]

		// числовая колонка (для -n), может быть отрицательной
		num := rnd.Intn(2_000_000) - 1_000_000 // [-1e6 .. +1e6]

		// текстовый ключ для -b (с хвостовыми пробелами)
		key := pad(fmt.Sprintf("key%03d", rnd.Intn(200)))

		// иногда «ломаем» строки, пропуская 3-ю или 5-ю колонку (для -k N + missing)
		row := ""
		switch {
		case rnd.Float64() < *missP:
			// пропустим 3-ю колонку (human-size)
			row = fmt.Sprintf("%s\t%s\t\t%d\t%s", id, mon, num, key)
		case rnd.Float64() < *missP:
			// пропустим 5-ю колонку (key)
			row = fmt.Sprintf("%s\t%s\t%s\t%d", id, mon, size, num)
		default:
			row = fmt.Sprintf("%s\t%s\t%s\t%d\t%s", id, mon, size, num, key)
		}

		_, _ = w.WriteString(row + "\n")
		prev = row
	}
	_ = w.Flush()
}
