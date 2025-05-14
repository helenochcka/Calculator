package use_cases

import (
	"Calculator/internal/arithmetic"
	"testing"
)

type mockResultService struct{}
type mockArithmeticService struct{}

func (mrs *mockResultService) PublishResult(result arithmetic.Result, queueName string) {}

func (mrs *mockResultService) PublishError(errMsg string, queueName string) {}

func (mar *mockArithmeticService) Sum(a, b int64) int64 {
	return 42
}

func (mar *mockArithmeticService) Multi(a, b int64) int64 {
	return 42
}

func (mar *mockArithmeticService) Sub(a, b int64) int64 {
	return 42
}

func TestExecute(t *testing.T) {
	//mockResultService := &mockResultService{}
	//mockArithmeticService := &mockArithmeticService{}
	//uc := use_cases.NewUseCase(mockResultService, mockArithmeticService)
	//
	//result := uc.Execute()
	//
	//if result != 42 {
	//	t.Errorf("expected 42, got %d", result)
	//}
}
