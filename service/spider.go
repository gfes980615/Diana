package service

import (
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/gfes980615/Diana/glob"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/models"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	baseRegexp = `<li><a href="/([a-z]+)/">(.{4,6})</a></li>`
	subRegexp  = `<li><h3><a href="(.{0,20})" title="(.{0,10})" target="_blank">`
	GBK        = "gbk"
	UTF8       = "utf-8"
)

var (
	bot          = glob.Bot
	baseURL      = "https://www.1juzi.com/"
	mapleBaseURL = "https://tw.beanfun.com/maplestory/"
)

func init() {
	injection.AutoRegister(&spiderService{})
}

type spiderService struct {
}

func (s *spiderService) GetPageSource(url string, code string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Http get err:", err)
	}
	if resp.StatusCode != 200 {
		fmt.Println("Http status code:", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error", err)
	}
	result := s.convertToString(string(body), code, "utf-8")
	s.writeToFile(body)

	return strings.Replace(result, "\n", "", -1)
}

func (s *spiderService) GetAllCount(page string) int {
	rp := regexp.MustCompile(rCount)
	items := rp.FindAllStringSubmatch(page, -1)
	count, _ := strconv.Atoi(items[0][1])
	return count
}

func (c *spiderService) convertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func (s *spiderService) writeToFile(b []byte) {
	// write the whole body at once
	err := ioutil.WriteFile("regexp.html", b, 0644)
	if err != nil {
		panic(err)
	}
}

func (s *spiderService) GetEveryDaySentence() string {
	juziSubURL := s.setURL(baseURL, baseRegexp, GBK)
	r := s.getRandomNumber(len(juziSubURL))
	subListURL := s.setURL(baseURL+juziSubURL[r].URL, subRegexp, GBK)
	lr := s.getRandomNumber(len(subListURL))
	result := s.GetPageSource(baseURL+subListURL[lr].URL, GBK)
	rp := regexp.MustCompile(`<p>([0-9]+)、(.*?)</p>`)
	items := rp.FindAllStringSubmatch(result, -1)
	ir := s.getRandomNumber(len(items))
	return fmt.Sprintf("%s > %s:\n\n%s", juziSubURL[r].Name, subListURL[lr].Name, items[ir][2])
}

func (s *spiderService) getRandomNumber(number int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Int() % number
}

func (s *spiderService) setURL(url string, regex string, code string) []models.URLStruct {
	subURL := []models.URLStruct{}
	result := s.GetPageSource(url, code)
	rp := regexp.MustCompile(regex)
	items := rp.FindAllStringSubmatch(result, -1)
	for _, item := range items {
		tmp := models.URLStruct{Name: item[2], URL: item[1]}
		subURL = append(subURL, tmp)
	}

	return subURL
}
