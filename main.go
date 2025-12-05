package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/RahulVervebot/go_cicd_demo/internal/calc"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go_cicd_demo <a> <b>")
		os.Exit(2)
	}

	a, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("invalid a:", err)
		os.Exit(2)
	}
	b, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("invalid b:", err)
		os.Exit(2)
	}

	fmt.Printf("%d + %d = %d\n", a, b, calc.Add(a, b))
}
