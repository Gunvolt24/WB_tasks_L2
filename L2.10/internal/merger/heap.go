package merger

import (
	"bufio"
	"os"
)

// stream - один входной поток (временный файл	+ сканер)
type stream struct {
	file    *os.File
	scanner *bufio.Scanner
	curStr  string
}

// minHeap - минимальная куча по функции less(curStr)
type minHeap struct {
	data     []*stream
	lessFunc func(a, b string) bool
}

// Len - длина кучи
func (h minHeap) Len() int {
	return len(h.data)
}

// Less - функция сравнения двух элементов кучи
func (h minHeap) Less(i, j int) bool {
	return h.lessFunc(h.data[i].curStr, h.data[j].curStr)
}

// Swap - обмен двумя элементами
func (h minHeap) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

// Push - добавление в кучу
func (h *minHeap) Push(x any) {
	h.data = append(h.data, x.(*stream))
}

// Pop - удаление из кучи
func (h *minHeap) Pop() any {
	old := h.data
	n := len(old)
	x := old[n-1]
	h.data = old[:n-1]
	return x
}
