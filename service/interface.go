package service

import (
	"github.com/gfes980615/Diana/models/dto"
	"github.com/line/line-bot-sdk-go/linebot"
)

type LineService interface {
	ReplyMessage(events []*linebot.Event) error
	PushMessage(message string)
}

type CurrencyService interface {
	GetMapleCurrencyMessage(mapleServer string) string
	GetMapleCurrencyChartData() (*dto.ReturnSlice, error)
	GetAllServerCurrency() string
	GetDailyMessage() (string, error)
}

type SpiderService interface {
	GetPageSource(string, string) string
	GetAllCount(string) int
	GetEveryDaySentence() string
}

type Maple8591ProductService interface {
	Get8591AllProduct()
}
