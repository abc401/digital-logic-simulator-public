package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Noop(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"operation": "noop"})
}
