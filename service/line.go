package service

import (
	"errors"
	"github.com/gfes980615/Diana/db"
	"github.com/gfes980615/Diana/repository/mysql"
	"regexp"
	"strings"

	"github.com/gfes980615/Diana/glob"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/line"
	"github.com/line/line-bot-sdk-go/linebot"
)

func init() {
	injection.AutoRegister(&lineService{})
}

type lineService struct {
	currencyService    CurrencyService          `injection:"currencyService"`
	lineUserRepository mysql.LineUserRepository `injection:"lineUserRepository"`
}

func (ls *lineService) ReplyMessage(events []*linebot.Event) error {
	if len(events) != 1 {
		errMessage := "except error"
		return errors.New(errMessage)
	}
	var err error
	event := events[0]

	DB := db.MysqlConn.Session()
	go ls.lineUserRepository.Create(DB, event.Source.UserID)

	switch event.Type {
	case linebot.EventTypeMessage:
		err = ls.eventTypeMessage(event)
	}
	return err
}

func (ls *lineService) eventTypeMessage(event *linebot.Event) error {
	message, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		errMessage := "message type is not linebot.TextMessage"
		return errors.New(errMessage)
	}

	keyword := strings.TrimSpace(message.Text)
	if keyword == "a" {
		daily := line.GetEveryDaySentence()
		glob.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(daily)).Do()
		return nil
	}

	if keyword == "楓谷幣值" {
		glob.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(glob.MapleCurrencyMessage)).Do()
		return nil
	}

	if _, ok := glob.MapleServerMap[keyword]; ok {
		currencyValue := ls.currencyService.GetMapleCurrencyMessage(keyword)
		glob.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(currencyValue)).Do()
		return nil
	}

	//web := strings.Split(keyword, " ")
	//if (len(web) == 2 && IsURL(web[1])) || (len(web) == 1 && IsURL((web[0]))) {
	//	message := "save successful"
	//	err := line.SaveWebsite(web)
	//	if err != nil {
	//		message = fmt.Sprintf("save failed : %v", err)
	//	}
	//	glob.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message)).Do()
	//	return nil
	//}

	glob.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("no work keyword")).Do()
	// if message.Text == "maple story" {
	// 	maple := line.GetMapleStoryAnnouncement()
	// 	glob.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(maple)).Do()
	// 	return
	// }

	//id, transferErr := strconv.ParseInt(message.Text, 10, 64)
	//text := line.GetGoogleExcelValueById(id)
	//if transferErr != nil {
	//	if _, err = glob.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(transferErr.Error())).Do(); err != nil {
	//		log.Print(err)
	//	}
	//	return nil
	//}
	//if _, err = glob.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
	//	log.Print(err)
	//}
	return nil
}

func IsURL(url string) bool {
	isurl, _ := regexp.MatchString("http([s]*)://(.*?)", url)
	return isurl
}
