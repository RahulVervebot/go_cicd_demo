package main

import "fmt"

func main() {
	a, b := 10, 5
	sum := add(a, b)
	fmt.Printf("Add: %d + %d = %d\n", a, b, sum)	
	mul := multiply(a, b)
	fmt.Printf("Mul: %d * %d = %d\n", a, b, mul)
	fmt.Printf("Add: %d + %d = %d\n", a, b, sum)
	div := divide(a, b)
	fmt.Printf("Div: %d / %d = %d\n", a, b, div)
}


func add(a, b int) int {
	return a + b
}

func subtract(a, b int) int {
	return a - b
}

// this is multiple
func multiply(a, b int) int {
	return a * b
}