package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("Hello, It Home!"))
	})

	router.POST("/callback", callbackHandler)

	router.Run()

	// defer func() {
	// 	if rc := recover(); rc != nil {
	// 		log.Printf("panic:\n%v\n", rc)
	// 	}
	// }()
	// port := os.Getenv("PORT")
	// addr := fmt.Sprintf(":%s", port)
	// go http.ListenAndServe(addr, nil)

	// var err error
	// bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	// log.Println("Bot:", bot, " err:", err)
	// http.HandleFunc("/callback", callbackHandler)
	// http.HandleFunc("/test", test)
}

func callbackHandler(c *gin.Context) {
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		log.Print(err.Error())
		if err == linebot.ErrInvalidSignature {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"reason": "values error.",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"reason": "values error.",
			})
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:

				if message.Text == "a" {
					daily := getEveryDaySentence()
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(daily)).Do()
					return
				}

				id, transferErr := strconv.ParseInt(message.Text, 10, 64)
				text := getGoogleExcelValueById(id)
				if transferErr != nil {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(transferErr.Error())).Do(); err != nil {
						log.Print(err)
					}
					return
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
	// fmt.Println(getEveryDaySentence())
}

func getGoogleExcelValueById(id int64) string {
	url := "https://script.google.com/macros/s/AKfycbzDtZfQHmr0YJF7F_m2ZfatU7Hu-FwTpBTwQfYXqZAv7P1JnHQ/exec?msg=" + fmt.Sprintf("%d", id)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("err:\n" + err.Error())
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read error", err)
		return ""
	}

	type Tmp struct {
		Msg interface{}
	}

	test := Tmp{}
	if err := json.Unmarshal(body, &test); err != nil {
		log.Print(err.Error())
		return ""
	}

	switch reflect.TypeOf(test.Msg).Kind() {
	case reflect.Int:
		return fmt.Sprintf("%d", test.Msg.(int))
	case reflect.Int8:
		return fmt.Sprintf("%d", test.Msg.(int8))
	case reflect.Int16:
		return fmt.Sprintf("%d", test.Msg.(int16))
	case reflect.Int32:
		return fmt.Sprintf("%d", test.Msg.(int32))
	case reflect.Int64:
		return fmt.Sprintf("%d", test.Msg.(int64))
	case reflect.String:
		return test.Msg.(string)
	case reflect.Float64:
		return fmt.Sprintf("%.f", test.Msg.(float64))
	case reflect.Float32:
		return fmt.Sprintf("%.f", test.Msg.(float32))
	default:
		fmt.Println(reflect.TypeOf(test.Msg).Kind())
		return "unknow type"
	}

	return "unexcept error"
}

var (
	baseURL = "https://www.1juzi.com/"
)

const (
	baseRegexp = `<li><a href="/([a-z]+)/">(.{4,6})</a></li>`
	subRegexp  = `<li><h3><a href="(.{0,20})" title="(.{0,10})" target="_blank">`
)

type URLStruct struct {
	CategoryName string
	URL          string
}

func getEveryDaySentence() string {
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
