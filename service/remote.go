package service

import (
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/utils"
	"github.com/gocolly/colly"
)

func init() {
	injection.AutoRegister(&remoteService{})
}

type remoteService struct {
}

func (w *remoteService) Simple(api string) {
	c := colly.NewCollector(colly.UserAgent(utils.RandomAgent()))
	c.Visit(api)
}
