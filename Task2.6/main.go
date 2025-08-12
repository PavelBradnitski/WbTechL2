package main

import (
	"fmt"
)

func main() {
	var s = []string{"1", "2", "3"}
	modifySlice(s)
	fmt.Println(s) // [3 2 3]
}

func modifySlice(i []string) {
	i[0] = "3"         // меняет слайс в main
	i = append(i, "4") // не меняет в main, i указывает на новый массив
	i[1] = "5"         // работа с новым массивом
	i = append(i, "6") // работа с новым массивом
}
