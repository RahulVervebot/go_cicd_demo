package main

import "fmt"

// Simple calculator demo to play with Git conflicts.
func main() {
	a, b := 10, 5
	sum := add(a, b)
	fmt.Printf("Add: %d + %d = %d\n", a, b, sum)
	// Multiplication added in feat/add-mul branch
}

func add(a, b int) int {
	return a + b
}

func subtract(a, b int) int {
	return a - b
}