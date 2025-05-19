package services

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestArithmeticService(t *testing.T) {
	as := NewArithmeticService()

	t.Run("Sum", func(t *testing.T) {
		result := as.Sum(int64(2), int64(3))
		require.Equal(t, int64(5), result)
	})

	t.Run("Sub", func(t *testing.T) {
		result := as.Sub(int64(10), int64(3))
		require.Equal(t, int64(7), result)
	})

	t.Run("Mul", func(t *testing.T) {
		result := as.Multi(int64(2), int64(3))
		require.Equal(t, int64(6), result)
	})
}
