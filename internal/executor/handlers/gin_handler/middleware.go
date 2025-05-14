package gin_handler

import (
	"Calculator/internal/executor/values"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ReqIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.NewString()
		c.Set(values.RequestIdKey, id)
		c.Next()
	}
}
