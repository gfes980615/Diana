package glob

import (
	"log"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	ChannelSecret      string
	ChannelAccessToken string
)

func init() {
	initEnvConfig()
	initLineBot()
}

func initEnvConfig() {
	ChannelSecret = os.Getenv("ChannelSecret")
	ChannelAccessToken = os.Getenv("ChannelAccessToken")
}

func initLineBot() {
	var err error
	Bot, err = linebot.New(ChannelSecret, ChannelAccessToken)
	log.Println("Bot:", Bot, " err:", err)
}
