package service

import (
	"fmt"
	"github.com/gfes980615/Diana/db"
	"github.com/gfes980615/Diana/glob/common/log"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/models/po"
	"github.com/gfes980615/Diana/repository/mysql"
	"github.com/gocolly/colly"
	"math/rand"
	"strings"
	"time"
)

func init() {
	injection.AutoRegister(&activityService{})
}

const (
	userAgent1 = "AppleWebKit/537.36 (KHTML, like Gecko)"
	userAgent2 = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
	userAgent3 = "Chrome/85.0.4183.102 Safari/537.36"
	userAgent4 = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36"

	kktix_root       = "https://kktix.com"
	kktix_exhibition = "https://kktix.com/events?category_id=11" // 展覽
	kktix_all        = "https://kktix.com/events"

	travelTaipei_root   = "https://www.travel.taipei"
	travelTaipei_show   = "https://www.travel.taipei/zh-tw/activity?sortby=Recently&page=1" // travel taipei 活動展演
	travelTaipei_themes = "https://www.travel.taipei/zh-tw/attraction/themes"
)

var (
	userAgent = map[int]string{
		1: userAgent1,
		2: userAgent2,
		3: userAgent3,
		4: userAgent4,
	}
)

type activityService struct {
	travelRepository mysql.TravelRepository `injection:""`
}

func (as *activityService) GetTravelTaipeiActivity(category string) []*po.TTActivity {
	switch category {
	case "exhibition":
		return as.getTravelTaipeiExhibitionItems()
	case "travel_list":
		as.getTravelTaipeiAllTravelList()
	}
	return []*po.TTActivity{}
}

func (as *activityService) GetKktixActivity(category string) []*po.KktixActivity {
	var url string
	switch category {
	case "exhibition":
		url = kktix_exhibition
	case "all":
		url = kktix_all
	}

	activityMap, forSortSlice := as.getKKtixActivityItems(url)
	result := []*po.KktixActivity{}
	number := 1
	for _, activity := range forSortSlice {
		if len(activity) == 0 {
			continue
		}
		//fmt.Println(number)
		//fmt.Println("活動名稱:", activityMap[activity].Title)
		//fmt.Println("報名網址:", activityMap[activity].URL)
		//fmt.Println("活動簡介:", activityMap[activity].Introduction)
		//fmt.Println("活動分類:", activityMap[activity].Category)
		//fmt.Println("參加人數:", activityMap[activity].ParticipateNumber)
		//fmt.Println("刊登日期:", activityMap[activity].CreateTime)
		//fmt.Println("票卷狀態:", activityMap[activity].TicketStatus)
		//fmt.Println("活動時間:", activityMap[activity].ActivityTime)
		number++

		result = append(result, activityMap[activity])
	}
	//fmt.Println("end ...")
	return result
}

func (as *activityService) getTravelTaipeiAllTravelList() []*po.TravelList {
	themesURL := as.getTravelTaipeiThemesList()

	c := colly.NewCollector(colly.UserAgent(as.randomAgent()))

	travelList := []*po.TravelList{}
	tmpCategory := ""
	c.OnHTML("ul.event-news-card-list > li.item", func(e *colly.HTMLElement) {
		tmp := &po.TravelList{
			Title:    e.ChildText("div.info-card-item > a.link > div.info-blk > h3.info-title"),
			Category: themesURL[tmpCategory],
			URL:      travelTaipei_root + e.ChildAttr("div.info-card-item > a.link", "href"),
		}
		travelList = append(travelList, tmp)
	})

	//c.OnHTML("div.page-bar > div.blk", func(e *colly.HTMLElement) {
	//	c.Visit(travelTaipei_root + e.ChildAttr("a.next-page", "href"))
	//})

	DB := db.MysqlConn.Session()
	for url, _ := range themesURL {
		tmpCategory = url
		c.Visit(url)
		if err := as.travelRepository.CreateTravelTaipeiTravelItem(DB, travelList); err != nil {
			log.Error(err)
		}
		travelList = []*po.TravelList{}
	}

	return travelList
}

func (as *activityService) getTravelTaipeiThemesList() map[string]string {
	c := colly.NewCollector(colly.UserAgent(as.randomAgent()))

	themesURLs := make(map[string]string)
	c.OnHTML("ul.d-flex > li.col-6", func(e *colly.HTMLElement) {
		title := e.ChildText("div.d-block a.text-decoration-none > h3.fz-16px")
		url := travelTaipei_root + e.ChildAttr("div.d-block a.text-decoration-none", "href")
		themesURLs[url] = title
	})

	c.Visit(travelTaipei_themes)

	return themesURLs
}

func (as *activityService) getTravelTaipeiExhibitionItems() []*po.TTActivity {
	c := colly.NewCollector(colly.UserAgent(as.randomAgent()))
	resultItems := []*po.TTActivity{}

	visitTag := true

	c.OnHTML("ul.event-news-card-list > li.item", func(e *colly.HTMLElement) {
		activityTime := e.ChildText("div.info-card-item > a.link > div.info-blk > span.duration")
		if as.activityIsEnd(activityTime) {
			visitTag = false
			return
		}
		ttItem := &po.TTActivity{
			Title:        e.ChildText("div.info-card-item > a.link > div.info-blk > h3.info-title"),
			URL:          travelTaipei_root + e.ChildAttr("div.info-card-item > a.link", "href"),
			ActivityTime: e.ChildText("div.info-card-item > a.link > div.info-blk > span.duration"),
			Viewers:      e.ChildText("div.info-card-item > a.link > div.info-blk > span.icon-view"),
		}
		resultItems = append(resultItems, ttItem)
	})

	c.OnHTML("div.page-bar > div.blk", func(e *colly.HTMLElement) {
		if visitTag {
			c.Visit(travelTaipei_root + e.ChildAttr("a.next-page", "href"))
		}
	})

	c.Visit(travelTaipei_show)

	for i, item := range resultItems {
		fmt.Printf("%d\n標題:%s\n網址:%s\n活動時間:%s\n查看人數:%s\n", i+1, item.Title, item.URL, item.ActivityTime, item.Viewers)
	}
	return resultItems
}

// 檢查活動是否結束
func (as *activityService) activityIsEnd(activityTime string) bool {
	activityTimes := strings.Split(activityTime, "～")
	var lastDate time.Time
	if len(activityTimes) > 1 {
		lastDate, _ = time.Parse("2006-01-02", activityTimes[1])
	} else {
		lastDate, _ = time.Parse("2006-01-02", activityTimes[0])
	}

	if lastDate.Before(time.Now()) {
		return true
	}

	return false
}
func (as *activityService) getKKtixActivityItems(url string) (map[string]*po.KktixActivity, []string) {
	c := colly.NewCollector(colly.UserAgent(as.randomAgent()))
	d := colly.NewCollector(colly.UserAgent(as.randomAgent()))

	activityMap := make(map[string]*po.KktixActivity)
	forSortSlice := []string{}

	c.OnHTML("ul.events", func(e *colly.HTMLElement) {
		e.ForEach("li", func(i int, ex *colly.HTMLElement) {
			url := ex.ChildAttr("a.cover", "href")

			if _, ok := activityMap[url]; ok {
				return
			}
			forSortSlice = append(forSortSlice, url)
			activityMap[url] = &po.KktixActivity{
				URL: url,
			}
			ex.ForEach("a.cover", func(j int, el *colly.HTMLElement) {
				activityMap[url].Title = el.ChildText("div.event-title")
				activityMap[url].Introduction = el.ChildText("div.introduction")
				activityMap[url].Category = el.ChildText("span.category")
				activityMap[url].CreateTime = el.ChildText("div.ft > span.date")
				activityMap[url].ParticipateNumber = el.ChildText("ul.groups")
				activityMap[url].TicketStatus = el.ChildText("div.ft > span.fake-btn")
			})
			d.Visit(ex.ChildAttr("a.cover", "href"))
		})
	})

	getActivityTimeMap := make(map[string]bool)

	d.OnHTML("ul.info", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		if _, ok := getActivityTimeMap[url]; ok {
			return
		}
		if _, ok := activityMap[url]; ok {
			t := e.ChildText("span.info-desc > span.timezoneSuffix")
			if len(t) > 0 {
				activityMap[url].ActivityTime = t
				getActivityTimeMap[url] = true
			}
		} else {
			log.Info("no such key")
		}
	})

	d.OnHTML("div.section", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		if _, ok := getActivityTimeMap[url]; ok {
			return
		}
		if _, ok := activityMap[url]; ok {
			t := e.ChildText("p > span.timezoneSuffix")
			if len(t) > 0 {
				activityMap[url].ActivityTime = t
				getActivityTimeMap[url] = true
			}
		} else {
			log.Info("no such key")
		}
	})

	for subUrl, _ := range as.getKktixPageByURL(url) {
		c.Visit(subUrl)
	}

	return activityMap, forSortSlice
}

func (as *activityService) getKktixPageByURL(url string) map[string]bool {
	pageColly := colly.NewCollector(
		colly.UserAgent(as.randomAgent()),
	)

	visitURL := make(map[string]bool)
	visitURL[url] = true
	pageColly.OnHTML("div.pagination", func(e *colly.HTMLElement) {
		urlSlice := []string{}
		lastURL := ""
		e.ForEach("ul > li", func(i int, el *colly.HTMLElement) {
			url := el.ChildAttr("a", "href")
			if url == "#" {
				return
			}
			if _, ok := visitURL[kktix_root+url]; !ok {
				visitURL[kktix_root+url] = true
				log.Info(kktix_root + url)
			}
			urlSlice = append(urlSlice, url)
			if i > 0 {
				lastURL = urlSlice[i-1]
			}
		})
		pageColly.Visit(kktix_root + lastURL)
	})

	pageColly.Visit(url)

	return visitURL
}

func (as *activityService) randomAgent() string {
	rn := random(1, 4)
	return userAgent[rn]
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
