package controller

import (
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/transport/http/common"

	"github.com/gin-gonic/gin"
)

func init() {
	injection.AutoRegister(&HealthController{})
}

type HealthController struct {
}

func (ctl *HealthController) SetupRouter(router *gin.Engine) {
	controller := router.Group("/diana")
	controller.GET("/health", ctl.HeathCheck)
}

func (ctl *HealthController) HeathCheck(ctx *gin.Context) {
	ctx.JSON(common.Success, "ok")
}
