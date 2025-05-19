package services

import (
	"Calculator/api/arithmeticpb"
	"Calculator/internal/arithmetic"
	"Calculator/internal/executor"
	"Calculator/internal/executor/dto"
	"Calculator/internal/executor/services/mocks"
	"Calculator/internal/infrastructure/rabbitmq"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestCommService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBrokerClient := mocks.NewMockBrokerClient(ctrl)
	mockArithmClient := mocks.NewMockArithmeticClient(ctrl)

	cs := &CommunicationService{
		ac: mockArithmClient,
		bc: mockBrokerClient,
	}

	t.Run("RequestCalculation", func(t *testing.T) {
		expectedReq := arithmeticpb.CalculationData{
			Var:       "x",
			Op:        "+",
			Left:      2,
			Right:     2,
			QueueName: "results",
		}

		expectedResp := arithmeticpb.Message{Text: fmt.Sprintf("Expression received %v",
			arithmetic.Expression{Variable: "x", Op: "+", Left: 2, Right: 2})}

		mockArithmClient.EXPECT().Calculate(context.Background(), &expectedReq).Return(&expectedResp, nil).Times(1)
		req := dto.CalculationData{
			Variable:  "x",
			Operation: "+",
			Left:      2,
			Right:     2,
			QueueName: "results",
		}
		cs.RequestCalculation(&req)
	})

	t.Run("DeclareResultsQueueSuccess", func(t *testing.T) {
		mockBrokerClient.EXPECT().DeclareQueue("results").Return(nil).Times(1)
		err := cs.DeclareResultsQueue("results")
		require.NoError(t, err)
	})

	t.Run("DeclareResultsQueueErr", func(t *testing.T) {
		mockBrokerClient.EXPECT().DeclareQueue("results").Return(errors.New("fail"))
		err := cs.DeclareResultsQueue("results")
		require.Error(t, err, errors.New("fail"))
	})

	rp := func(result executor.Result) bool {
		return true
	}
	t.Run("ConsumeResultsSuccess", func(t *testing.T) {
		mockBrokerClient.EXPECT().Consume("results", gomock.Any()).Return(nil).Times(1)

		err := cs.ConsumeResults("results", rp)
		require.NoError(t, err)
	})

	t.Run("ConsumeResultsErrCalculatingResult", func(t *testing.T) {
		mockBrokerClient.EXPECT().Consume("results", gomock.Any()).Return(rabbitmq.ErrCalculatingResult)
		err := cs.ConsumeResults("results", rp)
		require.Error(t, err, executor.ErrCalcExpression)
	})

	t.Run("ConsumeResultsErrUnmarshallingMsg", func(t *testing.T) {
		mockBrokerClient.EXPECT().Consume("results", gomock.Any()).Return(rabbitmq.ErrUnmarshallingMsg)
		err := cs.ConsumeResults("results", rp)
		require.Error(t, err, executor.ErrConsumingResult)
	})

	t.Run("ConsumeResultsErrConsumingMsgs", func(t *testing.T) {
		mockBrokerClient.EXPECT().Consume("results", gomock.Any()).Return(rabbitmq.ErrConsumingMsgs)
		err := cs.ConsumeResults("results", rp)
		require.Error(t, err, executor.ErrConsumingResult)
	})

	t.Run("ConsumeResultsAnotherErr", func(t *testing.T) {
		mockBrokerClient.EXPECT().Consume("results", gomock.Any()).Return(errors.New("some error"))
		err := cs.ConsumeResults("results", rp)
		require.Error(t, err, errors.New("some error"))
	})
}
