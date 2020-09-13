package service

import (
	"fmt"
	"github.com/gfes980615/Diana/db"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/models/po"
	"github.com/gfes980615/Diana/repository/mysql"
	"github.com/gocolly/colly"
	"strings"
	"time"
	"unicode"
)

func init() {
	injection.AutoRegister(&maple8591ProductService{})
}

type maple8591ProductService struct {
	mapleProductRepository mysql.Maple8591ProductRepository `injection:"mapleProductRepository"`
}

func (mp *maple8591ProductService) Get8591AllProduct() {
	s := time.Now()
	DB := db.MysqlConn.Session()
	if err := mp.mapleProductRepository.CreateTable(DB); err != nil {
		return
	}

	products := make([]*po.Maple8591Product, 1000000)
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64)"),
	)
	number := 1
	c.OnHTML("ul#wc_list.clearfix > li", func(e *colly.HTMLElement) {
		t := time.Now()
		url := e.ChildAttr("div > a.detail_link", "href")
		url = strings.Replace(url, "?", "#", 1)
		title := removeExtraChar(e.ChildText("div > a.detail_link > span.ml-item-title"))
		tmp := &po.Maple8591Product{
			Title:     title,
			URL:       url,
			Pageviews: e.ChildText("div.creatTime > span.ListHour"),
		}
		//fmt.Println("標題:\t", e.ChildText("div > a.detail_link > span.ml-item-title"))
		//fmt.Println("網址:\t", e.ChildAttr("div > a.detail_link", "href"))
		e.ForEach("div.other", func(i int, el *colly.HTMLElement) {
			switch i {
			case 0:
				tmp.Amount = el.Text
				//fmt.Println("金額:\t", el.Text)
			case 1:
				tmp.Number = el.Text
				//fmt.Println("庫存:\t", el.Text)
			}
		})
		//fmt.Println("瀏覽量:\t", e.ChildText("div.creatTime > span.ListHour"))
		//fmt.Println(number)

		products = append(products, tmp)
		fmt.Printf("number %d done for %v time\n", number, time.Since(t))
		number++
	})
	// 下一页
	c.OnHTML("a[href].pageNum", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	// 启动
	c.Visit("https://www.8591.com.tw/mallList-list.html?id=859&%251=&gst=1&searchGame=859")

	if err := mp.mapleProductRepository.Insert(DB, products); err != nil {
		return
	}

	fmt.Println(time.Since(s))
}

func removeExtraChar(title string) string {
	var s []int32
	for _, t := range title {
		if unicode.Is(unicode.Han, t) || unicode.IsDigit(t) || unicode.IsLetter(t) {
			s = append(s, t)
		}
	}
	return string(s)
}

func replaceQuestionMark(url string) string {
	return strings.Replace(url, "?", "^", 1)
}

func recoverQuestionMark(url string) string {
	return strings.Replace(url, "^", "?", 1)
}
