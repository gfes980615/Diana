package line

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/gfes980615/Diana/db"
	"github.com/gfes980615/Diana/glob"
	"github.com/gfes980615/Diana/model"
)

const (
	raw            = 21
	root8591       = "https://www.8591.com.tw"
	URL8591        = "https://www.8591.com.tw/mallList-list.html?&group=1&searchType=0&priceSort=0&ratios=0&searchGame=859&searchServer=%d&firstRow=%d"
	regexpA        = `<a class="detail_link" href="(.*?)" target="_blank" class="goods_link" title="(.*?)" data-key="0">`
	rCount         = `<span class="R">([0-9]+)</span>`
	rCurrencyValue = `【([0-9]+):(.*?)萬】`
)

func getAllCount(page string) int {
	rp := regexp.MustCompile(rCount)
	items := rp.FindAllStringSubmatch(page, -1)
	count, _ := strconv.Atoi(items[0][1])
	return count
}

// GetMapleCurrencyMessage ...
func GetMapleCurrencyMessage(mapleServer string) string {
	value, err := getLastAvgValue()
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	fmt.Println(value)
	// currencySlice, count := get8591CurrencyValueTop5(mapleServer)

	// message := fmt.Sprintf("85幣值前五(共 %d 筆):\n", count)
	// for i := 0; i < len(currencySlice); i++ {
	// 	message += fmt.Sprintf("1 : %.f萬\n", currencySlice[i])
	// }
	// return message
	return ""
}

// get8591CurrencyValueTop5 ...
func get8591CurrencyValueTop5(mapleServer string) ([]float64, int) {
	// 取得所有數量
	resultPage := getPageSource(fmt.Sprintf(URL8591, glob.MapleServerMap[mapleServer], 0), UTF8)
	count := getAllCount(resultPage)
	titleArray := []URLStruct{}
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
				title := URLStruct{
					URL:  root8591 + item[1],
					Name: item[2],
				}
				titleArray = append(titleArray, title)
			}
		}
	}

	// 將幣值取出，string -> float 存成 array
	currencySlice := []model.Currency{}
	for _, title := range titleArray {
		rp := regexp.MustCompile(rCurrencyValue)
		items := rp.FindAllStringSubmatch(title.Name, -1)
		var value float64
		for _, item := range items {
			a, _ := strconv.ParseFloat(item[1], 64)
			b, _ := strconv.ParseFloat(item[2], 64)
			value = b / a
		}
		currencySlice = append(currencySlice, model.Currency{Value: value, Server: mapleServer, Title: title.Name, URL: title.URL})
	}

	sort.Slice(currencySlice, func(i, j int) bool {
		return currencySlice[i].Value > currencySlice[j].Value
	})

	avgValue, err := getLastAvgValue()
	if err != nil {
		return nil, 0
	}

	reuslt := []float64{}
	if len(currencySlice) < 5 {
		err := addCurrency(currencySlice, avgValue)
		if err != nil {
			log.Println("addCurrency error: ", err)
		}
		for _, c := range currencySlice {
			reuslt = append(reuslt, c.Value)
		}
		return reuslt, count
	}

	// 存入MYSQL
	err := addCurrency(currencySlice[0:5], avgValue)
	if err != nil {
		log.Println("addCurrency error: ", err)
	}

	for _, c := range currencySlice[0:5] {
		reuslt = append(reuslt, c.Value)
	}

	return reuslt, count
}

// addCurrency 存入MYSQL
func addCurrency(currencySlice []model.Currency, avgValue float64) error {
	// TODO: 對DB的操作移到另外的package
	mysql, err := db.NewMySQL(glob.DataBase)
	if err != nil {
		return err
	}
	defer mysql.Close()

	for _, c := range currencySlice {
		abnormal := 0
		if c.Value >= (avgValue * 2) {
			abnormal = 1
		}
		err := mysql.DB.Exec("INSERT IGNORE INTO `currency` (`added_time`,`value`,`server`,`title`,`url`,`abnormal`) VALUES (NOW(),?,?,?,?,?)", c.Value, c.Server, c.Title, c.URL, abnormal)
		if err.Error != nil {
			return err.Error
		}
	}

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

func getLastAvgValue() (float64, error) {
	mysql, err := db.NewMySQL(glob.DataBase)
	if err != nil {
		return 0, err
	}
	defer mysql.Close()
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	sql := fmt.Sprintf("SELECT avg(value) as `value` FROM `currency` where abnormal = 0 and added_time = '%s'", yesterday)
	type tmpValue struct {
		Value float64 `gorm:"column:value"`
	}
	value := []tmpValue{}
	result := mysql.DB.Raw(sql).Scan(&value)
	if result.Error != nil {
		return 0, result.Error
	}

	if len(value) == 0 {
		return 0, errors.New("no avg value")
	}

	return value[0].Value, nil
}

// GetMapleCurrencyChartData ...
func GetMapleCurrencyChartData(subFunc string) (model.ReturnSlice, error) {
	r := model.ReturnSlice{}
	mysql, err := db.NewMySQL(glob.DataBase)
	if err != nil {
		return r, err
	}
	defer mysql.Close()

	sql := fmt.Sprintf("SELECT `added_time`, `server`, %s(value) as `value` FROM `currency` GROUP BY `added_time`, `server` ORDER BY `added_time` ASC", subFunc)
	currency := []*model.Currency{}
	result := mysql.DB.Raw(sql).Scan(&currency)
	if result.Error != nil {
		return r, result.Error
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
