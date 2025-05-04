package gin_handlers

import (
	"Calculator/internal/executor"
	"Calculator/internal/executor/dto"
	"Calculator/internal/executor/use_cases"
	"Calculator/internal/executor/values"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GinHandler struct {
	uc *use_cases.UseCase
}

func NewGinHandler(uc *use_cases.UseCase) *GinHandler {
	return &GinHandler{uc: uc}
}

func (gh *GinHandler) Execute(gctx *gin.Context) {
	var insts []Instruction
	if err := gctx.ShouldBindJSON(&insts); err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body, " + err.Error()})
		return
	}

	gi := dto.GroupedInstructions{
		Expressions: make([]executor.Expression, 0),
		VarsToPrint: make(map[string]bool),
	}

	err := gh.groupInstructions(&insts, &gi)
	if err != nil {
		gh.mapExecutorErrToHTTPErr(err, gctx)
		return
	}

	reqId, exists := gctx.Get(values.RequestIdKey)
	if !exists {
		gctx.JSON(http.StatusInternalServerError, gin.H{"error": "request id is missing in the context"})
		return
	}

	ctx := gctx.Request.Context()
	ctx = context.WithValue(ctx, values.RequestIdKey, reqId)

	results, err := gh.uc.Execute(ctx, &gi)
	if err != nil {
		gh.mapExecutorErrToHTTPErr(err, gctx)
		return
	}

	gctx.JSON(http.StatusOK, gin.H{"items": gh.resultsToItems(&results)})
}

func (gh *GinHandler) resultsToItems(results *[]executor.Result) *[]Item {
	items := make([]Item, len(*results))
	for i, result := range *results {
		items[i] = Item{
			Var:   result.Key,
			Value: result.Value,
		}
	}
	return &items
}

func (gh *GinHandler) groupInstructions(instructions *[]Instruction, gi *dto.GroupedInstructions) error {
	for _, instruction := range *instructions {
		if instruction.Type == values.Calculate {
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
			gi.Expressions = append(gi.Expressions, expression)
			continue
		} else if instruction.Type == values.Print {
			gi.VarsToPrint[instruction.Variable] = false
			continue
		}
		return fmt.Errorf("%w: %v", executor.ErrUnknownInstructionType, instruction.Type)
	}
	return nil
}

func (gh *GinHandler) mapExecutorErrToHTTPErr(err error, c *gin.Context) {
	switch {
	case errors.Is(err, executor.ErrCyclicDependency) ||
		errors.Is(err, executor.ErrCalcExpression) ||
		errors.Is(err, executor.ErrVarNeverBeCalc):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, executor.ErrVarAlreadyUsed):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, executor.ErrVarToPrintNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, executor.ErrUnknownInstructionType):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
