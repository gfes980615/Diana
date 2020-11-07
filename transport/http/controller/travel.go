package controller

import (
	"github.com/gfes980615/Diana/glob/common/log"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/models/dto"
	"github.com/gfes980615/Diana/service"
	"github.com/gfes980615/Diana/transport/http/common"
	"github.com/gfes980615/Diana/utils"
	"github.com/gin-gonic/gin"
	"net/http"
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
		travel.GET("/get_lat_lng", ctl.getLatLng)
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
	start := utils.TraceMemStats()
	result, err := ctl.travelService.GetTravelPlaceByArea(req.Country, req.Location)
	if err != nil {
		common.Error(ctx, err)
		return
	}
	end := utils.TraceMemStats()
	log.Infof("size of tourist items :%d (bytes)", end-start)
	common.Send(ctx, result)
}

func (ctl *travelController) getLatLng(ctx *gin.Context) {
	url := "http://127.0.0.1:8000/fenrir/test"
	_, err := http.Get(url)
	if err != nil {
		common.Error(ctx, err)
		return
	}
	common.Send(ctx, "success")
}
