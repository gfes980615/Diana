package mysql

import (
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/models/po"
	"gorm.io/gorm"
	"strings"
)

func init() {
	injection.AutoRegister(&travelRepository{})
}

type travelRepository struct {
}

func (tr *travelRepository) CreateTravelTaipeiTravelItem(DB *gorm.DB, items []*po.TravelList) error {
	for _, item := range items {
		err := DB.Create(item).Error
		if err == nil {
			continue
		}
		if !strings.HasPrefix(err.Error(), "Error 1062:") {
			return err
		}
	}
	return nil
}
