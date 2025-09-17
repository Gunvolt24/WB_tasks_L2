package unpackstr

import (
	"errors"
	"strings"
	"unicode"
)

// ErrAllDigits ошибка - строка состоит только из цифр
var ErrAllDigits = errors.New("строка состоит только из цифр")

// ErrTrailingSlash ошибка - недопустимый символ после экранирования
var ErrTrailingSlash = errors.New("недопустимый символ после экранирования")

// UnpackString распаковывает строку.
// Поддерживает экранирование символов с помощью '\'.
// Возвращает ошибку в случаях некорректности (например, все символы - цифры).
func UnpackString(s string) (string, error) {
	runes := []rune(s)
	var result []rune
	var escape []bool

	// Снимаем экранирование: удаляем '\', следующий символ помечаем как экранированный
	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' {
			// Если есть следуюший символ - добавляем его в результат и помечаем как экранированный
			if i+1 < len(runes) {
				result = append(result, runes[i+1])
				escape = append(escape, true)
				i++ // пропускаем следующий символ, т.к. он уже добавлен
			} else {
				return "", ErrTrailingSlash
			}
			// Если символ не экранированный - добавляем его в результат
		} else {
			result = append(result, runes[i])
			escape = append(escape, false)
		}
	}

	// Если после снятия экранирования строка пустая - возвращаем пустую строку
	if len(result) == 0 {
		return "", nil
	}

	// Если все символы - цифры - возвращаем ошибку
	allDigits := true
	for i := 0; i < len(result); i++ {
		if !unicode.IsDigit(result[i]) {
			allDigits = false
			break
		}
	}
	if allDigits {
		return "", ErrAllDigits
	}

	// Распаковываем строку
	var b strings.Builder
	for i := 0; i < len(result); {
		char := result[i]
		// Если следующий символ - цифра - умножаем текущий символ на количество повторений
		if i+1 < len(result) && unicode.IsDigit(result[i+1]) && !escape[i+1] {
			// Получаем количество повторений
			num := int(result[i+1] - '0')
			if num > 0 {
				b.WriteString(strings.Repeat(string(char), num))
			}
			// Пропускаем количество повторений
			i += 2
		} else {
			// Если следующий символ - не цифра - добавляем его в результат
			b.WriteRune(char)
			i++
		}
	}

	return b.String(), nil
}
