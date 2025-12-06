package main

import "fmt"

// Simple calculator demo to play with Git conflicts.
func main() {
	a, b := 10, 5

	sum := add(a, b)
	fmt.Printf("Add: %d + %d = %d\n", a, b, sum)

	// NOTE:
	// - Branch feat/add-mul will add multiplication here.
	// - Branch feat/add-div will add division here.
	// Both branches will modify this SAME area, so we can see merge conflicts.
}

func add(a, b int) int {
	return a + b
}

// You can also later add subtraction on main branch if you want:
func subtract(a, b int) int {
	return a - b
}
