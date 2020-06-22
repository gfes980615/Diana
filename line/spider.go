package line

import (
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/gfes980615/Diana/glob"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	baseRegexp = `<li><a href="/([a-z]+)/">(.{4,6})</a></li>`
	subRegexp  = `<li><h3><a href="(.{0,20})" title="(.{0,10})" target="_blank">`
)

var (
	bot     = glob.Bot
	baseURL = "https://www.1juzi.com/"
)

func GetEveryDaySentence() string {
	juziSubURL := setJuziURL(baseURL, baseRegexp)
	r := getRandomNumber(len(juziSubURL))
	subListURL := setJuziURL(baseURL+juziSubURL[r].URL, subRegexp)
	lr := getRandomNumber(len(subListURL))
	result := getPageSource(baseURL + subListURL[lr].URL)
	rp := regexp.MustCompile(`<p>([0-9]+)、(.*?)</p>`)
	items := rp.FindAllStringSubmatch(result, -1)
	ir := getRandomNumber(len(items))
	return fmt.Sprintf("%s > %s:\n\n%s", juziSubURL[r].CategoryName, subListURL[lr].CategoryName, items[ir][2])
}

func getRandomNumber(number int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Int() % number
}

func setJuziURL(url string, regex string) []URLStruct {
	juziSubURL := []URLStruct{}
	result := getPageSource(url)
	rp := regexp.MustCompile(regex)
	items := rp.FindAllStringSubmatch(result, -1)
	for _, item := range items {
		tmp := URLStruct{CategoryName: item[2], URL: item[1]}
		juziSubURL = append(juziSubURL, tmp)
	}

	return juziSubURL
}

func getPageSource(url string) string {
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
	result := ConvertToString(string(body), "gbk", "utf-8")

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
