package mysql

import (
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/models/po"
	"gorm.io/gorm"
	"strings"
)

func init() {
	injection.AutoRegister(&activityRepository{})
}

type activityRepository struct {
}

func (tr *activityRepository) CreateKKtixActivityItem(DB *gorm.DB, items []*po.KktixActivity) error {
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
