package services

import (
	"Calculator/internal/executor"
	"errors"
	"testing"
)

func TestValidateService(t *testing.T) {
	vs := NewValidationService()

	tests1 := []struct {
		name string
		arg  interface{}
		want error
	}{
		{"x", "x", nil},
		{"2", 2, nil},
		{"custom type", struct {
			int    int
			string string
		}{0, ""}, executor.ErrUnsupportedArgType},
	}
	for _, tt := range tests1 {
		t.Run(tt.name, func(t *testing.T) {
			result := vs.ValidateArgType(tt.arg)
			if !errors.Is(result, tt.want) {
				t.Errorf("expected %d, got %d", tt.want, result)
			}
		})
	}

	tests2 := []struct {
		name           string
		expressionVars map[string]bool
		variable       string
		want           error
	}{
		{"x already used", map[string]bool{"x": true, "y": true}, "x", executor.ErrVarAlreadyUsed},
		{"x already used", map[string]bool{"x": false, "y": true}, "x", executor.ErrVarAlreadyUsed},
		{"x not used", map[string]bool{"y": true}, "x", nil},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			result := vs.CheckIfVarAlreadyUsed(tt.expressionVars, tt.variable)
			if !errors.Is(result, tt.want) {
				t.Errorf("expected %d, got %d", tt.want, result)
			}
		})
	}

	tests3 := []struct {
		name        string
		expressions []executor.Expression
		variable    string
		want        error
	}{
		{"cyclic dependency true", []executor.Expression{{Variable: "z"}, {Variable: "y"}}, "y", executor.ErrCyclicDependency},
		{"cyclic dependency false", []executor.Expression{{Variable: "z"}}, "y", nil},
	}
	for _, tt := range tests3 {
		t.Run(tt.name, func(t *testing.T) {
			result := vs.CheckCyclicDependency(tt.expressions, tt.variable)
			if !errors.Is(result, tt.want) {
				t.Errorf("expected %d, got %d", tt.want, result)
			}
		})
	}

	tests4 := []struct {
		name           string
		dependencyMap  map[string][]executor.Expression
		expressionVars map[string]bool
		want           error
	}{
		{"x never calculated", map[string][]executor.Expression{"x": {}, "y": {}}, map[string]bool{"y": true}, executor.ErrVarWillNeverBeCalc},
		{"all args will calculated", map[string][]executor.Expression{"x": {}, "y": {}}, map[string]bool{"y": true, "x": true}, nil},
	}
	for _, tt := range tests4 {
		t.Run(tt.name, func(t *testing.T) {
			result := vs.CheckIfArgNeverCalculated(tt.dependencyMap, tt.expressionVars)
			if !errors.Is(result, tt.want) {
				t.Errorf("expected %d, got %d", tt.want, result)
			}
		})
	}

	tests5 := []struct {
		name           string
		varsToPrint    map[string]bool
		expressionVars map[string]bool
		want           error
	}{
		{"x never calculated", map[string]bool{"x": true, "y": true}, map[string]bool{"y": true}, executor.ErrVarWillNeverBeCalc},
		{"all vars will calculated", map[string]bool{"x": true, "y": true}, map[string]bool{"y": true, "x": true}, nil},
	}
	for _, tt := range tests5 {
		t.Run(tt.name, func(t *testing.T) {
			result := vs.CheckIfPrintVarNeverCalculated(tt.varsToPrint, tt.expressionVars)
			if !errors.Is(result, tt.want) {
				t.Errorf("expected %d, got %d", tt.want, result)
			}
		})
	}
}
