package fields

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Gunvolt24/wb_l2/L2.13/internal/dto"
)

// задаем максимальное значение int
const maxInt = int(^uint(0) >> 1)

// DefineFields парсит спецификацию -f (например: "1,3-5,-2,7-")
// и возвращает список интервалов в 1-базной нумерации.
func DefineFields(flag string) ([]dto.Fields, error) {
	flag = strings.TrimSpace(flag) // обрезаем пробелы

	// если пользователь не указал диапазон, то выводим ошибку
	if flag == "" {
		return nil, fmt.Errorf("флаг -f обязателен")
	}

	// разбиваем строку на диапазоны
	parts := strings.Split(flag, ",") // [1, 2, 3]
	// создаем слайс диапазонов
	out := make([]dto.Fields, 0, len(parts))

	// проходим по каждому диапазону и проверяем его
	for _, part := range parts {
		part = strings.TrimSpace(part) // обрезаем пробелы

		// если диапазон пустой, то пропускаем
		if part == "" {
			continue
		}

		// если диапазон содержит "-", то это диапазон
		if strings.Contains(part, "-") {
			// Диапазон: A-B, -B (от 1 до B), A- (от A до бесконечности)
			sides := strings.SplitN(part, "-", 2)                                  // разбиваем на два элемента, например [1, 3]
			start, end := strings.TrimSpace(sides[0]), strings.TrimSpace(sides[1]) // определяем начало и конец и обрезаем пробелы

			var (
				startInt = 1      // начало диапазона
				endInt   = maxInt // конец диапазона
				err      error    // ошибка
			)

			// проверяем начало и конец
			if start != "" {
				startInt, err = strconv.Atoi(start)
				if err != nil || startInt < 1 {
					return nil, fmt.Errorf("не корректное начало диапазона: %q", part)
				}
			}

			if end != "" {
				endInt, err = strconv.Atoi(end)
				if err != nil || endInt < 1 {
					return nil, fmt.Errorf("не корректный конец диапазона: %q", part)
				}
			}

			// если начало больше конца, то возвращаем ошибку
			if startInt > endInt {
				return nil, fmt.Errorf("начало диапазона больше конца")
			}

			// формируем диапазон
			out = append(out, dto.Fields{Start: startInt, End: endInt})

		} else {
			// если диапазон не содержит "-", то это одно поле
			oneNum, err := strconv.Atoi(part)
			if err != nil || oneNum < 1 {
				return nil, fmt.Errorf("не корректное поле: %q", part)
			}
			// формируем одно поле
			out = append(out, dto.Fields{Start: oneNum, End: oneNum})
		}
	}

	return out, nil
}
