package calc

import "testing"

func TestMul(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"positive", 2, 3, 6},
		{"with zero", 10, 0, 0},
		{"negative", -2, 4, -8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mul(tt.a, tt.b); got != tt.want {
				t.Fatalf("Mul(%d,%d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
