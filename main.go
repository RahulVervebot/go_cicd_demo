package main

import "fmt"

func main() {
	a, b := 10, 5
	sum := add(a, b)
	fmt.Printf("Add: %d + %d = %d\n", a, b, sum)
	mul := multiply(a, b)
	fmt.Printf("Mul: %d * %d = %d\n", a, b, mul)
}

func add(a, b int) int {
	return a + b
}

func subtract(a, b int) int {
	return a - b
}



// this should be added
func multiply(a, b int) int {
	return a * b
}