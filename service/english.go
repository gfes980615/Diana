package service

import (
	"fmt"
	"github.com/gfes980615/Diana/db"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/repository/mysql"
)

func init() {
	injection.AutoRegister(&englishService{})
}

const (
	managementEnglishPage = "https://www.managertoday.com.tw/quotes?page=%d" // 經理人每日一句學管理
)

type englishService struct {
	lineService       LineService             `injection:"lineService"`
	englishRepository mysql.EnglishRepository `injection:"englishRepository"`
}

func (es *englishService) SendDailyMessage() error {
	DB := db.MysqlConn.Session()
	content, err := es.englishRepository.Search(DB)
	if err != nil {
		return err
	}
	message := fmt.Sprintf("%s\n%s\n%s", content.Date.Format("2006-01-02"), content.Content, content.Translation)
	es.lineService.PushMessage(message)
	return nil
}
