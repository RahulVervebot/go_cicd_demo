package main

import (
	"testing"
)

func TestAdd(t *testing.T) {
	got := add(2, 3)
	if got != 5 {
		t.Fatalf("add(2, 3) = %d, want 5", got)
	}
}

func TestSubtract(t *testing.T) {
	got := subtract(7, 4)
	if got != 3 {
		t.Fatalf("subtract(7, 4) = %d, want 3", got)
	}
}

func TestDivide(t *testing.T) {
	got := divide(10, 5)
	if got != 2 {
		t.Fatalf("divide(10, 5) = %d, want 2", got)
	}
}

func TestDivideByZeroPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("divide did not panic on divide by zero")
		}
	}()
	divide(1, 0)
}