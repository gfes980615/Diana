package service

import (
	"fmt"
	"github.com/gfes980615/Diana/glob/common/log"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/repository/redis"
	"github.com/gfes980615/Diana/utils"
)

func init() {
	injection.AutoRegister(&mapleBulletinService{})
}

type mapleBulletinService struct {
	mapleRedisRepository redis.MapleRedisRepository `injection:"mapleRedisRepository"`
	lineService          LineService                `injection:"lineService"`
}

func (mbs *mapleBulletinService) GetBulletinMessage() (string, error) {
	bulletin, err := mbs.mapleRedisRepository.GetBulletinData()
	if err != nil {
		return "", err
	}
	message := ""
	for _, b := range bulletin {
		message += fmt.Sprintf("日期:%s\n標題:%s\n分類:%s\n網址:%s\n----------------\n", b.Date, b.Title, b.Category, b.URL)
	}
	return message, nil
}

func (mbs *mapleBulletinService) PushToLine() {
	msg, err := mbs.GetBulletinMessage()
	if err != nil {
		log.Error(err)
		return
	}
	if utils.EmptyString(msg) {
		return
	}
	mbs.lineService.PushMessage(msg)
}
