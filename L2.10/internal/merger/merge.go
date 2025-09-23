package merger

import (
	"bufio"
	"container/heap"
	"io"
	"os"
)

// KWayMerge - выполняет слияние заранее отсортированных временных файлов в поток out.
// Если Unique=true, то удаляет дублирующиеся строки по ключу на стыках разных чанков.
func KWayMerge(
	paths []string,
	less func(a, b string) bool,
	equal func(a, b string) bool,
	out io.Writer,
	unique bool,
) error {
	var streams []*stream
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(file)
		buffer := make([]byte, 64*1024)
		scanner.Buffer(buffer, 4*1024*1024) // до 4 МБ на одну строку
		s := &stream{
			file:    file,
			scanner: scanner,
		}
		if scanner.Scan() {
			s.curStr = scanner.Text()
			streams = append(streams, s)
		} else {
			_ = file.Close()
			if err := scanner.Err(); err != nil {
				return err
			}
		}
	}

	h := &minHeap{
		lessFunc: less,
	}
	for _, stream := range streams {
		h.data = append(h.data, stream)
	}
	heap.Init(h)

	writer := bufio.NewWriter(out)

	var (
		prev string
		ok   bool
	)

	for h.Len() > 0 {
		s := heap.Pop(h).(*stream)

		// -u: подавление дубликатов на стыках чанков
		if !unique || !ok || !equal(prev, s.curStr) {
			_, _ = writer.WriteString(s.curStr)
			_ = writer.WriteByte('\n')
			prev = s.curStr
			ok = true
		}

		if s.scanner.Scan() {
			s.curStr = s.scanner.Text()
			heap.Push(h, s)
		} else {
			_ = s.file.Close()
			if err := s.scanner.Err(); err != nil {
				_ = writer.Flush()
				return err
			}
		}
	}

	return writer.Flush()
}
