package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

func main() {
	// Запрашиваем точное время с NTP-сервера
	curTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		// Выводим ошибку в stderr и возвращаем ненулевой код
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	// Выводим UTC и локальное точное время
	fmt.Println("UTC time:", curTime.UTC().Format("2006-01-02 15:04:05"))
	fmt.Println("Local time:", curTime.Local().Format("2006-01-02 15:04:05"))
}
