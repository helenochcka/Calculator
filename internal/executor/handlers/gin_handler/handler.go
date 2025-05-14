package gin_handler

import (
	"Calculator/internal/executor"
	"Calculator/internal/executor/dto"
	"Calculator/internal/executor/use_cases"
	"Calculator/internal/executor/values"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GinHandler struct {
	uc *use_cases.UseCase
}

func NewGinHandler(uc *use_cases.UseCase) *GinHandler {
	return &GinHandler{uc: uc}
}

// Execute godoc
//
//	@Summary		Execute instructions
//	@Description	Calculate expressions and print results of requested instructions.
//	@Description	The result can only be written to the same variable once.
//	@Produce		json
//	@Param			instructions	body		[]Instruction	true	"instructions to calculate"
//	@Success		200				{object}	[]Item
//	@Failure		400				{object}	HTTPError		"possible error codes: INVALID_JSON_BODY, UNKNOWN_TYPE_OF_INSTRUCTION, FIELD_MISSING_IN_CALC, CYCLIC_DEPENDENCY, ERR_CALC_EXPRESSION, VAR_WILL_NEVER_BE_CALC"
//	@Failure		409				{object}	HTTPError		"possible error codes: VAR_ALREADY_USED"
//	@Failure		500				{object}	HTTPError		"possible error codes: REQUEST_ID_MISSING, INTERNAL_SERVER_ERROR"
//	@Router			/execute		[post]
func (gh *GinHandler) Execute(gctx *gin.Context) {
	var insts []Instruction
	if err := gctx.ShouldBindJSON(&insts); err != nil {
		gctx.JSON(http.StatusBadRequest, HTTPError{Code: "INVALID_JSON_BODY", Message: "invalid json body, " + err.Error()})
		return
	}

	gi := dto.GroupedInstructions{
		Expressions: make([]executor.Expression, 0),
		VarsToPrint: make(map[string]bool),
	}

	httpErr := gh.groupInstructions(&insts, &gi)
	if httpErr != nil {
		gctx.JSON(http.StatusBadRequest, httpErr)
		return
	}

	reqId, exists := gctx.Get(values.RequestIdKey)
	if !exists {
		gctx.JSON(http.StatusInternalServerError, HTTPError{Code: "REQUEST_ID_MISSING", Message: executor.ErrReqIdMissing.Error()})
		return
	}

	ctx := gctx.Request.Context()
	ctx = context.WithValue(ctx, values.RequestIdKey, reqId)

	results, err := gh.uc.Execute(ctx, &gi)
	if err != nil {
		gh.mapExecutorErrToHTTPErr(err, gctx)
		return
	}

	gctx.JSON(http.StatusOK, gh.resultsToItems(&results))
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

func (gh *GinHandler) groupInstructions(instructions *[]Instruction, gi *dto.GroupedInstructions) *HTTPError {
	for _, instruction := range *instructions {
		if instruction.Type == values.Calculate {
			httpErr := gh.validateCalcInst(&instruction)
			if httpErr != nil {
				return httpErr
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
		return &HTTPError{Code: "UNKNOWN_TYPE_OF_INSTRUCTION", Message: "unknown type of instruction (" + instruction.Type + ")"}
	}
	return nil
}

func (gh *GinHandler) validateCalcInst(instruction *Instruction) *HTTPError {
	if instruction.Operation == nil {
		return &HTTPError{Code: "FIELD_MISSING_IN_CALC", Message: "field 'op' is missing in calculate instruction"}
	}
	if instruction.Left == nil {
		return &HTTPError{Code: "FIELD_MISSING_IN_CALC", Message: "field 'left' is missing in calculate instruction"}
	}
	if instruction.Right == nil {
		return &HTTPError{Code: "FIELD_MISSING_IN_CALC", Message: "field 'right' is missing in calculate instruction"}
	}
	return nil
}

func (gh *GinHandler) mapExecutorErrToHTTPErr(err error, c *gin.Context) {
	switch {
	case errors.Is(err, executor.ErrCyclicDependency):
		c.JSON(http.StatusBadRequest, HTTPError{Code: "CYCLIC_DEPENDENCY", Message: err.Error()})
	case errors.Is(err, executor.ErrCalcExpression):
		c.JSON(http.StatusBadRequest, HTTPError{Code: "ERR_CALC_EXPRESSION", Message: err.Error()})
	case errors.Is(err, executor.ErrVarWillNeverBeCalc):
		c.JSON(http.StatusBadRequest, HTTPError{Code: "VAR_WILL_NEVER_BE_CALC", Message: err.Error()})
	case errors.Is(err, executor.ErrVarAlreadyUsed):
		c.JSON(http.StatusConflict, HTTPError{Code: "VAR_ALREADY_USED", Message: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, HTTPError{Code: "INTERNAL_SERVER_ERROR", Message: err.Error()})
	}
}
