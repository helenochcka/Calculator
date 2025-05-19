package services

import (
	"Calculator/api/arithmeticpb"
	"Calculator/internal/arithmetic"
	"Calculator/internal/arithmetic/services/mocks"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestResultService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBrokerClient := mocks.NewMockBrokerClient(ctrl)

	rs := &ResultService{
		bc: mockBrokerClient,
	}

	t.Run("SuccessPublishResult", func(t *testing.T) {
		body, _ := proto.Marshal(&arithmeticpb.Result{Key: ptrString("x"), Value: ptrInt64(4)})
		mockBrokerClient.EXPECT().Publish("result", body).Return(nil).Times(1)
		rs.PublishResult(arithmetic.Result{Key: "x", Value: 4}, "result")
	})

	t.Run("PublishError", func(t *testing.T) {
		body, _ := proto.Marshal(&arithmeticpb.Result{ErrMsg: ptrString("operation '/' is not supported")})
		mockBrokerClient.EXPECT().Publish("result", body).Times(1)
		rs.PublishError("operation '/' is not supported", "result")
	})
}

func ptrString(s string) *string {
	return &s
}

func ptrInt64(i int64) *int64 {
	return &i
}
