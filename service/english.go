package service

import (
	"fmt"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/utils"
	"github.com/gocolly/colly"
)

func init() {
	injection.AutoRegister(&englishService{})
}

const (
	managementEnglishPage = "https://www.managertoday.com.tw/quotes?page=%d" // 經理人每日一句學管理
)

type englishService struct {
}

func (es *englishService) GetDailySentence() {
	es.dailySentence("", 0)
}

//> div.container > div.my-rwd-container > div.row > div.col
func (es *englishService) dailySentence(url string, page int) {
	c := colly.NewCollector(colly.UserAgent(utils.RandomAgent()))
	c.OnHTML("div.body-container > div.v-application", func(e *colly.HTMLElement) {
		//fmt.Println(e.ChildText("a.quote_box > div.pa-4 > h2.title_sty02"))
		//fmt.Println(e.ChildText("a.quote_box > div.pa-4 > span"))
		//fmt.Println(e.ChildText("a.quote_box > div.pa-4 > div.d-flex > div.text-right > div.title_syt04"))
		//fmt.Println(e.ChildText("a.quote_box > div.pa-4 > div.d-flex > div.text-right > span"))
		fmt.Println(e.Text)
	})

	c.Visit("https://www.managertoday.com.tw/quotes?page=1")
}
