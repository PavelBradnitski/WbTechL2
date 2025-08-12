package main

import "fmt"

func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return // defer выполняется и x инкрементируется
}

func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x // сохраняется значение x = 1 и вызов defer не меняет возвращаемое значение
}

func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
