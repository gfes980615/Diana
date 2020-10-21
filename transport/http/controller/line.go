package controller

import (
	"log"

	"github.com/gfes980615/Diana/models/dto"
	"github.com/gfes980615/Diana/transport/http/common"

	"github.com/gfes980615/Diana/service"

	"github.com/gfes980615/Diana/glob"
	"github.com/gfes980615/Diana/injection"

	"github.com/gin-gonic/gin"
)

func init() {
	injection.AutoRegister(&LineController{})
}

type LineController struct {
	currencyService service.CurrencyService `injection:"currencyService"`
	lineService     service.LineService     `injection:"lineService"`
	spiderService   service.SpiderService   `injection:"spiderService"`
}

func (lc *LineController) SetupRouter(router *gin.Engine) {
	router.GET("/callback", lc.callbackHandler)
	router.GET("/daily/sentence", lc.Daily)
	router.GET("/daily/currency/message", lc.dailyCurrencyMessage)
	router.GET("/debug", lc.Debug)
}

func (lc *LineController) callbackHandler(ctx *gin.Context) {
	events, err := glob.Bot.ParseRequest(ctx.Request)
	if err != nil {
		log.Print(err.Error())
		return
	}
	lc.lineService.ReplyMessage(events)
}

func (lc *LineController) Daily(ctx *gin.Context) {
	common.Send(ctx, lc.spiderService.GetEveryDaySentence())
}

func (lc *LineController) dailyCurrencyMessage(ctx *gin.Context) {
	message, err := lc.currencyService.GetDailyMessage()
	if err != nil {
		common.Send(ctx, err)
		return
	}

	testMessage := &dto.Message{
		Message: message,
	}
	common.Send(ctx, testMessage)
}

func (lc *LineController) Debug(ctx *gin.Context) {
	common.Send(ctx, lc.lineService.GetActivityMessage())
}
