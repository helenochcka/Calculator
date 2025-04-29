package gin_handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ReqIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.NewString()
		c.Set("request_id", id)
		c.Next()
	}
}
