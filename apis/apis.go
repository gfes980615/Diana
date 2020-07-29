package apis

import (
	"github.com/gfes980615/Diana/glob"
	"github.com/gfes980615/Diana/line"
	"github.com/gfes980615/Diana/line/message"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func MainApis() {
	// line.GetMapleCurrencyMessage("izr")

	router := gin.Default()
	router.LoadHTMLGlob("template/*")
	router.GET("/hello", func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("Hello, It Home!"))
	})
	router.POST("/callback", callbackHandler)
	router.POST("/currency", addCurrency)
	router.GET("/currency/chart/:subfunc", currencyChart)

	router.Run()
}

func currencyChart(c *gin.Context) {
	subfunc := c.Param("subfunc")
	result, err := line.GetMapleCurrencyChartData(subfunc)
	if err != nil {
		c.JSON(400, err)
		return
	}
	var category string
	switch subfunc {
	case "avg":
		category = "每日平均幣值"
	case "max":
		category = "每日最高幣值"
	default:

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
		"subfunc": category,
	}
	c.HTML(http.StatusOK, "maple_story.html", chartData)
}

func addCurrency(c *gin.Context) {
	line.AddAllServerCurrency()
	c.JSON(200, "開始蒐集資料")
}

func callbackHandler(c *gin.Context) {
	events, err := glob.Bot.ParseRequest(c.Request)
	if err != nil {
		log.Print(err.Error())
		return
	}
	message.Message(events)
}
