package main

import "fmt"

func main() {
	a, b := 10, 5
	sum := add(a, b)
	subtract := subtract(a, b)
	multiply := multiply(a, b)
	fmt.Printf("Add: %d + %d = %d\n", a, b, sum)
	fmt.Printf("Sub: %d - %d = %d\n", a, b, subtract)
	fmt.Printf("Mul: %d * %d = %d\n", a, b, multiply)
}

func add(a, b int) int {
	return a + b
}

func subtract(a, b int) int {
	return a - b
}
// changes
func multiply(a, b int) int {
	return a * b
}

func divide(a, b int) int {
	return a / b
}