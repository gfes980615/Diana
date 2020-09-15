package service

import (
	"fmt"
	"github.com/gfes980615/Diana/db"
	"github.com/gfes980615/Diana/glob"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/models"
	"github.com/gfes980615/Diana/models/dto"
	"github.com/gfes980615/Diana/models/po"
	"github.com/gfes980615/Diana/repository/mysql"
	"github.com/gocolly/colly"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"regexp"
	"sort"
	"strconv"
	"sync"
)

func init() {
	injection.AutoRegister(&currencyService{})
}

const (
	raw            = 21
	root8591       = "https://www.8591.com.tw"
	URL8591        = "https://www.8591.com.tw/mallList-list.html?&group=1&searchType=0&priceSort=0&ratios=0&searchGame=859&searchServer=%d&firstRow=%d"
	regexpA        = `<a class="detail_link" href="(.*?)" target="_blank" class="goods_link" title="(.*?)" data-key="0">`
	rCount         = `<span class="R">([0-9]+)</span>`
	rCurrencyValue = `【([0-9]+):(.*?)萬】`
)

type currencyService struct {
	currencyRepository mysql.CurrencyRepository `injection:"currencyRepository"`
	lineUserRepository mysql.LineUserRepository `injection:"lineUserRepository"`
	spiderService      SpiderService            `injection:"spiderService"`
	lineService        LineService              `injection:"lineService"`
}

func (cs *currencyService) GetMapleCurrencyMessage(mapleServer string) string {
	currencySlice, count := cs.get8591CurrencyValueTop5V2(mapleServer)

	message := fmt.Sprintf("%s\n85幣值前五(共 %d 筆):\n", mapleServer, count)
	for i := 0; i < len(currencySlice); i++ {
		message += fmt.Sprintf("1 : %.f萬\n", currencySlice[i])
	}
	return message
	//cs.get8591CurrencyValueTop5V2(mapleServer)

	//return ""
}

func (cs *currencyService) get8591CurrencyValueTop5V2(mapleServer string) ([]float64, int) {
	products := []*po.Maple8591Product{}

	urlMap := make(map[string]bool)

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64)"),
	)

	c.OnHTML("ul#wc_list.clearfix > li", func(e *colly.HTMLElement) {
		url := e.ChildAttr("div > a.detail_link", "href")
		if _, ok := urlMap[url]; ok {
			return
		} else {
			urlMap[url] = true
		}

		tmp := &po.Maple8591Product{
			Title:     e.ChildText("div > a.detail_link > span.ml-item-title"),
			URL:       url,
			Pageviews: e.ChildText("div.creatTime > span.ListHour"),
		}

		e.ForEach("div.other", func(i int, el *colly.HTMLElement) {
			switch i {
			case 0:
				tmp.Amount = el.Text
			case 1:
				tmp.Number = el.Text
			}
		})
		tmp.Server = mapleServer

		products = append(products, tmp)
	})

	visitUrlMap := make(map[string]bool)
	c.OnHTML("a[href].pageNum", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		if _, ok := visitUrlMap[url]; ok {
			return
		} else {
			visitUrlMap[url] = true
		}

		e.Request.Visit(url)
	})
	count := 0
	c.OnHTML("span.R", func(e *colly.HTMLElement) {
		c, _ := strconv.Atoi(e.Text)
		count = c
	})

	// 启动
	c.Visit(fmt.Sprintf(URL8591, glob.MapleServerMap[mapleServer], 0))

	currencySlice := cs.setProductToCurrency(products)

	defaultSize := 5
	if len(currencySlice) < defaultSize {
		defaultSize = len(currencySlice)
	}
	cResultSlice := currencySlice[0:defaultSize]
	// 存入MYSQL
	go cs.insertCurrency(cResultSlice)

	result := []float64{}
	for _, c := range cResultSlice {
		result = append(result, c.Value)
	}

	return result, count
}

func (cs *currencyService) setProductToCurrency(products []*po.Maple8591Product) []*po.Currency {
	currencySlice := []*po.Currency{}
	for _, product := range products {
		rp := regexp.MustCompile(rCurrencyValue)
		items := rp.FindAllStringSubmatch(product.Title, -1)
		var value float64
		for _, item := range items {
			a, _ := strconv.ParseFloat(item[1], 64)
			b, _ := strconv.ParseFloat(item[2], 64)
			value = b / a
		}
		currencySlice = append(currencySlice, &po.Currency{Value: value, Server: product.Server, Title: removeExtraChar(product.Title), URL: root8591 + product.URL})
	}

	sort.Slice(currencySlice, func(i, j int) bool {
		return currencySlice[i].Value > currencySlice[j].Value
	})

	return currencySlice
}

// get8591CurrencyValueTop5 ...
func (cs *currencyService) get8591CurrencyValueTop5(mapleServer string) ([]float64, int) {
	// 取得所有數量
	resultPage := cs.spiderService.GetPageSource(fmt.Sprintf(URL8591, glob.MapleServerMap[mapleServer], 0), UTF8)
	count := cs.spiderService.GetAllCount(resultPage)
	titleArray := []models.URLStruct{}
	pageResult := make(chan string, count/raw+1)
	// 取得每條商品資訊
	for i := 0; i < count; i = i + raw {
		go func(index int) {
			tmpResult := cs.spiderService.GetPageSource(fmt.Sprintf(URL8591, glob.MapleServerMap[mapleServer], index), UTF8)
			pageResult <- tmpResult
		}(i)
	}

	for i := 0; i < count; i = i + raw {
		select {
		case tmp := <-pageResult:
			rp := regexp.MustCompile(regexpA)
			items := rp.FindAllStringSubmatch(tmp, -1)
			for _, item := range items {
				title := models.URLStruct{
					URL:  root8591 + item[1],
					Name: item[2],
				}
				titleArray = append(titleArray, title)
			}
		}
	}

	// 將幣值取出，string -> float 存成 array
	currencySlice := []*po.Currency{}
	for _, title := range titleArray {
		rp := regexp.MustCompile(rCurrencyValue)
		items := rp.FindAllStringSubmatch(title.Name, -1)
		var value float64
		for _, item := range items {
			a, _ := strconv.ParseFloat(item[1], 64)
			b, _ := strconv.ParseFloat(item[2], 64)
			value = b / a
		}
		name := removeExtraChar(title.Name)
		url := replaceQuestionMark(title.URL)
		currencySlice = append(currencySlice, &po.Currency{Value: value, Server: mapleServer, Title: name, URL: url})
	}

	sort.Slice(currencySlice, func(i, j int) bool {
		return currencySlice[i].Value > currencySlice[j].Value
	})

	defaultSize := 5
	if len(currencySlice) < defaultSize {
		defaultSize = len(currencySlice)
	}
	cResultSlice := currencySlice[0:defaultSize]
	// 存入MYSQL
	go cs.insertCurrency(cResultSlice)

	result := []float64{}
	for _, c := range cResultSlice {
		result = append(result, c.Value)
	}

	return result, count
}

// 幣值異常時通知用戶
func (cs currencyService) insertCurrency(cresult []*po.Currency) {
	DB := db.MysqlConn.Session()
	lineUsers, err := cs.lineUserRepository.GetAllUser(DB)
	if err != nil {
		log.Println("Get line user error: ", err)
		return
	}
	err = cs.insertAndWarning(cresult, lineUsers)
	if err != nil {
		log.Println("Insert error: ", err)
	}
}

func (cs currencyService) insertAndWarning(currencySlice []*po.Currency, users []*po.LineUser) error {
	DB := db.MysqlConn.Session()

	avgValue, err := cs.currencyRepository.GetLastDayAvgValue(DB)
	if err != nil {
		return err
	}

	for _, c := range currencySlice {
		abnormal := 0
		if c.Value >= (avgValue * 2) {
			cs.pushAbnormalCurrency(c, users)
			abnormal = 1
		}
		c.Abnormal = abnormal
		tmp := []*po.Currency{}
		tmp = append(tmp, c)
		err := cs.currencyRepository.Insert(DB, tmp)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cs currencyService) pushAbnormalCurrency(product *po.Currency, users []*po.LineUser) {
	messageFormat := "%s : 標題幣值異常\n1:%.f\n%s"
	replyMessage := fmt.Sprintf(messageFormat, glob.MapleServerZH[product.Server], product.Value, product.URL)

	for _, u := range users {
		glob.Bot.PushMessage(u.UserID, linebot.NewTextMessage(replyMessage)).Do()
	}
}

// GetMapleCurrencyChartData ...
func (cs currencyService) GetMapleCurrencyChartData() (*dto.ReturnSlice, error) {
	r := &dto.ReturnSlice{}
	DB := db.MysqlConn.Session()
	currency, err := cs.currencyRepository.GetCurrencyChartData(DB)
	if err != nil {
		return nil, err
	}

	tmpDateMap := make(map[string]bool)
	min := currency[0].Value
	max := currency[0].Value
	for _, item := range currency {
		date := item.AddedTime.Format("2006-01-02")
		if _, exist := tmpDateMap[date]; !exist {
			r.Date = append(r.Date, date)
			tmpDateMap[date] = true
		}
		if item.Value > max {
			max = item.Value
		}
		if item.Value < min {
			min = item.Value
		}
		switch item.Server {
		case "izcr":
			r.Izcr = append(r.Izcr, glob.FloatRound(item.Value))
		case "izr":
			r.Izr = append(r.Izr, glob.FloatRound(item.Value))
		case "ld":
			r.Ld = append(r.Ld, glob.FloatRound(item.Value))
		case "plt":
			r.Plt = append(r.Plt, glob.FloatRound(item.Value))
		case "slc":
			r.Slc = append(r.Slc, glob.FloatRound(item.Value))
		case "yen":
			r.Yen = append(r.Yen, glob.FloatRound(item.Value))
		}
	}

	r.YMax = int(max/10)*10 + 10
	r.YMin = int(min/10)*10 - 10

	return r, nil
}

func (cs currencyService) GetAllServerCurrency() string {
	serverMap := make(map[string]string)
	var wg sync.WaitGroup
	for server, _ := range glob.MapleServerMap {
		wg.Add(1)
		go func(server string) {
			defer wg.Done()
			message := cs.GetMapleCurrencyMessage(server)
			if _, ok := serverMap[server]; !ok {
				serverMap[server] = message
			}
		}(server)
	}
	wg.Wait()

	var resultMessage string
	for server, message := range serverMap {
		resultMessage += server + " : " + message
	}

	return resultMessage
}

func (cs currencyService) GetDailyMessage() (string, error) {
	DB := db.MysqlConn.Session()
	dailyItems, err := cs.currencyRepository.GetDailyItems(DB)
	if err != nil {
		return "", err
	}
	returnMessage := "每日提醒:\n"
	for key, items := range dailyItems {
		dailyMessageFormat := "%s : %.f\n"
		switch key {
		case "max":
			returnMessage += "\n最大幣值:\n"
			for _, item := range items {
				returnMessage += fmt.Sprintf(dailyMessageFormat, glob.MapleServerZH[item.Server], item.Value)
			}
		case "avg":
			returnMessage += "\n平均幣值:\n"
			for _, item := range items {
				returnMessage += fmt.Sprintf(dailyMessageFormat, glob.MapleServerZH[item.Server], item.Value)
			}
		}
	}
	cs.lineService.PushMessage(returnMessage)
	return returnMessage, nil
}
