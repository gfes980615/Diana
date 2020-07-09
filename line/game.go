package line

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/gfes980615/Diana/glob"
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

func Get8591CurrencyValue(mapleServer string) string {
	// 取得所有數量
	result := getPageSource(fmt.Sprintf(URL8591, glob.MapleServerMap[mapleServer], 0), UTF8)
	count := getAllCount(result)
	tmpArray := []string{}
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
				tmpArray = append(tmpArray, item[2])
			}
		}
	}
	// // 取得每條商品資訊
	// for i := 0; i < count; i = i + raw {
	// 	result := getPageSource(fmt.Sprintf(URL8591, i), UTF8)
	// 	rp := regexp.MustCompile(regexpA)
	// 	items := rp.FindAllStringSubmatch(result, -1)
	// 	for _, item := range items {
	// 		tmpArray = append(tmpArray, item[2])
	// 	}
	// }
	// 將幣值取出，string -> int 存成 array
	currencyValueArray := []float64{}
	for _, t := range tmpArray {
		rp := regexp.MustCompile(rCurrencyValue)
		items := rp.FindAllStringSubmatch(t, -1)
		for _, item := range items {
			a, _ := strconv.Atoi(item[1])
			b, _ := strconv.Atoi(item[2])
			currencyValueArray = append(currencyValueArray, float64(b)/float64(a))
		}
	}
	sort.Slice(currencyValueArray, func(i, j int) bool {
		return currencyValueArray[i] > currencyValueArray[j]
	})

	message := fmt.Sprintf("85幣值前五(共 %d 筆):\n", count)
	for i := 0; i < 5; i++ {
		message += fmt.Sprintf("1 : %.f萬\n", currencyValueArray[i])
	}
	return message
}
