package line

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"

	"github.com/gfes980615/Diana/db"
	"github.com/gfes980615/Diana/glob"
	"github.com/gfes980615/Diana/model"
)

const (
	raw            = 21
	URL8591        = "https://www.8591.com.tw/mallList-list.html?&group=1&searchType=0&priceSort=0&ratios=0&searchGame=859&searchServer=%d&firstRow=%d"
	regexpA        = `<a class="detail_link" href="(.*?)" target="_blank" class="goods_link" title="(.*?)" data-key="0">`
	rCount         = `<span class="R">([0-9]+)</span>`
	rCurrencyValue = `【([0-9]+):([0-9]+)萬】`
)

func getAllCount(page string) int {
	rp := regexp.MustCompile(rCount)
	items := rp.FindAllStringSubmatch(page, -1)
	count, _ := strconv.Atoi(items[0][1])
	return count
}

// GetMapleCurrencyMessage ...
func GetMapleCurrencyMessage(mapleServer string) string {
	currencySlice, count := get8591CurrencyValueTop5(mapleServer)

	message := fmt.Sprintf("85幣值前五(共 %d 筆):\n", count)
	for i := 0; i < len(currencySlice); i++ {
		message += fmt.Sprintf("1 : %.f萬\n", currencySlice[i])
	}
	return message
}

// get8591CurrencyValueTop5 ...
func get8591CurrencyValueTop5(mapleServer string) ([]float64, int) {
	// 取得所有數量
	resultPage := getPageSource(fmt.Sprintf(URL8591, glob.MapleServerMap[mapleServer], 0), UTF8)
	count := getAllCount(resultPage)
	titleArray := []string{}
	pageResult := make(chan string, count/raw+1)
	// 取得每條商品資訊
	for i := 0; i < count; i = i + raw {
		go func(index int) {
			tmpResult := getPageSource(fmt.Sprintf(URL8591, glob.MapleServerMap[mapleServer], index), UTF8)
			pageResult <- tmpResult
		}(i)
	}

	for i := 0; i < count; i = i + raw {
		select {
		case tmp := <-pageResult:
			rp := regexp.MustCompile(regexpA)
			items := rp.FindAllStringSubmatch(tmp, -1)
			for _, item := range items {
				titleArray = append(titleArray, item[2])
			}
		}
	}

	// 將幣值取出，string -> int 存成 array
	currencySlice := []model.Currency{}
	for _, title := range titleArray {
		rp := regexp.MustCompile(rCurrencyValue)
		items := rp.FindAllStringSubmatch(title, -1)
		var value float64
		for _, item := range items {
			a, _ := strconv.Atoi(item[1])
			b, _ := strconv.Atoi(item[2])
			value = float64(b) / float64(a)
		}
		currencySlice = append(currencySlice, model.Currency{Value: value, Server: mapleServer, Title: title})
	}

	sort.Slice(currencySlice, func(i, j int) bool {
		return currencySlice[i].Value > currencySlice[j].Value
	})

	reuslt := []float64{}
	if len(currencySlice) < 5 {
		err := addCurrency(currencySlice)
		if err != nil {
			log.Println("addCurrency error: ", err)
		}
		for _, c := range currencySlice {
			reuslt = append(reuslt, c.Value)
		}
		return reuslt, count
	}

	// 存入MYSQL
	err := addCurrency(currencySlice[0:5])
	if err != nil {
		log.Println("addCurrency error: ", err)
	}

	for _, c := range currencySlice {
		reuslt = append(reuslt, c.Value)
	}

	return reuslt, count
}

// addCurrency 存入MYSQL
func addCurrency(currencySlice []model.Currency) error {
	mysql, err := db.NewMySQL(glob.DataBase)
	if err != nil {
		return err
	}

	for _, c := range currencySlice {
		err := mysql.DB.Exec("INSERT IGNORE INTO `currency` (`added_time`,`value`,`server`,`title`) VALUES (NOW(),?,?,?)", c.Value, c.Server, c.Title)
		if err.Error != nil {
			return err.Error
		}
	}
	defer mysql.Close()

	return nil
}

// AddAllServerCurrency 存入所有品牌幣值
func AddAllServerCurrency() {
	for server, _ := range glob.MapleServerMap {
		go func(server string) {
			get8591CurrencyValueTop5(server)
		}(server)
	}
}

// GetMapleCurrencyChartData ...
func GetMapleCurrencyChartData() (model.ReturnSlice, error) {
	r := model.ReturnSlice{}
	mysql, err := db.NewMySQL(glob.DataBase)
	if err != nil {
		return r, err
	}
	currency := []*model.Currency{}
	result := mysql.DB.Raw("select added_time, server, avg(value) as value from currency group by added_time, server order by added_time asc").Scan(&currency)
	if result.Error != nil {
		return r, result.Error
	}

	defer mysql.Close()
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
