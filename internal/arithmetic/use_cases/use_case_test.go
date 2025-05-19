package use_cases

import (
	"Calculator/internal/arithmetic"
	"Calculator/internal/arithmetic/services"
	"Calculator/internal/arithmetic/use_cases/mocks"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockResultService := mocks.NewMockResultService(ctrl)

	uc := &UseCase{
		rs: mockResultService,
		as: &services.ArithmeticService{},
	}

	t.Run("ExecuteSuccess", func(t *testing.T) {
		mockResultService.EXPECT().PublishResult(arithmetic.Result{Key: "x", Value: 4}, "results")
		uc.Execute(arithmetic.Expression{Op: "+", Variable: "x", Left: 2, Right: 2}, "results")
	})

	t.Run("ExecuteError", func(t *testing.T) {
		mockResultService.EXPECT().PublishError(gomock.Any(), "results")
		uc.Execute(arithmetic.Expression{Op: "/", Variable: "x", Left: 2, Right: 2}, "results")
	})
}
