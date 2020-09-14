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
