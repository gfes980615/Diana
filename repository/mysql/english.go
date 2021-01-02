package mysql

import (
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/models/po"
	"gorm.io/gorm"
	"time"
)

func init() {
	injection.AutoRegister(&englishRepository{})
}

type englishRepository struct {
}

func (cr *englishRepository) Search(DB *gorm.DB) (*po.EnglishSentence, error) {
	content := &po.EnglishSentence{}
	err := DB.Where("date = ?", time.Now().Format("2006-01-02")).Find(content).Error
	return content, err
}
