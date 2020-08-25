package controller

import (
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/service"
	"github.com/gfes980615/Diana/transport/http/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	injection.AutoRegister(&CurrencyController{})
}

type CurrencyController struct {
	currencyService service.CurrencyService `injection:"currencyService"`
}

func (ctl *CurrencyController) SetupRouter(router *gin.Engine) {
	controller := router.Group("/diana")
	controller.GET("/currency/chart", ctl.currencyChart)
	controller.GET("/currency/value", ctl.currencyValue)
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

	common.Send(ctx, message)
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
