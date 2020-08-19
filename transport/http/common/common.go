package common

import (
	"github.com/gin-gonic/gin"
)

type Resp struct {
	Data interface{}
	Err  error
}

const (
	Success      = 200
	NotFound     = 404
	Unauthorized = 401
	ServerError  = 500
)

// Send ...
func Send(ctx *gin.Context, ret interface{}) {
	// convertutils.FloatRound(ret)

	ctx.JSON(Success, ret)
}

// Error ...
func Error(ctx *gin.Context, err error) {
	ctx.JSON(ServerError, map[string]string{"error": err.Error()})
}
