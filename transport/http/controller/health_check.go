package controller

import (
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/service"
	"github.com/gfes980615/Diana/transport/http/common"

	"github.com/gin-gonic/gin"
)

func init() {
	injection.AutoRegister(&HealthController{})
}

type HealthController struct {
	mapleBulletinService service.MapleBulletinService `injection:"mapleBulletinService"`
}

func (ctl *HealthController) SetupRouter(router *gin.Engine) {
	controller := router.Group("/diana")
	controller.GET("/health", ctl.HeathCheck)
	controller.GET("/test", ctl.forTest)
}

func (ctl *HealthController) HeathCheck(ctx *gin.Context) {
	ctx.JSON(common.Success, "ok")
}

func (ctl *HealthController) forTest(ctx *gin.Context) {
	message, err := ctl.mapleBulletinService.GetBulletinMessage()
	if err != nil {
		common.Error(ctx, err)
		return
	}
	ctx.JSON(common.Success, message)
}
