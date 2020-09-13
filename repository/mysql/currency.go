package mysql

import (
	"fmt"
	"github.com/gfes980615/Diana/glob/common/log"
	"github.com/gfes980615/Diana/injection"
	"time"

	"github.com/gfes980615/Diana/models/po"
	"gorm.io/gorm"
)

func init() {
	injection.AutoRegister(&currencyRepository{})
}

type currencyRepository struct {
}

func (cr *currencyRepository) Insert(DB *gorm.DB, currencySlice []*po.Currency) error {
	for _, cur := range currencySlice {
		cur.AddedTime = time.Now()
		err := DB.Table("currency").Create(cur).Error
		if err != nil {
			log.Println("insert error: ", err)
		}
	}
	return nil
}

// getLastDayAvgValue 取得前一日平均幣值
func (cr *currencyRepository) GetLastDayAvgValue(DB *gorm.DB) (float64, error) {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	value := &po.Currency{}
	err := DB.Table("currency").Select("AVG(`value`) AS `value`").Where("`abnormal` = ? AND `added_time` = ?", 0, yesterday).Take(value).Error
	if err != nil {
		return 0, err
	}

	return value.Value, nil
}

func (cr *currencyRepository) GetCurrencyChartData(DB *gorm.DB) ([]*po.Currency, error) {
	//DB.Table("currency").Select("`added_time`,`server`,AVG(`value`) AS `value`")
	sql := fmt.Sprintf("SELECT `added_time`, `server`, avg(value) as `value` FROM `currency` WHERE `abnormal` = 0 GROUP BY `added_time`, `server` ORDER BY `added_time` ASC")
	currency := []*po.Currency{}
	err := DB.Raw(sql).Scan(&currency).Error
	if err != nil {
		return nil, err
	}

	return currency, nil
}

// Insert maybe can use Create()?
// Insert 存入MYSQL
//func (cr currencyRepository) InsertAndWarning(DB *gorm.DB, currencySlice []po.Currency, users []po.LineUser) error {
//	mysql, err := db.NewMySQL(glob.DataBase)
//	if err != nil {
//		log.Print(err)
//		return err
//	}
//	defer mysql.Close()
//
//	avgValue, err := cr.getLastDayAvgValue(mysql)
//	if err != nil {
//		log.Print(err)
//		return err
//	}
//
//	for _, c := range currencySlice {
//		abnormal := 0
//		if c.Value >= (avgValue * 2) {
//			pushAbnormalCurrency(c.URL, users)
//			abnormal = 1
//		}
//		err := mysql.DB.Exec("INSERT IGNORE INTO `currency` (`added_time`,`value`,`server`,`title`,`url`,`abnormal`) VALUES (NOW(),?,?,?,?,?)", c.Value, c.Server, c.Title, c.URL, abnormal)
//		if err.Error != nil {
//			return err.Error
//		}
//	}
//
//	return nil
//}

//// getLastDayAvgValue 取得前一日平均幣值
//func (cr currencyRepository) getLastDayAvgValue(mysql *db.MySQL) (float64, error) {
//	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
//	sql := fmt.Sprintf("SELECT avg(value) as `value` FROM `currency` where `abnormal` = 0 AND `added_time` = '%s'", yesterday)
//
//	value := []models.Currency{}
//	result := mysql.DB.Raw(sql).Scan(&value)
//	if result.Error != nil {
//		log.Print(result.Error)
//		return 0, result.Error
//	}
//
//	if len(value) == 0 {
//		log.Print(errors.New("no avg value"))
//		return 0, errors.New("no avg value")
//	}
//
//	return value[0].Value, nil
//}

//func pushAbnormalCurrency(url string, users []model.LineUser) {
//	for _, u := range users {
//		message := "幣值異常，趕快來看看\n"
//		message += url
//		glob.Bot.PushMessage(u.UserID, linebot.NewTextMessage(message)).Do()
//	}
//}
