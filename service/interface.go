package service

import "github.com/line/line-bot-sdk-go/linebot"

// LineService ...
type LineService interface {
	ReplyMessage(events []*linebot.Event) error
}
