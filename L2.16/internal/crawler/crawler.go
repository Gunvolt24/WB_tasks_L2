package crawler

import (
	"fmt"
	"sync"

	"github.com/Gunvolt24/wb_l2/L2.16/internal/downloader"
	"github.com/Gunvolt24/wb_l2/L2.16/internal/parser"
	"github.com/Gunvolt24/wb_l2/L2.16/internal/saver"
)

// Crawler хранит состояние рекурсивного скачивания
type Crawler struct {
	visitedPages     map[string]bool
	visitedResources map[string]bool
	mu               sync.Mutex
	wg               sync.WaitGroup
	sema             chan struct{}
	baseDir          string
}

// NewCrawler создаёт новый Crawler
func NewCrawler(concurrency int, baseDir string) *Crawler {
	return &Crawler{
		visitedPages:     make(map[string]bool),
		visitedResources: make(map[string]bool),
		sema:             make(chan struct{}, concurrency),
		baseDir:          baseDir,
	}
}

// Start запускает скачивание с указанного URL и глубины
func (c *Crawler) Start(url string, depth int) {
	c.crawl(url, depth)
	c.wg.Wait()
}

// crawl рекурсивно скачивает страницу и ресурсы
func (c *Crawler) crawl(pageURL string, depth int) {
	if depth <= 0 {
		return
	}

	c.mu.Lock()
	if c.visitedPages[pageURL] {
		c.mu.Unlock()
		return
	}
	c.visitedPages[pageURL] = true
	c.mu.Unlock()

	c.wg.Add(1)
	go func(url string, d int) {
		defer c.wg.Done()
		c.sema <- struct{}{}
		defer func() { <-c.sema }()

		fmt.Println("Crawling:", url)

		body, err := downloader.Download(url)
		if err != nil {
			fmt.Println("Error downloading page:", url, err)
			return
		}

		if err := saver.SaveResource(url, c.baseDir, body); err != nil {
			fmt.Println("Error saving page:", url, err)
		}

		links, resources, err := parser.ExtractLinksAndResources(body, url)
		if err != nil {
			fmt.Println("Error parsing page:", url, err)
			return
		}

		// Скачиваем ресурсы параллельно
		for _, res := range resources {
			c.mu.Lock()
			if c.visitedResources[res] {
				c.mu.Unlock()
				continue
			}
			c.visitedResources[res] = true
			c.mu.Unlock()

			c.wg.Add(1)
			go func(r string) {
				defer c.wg.Done()
				c.sema <- struct{}{}
				defer func() { <-c.sema }()

				data, err := downloader.Download(r)
				if err != nil {
					fmt.Println("Error downloading resource:", r, err)
					return
				}
				if err := saver.SaveResource(r, c.baseDir, data); err != nil {
					fmt.Println("Error saving resource:", r, err)
				}
			}(res)
		}

		// Рекурсивно обрабатываем ссылки
		for _, link := range links {
			c.crawl(link, d-1)
		}
	}(pageURL, depth)
}
