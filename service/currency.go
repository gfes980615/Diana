package service

import (
	"fmt"
	"github.com/gfes980615/Diana/glob"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/models"
	"github.com/gfes980615/Diana/models/po"
	"github.com/gfes980615/Diana/repository/mysql"
	"regexp"
	"sort"
	"strconv"
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
}

func (cs *currencyService) GetMapleCurrencyMessage(mapleServer string) string {
	currencySlice, count := cs.get8591CurrencyValueTop5(mapleServer)

	message := fmt.Sprintf("85幣值前五(共 %d 筆):\n", count)
	for i := 0; i < len(currencySlice); i++ {
		message += fmt.Sprintf("1 : %.f萬\n", currencySlice[i])
	}
	return message
}


// get8591CurrencyValueTop5 ...
func (cs *currencyService) get8591CurrencyValueTop5(mapleServer string) ([]float64, int) {
	// 取得所有數量
	resultPage := GetPageSource(fmt.Sprintf(URL8591, glob.MapleServerMap[mapleServer], 0), UTF8)
	count := GetAllCount(resultPage)
	titleArray := []models.URLStruct{}
	pageResult := make(chan string, count/raw+1)
	// 取得每條商品資訊
	for i := 0; i < count; i = i + raw {
		go func(index int) {
			tmpResult := GetPageSource(fmt.Sprintf(URL8591, glob.MapleServerMap[mapleServer], index), UTF8)
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
		currencySlice = append(currencySlice, &po.Currency{Value: value, Server: mapleServer, Title: title.Name, URL: title.URL})
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
	go insertCurrency(cResultSlice)

	result := []float64{}
	for _, c := range currencySlice[0:defaultSize] {
		result = append(result, c.Value)
	}

	return result, count
}

// 幣值異常時通知用戶
func insertCurrency(cresult []models.Currency) {
	lineUser := line_user.LineUserRepository{}
	c := currency.CurrencyRepository{}
	err := c.InsertAndWarning(cresult, lineUser.GetAllUser())
	if err != nil {
		log.Println("Insert error: ", err)
	}
}
