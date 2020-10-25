package controller

import (
	"encoding/json"
	"github.com/line/line-bot-sdk-go/linebot"
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
	router.GET("/callback_lineTemplate", lc.callbackLineTemplateHandler)
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

type LineMessage struct {
	Events string `json:"events" form:"events"`
}

func (lc *LineController) callbackLineTemplateHandler(ctx *gin.Context) {
	conds := &LineMessage{}
	err := ctx.ShouldBind(conds)
	if err != nil {
		common.Error(ctx, err)
		log.Print(err)
		return
	}
	events := []*linebot.Event{}
	err = json.Unmarshal([]byte(conds.Events), &events)
	if err != nil {
		common.Error(ctx, err)
		log.Print(err)
		return
	}

	err = lc.lineService.ReplyMessage(events)
	if err != nil {
		common.Error(ctx, err)
		log.Print(err)
		return
	}
	common.Send(ctx, "ok")
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
