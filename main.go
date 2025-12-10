package main
import "fmt"

func main() {
	a, b := 10, 5
	sum := add(a, b)
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

func divide(a, b int) int {
	if b == 0 {
		panic("cannot divide by zero")
	}
	return a / b
}