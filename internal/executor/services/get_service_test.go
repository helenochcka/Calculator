package services

import (
	"Calculator/internal/executor"
	"Calculator/internal/executor/values"
	"context"
	"errors"
	"testing"
)

func TestGetterService(t *testing.T) {
	gs := NewGetterService()

	tests1 := []struct {
		name        string
		ctx         context.Context
		want        *string
		expectErr   bool
		expectedErr error
	}{
		{"expected error", context.Background(), nil, true, executor.ErrReqIdMissing},
		{"expected no error and 12345", context.WithValue(context.Background(), values.RequestIdKey, "12345"), ptrString("12345"), false, nil},
	}
	for _, tt := range tests1 {
		t.Run(tt.name, func(t *testing.T) {
			result, err := gs.GetReqIdFromCtx(tt.ctx)

			if tt.expectErr && (err == nil || !errors.Is(err, tt.expectedErr)) {
				t.Errorf("expected error %v but got %v", tt.expectedErr, err)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("expected no error but got %v", err)
			}

			if tt.want != nil && (result == nil || *result != *tt.want) {
				t.Errorf("expected %v but got %v", *tt.want, result)
			}

			if tt.want == nil && result != nil {
				t.Errorf("expected nil but got %v", result)
			}
		})
	}

	tests2 := []struct {
		name      string
		variable  interface{}
		resultMap map[string]int
		want      *int
		exists    bool
	}{
		{"value for x does not exist", "x", map[string]int{}, nil, false},
		{"value for x exist", "x", map[string]int{"x": 1}, ptrInt(1), true},
		{"value for 1", 1, map[string]int{}, ptrInt(1), true},
		{"nil for custom struct", struct {
			int    int
			string string
		}{0, ""}, map[string]int{}, nil, false},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			result, exists := gs.GetVarValue(tt.variable, tt.resultMap)

			if !tt.exists && (exists || result != nil) {
				t.Errorf("expected not exists but exists or got %d", result)
			}

			if tt.exists && (result == nil || *result != *tt.want) {
				t.Errorf("expected %d but not exists or got %d", *tt.want, result)
			}

		})
	}
}

func ptrString(s string) *string {
	return &s
}

func ptrInt(i int) *int {
	return &i
}
