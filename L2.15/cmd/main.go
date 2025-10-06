package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"

	runcmdctx "github.com/Gunvolt24/wb_l2/L2.15/internal/runCmdCtx"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		dir, _ := os.Getwd()
		fmt.Printf("$ \x1b[32m%s\x1b[0m> ", dir)

		line, err := reader.ReadString('\n')
		// EOF / Ctrl+D - выходим
		if err != nil {
			fmt.Println("\x1b[1;33m\n[+] Goodbye!\x1b[0m")
			return
		}

		// контекст с отменой по сигналу прерывания (Ctrl+C)
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		err = runcmdctx.RunCmdCtx(ctx, line)
		cancel()

		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
	}
}
