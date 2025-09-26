package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// config - конфигурация командной строки
// используется для передачи параметров в функцию grep
type config struct {
	after  int
	before int
	circle int

	count  bool
	ignore bool
	invert bool
	number bool
	fixed  bool

	pattern string
}

// createMatcher создает функцию, которая проверяет, есть ли в строке шаблон
func createMatcher(pattern string, fixed, ignore bool) (func(string) bool, error) {
	// если флаг -F установлен, то шаблон должен быть фиксированным
	if fixed {
		// если флаг -i установлен, то регистр игнорируется
		if ignore {
			pat := strings.ToLower(pattern)
			// возвращаем функцию, которая проверяет, есть ли в строке шаблон
			return func(s string) bool {
				return strings.Contains(strings.ToLower(s), pat)
			}, nil
		}
		// возвращаем функцию, которая проверяет, есть ли в строке шаблон
		return func(s string) bool {
			return strings.Contains(s, pattern)
		}, nil
	}

	// если флаг -i установлен, то паттерн должен игнорировать регистр
	if ignore {
		pattern = "(?i)" + pattern
	}

	// компилируем регулярное выражение
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("compile regexp: %w", err)
	}

	// возвращаем функцию, которая проверяет, есть ли в строке шаблон
	return re.MatchString, nil
}

// grep - реализация утилиты grep.
// Работа с флагами -c, -i, -v, -n, -F, -A, -B, -C
func grep(r io.Reader, w io.Writer, c config) (int, error) {
	// объявлем мэтчер, который проверяет, есть ли в строке шаблон согласно флагам -F, -i
	match, err := createMatcher(c.pattern, c.fixed, c.ignore)
	if err != nil {
		return 0, err
	}

	// создаем структуру для хранения данных предыдущих строк
	type record struct {
		num int
		str string
	}
	// создаем буфер для хранения предыдущих строк
	bufBefore := make([]record, 0, c.before)

	scanner := bufio.NewScanner(r)
	buf := make([]byte, 64*1024)     // устанавливаем буфер для сканирования
	scanner.Buffer(buf, 4*1024*1024) // установка максимального размера буфера

	lineN := 0        // номер текущей строки
	remainAfter := 0  // количество строк, которые остались после текущей строки
	lastPrinted := -1 // номер последней распечатанной строки
	matches := 0      // количество совпадений

	// функция для распечатки строки
	printLine := func(n int, s string) {
		// если режим -c, то распечатываем только количество совпадений
		if c.count {
			return // поскольку это режим -c - нам не нужно распечатывать строку
		}
		// если режим -n, то распечатываем номер строки и строку
		if c.number {
			fmt.Fprintf(w, "%d:%s\n", n+1, s)
		} else {
			// если режим не -n, то распечатываем только строку
			fmt.Fprintln(w, s)
		}
		lastPrinted = n // обновляем номер последней распечатанной строки
	}

	for scanner.Scan() {
		line := scanner.Text() // считываем строку

		// если режим -A установлен, то распечатываем количество строк, которые остались после текущей, без дублирования
		if remainAfter > 0 && lineN > lastPrinted {
			printLine(lineN, line)
			remainAfter--
		}

		// проверяем, есть ли совпадение
		ok := match(line)
		// если режим -v, то инвертируем результат
		if c.invert {
			ok = !ok
		}

		// если совпадение найдено - увеличиваем количество совпадений
		if ok {
			matches++

			// если флаг -B установлен, то распечатываем предыдущие строки из буфера, без дублирования
			if c.before > 0 && len(bufBefore) > 0 {
				start := 0
				if len(bufBefore) > c.before {
					start = len(bufBefore) - c.before
				}

				for _, r := range bufBefore[start:] {
					if r.num > lastPrinted {
						printLine(r.num, r.str)
					}
				}
			}

			// распечатываем текущую строку, если она еще не распечатана
			if lineN > lastPrinted {
				printLine(lineN, line)
			}

			// устанавливаем количество строк, которые остались после текущей
			if c.after > remainAfter {
				remainAfter = c.after
			}
		}

		// обновляем буфер предыдущих строк
		if c.before > 0 {
			if len(bufBefore) == c.before {
				bufBefore = bufBefore[1:]
			}
			bufBefore = append(bufBefore, record{num: lineN, str: line})
		}

		lineN++
	}

	if err := scanner.Err(); err != nil {
		return matches, err
	}

	// если флаг -c - печатаем количество совпадений
	if c.count {
		fmt.Fprintln(w, matches)
	}

	return matches, nil
}

func main() {
	// устанавливаем флаги
	var (
		countStr     = flag.Bool("c", false, "вывод количества строк, что совпадают с шаблоном")
		ignoreCase   = flag.Bool("i", false, "игнорировать регистр")
		invertFilter = flag.Bool("v", false, "выводить строки, не совпадающие с шаблоном")
		numStr       = flag.Bool("n", false, "выводить номер строки")
		fixedStr     = flag.Bool("F", false, "шаблон - фиксированная строка")
		afterStr     = flag.Int("A", 0, "выводить строки, следующие после шаблона")
		beforeStr    = flag.Int("B", 0, "выводить строки, предшествующие шаблону")
		circleStr    = flag.Int("C", 0, "выводить строки, встречающиеся вокруг шаблона")
	)

	// устанавливаем обработчик ошибок
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [FLAGS] PATTERN [FILE|-]\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	// парсим флаги
	flag.Parse()

	// позиционные аргументы: PATTERN [FILE]
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "error: missing pattern")
		os.Exit(2)
	}
	pattern := args[0]

	// Открытие входного потока
	var in io.Reader = os.Stdin
	if len(args) > 1 && args[1] != "-" {
		file, err := os.Open(args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "open:", err)
			os.Exit(1)
		}
		defer file.Close()
		in = file
	}

	// находим максимальное значение A и B и нормализуем C
	A, B := *afterStr, *beforeStr
	if *circleStr > A {
		A = *circleStr
	}
	if *circleStr > B {
		B = *circleStr
	}

	// устанавливаем конфигурацию
	conf := config{
		after:  A,
		before: B,
		circle: *circleStr,

		count:  *countStr,
		ignore: *ignoreCase,
		invert: *invertFilter,
		number: *numStr,
		fixed:  *fixedStr,

		pattern: pattern,
	}

	// выполняем поиск
	matches, err := grep(in, os.Stdout, conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}

	// код возврата: 0 - совпадение, 1 - не совпадение, 2 - ошибка
	if matches > 0 {
		os.Exit(0)
	}

	os.Exit(1)
}
