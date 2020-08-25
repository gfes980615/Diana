package mysql

import (
	"github.com/gfes980615/Diana/models/po"
	"github.com/jinzhu/gorm"
)

type CurrencyRepository interface {
	Insert(DB *gorm.DB, currencySlice []*po.Currency) error
	GetLastDayAvgValue(DB *gorm.DB) (float64, error)
	GetCurrencyChartData(DB *gorm.DB) ([]*po.Currency, error)
}

type LineUserRepository interface {
	GetAllUser(DB *gorm.DB) ([]*po.LineUser, error)
	Create(DB *gorm.DB, id string)
}
