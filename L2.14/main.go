// Источник информации: https://ndukwearmstrong.medium.com/the-or-channel-pattern-w-dynamic-selects-ae2c52fe1cfd

package main

import (
	"fmt"
	"reflect"
	"time"
)

// or - возвращает канал, который закроется, когда закроются все переданные каналы
// использует динамический reflect.Select (Intuitive OR-Channel Pattern)
func or(channels ...<-chan any) <-chan any {
	done := make(chan any)
	// создаем слайс SelectCase - используется, когда есть неизвестное количество cases
	cases := make([]reflect.SelectCase, 0, len(channels))

	// итерируемся по списку каналов и создаем case для каждого из них,
	// коротый сам по себе является структурой SelectCase{Dir: для указания направления канала, Chan: сам канал, но типа reflect.ValueOf(ch)}
	for _, ch := range channels {
		// собираем только не nil каналы
		if ch == nil {
			continue
		}
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}

	// если каналов нет (или все nil) - вернуть сразу закрытый done
	if len(cases) == 0 {
		close(done)
		return done
	}

	// создаем горутину для проверки каналов и закрытия done
	go func() {
		defer close(done)
		for {
			_, _, ok := reflect.Select(cases)
			if !ok { // если канал закрыт - закрываем done
				return
			}
			// если пришло значение (ok == true) - ждем пока закроются все каналы
		}
	}()
	return done

}

func main() {
	signal := func(after time.Duration) <-chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	<-or(
		signal(5*time.Second),
		signal(5*time.Minute),
		signal(3*time.Second),
		signal(1*time.Second),
		signal(3*time.Minute),
		signal(2*time.Second),
	)

	fmt.Printf("done after %v", time.Since(start))
}
