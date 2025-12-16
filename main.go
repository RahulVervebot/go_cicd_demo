package main

import "fmt"

func main() {
	a, b := 10, 5
	sum := add(a, b)
	fmt.Printf("Add: %d + %d = %d\n", a, b, sum)

	diff := subtract(a, b)
	fmt.Printf("Sub: %d - %d = %d\n", a, b, diff)

	prod := multiply(a, b)
	fmt.Printf("Mul: %d * %d = %d\n", a, b, prod)

	quot := divide(a, b)
	fmt.Printf("Div: %d / %d = %d\n", a, b, quot)
}

func add(a, b int) int {
	return a + b
}

func subtract(a, b int) int {
	return a - b
}

func multiply(a, b int) int {
	return a * b
}

func divide(a, b int) int {
	if b == 0 {
		panic("cannot divide by zero")
	}
	return a / b
}
