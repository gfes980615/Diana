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

func (tr *travelRepository) CreateTaoyuanTravelItem(DB *gorm.DB, items []*po.TouristAttractionList) error {
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

func (tr *travelRepository) GetTravelListByArea(DB *gorm.DB, country, loction string) ([]*po.TouristAttractionList, error) {
	items := []*po.TouristAttractionList{}

	result := DB.Debug().Where("country = ? and location = ?", country, loction).Find(&items)
	if err := result.Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (tr *travelRepository) GetAllTravelList(DB *gorm.DB) ([]*po.TouristAttractionList, error) {
	items := []*po.TouristAttractionList{}

	result := DB.Find(&items)
	if err := result.Error; err != nil {
		return nil, err
	}

	return items, nil
}
