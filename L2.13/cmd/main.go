package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Gunvolt24/wb_l2/L2.13/internal/config"
	"github.com/Gunvolt24/wb_l2/L2.13/internal/cut"
)

func main() {
	// устанавливаем флаги
	var (
		fields    = flag.String("f", "", "список полей")
		delimeter = flag.String("d", "\t", "разделитель полей")
		separator = flag.Bool("s", false, "выводить только строки, содержащие разделитель")
	)

	// устанавливаем обработчик ошибок
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s -f FIELDS [-d DELIM] [-s] [FILE|-]\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	// парсим флаги
	flag.Parse()

	// определяем конфигурацию
	cfg := config.Config{
		Fields:    *fields,
		Delimiter: *delimeter,
		Separator: *separator,
	}

	// если флаг "-f" не задан, то выводим ошибку
	if cfg.Fields == "" {
		fmt.Fprintln(os.Stderr, "error: missing fields")
		os.Exit(2)
	}

	// открытие входного потока
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

	// если флаг "-f" задан, то вызываем функцию
	if err := cut.Cut(in, os.Stdout, cfg); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
