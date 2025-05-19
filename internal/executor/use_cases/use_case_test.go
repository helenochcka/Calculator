package use_cases

import (
	"Calculator/internal/executor"
	"Calculator/internal/executor/dto"
	"Calculator/internal/executor/services"
	"Calculator/internal/executor/use_cases/mocks"
	"Calculator/internal/executor/values"
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCS := mocks.NewMockCommunicationService(ctrl)

	uc := &UseCase{
		cs: mockCS,
		vs: &services.ValidationService{},
		gs: &services.GetterService{},
	}

	t.Run("ReqIdMissing", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), &dto.GroupedInstructions{})

		require.Equal(t, err, executor.ErrReqIdMissing)
	})

	ctx := context.WithValue(context.Background(), values.RequestIdKey, "12345")
	t.Run("ErrDeclaringQueue", func(t *testing.T) {
		mockCS.EXPECT().DeclareResultsQueue("12345").Return(executor.ErrDeclaringQueue).Times(1)
		_, err := uc.Execute(ctx, &dto.GroupedInstructions{})

		require.Equal(t, err, executor.ErrDeclaringQueue)
	})

	t.Run("VarAlreadyUsed", func(t *testing.T) {
		mockCS.EXPECT().DeclareResultsQueue("12345").Return(nil).Times(1)
		mockCS.EXPECT().RequestCalculation(gomock.Any()).AnyTimes()
		gi := dto.GroupedInstructions{Expressions: []executor.Expression{{Variable: "x", Left: 2, Right: 2}, {Variable: "x"}}}
		_, err := uc.Execute(ctx, &gi)
		if !errors.Is(err, executor.ErrVarAlreadyUsed) {
			t.Errorf("want err: executor.ErrVarAlreadyUsed, got %v", err)
		}
	})

	t.Run("UnsupportedLeftArgType", func(t *testing.T) {
		mockCS.EXPECT().DeclareResultsQueue("12345").Return(nil).Times(1)
		gi := dto.GroupedInstructions{Expressions: []executor.Expression{{Left: struct {
			int    int
			string string
		}{0, ""}, Variable: "x", Right: 2}}}
		_, err := uc.Execute(ctx, &gi)
		if !errors.Is(err, executor.ErrUnsupportedArgType) {
			t.Errorf("want err: executor.ErrUnsupportedArgType, got %v", err)
		}
	})

	t.Run("UnsupportedRightArgType", func(t *testing.T) {
		mockCS.EXPECT().DeclareResultsQueue("12345").Return(nil).Times(1)
		gi := dto.GroupedInstructions{Expressions: []executor.Expression{{Right: struct {
			int    int
			string string
		}{0, ""}, Variable: "x", Left: 2}}}
		_, err := uc.Execute(ctx, &gi)
		if !errors.Is(err, executor.ErrUnsupportedArgType) {
			t.Errorf("want err: executor.ErrUnsupportedArgType, got %v", err)
		}
	})

	t.Run("CyclicDependencyLeft", func(t *testing.T) {
		mockCS.EXPECT().DeclareResultsQueue("12345").Return(nil).Times(1)
		mockCS.EXPECT().RequestCalculation(gomock.Any()).AnyTimes()
		gi := dto.GroupedInstructions{Expressions: []executor.Expression{{Variable: "x", Left: 2, Right: "y"}, {Variable: "y", Left: "x", Right: 2}}}
		_, err := uc.Execute(ctx, &gi)
		require.Equal(t, err, executor.ErrCyclicDependency)
	})

	t.Run("CyclicDependencyRight", func(t *testing.T) {
		mockCS.EXPECT().DeclareResultsQueue("12345").Return(nil).Times(1)
		mockCS.EXPECT().RequestCalculation(gomock.Any()).AnyTimes()
		gi := dto.GroupedInstructions{Expressions: []executor.Expression{{Variable: "x", Left: 2, Right: "y"}, {Variable: "y", Left: 2, Right: "x"}}}
		_, err := uc.Execute(ctx, &gi)
		require.Equal(t, err, executor.ErrCyclicDependency)
	})

	t.Run("ArgNeverCalculated", func(t *testing.T) {
		mockCS.EXPECT().DeclareResultsQueue("12345").Return(nil).Times(1)
		gi := dto.GroupedInstructions{Expressions: []executor.Expression{
			{Variable: "x", Left: 2, Right: "y"}}}
		_, err := uc.Execute(ctx, &gi)
		if !errors.Is(err, executor.ErrVarWillNeverBeCalc) {
			t.Errorf("want err: executor.ErrVarWillNeverBeCalc, got %v", err)
		}
	})

	t.Run("PrintVarNeverCalculated", func(t *testing.T) {
		mockCS.EXPECT().DeclareResultsQueue("12345").Return(nil).Times(1)
		mockCS.EXPECT().RequestCalculation(gomock.Any()).AnyTimes()
		gi := dto.GroupedInstructions{Expressions: []executor.Expression{
			{Variable: "x", Left: 2, Right: 2}}, VarsToPrint: map[string]bool{"z": true}}
		_, err := uc.Execute(ctx, &gi)
		if !errors.Is(err, executor.ErrVarWillNeverBeCalc) {
			t.Errorf("want err: executor.ErrVarWillNeverBeCalc, got %v", err)
		}
	})

	t.Run("ErrConsumingResult", func(t *testing.T) {
		mockCS.EXPECT().DeclareResultsQueue("12345").Return(nil).Times(1)
		mockCS.EXPECT().RequestCalculation(gomock.Any()).AnyTimes()
		mockCS.EXPECT().ConsumeResults(gomock.Any(), gomock.Any()).Return(executor.ErrConsumingResult)
		gi := dto.GroupedInstructions{Expressions: []executor.Expression{
			{Variable: "x", Left: 2, Right: 2}}, VarsToPrint: map[string]bool{"x": true}}
		_, err := uc.Execute(ctx, &gi)
		if !errors.Is(err, executor.ErrConsumingResult) {
			t.Errorf("want err: executor.ErrConsumingResult, got %v", err)
		}
	})

	t.Run("ErrCalcExpression", func(t *testing.T) {
		mockCS.EXPECT().DeclareResultsQueue("12345").Return(nil).Times(1)
		mockCS.EXPECT().RequestCalculation(gomock.Any()).AnyTimes()
		mockCS.EXPECT().ConsumeResults(gomock.Any(), gomock.Any()).Return(executor.ErrCalcExpression)
		gi := dto.GroupedInstructions{Expressions: []executor.Expression{
			{Variable: "x", Left: 2, Right: 2}}, VarsToPrint: map[string]bool{"x": true}}
		_, err := uc.Execute(ctx, &gi)
		if !errors.Is(err, executor.ErrCalcExpression) {
			t.Errorf("want err: executor.ErrCalcExpression, got %v", err)
		}
	})

	t.Run("SuccessCalc", func(t *testing.T) {
		mockCS.EXPECT().DeclareResultsQueue("12345").Return(nil).Times(1)
		mockCS.EXPECT().RequestCalculation(gomock.Any()).AnyTimes()
		mockCS.EXPECT().ConsumeResults(gomock.Any(), gomock.Any()).Return(nil)
		gi := dto.GroupedInstructions{Expressions: []executor.Expression{
			{Variable: "x", Left: "y", Right: "z"},
			{Variable: "y", Left: 2, Right: "z"},
			{Variable: "z", Left: 3, Right: 4},
		}, VarsToPrint: map[string]bool{
			"y": true,
			"x": true,
			"z": true,
		}}

		_, err := uc.Execute(ctx, &gi)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
		}
	})
}
