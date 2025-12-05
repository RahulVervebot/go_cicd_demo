package calc

import "testing"

func TestAdd(t *testing.T) {
	if got := Add(2, 3); got != 5 {
		t.Fatalf("Add(2,3) = %d, want 5", got)
	}
}

func TestSub(t *testing.T) {
	if got := Sub(10, 4); got != 6 {
		t.Fatalf("Sub(10,4) = %d, want 6", got)
	}
}
