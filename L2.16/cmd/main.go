package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Gunvolt24/wb_l2/L2.16/internal/crawler"
)

func main() {

	// Парсим флаги командной строки
	pUrl := flag.String("url", "", "URL to crawl")
	pDepth := flag.Int("depth", 1, "Crawl depth")
	pConcurrency := flag.Int("concurrency", 5, "Number of parallel downloads")
	flag.Parse()

	// Проверяем обязательный параметр URL
	if *pUrl == "" {
		fmt.Fprintln(os.Stderr, "Error: empty URL")
		return
	}

	fmt.Printf("Starting crawl of %s with depth %d\n", *pUrl, *pDepth)

	// Создаём и запускаем краулер
	crawl := crawler.NewCrawler(*pConcurrency, "download")
	crawl.Start(*pUrl, *pDepth)
}
