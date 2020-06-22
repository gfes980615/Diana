package apis

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gfes980615/Diana/glob"
	"github.com/gfes980615/Diana/line"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

func MainApis() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("Hello, It Home!"))
	})
	router.POST("/callback", callbackHandler)

	router.Run()
}

var bot = glob.Bot

func callbackHandler(c *gin.Context) {
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		log.Print(err.Error())
		if err == linebot.ErrInvalidSignature {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"reason": "values error.",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"reason": "values error.",
			})
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:

				if message.Text == "a" {
					daily := line.GetEveryDaySentence()
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(daily)).Do()
					return
				}

				id, transferErr := strconv.ParseInt(message.Text, 10, 64)
				text := line.GetGoogleExcelValueById(id)
				if transferErr != nil {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(transferErr.Error())).Do(); err != nil {
						log.Print(err)
					}
					return
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
