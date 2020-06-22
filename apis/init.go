package apis

import (
	"log"
	"os"

	"github.com/gfes980615/Diana/glob"
	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	ChannelSecret      string
	ChannelAccessToken string
)

func initEnvConfig() {
	ChannelSecret = os.Getenv("ChannelSecret")
	ChannelAccessToken = os.Getenv("ChannelAccessToken")
}

func initLineBot() {
	var err error
	glob.Bot, err = linebot.New(ChannelSecret, ChannelAccessToken)
	log.Println("Bot:", glob.Bot, " err:", err)
}

// InitSetting ...
func Init() {
	initEnvConfig()
	initLineBot()
}
