package downloader

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Download делает GET запрос и возвращает тело ответа
func Download(url string) ([]byte, error) {
	fmt.Println("Downloading", url)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad status code %d for %s", resp.StatusCode, url)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
