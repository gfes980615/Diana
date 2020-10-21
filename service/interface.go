package service

import (
	"github.com/gfes980615/Diana/models/dto"
	"github.com/gfes980615/Diana/models/po"
	"github.com/line/line-bot-sdk-go/linebot"
)

type LineService interface {
	ReplyMessage(events []*linebot.Event) error
	GetActivityMessage() string
}

type CurrencyService interface {
	GetMapleCurrencyMessage(mapleServer string) string
	GetMapleCurrencyChartData() (*dto.ReturnSlice, error)
	GetAllServerCurrency() string
}

type SpiderService interface {
	GetPageSource(string, string) string
	GetAllCount(string) int
	GetEveryDaySentence() string
}

type Maple8591ProductService interface {
	Get8591AllProduct()
}

type ActivityService interface {
	GetKktixActivity(category string) []*po.KktixActivity
	GetTravelTaipeiActivity(category string) []*po.TTActivity
}
