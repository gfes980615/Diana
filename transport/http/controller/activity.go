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
	controller.GET("/kktix", ctl.getKktixActivity)
	controller.GET("/travel_taipei", ctl.getTravelTaipeiActivity)
}

func (ctl *ActivityController) getKktixActivity(ctx *gin.Context) {
	ctl.activityService.GetKktixActivity()
}

func (ctl *ActivityController) getTravelTaipeiActivity(ctx *gin.Context) {
	ctl.activityService.GetTravelTaipeiActivity()
}
