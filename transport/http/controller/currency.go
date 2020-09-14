package controller

import (
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/models/dto"
	"github.com/gfes980615/Diana/service"
	"github.com/gfes980615/Diana/transport/http/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	injection.AutoRegister(&CurrencyController{})
}

type CurrencyController struct {
	currencyService         service.CurrencyService         `injection:"currencyService"`
	maple8591ProductService service.Maple8591ProductService `injection:"maple8591ProductService"`
}

func (ctl *CurrencyController) SetupRouter(router *gin.Engine) {
	router.LoadHTMLGlob("template/*")
	controller := router.Group("/diana/currency")
	controller.GET("/chart", ctl.currencyChart)
	controller.GET("/value", ctl.currencyValue)
	controller.GET("/test", ctl.get8591AllProduct)
}

//func (ctl *CurrencyController) currencyChart(ctx *gin.Context) {
//	ctx.JSON(common.Success, "ok")
//}

func (ctl *CurrencyController) currencyValue(ctx *gin.Context) {
	type Tmp struct {
		Server string `form:"server"`
	}
	t := &Tmp{}
	ctx.Bind(t)

	var message string
	switch t.Server {
	case "all":
		message = ctl.currencyService.GetAllServerCurrency()
	default:
		message = ctl.currencyService.GetMapleCurrencyMessage(t.Server)
	}

	testMessage := &dto.Message{
		Message: message,
	}

	common.Send(ctx, testMessage)
}

func (ctl *CurrencyController) currencyChart(c *gin.Context) {
	result, err := ctl.currencyService.GetMapleCurrencyChartData()
	if err != nil {
		c.JSON(400, err)
		return
	}

	chartData := map[string]interface{}{
		"date":    result.Date,
		"izcr":    result.Izcr,
		"izr":     result.Izr,
		"ld":      result.Ld,
		"plt":     result.Plt,
		"slc":     result.Slc,
		"yen":     result.Yen,
		"ymax":    result.YMax,
		"ymin":    result.YMin,
		"subfunc": "每日平均幣值",
	}
	c.HTML(http.StatusOK, "maple_story.html", chartData)
}

func (ctl *CurrencyController) get8591AllProduct(c *gin.Context) {
	ctl.maple8591ProductService.Get8591AllProduct()
	common.Send(c, "ok")
}