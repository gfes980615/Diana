package line

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	"github.com/gfes980615/Diana/glob"
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

func GetEveryDaySentence() string {
	juziSubURL := setURL(baseURL, baseRegexp, GBK)
	r := getRandomNumber(len(juziSubURL))
	subListURL := setURL(baseURL+juziSubURL[r].URL, subRegexp, GBK)
	lr := getRandomNumber(len(subListURL))
	result := getPageSource(baseURL+subListURL[lr].URL, GBK)
	rp := regexp.MustCompile(`<p>([0-9]+)、(.*?)</p>`)
	items := rp.FindAllStringSubmatch(result, -1)
	ir := getRandomNumber(len(items))
	return fmt.Sprintf("%s > %s:\n\n%s", juziSubURL[r].Name, subListURL[lr].Name, items[ir][2])
}

func GetMapleStoryAnnouncement() string {
	url := "https://tw.beanfun.com/maplestory/BullentinList.aspx?cate=71"
	tmpRegex := `<TD class="maple01"><a href=(.*?)>(.*?)</a></TD>`
	items := setURL(url, tmpRegex, UTF8)
	announcement := ""
	for _, item := range items {
		announcement += fmt.Sprintf("%s\n%s%s\n\n", item.Name, mapleBaseURL, item.URL)
	}

	return announcement
}

func getRandomNumber(number int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Int() % number
}

func setURL(url string, regex string, code string) []URLStruct {
	subURL := []URLStruct{}
	result := getPageSource(url, code)
	rp := regexp.MustCompile(regex)
	items := rp.FindAllStringSubmatch(result, -1)
	for _, item := range items {
		tmp := URLStruct{Name: item[2], URL: item[1]}
		subURL = append(subURL, tmp)
	}

	return subURL
}

func getPageSource(url string, code string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
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
	result := ConvertToString(string(body), code, "utf-8")

	return strings.Replace(result, "\n", "", -1)
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}
