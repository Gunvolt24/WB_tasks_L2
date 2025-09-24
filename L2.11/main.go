package main

import (
	"fmt"
	"slices"
	"strings"
)

// sortRunes - функция для сортировки рун, чтобы они были в алфавитном порядке
func sortRunes(s string) string {
	runes := []rune(strings.ToLower(s))
	slices.Sort(runes)
	return string(runes)
}

// findAnagramms - функция для поиска анаграмм
func findAnagramms(input []string) map[string][]string {
	group := make(map[string][]string)  // map[signature]words - сигнатура - все слова с такой сигнатурой
	firstKey := make(map[string]string) // map[signature]firstWord - первое слово с такой сигнатурой

	// группируем слова по сигнатуре
	for _, word := range input {
		lowCase := strings.ToLower(word) // приводим к нижнему регистру
		sig := sortRunes(lowCase)        // сортируем буквы в слове

		// запоминаем первое слово с такой сигнатурой
		if _, seen := firstKey[sig]; !seen {
			firstKey[sig] = lowCase
		}

		// добавляем слово в группу
		group[sig] = append(group[sig], lowCase)
	}

	// собираем анаграммы
	result := make(map[string][]string)

	// проходим по всем группам
	for sig, words := range group {
		seen := make(map[string]struct{}, len(words)) // мапа для отслеживания уникальных слов
		unique := make([]string, 0, len(words))       // массив для хранения уникальных слов

		for _, word := range words {
			if _, ok := seen[word]; ok { // если слово уже встречалось - пропускаем
				continue
			}

			seen[word] = struct{}{}       // отмечаем слово как уникальное
			unique = append(unique, word) // добавляем слово в уникальные слова
		}

		// если уникальных слов меньше двух - пропускаем
		if len(unique) < 2 {
			continue
		}

		slices.Sort(unique)            // сортируем уникальные слова
		result[firstKey[sig]] = unique // добавляем уникальные слова в анаграммы
	}

	return result
}

func main() {
	input := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол",
		"карта", "тарка", "слово", "карат", "катар", "голанг", "клоун", "калун", "уклон", "кулон"}

	output := findAnagramms(input)

	// keys - сортированный список ключей
	keys := make([]string, 0, len(output))
	for key := range output {
		keys = append(keys, key) // добавляем ключ
	}
	slices.Sort(keys)

	for _, key := range keys {
		fmt.Printf("%s: %v\n", key, output[key])
	}
}
