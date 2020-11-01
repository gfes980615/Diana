package controller

import (
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/models/dto"
	"github.com/gfes980615/Diana/service"
	"github.com/gfes980615/Diana/transport/http/common"
	"github.com/gin-gonic/gin"
)

func init() {
	injection.AutoRegister(&travelController{})
}

type travelController struct {
	travelService service.TravelService `injection:"travelService"`
}

func (ctl *travelController) SetupRouter(router *gin.Engine) {
	travel := router.Group("/diana/travel")
	{
		// 取得桃園景點存進資料庫
		travel.GET("/taoyuan/list", ctl.getTaoyuanPlace)
		// 根據輸入的縣市從資料庫取得景點
		travel.GET("/area", ctl.getTravelPlaceByArea)
	}
}

func (ctl *travelController) getTaoyuanPlace(ctx *gin.Context) {
	if err := ctl.travelService.TaoyuanTravelPlace(); err != nil {
		common.Error(ctx, err)
		return
	}
	common.Send(ctx, "success")
}

func (ctl *travelController) getTravelPlaceByArea(ctx *gin.Context) {
	req := &dto.TravelReq{}
	if err := req.Init(ctx); err != nil {
		common.Error(ctx, err)
		return
	}
	result, err := ctl.travelService.GetTravelPlaceByArea(req.Country, req.Location)
	if err != nil {
		common.Error(ctx, err)
		return
	}
	common.Send(ctx, result)
}
