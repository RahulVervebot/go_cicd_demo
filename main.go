package main

import "fmt"

func main() {
	a, b := 10, 5
	sum := add(a, b)
	subtract := subtract(a, b)
	fmt.Printf("Add: %d + %d = %d\n", a, b, sum)

	diff := subtract(a, b)
	fmt.Printf("Sub: %d - %d = %d\n", a, b, diff)
}

func add(a, b int) int {
	return a + b
}

func subtract(a, b int) int {
	return a - b
}
