package mysql

import "github.com/jinzhu/gorm"

type CurrencyRepository interface {
	Insert(DB *gorm.DB) error
}
