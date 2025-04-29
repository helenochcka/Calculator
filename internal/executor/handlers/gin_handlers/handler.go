package gin_handlers

import (
	"Calculator/internal/executor"
	"Calculator/internal/executor/use_cases"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerGin struct {
	useCase use_cases.UseCase
}

func NewHandlerGin(uc use_cases.UseCase) HandlerGin {
	return HandlerGin{useCase: uc}
}

func (hg *HandlerGin) Execute(c *gin.Context) {
	var instructions []executor.Instruction
	if err := c.ShouldBindJSON(&instructions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body, " + err.Error()})
		return
	}

	varsToPrint := make(map[string]bool)
	var expressions []executor.Expression

	for _, instruction := range instructions {
		expression, err := hg.distributeInstructions(instruction, varsToPrint)
		if err != nil {
			hg.mapExecErrToHTTPErr(err, c)
			return
		}
		if expression != nil {
			expressions = append(expressions, *expression)
		}
	}
	reqId, exists := c.Get("request_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "request id is missing in context"})
		return
	}

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "request_id", reqId)

	items, err := hg.useCase.Execute(ctx, expressions, varsToPrint)
	if err != nil {
		hg.mapExecErrToHTTPErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (hg *HandlerGin) distributeInstructions(
	instruction executor.Instruction,
	varsToPrint map[string]bool) (*executor.Expression, error) {
	if instruction.Type == "calc" {
		expression := executor.Expression{
			Type:      instruction.Type,
			Operation: *instruction.Operation,
			Variable:  instruction.Variable,
		}
		switch right := instruction.Right.(type) {
		case float64:
			expression.Right = int(right)
		case string:
			expression.Right = right
		}

		switch left := instruction.Left.(type) {
		case float64:
			expression.Left = int(left)
		case string:
			expression.Left = left
		}
		return &expression, nil
	} else if instruction.Type == "print" {
		varsToPrint[instruction.Variable] = false
		return nil, nil
	}
	return nil, fmt.Errorf("%w: %v", executor.UnknownTypeOfInstruction, instruction.Type)
}

func (hg *HandlerGin) mapExecErrToHTTPErr(err error, c *gin.Context) {
	switch {
	case errors.Is(err, executor.CyclicDependency) ||
		errors.Is(err, executor.ErrCalcExpression) ||
		errors.Is(err, executor.VarNeverBeCalc):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, executor.VarIsAlreadyUsed):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, executor.VarToPrintNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, executor.UnknownTypeOfInstruction):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
