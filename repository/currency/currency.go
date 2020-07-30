package currency

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gfes980615/Diana/db"
	"github.com/gfes980615/Diana/glob"
	"github.com/gfes980615/Diana/model"
	"github.com/line/line-bot-sdk-go/linebot"
)

// CurrencyRepository ...
type CurrencyRepository struct {
}

// Insert maybe can use Create()?
// Insert 存入MYSQL
func (cr CurrencyRepository) InsertAndWarning(currencySlice []model.Currency, users []model.LineUser) error {
	mysql, err := db.NewMySQL(glob.DataBase)
	if err != nil {
		log.Print(err)
		return err
	}
	defer mysql.Close()

	avgValue, err := cr.getLastDayAvgValue(mysql)
	if err != nil {
		log.Print(err)
		return err
	}

	for _, c := range currencySlice {
		abnormal := 0
		if c.Value >= (avgValue * 2) {
			pushAbnormalCurrency(c.URL, users)
			abnormal = 1
		}
		err := mysql.DB.Exec("INSERT IGNORE INTO `currency` (`added_time`,`value`,`server`,`title`,`url`,`abnormal`) VALUES (NOW(),?,?,?,?,?)", c.Value, c.Server, c.Title, c.URL, abnormal)
		if err.Error != nil {
			return err.Error
		}
	}

	return nil
}

func (cr CurrencyRepository) GetChartData(subFunc string) ([]model.Currency, error) {
	mysql, err := db.NewMySQL(glob.DataBase)
	if err != nil {
		return nil, err
	}
	defer mysql.Close()

	sql := fmt.Sprintf("SELECT `added_time`, `server`, %s(value) as `value` FROM `currency` WHERE `abnormal` = 0 GROUP BY `added_time`, `server` ORDER BY `added_time` ASC", subFunc)
	currency := []model.Currency{}
	result := mysql.DB.Raw(sql).Scan(&currency)
	if result.Error != nil {
		return nil, result.Error
	}

	return currency, nil
}

// getLastDayAvgValue 取得前一日平均幣值
func (cr CurrencyRepository) getLastDayAvgValue(mysql *db.MySQL) (float64, error) {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	sql := fmt.Sprintf("SELECT avg(value) as `value` FROM `currency` where `abnormal` = 0 AND `added_time` = '%s'", yesterday)

	value := []model.Currency{}
	result := mysql.DB.Raw(sql).Scan(&value)
	if result.Error != nil {
		log.Print(result.Error)
		return 0, result.Error
	}

	if len(value) == 0 {
		log.Print(errors.New("no avg value"))
		return 0, errors.New("no avg value")
	}

	return value[0].Value, nil
}

func pushAbnormalCurrency(url string, users []model.LineUser) {
	for _, u := range users {
		message := "幣值異常，趕快來看看\n"
		message += url
		glob.Bot.PushMessage(u.UserID, linebot.NewTextMessage(message)).Do()
	}
}
