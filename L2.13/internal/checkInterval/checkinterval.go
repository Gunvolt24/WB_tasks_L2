package checkinterval

import "github.com/Gunvolt24/wb_l2/L2.13/internal/dto"

// IsInInterval - проверяет, находится ли число в диапазоне, чтобы не выводить лишние поля
func IsInInterval(num int, fields []dto.Fields) bool {
	for _, field := range fields {
		if num >= field.Start && num <= field.End {
			return true
		}
	}
	return false
}
