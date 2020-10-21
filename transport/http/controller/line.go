package controller

import (
	"github.com/gfes980615/Diana/transport/http/common"
	"log"

	"github.com/gfes980615/Diana/service"

	"github.com/gfes980615/Diana/glob"
	"github.com/gfes980615/Diana/injection"

	"github.com/gin-gonic/gin"
)

func init() {
	injection.AutoRegister(&LineController{})
}

type LineController struct {
	lineService   service.LineService   `injection:"lineService"`
	spiderService service.SpiderService `injection:"spiderService"`
}

func (lc *LineController) SetupRouter(router *gin.Engine) {
	router.GET("/callback", lc.callbackHandler)
	router.GET("/daily/sentence", lc.Daily)
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

func (lc *LineController) Debug(ctx *gin.Context) {
	common.Send(ctx, lc.lineService.GetActivityMessage())
}
