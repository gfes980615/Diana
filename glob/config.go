package glob

import (
	"log"
	"os"

	"github.com/gfes980615/Diana/models"
	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	ChannelSecret      string
	ChannelAccessToken string
	DataBase           models.DataBaseConfig
)

func init() {
	initEnvConfig()
	initLineBot()
	initDB()
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

func initDB() {
	DataBase.Username = os.Getenv("mysql_username")
	DataBase.Password = os.Getenv("mysql_password")
	DataBase.Address = os.Getenv("mysql_address")
	DataBase.Database = os.Getenv("mysql_database")
}
