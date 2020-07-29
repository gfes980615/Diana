package message

import (
	"errors"
	"github.com/gfes980615/Diana/glob"
	"github.com/gfes980615/Diana/line"
	"github.com/line/line-bot-sdk-go/linebot"
	"strings"
)

func Message(events []*linebot.Event) error {
	if len(events) != 1 {
		errMessage := "except error"
		return errors.New(errMessage)
	}
	var err error
	event := events[0]
	go line.SaveUserID(event.Source.UserID)
	switch event.Type {
	case linebot.EventTypeMessage:
		err = eventTypeMessage(event)
	}
	return err
}

func eventTypeMessage(event *linebot.Event) error {
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

	if _, ok := glob.MapleServerMap[message.Text]; ok {
		currencyValue := line.GetMapleCurrencyMessage(message.Text)
		glob.Bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(currencyValue)).Do()
		return nil
	}
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
