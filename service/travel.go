package service

import (
	"errors"
	"fmt"
	"github.com/gfes980615/Diana/db"
	"github.com/gfes980615/Diana/glob/common/log"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/models/bo"
	"github.com/gfes980615/Diana/models/po"
	"github.com/gfes980615/Diana/repository/mysql"
	"github.com/gfes980615/Diana/utils"
	"github.com/gocolly/colly"
	"sort"
	"strconv"
	"strings"
)

func init() {
	injection.AutoRegister(&travelService{})
}

const (
	taoyuan_root = "https://travel.tycg.gov.tw"
	taoyuan_hits = "https://travel.tycg.gov.tw/zh-tw/travel?sortby=Hits"
)

type travelService struct {
	travelRepository mysql.TravelRepository `injection:"travelRepository"`
}

func (ts *travelService) TaoyuanTravelPlace() error {
	c := colly.NewCollector(colly.UserAgent(utils.RandomAgent()), colly.Async(true))
	travelList := []*po.TouristAttractionList{}
	c.OnHTML("ul.info-card-list.mode-switch > li.item", func(e *colly.HTMLElement) {
		address := e.ChildText("a > div.info-blk.w-100 > p.icon-location")
		country, location := ts.parseAddress(address)
		tmp := &po.TouristAttractionList{
			Place:        e.ChildText("a > div.info-blk.w-100 > h3"),
			URL:          taoyuan_root + e.ChildAttr("a", "href"),
			ActivityTime: e.ChildText("a > div.info-blk.w-100 > p.opening-status.open"),
			Address:      address,
			Country:      country,
			Location:     location,
		}
		travelList = append(travelList, tmp)
	})

	pages, err := ts.getTaoyuanTravelPages()
	if err != nil {
		return err
	}
	for _, page := range pages {
		c.Visit(page)
	}
	c.Wait()

	return ts.travelRepository.CreateTaoyuanTravelItem(db.MysqlConn.Session(), travelList)
}

func (ts *travelService) getTaoyuanTravelPages() ([]string, error) {
	c := colly.NewCollector(colly.UserAgent(utils.RandomAgent()))
	lastPage := ""
	c.OnHTML("div.page-bar > div.blk.next-blk", func(e *colly.HTMLElement) {
		e.ForEach("a", func(index int, el *colly.HTMLElement) {
			if index == 1 {
				lastPage = taoyuan_root + el.Attr("href")
			}
		})
	})
	c.Visit(taoyuan_hits)

	return ts.getPageArray(lastPage)
}

func (ts *travelService) getPageArray(lastPage string) ([]string, error) {
	pages := []string{}
	pageString := strings.Split(lastPage, "&page=")
	if len(pageString) != 2 {
		return nil, errors.New("get page string can't resolve")
	}
	lastPageNumber, err := strconv.Atoi(pageString[1])
	if err != nil {
		return nil, err
	}
	rootURL := pageString[0] + "&page=%d"

	for i := 1; i <= lastPageNumber; i++ {
		pages = append(pages, fmt.Sprintf(rootURL, i))
	}
	return pages, nil
}

func (ts *travelService) parseAddress(address string) (string, string) {
	str := strings.Split(address, " ")
	if len(str) != 2 {
		log.Errorf("can't parse address: %s", address)
		return "", ""
	}
	if len(str[1]) < 18 {
		log.Errorf("can't parse address: %s", address)
		return "", ""
	}
	country := str[1][0:9]
	location := str[1][9:18]
	return country, location
}

func (ts *travelService) GetTravelPlaceByArea(country, location string) ([]*po.TouristAttractionList, error) {
	result, err := ts.travelRepository.GetTravelListByArea(db.MysqlConn.Session(), country, location)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *travelService) GetPlaceLatLong() {

}

func (ts *travelService) GetClosestTravelPlaceListTop5(latlng *bo.LatLong) ([]*po.TouristAttractionList, error) {
	travelList, err := ts.travelRepository.GetAllTravelList(db.MysqlConn.Session())
	if err != nil {
		return nil, err
	}
	distanceList := []float64{}
	distanceMap := make(map[float64][]*po.TouristAttractionList)
	for _, list := range travelList {
		distance := utils.EarthDistance(latlng.Lat, latlng.Lng, list.Latitude, list.Longitude)
		distanceList = append(distanceList, distance)
		if _, ok := distanceMap[distance]; !ok {
			distanceMap[distance] = []*po.TouristAttractionList{}
		}
		distanceMap[distance] = append(distanceMap[distance], list)
	}

	sort.Slice(distanceList, func(i, j int) bool {
		return distanceList[i] > distanceList[j]
	})

	result := []*po.TouristAttractionList{}
	for _, list := range distanceList {
		result = append(result, distanceMap[list]...)
		if len(result) > 5 {
			break
		}
	}
	return result, nil
}
