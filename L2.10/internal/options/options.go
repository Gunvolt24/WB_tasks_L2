package options

// Options - опции сравнения
type Options struct {
	Column        int    // 0 = вся стрка; 1..N - номер столбца
	NumSort       bool   // флаг: -n: числовая сортировка
	Reverse       bool   // флаг: -r: сортировка в обратном порядке
	Unique        bool   // флаг: -u: вывод уникальных строк
	MonthNames    bool   // флаг: -M: сортировать по названию месяца
	TrailBlanks   bool   // флаг: -b: игнорировать хвостовые пробелы (trailing blanks)
	HumanReadable bool   // флаг: -h: человекочитаемая сортировка по числовому значению с учетом суффиксов
	Splitter      string // флаг: -d: разделитель столбцов
}
