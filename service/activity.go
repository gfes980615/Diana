package service

import (
	"fmt"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/models/po"
	"github.com/gocolly/colly"
	"math/rand"
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
}

func (as *activityService) GetKktixActivity() {
	rn := random(1, 4)
	agent := userAgent[rn]

	c := colly.NewCollector(
		colly.UserAgent(agent),
	)
	d := colly.NewCollector(
		colly.UserAgent(agent),
	)

	activityMap := make(map[string]*po.Activity)
	c.OnHTML("ul.events > li.type-selling", func(e *colly.HTMLElement) {
		url := e.ChildAttr("a.cover", "href")

		if _, ok := activityMap[url]; ok {
			return
		}
		activityMap[url] = &po.Activity{
			URL: url,
		}
		e.ForEach("a.cover", func(i int, el *colly.HTMLElement) {
			activityMap[url].Title = el.ChildText("div.event-title")
			activityMap[url].Introduction = el.ChildText("div.introduction")
			activityMap[url].Category = el.ChildText("span.category")
			activityMap[url].CreateTime = el.ChildText("div.ft > span.date")
			activityMap[url].ParticipateNumber = el.ChildText("ul.groups")
			activityMap[url].TicketStatus = el.ChildText("div.ft > span.fake-btn")
		})
		d.Visit(e.ChildAttr("a.cover", "href"))
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
			fmt.Println("no such key")
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
			fmt.Println("no such key")
		}
	})

	for url, _ := range as.getKktixAllPage(agent) {
		c.Visit("https://kktix.com" + url)
	}

	number := 1
	for _, activity := range activityMap {
		fmt.Println(number)
		fmt.Println("活動名稱:", activity.Title)
		fmt.Println("報名網址:", activity.URL)
		fmt.Println("活動簡介:", activity.Introduction)
		fmt.Println("活動分類:", activity.Category)
		fmt.Println("參加人數:", activity.ParticipateNumber)
		fmt.Println("刊登日期:", activity.CreateTime)
		fmt.Println("票卷狀態:", activity.TicketStatus)
		fmt.Println("活動時間:", activity.ActivityTime)
		number++
	}
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func (as *activityService) getKktixAllPage(agent string) map[string]bool {
	b := colly.NewCollector(
		colly.UserAgent(agent),
	)
	visitURL := make(map[string]bool)
	b.OnHTML("div.pagination", func(e *colly.HTMLElement) {
		urlSlice := []string{}
		lastURL := ""
		e.ForEach("ul > li", func(i int, el *colly.HTMLElement) {
			url := el.ChildAttr("a", "href")
			if url == "#" {
				return
			}
			if _, ok := visitURL[url]; !ok {
				visitURL[url] = true
				fmt.Println(url)
			}
			urlSlice = append(urlSlice, url)
			if i > 0 {
				lastURL = urlSlice[i-1]
			}
		})
		b.Visit("https://kktix.com" + lastURL)
	})

	b.Visit("https://kktix.com/events?category_id=2")

	return visitURL
}
