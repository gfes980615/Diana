package service

import (
	"fmt"
	"github.com/gfes980615/Diana/db"
	"github.com/gfes980615/Diana/glob/common/log"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/models/po"
	"github.com/gfes980615/Diana/repository/mysql"
	"github.com/gfes980615/Diana/utils"
	"github.com/gocolly/colly"
	"strings"
	"time"
)

func init() {
	injection.AutoRegister(&activityService{})
}

const (
	kktix_root       = "https://kktix.com"
	kktix_exhibition = "https://kktix.com/events?category_id=11" // 展覽
	kktix_all        = "https://kktix.com/events"

	travelTaipei_root   = "https://www.travel.taipei"
	travelTaipei_show   = "https://www.travel.taipei/zh-tw/activity?sortby=Recently&page=1" // travel taipei 活動展演
	travelTaipei_themes = "https://www.travel.taipei/zh-tw/attraction/themes"
)

type activityService struct {
	travelRepository   mysql.TravelRepository   `injection:"travelRepository"`
	activityRepository mysql.ActivityRepository `injection:"activityRepository"`
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
		number++

		result = append(result, activityMap[activity])
	}
	//fmt.Println("end ...")
	err := utils.CreateTable(db.MysqlConn.Session(), "kktix_activity", &po.KktixActivity{})
	if err != nil {
		log.Error(err)
	}
	err = as.activityRepository.CreateKKtixActivityItem(db.MysqlConn.Session(), result)
	if err != nil {
		log.Error(err)
	}

	return result
}

func (as *activityService) getTravelTaipeiAllTravelList() []*po.TravelList {
	themesURL := as.getTravelTaipeiThemesList()

	c := colly.NewCollector(colly.UserAgent(utils.RandomAgent()))

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
	c := colly.NewCollector(colly.UserAgent(utils.RandomAgent()))

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
	c := colly.NewCollector(colly.UserAgent(utils.RandomAgent()))
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
	c := colly.NewCollector(colly.UserAgent(utils.RandomAgent()))
	d := colly.NewCollector(colly.UserAgent(utils.RandomAgent()))

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
				activityMap[url].Title = utils.RemoveExtraChar(el.ChildText("div.event-title"))
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
		colly.UserAgent(utils.RandomAgent()),
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
