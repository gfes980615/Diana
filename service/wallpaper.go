package service

import (
	"bytes"
	"fmt"
	"github.com/gfes980615/Diana/glob/common/log"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/utils"
	"github.com/gocolly/colly"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	injection.AutoRegister(&wallPaperService{})
}

const (
	animeWallPaperPage = "https://www.zdqx.com/list-9-0-3444-0-0-0-%d.html"
)

type wallPaperService struct {
}

func (w *wallPaperService) GetWallPaper() {
	for i := 10; i <= 69; i++ {
		w.categoryPicture(fmt.Sprintf(animeWallPaperPage, i), i)
		time.Sleep(10 * time.Second)
	}
}

func (w *wallPaperService) categoryPicture(url string, page int) {
	c := colly.NewCollector(colly.UserAgent(utils.RandomAgent()))
	c.OnHTML("div.main > div.piclist > ul.clearfix > li > div.listbox", func(e *colly.HTMLElement) {
		pageUrl := e.ChildAttr("a", "href")
		imgUrl := e.ChildAttr("img", "src")
		w.downloadPicture(imgUrl, page)
		r := regexp.MustCompile(`共(\d+)张`)
		rCount := r.FindStringSubmatch(e.ChildText("em.page_num"))
		if len(rCount) != 2 {
			log.Errorf("regexp is out of expectation ,url:%s,pageUrl:%s", url, pageUrl)
			return
		}
		pictureCount := rCount[1]
		count, err := strconv.Atoi(pictureCount)
		if err != nil {
			log.Error(err)
			return
		}
		w.pagePicture(pageUrl, count, page)
	})

	c.Visit(url)
}

func (w *wallPaperService) pagePicture(url string, count, page int) {
	c := colly.NewCollector(colly.UserAgent(utils.RandomAgent()))
	c.OnHTML("div.main_center", func(e *colly.HTMLElement) {
		pictureUrl := e.ChildAttr("div.cb > div.fr > a", "href")
		go w.downloadPicture(pictureUrl, page)
	})
	for _, pageURL := range w.setPagePictureURL(url, count) {
		c.Visit(pageURL)
	}
}

func (w *wallPaperService) setPagePictureURL(url string, count int) []string {
	urlSlice := make([]string, 0, count)
	urlSlice = append(urlSlice, url)
	resetUrl := strings.Split(url, ".html")
	for i := 2; i < count; i++ {
		urlSlice = append(urlSlice, fmt.Sprintf("%s_%d.html", resetUrl[0], i))
	}
	return urlSlice
}

func (w *wallPaperService) downloadPicture(url string, pageFolder int) {
	imagPath := "http:" + url
	//图片正则
	reg := regexp.MustCompile(`(\w|\d|_)*.(jpg|png+)`)
	regPath := reg.FindStringSubmatch(imagPath)
	if len(regPath) < 3 {
		log.Errorf("regexp is out of expectation ,imagePath:%s", imagPath)
		return
	}
	name := regPath[0]
	extension := regPath[2]
	log.Infof("download picture %s", name)
	//通过http请求获取图片的流文件
	resp, err := http.Get(imagPath)
	if err != nil {
		log.Errorf("[http get error]: %v", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[read response error]: %v", err)
		return
	}
	folder := fmt.Sprintf("wallpaper/%d", pageFolder)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.Mkdir(folder, os.ModePerm)
		if err != nil {
			log.Errorf("[mkdir error]: %v", err)
			return
		}
	}
	pictureName := fmt.Sprintf("%s/%v.%s", folder, time.Now().UnixNano(), extension)
	out, err := os.Create(pictureName)
	defer out.Close()
	if err != nil {
		log.Errorf("[create file error]: %v", err)
		return
	}
	_, err = io.Copy(out, bytes.NewReader(body))
	if err != nil {
		log.Errorf("[copy error]: %v", err)
		return
	}
	log.Infof("download picture from: %s\n name: %s", imagPath, pictureName)
}
