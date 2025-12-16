package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {

	// Задаем флаг для установки таймаута
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "Timeout for the operation")
	flag.Parse()

	log.SetOutput(os.Stderr)
	log.SetFlags(0)
	if len(flag.Args()) < 2 {
		log.Fatalf("Usage: %s [--timeout duration] <host> <port>\n", os.Args[0])
	}

	// Создаем соединение host:port из аргументов
	addr := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %v\n", addr, err)
	}

	log.Printf("Connected to %s\n", addr)

	// Создаем каналы для завершения
	doneRead := make(chan struct{})
	doneWrite := make(chan struct{})

	// Запускаем горутину для чтения из сокета
	go func() {
		reader := bufio.NewReader(conn)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Printf("Error reading from connection: %v\n", err)
				}
				break
			}

			os.Stdout.WriteString(line)
		}
		log.Println("Connection closed from server side")
		close(doneRead)
	}()

	// Запускаем горутину для записи в сокет
	go func() {
		stdin := bufio.NewReader(os.Stdin)
		for {
			line, err := stdin.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Printf("Error reading from stdin: %v\n", err)
				}
				break
			}

			_, err = conn.Write([]byte(line))
			if err != nil {
				log.Printf("Error writing to connection: %v\n", err)
				break
			}
		}

		log.Println("Connection closed from client side")
		conn.Close()
		close(doneWrite)
	}()

	select {
	case <-doneRead:
	case <-doneWrite:
	}

	<-doneRead
	<-doneWrite
}
