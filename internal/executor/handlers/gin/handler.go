package gin

import (
	"Calculator/internal/executor"
	"Calculator/internal/executor/use_case"
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
	var instructions []executor.Instruction
	if err := c.ShouldBindJSON(&instructions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items, err := hg.useCase.Execute(instructions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}
