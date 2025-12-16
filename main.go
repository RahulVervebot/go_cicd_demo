package main

import "fmt"

func main() {
	a, b := 10, 5
	sum := add(a, b)
	subtract := subtract(a, b)
	fmt.Printf("Add: %d + %d = %d\n", a, b, sum)
  fmt.Printf("Sub: %d - %d = %d\n", a, b, subtract)
}

func add(a, b int) int {
	return a + b
}

func subtract(a, b int) int {
	return a - b
}