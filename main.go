package main

import "fmt"

// Simple calculator demo to play with Git conflicts.
func main() {
	a, b := 10, 5

	sum := add(a, b)
	fmt.Printf("Add: %d + %d = %d\n", a, b, sum)

	// Division added in feat/add-div branch
	div := divide(a, b)
	fmt.Printf("Div: %d / %d = %d\n", a, b, div)

	// NOTE:
	// Multiplication was added in feat/add-mul on the SAME block.
	// When we merge this branch to main, Git will see different changes
}

func add(a, b int) int {
	return a + b
}

func subtract(a, b int) int {
	return a - b
}

func divide(a, b int) int {
	if b == 0 {
		panic("cannot divide by zero")
	}
	return a / b
}
