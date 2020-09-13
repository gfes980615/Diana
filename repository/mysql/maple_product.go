package mysql

import (
	"errors"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/models/po"
	"gorm.io/gorm"
)

func init() {
	injection.AutoRegister(&mapleProductRepository{})
}

type mapleProductRepository struct {
}

func (mp *mapleProductRepository) Insert(DB *gorm.DB, products []*po.Maple8591Product) error {
	return DB.Create(products).Error
}

func (mp *mapleProductRepository) CreateTable(DB *gorm.DB) error {
	var err string
	table := &po.Maple8591Product{}
	db := DB.Migrator()
	if !db.HasTable(table) {
		err = db.CreateTable(table).Error()
	}
	if len(err) > 0 {
		return errors.New(err)
	}
	return nil
}
