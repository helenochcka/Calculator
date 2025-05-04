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
		gctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	gi := dto.GroupedInstructions{
		Expressions: make([]executor.Expression, 0),
		VarsToPrint: make(map[string]bool),
	}

	err := gh.groupInstructions(&insts, &gi)
	if err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	reqId, exists := gctx.Get(values.RequestIdKey)
	if !exists {
		gctx.JSON(http.StatusInternalServerError, gin.H{"error": executor.ErrReqIdMissing.Error()})
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
			err := gh.validateCalcInst(&instruction)
			if err != nil {
				return err
			}
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
			gi.VarsToPrint[instruction.Variable] = true
			continue
		}
		return fmt.Errorf("%v (%v)", "unknown type of instruction", instruction.Type)
	}
	return nil
}

func (gh *GinHandler) validateCalcInst(instruction *Instruction) error {
	if instruction.Operation == nil {
		return errors.New("field 'op' is missing in calculate instruction")
	}
	if instruction.Left == nil {
		return errors.New("field 'left' is missing in calculate instruction")
	}
	if instruction.Right == nil {
		return errors.New("field 'right' is missing in calculate instruction")
	}
	return nil
}

func (gh *GinHandler) mapExecutorErrToHTTPErr(err error, c *gin.Context) {
	switch {
	case errors.Is(err, executor.ErrCyclicDependency) ||
		errors.Is(err, executor.ErrCalcExpression) ||
		errors.Is(err, executor.ErrVarWillNeverBeCalc):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, executor.ErrVarAlreadyUsed):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
