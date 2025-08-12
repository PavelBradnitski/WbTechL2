package main

import "fmt"

func main() {
	a := [5]int{76, 77, 78, 79, 80}
	var b []int = a[1:4] // b = a[1] a[2] a[3]
	fmt.Println(b)       // [77 78 79]
}
