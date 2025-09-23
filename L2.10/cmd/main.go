package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Gunvolt24/wb_l2/L2.10/internal/checker"
	chunker "github.com/Gunvolt24/wb_l2/L2.10/internal/chunk"
	"github.com/Gunvolt24/wb_l2/L2.10/internal/comparator"
	inputoutput "github.com/Gunvolt24/wb_l2/L2.10/internal/inputOutput"
	"github.com/Gunvolt24/wb_l2/L2.10/internal/merger"
	"github.com/Gunvolt24/wb_l2/L2.10/internal/options"
)

func main() {
	// Установка опций командной строки и парсинг аргументов
	var (
		k      = flag.Int("k", 0, "колонка (1..N), 0 - вся строка")
		n      = flag.Bool("n", false, "числовая сортировка")
		r      = flag.Bool("r", false, "обратный порядок")
		u      = flag.Bool("u", false, "только уникальные строки по ключу")
		M      = flag.Bool("M", false, "месяцы (Jan..Dec)")
		b      = flag.Bool("b", false, "игнорировать хвостовые пробелы при сравнении")
		c      = flag.Bool("c", false, "только проверить отсортированность и выйти")
		h      = flag.Bool("h", false, "человекочитаемые размеры (1K 2M ...)")
		d      = flag.String("d", "\t", "разделитель колонок")
		chunk  = flag.Int("chunk", 1_000, "максимум строк в одном чанке перед сбросом на диск")
		tmpdir = flag.String("tmpdir", os.TempDir(), "директория для временных файлов")
	)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] [FILE|-]\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// Открытие входного потока
	var in io.Reader = os.Stdin
	if flag.NArg() > 0 && flag.Arg(0) != "-" {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "open:", err)
			os.Exit(1)
		}
		defer file.Close()
		in = file
	}

	// Установка опций
	opts := options.Options{
		Column:        *k,
		NumSort:       *n,
		Reverse:       *r,
		Unique:        *u,
		MonthNames:    *M,
		TrailBlanks:   *b,
		HumanReadable: *h,
		Splitter:      *d,
	}

	less, equal := comparator.CreateComparator(opts)

	// Режим проверки (-c): читаем поток и проверяем на лету.
	if *c {
		ok, i, j, err := checker.CheckSorted(in, less)
		if err != nil {
			fmt.Fprintln(os.Stderr, "read:", err)
			os.Exit(1)
		}
		if !ok {
			fmt.Fprintf(os.Stderr, "check: data not sorted at lines %d and %d\n", i+1, j+1)
			os.Exit(1)
		}
		return
	}

	// Этап 1: нарезка входа на чанки, сортировка, запись во временные файлы.
	paths, err := chunker.SplitSort(in, *chunk, less, equal, *tmpdir, opts.Unique)
	if err != nil {
		fmt.Fprintln(os.Stderr, "split/sort:", err)
		os.Exit(1)
	}
	defer func() {
		for _, p := range paths {
			_ = os.Remove(p)
		}
	}()

	// Быстрый путь: один файл - просто выводим.
	if len(paths) == 1 {
		if err := inputoutput.CopyFileToStdout(paths[0]); err != nil {
			fmt.Fprintln(os.Stderr, "copy:", err)
			os.Exit(1)
		}
		return
	}

	// Этап 2: k-way merge всех временных файлов в stdout.
	if err := merger.KWayMerge(paths, less, equal, os.Stdout, opts.Unique); err != nil {
		fmt.Fprintln(os.Stderr, "merge:", err)
		os.Exit(1)
	}

}
