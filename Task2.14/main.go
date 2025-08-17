package main

import (
	"fmt"
	"sync"
	"time"
)

// Функция объединяет один или более каналов done (каналов сигнала завершения) в один
// Возвращаемый канал закрыватется, как только закроется любой из исходных каналов.
func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		// Если нет каналов, возвращаем закрытый канал
		out := make(chan interface{})
		close(out)
		return out
	case 1:
		// Если канал один возвращаем его
		return channels[0]
	}
	// Создаем новый канал, который будет использоваться для объединения входных каналов
	out := make(chan interface{})
	var once sync.Once

	launch := func(ch <-chan interface{}) {
		go func() {
			<-ch
			once.Do(func() {
				close(out)
			})
		}()
	}

	for _, ch := range channels {
		launch(ch)
	}

	return out
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()

		return c
	}
	_ = sig(1 * time.Second)
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}
