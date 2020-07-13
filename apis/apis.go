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
	router.LoadHTMLGlob("template/*")
	router.GET("/hello", func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("Hello, It Home!"))
	})
	router.POST("/callback", callbackHandler)
	router.POST("/currency", addCurrency)
	router.GET("/currency/chart", currencyChart)

	router.Run()
	// line.GetMapleCurrencyMessage("izr")
}

func currencyChart(c *gin.Context) {
	result, err := line.GetMapleCurrencyChartData()
	if err != nil {
		c.JSON(400, "error")
		return
	}
	c.HTML(http.StatusOK, "echarts.html", gin.H{
		"date": result.Date,
		"izcr": result.Izcr,
		"izr":  result.Izr,
		"ld":   result.Ld,
		"plt":  result.Plt,
		"slc":  result.Slc,
		"yen":  result.Yen,
		"ymax": result.YMax,
		"ymin": result.YMin,
	})
}

func addCurrency(c *gin.Context) {
	line.AddAllServerCurrency()
	c.JSON(200, "開始蒐集資料")
}

func callbackHandler(c *gin.Context) {
	events, err := glob.Bot.ParseRequest(c.Request)
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
					glob.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(daily)).Do()
					return
				}

				if message.Text == "maple story" {
					maple := line.GetMapleStoryAnnouncement()
					glob.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(maple)).Do()
					return
				}

				if message.Text == "楓谷幣值" {
					glob.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(glob.MapleCurrencyMessage)).Do()
					return
				}

				if _, ok := glob.MapleServerMap[message.Text]; ok {
					currencyValue := line.GetMapleCurrencyMessage(message.Text)
					glob.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(currencyValue)).Do()
					return
				}

				id, transferErr := strconv.ParseInt(message.Text, 10, 64)
				text := line.GetGoogleExcelValueById(id)
				if transferErr != nil {
					if _, err = glob.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(transferErr.Error())).Do(); err != nil {
						log.Print(err)
					}
					return
				}
				if _, err = glob.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
