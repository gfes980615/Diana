package controller

import (
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/service"
	"github.com/gin-gonic/gin"
)

func init() {
	injection.AutoRegister(&ActivityController{})
}

type ActivityController struct {
	activityService service.ActivityService `injection:"activityService"`
}

func (ctl *ActivityController) SetupRouter(router *gin.Engine) {
	controller := router.Group("/diana/activity")
	controller.GET("/kktix/:category", ctl.getKktixActivity)
	controller.GET("/travel_taipei/:category", ctl.getTravelTaipeiActivity)
}

func (ctl *ActivityController) getKktixActivity(ctx *gin.Context) {
	category:=ctx.Param("category")
	ctl.activityService.GetKktixActivity(category)
}

func (ctl *ActivityController) getTravelTaipeiActivity(ctx *gin.Context) {
	category:=ctx.Param("category")
	ctl.activityService.GetTravelTaipeiActivity(category)
}
