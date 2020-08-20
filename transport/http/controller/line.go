package controller

import (
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
	lineService service.LineService `injection:"lineService"`
}

func (ctl *LineController) SetupRouter(router *gin.Engine) {
	router.GET("/callback", ctl.callbackHandler)
}

func (lc *LineController) callbackHandler(ctx *gin.Context) {
	events, err := glob.Bot.ParseRequest(ctx.Request)
	if err != nil {
		log.Print(err.Error())
		return
	}
	lc.lineService.ReplyMessage(events)
}
