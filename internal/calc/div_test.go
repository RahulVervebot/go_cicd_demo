package calc

import "testing"

func TestDiv(t *testing.T) {
	got, err := Div(10, 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != 5 {
		t.Fatalf("Div(10,2) = %d, want 5", got)
	}
}

func TestDivByZero(t *testing.T) {
	_, err := Div(10, 0)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != ErrDivideByZero {
		t.Fatalf("expected ErrDivideByZero, got %v", err)
	}
}
