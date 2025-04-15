package gin

import (
	"Calculator/core"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerGin struct {
	useCase core.UseCase
}

func NewHandlerGin(uc core.UseCase) HandlerGin {
	return HandlerGin{useCase: uc}
}

func (hg *HandlerGin) Calculate(c *gin.Context) {
	var instructions []core.Instruction
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
