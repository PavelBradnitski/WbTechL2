package main

func main() {
	ch := make(chan int)
	go func() {
		// для исправления нужно раскомментировать эту строку
		//defer close(ch) // Закрываем канал после отправки всех значений
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()
	for n := range ch {
		// выведет от 0 до 9, а позже завершится ошибкой "fatal error: all goroutines are asleep - deadlock!".
		// из за того что не закрываем канал
		println(n)
	}
}
