package main

import (
	"fmt"
	"math/rand"
	"time"
)

// отправляем значения в канал с рандомной паузой, а потом закрываем его
func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

// объединяем значения из двух каналов в 3-ий
func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select { // если оба канала готовы, выбирается случайный
			case v, ok := <-a:
				if ok { // канал открыт - передаем значение в c
					c <- v
				} else { // канал закрыт - присваем nil чтобы убрать ветку с проверкой этого канала
					a = nil
				}
			case v, ok := <-b:
				if ok { // канал открыт - передаем значение в c
					c <- v
				} else {
					b = nil // канал закрыт - присваем nil чтобы убрать ветку с проверкой этого канала
				}
			}
			if a == nil && b == nil { // оба канала закрыты - закрываем и итоговый канал
				close(c)
				return
			}
		}
	}()
	return c
}

func main() {
	rand.Seed(time.Now().Unix())
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Print(v) // вывод вариативен, сохраняется только относительный порядок элементов из канала в отдельности. Например 12354687
	}
}
