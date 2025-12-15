package saver

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// sanitizeFilename делает имя файла валидным для Windows
func sanitizeFilename(name string) string {
	reg := regexp.MustCompile(`[<>:"|?*]`)
	return reg.ReplaceAllString(name, "_")
}

// SaveResource сохраняет данные ресурса по URL в локальную папку
func SaveResource(rawurl string, baseDir string, data []byte) error {
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	// Разбиваем путь на части и очищаем имена
	parts := strings.Split(u.Path, "/")
	for i, p := range parts {
		parts[i] = sanitizeFilename(p)
	}

	// Если путь пустой или "/", используем домен как имя файла
	var localPath string
	if u.Path == "" || u.Path == "/" || filepath.Ext(u.Path) == "" {
		localPath = filepath.Join(baseDir, u.Host)
		// Добавляем имя файла: домен + .html
		localPath = filepath.Join(localPath, sanitizeFilename(u.Host)+".html")
	} else {
		localPath = filepath.Join(baseDir, u.Host, filepath.Join(parts...))
	}

	// Создаём папки
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("cannot create directory %s: %w", dir, err)
	}

	// Сохраняем файл
	f, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return err
	}

	fmt.Println("Saved:", localPath)
	return nil
}
