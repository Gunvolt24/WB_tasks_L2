package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	unpackstr "github.com/Gunvolt24/wb_l2/L2.9/unpackstr"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите строку: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка чтения: %v\n", err)
		os.Exit(1)
	}
	input = strings.TrimRight(input, "\r\n") // кроссплатформенно удаляем перенос строки

	result, err := unpackstr.UnpackString(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Распакованная строка: %q\n", result)
}
