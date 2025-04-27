package gin

import (
	"Calculator/internal/executor"
	"Calculator/internal/executor/use_case"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerGin struct {
	useCase use_case.UseCase
}

func NewHandlerGin(uc use_case.UseCase) HandlerGin {
	return HandlerGin{useCase: uc}
}

func (hg *HandlerGin) Calculate(c *gin.Context) {
	var instructions []executor.InstructionGin
	if err := c.ShouldBindJSON(&instructions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	varsToPrint := make(map[string]bool)
	var expressions []executor.Expression

	for _, instruction := range instructions {
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

			expressions = append(expressions, expression)
		} else if instruction.Type == "print" {
			varsToPrint[instruction.Variable] = true
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("unknown type: " + instruction.Type)})
			return
		}
	}

	items, err := hg.useCase.Execute(expressions, varsToPrint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}
