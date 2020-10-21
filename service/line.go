package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/gfes980615/Diana/db"
	"github.com/gfes980615/Diana/glob/common/log"
	"github.com/gfes980615/Diana/models/dto"
	"github.com/gfes980615/Diana/repository/mysql"

	"github.com/gfes980615/Diana/glob"
	"github.com/gfes980615/Diana/injection"
	"github.com/line/line-bot-sdk-go/linebot"
)

func init() {
	injection.AutoRegister(&lineService{})
}

type lineService struct {
	currencyService    CurrencyService          `injection:"currencyService"`
	spiderService      SpiderService            `injection:"spiderService"`
	lineUserRepository mysql.LineUserRepository `injection:"lineUserRepository"`
	activityService    ActivityService          `injection:"activityService"`
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
		daily := ls.spiderService.GetEveryDaySentence()
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

	if keyword == "活動" {
		glob.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(ls.GetActivityMessage())).Do()
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

func (ls *lineService) GetActivityMessage() string {
	activitys := []*dto.Activity{}
	ttActivity := ls.activityService.GetTravelTaipeiActivity("exhibition")
	for _, a := range ttActivity {
		tmp := &dto.Activity{
			Title: a.Title,
			URL:   a.URL,
			Time:  a.ActivityTime,
		}
		activitys = append(activitys, tmp)
	}
	kkActivity := ls.activityService.GetKktixActivity("exhibition")
	for _, a := range kkActivity {
		tmp := &dto.Activity{
			Title: a.Title,
			URL:   a.URL,
			Time:  a.ActivityTime,
		}
		activitys = append(activitys, tmp)
	}
	message := ""
	for _, a := range activitys {
		message += fmt.Sprintf("活動:%s\n時間:%s\n網站:%s\n\n", a.Title, a.Time, a.URL)
	}
	return message
}

func IsURL(url string) bool {
	isurl, _ := regexp.MatchString("http([s]*)://(.*?)", url)
	return isurl
}

func (ls *lineService) PushMessage(message string) {
	DB := db.MysqlConn.Session()
	lineUsers, err := ls.lineUserRepository.GetAllUser(DB)
	if err != nil {
		log.Errorf("get line user error %v", err)
		return
	}

	for _, u := range lineUsers {
		glob.Bot.PushMessage(u.UserID, linebot.NewTextMessage(message)).Do()
	}
}
