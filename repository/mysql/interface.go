package mysql

import (
	"github.com/gfes980615/Diana/models/po"
	"gorm.io/gorm"
)

type CurrencyRepository interface {
	Insert(DB *gorm.DB, currencySlice []*po.Currency) error
	GetLastDayAvgValue(DB *gorm.DB) (float64, error)
	GetCurrencyChartData(DB *gorm.DB) ([]*po.Currency, error)
	GetDailyItems(DB *gorm.DB) (map[string][]*po.Currency, error)
}

type LineUserRepository interface {
	GetAllUser(DB *gorm.DB) ([]*po.LineUser, error)
	Create(DB *gorm.DB, id string)
}

type Maple8591ProductRepository interface {
	Insert(DB *gorm.DB, products []*po.Maple8591Product) error
	CreateTable(DB *gorm.DB) error
}

type TravelRepository interface {
	CreateTravelTaipeiTravelItem(DB *gorm.DB, items []*po.TravelList) error
	CreateTaoyuanTravelItem(DB *gorm.DB, items []*po.TouristAttractionList) error
	GetTravelListByArea(DB *gorm.DB, country, loction string) ([]*po.TouristAttractionList, error)
	GetAllTravelList(DB *gorm.DB) ([]*po.TouristAttractionList, error)
}

type ActivityRepository interface {
	CreateKKtixActivityItem(DB *gorm.DB, items []*po.KktixActivity) error
}

type EnglishRepository interface {
	Search(DB *gorm.DB) (*po.EnglishSentence, error)
}
