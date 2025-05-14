package services

import (
	"Calculator/internal/arithmetic/services"
	"testing"
)

func TestArithmeticService(t *testing.T) {
	as := services.NewArithmeticService()

	t.Run("Sum", func(t *testing.T) {
		result := as.Sum(2, 3)
		if result != 5 {
			t.Errorf("expected 5, got %d", result)
		}
	})

	t.Run("Sub", func(t *testing.T) {
		result := as.Sub(10, 3)
		if result != 7 {
			t.Errorf("expected 7, got %d", result)
		}
	})

	t.Run("Mul", func(t *testing.T) {
		result := as.Multi(2, 3)
		if result != 6 {
			t.Errorf("expected 6, got %d", result)
		}
	})
}
