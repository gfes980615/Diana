package mysql

import (
	"github.com/gfes980615/Diana/glob/common/log"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/models/po"
	"github.com/jinzhu/gorm"
	"time"
)

func init() {
	injection.AutoRegister(&lineUserRepository{})
}

type lineUserRepository struct {
}

// GetAllUser ...
func (lr *lineUserRepository) GetAllUser(DB *gorm.DB) ([]*po.LineUser, error) {
	users := []*po.LineUser{}

	err := DB.Table("line_user").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (lr *lineUserRepository) Create(DB *gorm.DB, id string) {
	user := &po.LineUser{
		UserID:    id,
		AddedTime: time.Now(),
	}
	err := DB.Table("line_user").Create(user).Error
	log.Println(err)
	return
}
