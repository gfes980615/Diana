package dto

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

// ReturnSlice ...
type ReturnSlice struct {
	Date []string
	Izcr []float64
	Izr  []float64
	Ld   []float64
	Plt  []float64
	Slc  []float64
	Yen  []float64
	YMax int
	YMin int
}

type Message struct {
	Message string `json:"message"`
}

type Activity struct {
	Title string
	URL   string
	Time  string
}

func (a Activity) GetStruct() interface{} {
	return Activity{}
}

func (a Activity) GetStructPtr() interface{} {
	return &Activity{}
}

type LineMessage struct {
	Events     string `json:"events" form:"events"`
	LineEvents []*linebot.Event
}

func (lm *LineMessage) Init(ctx *gin.Context) error {
	err := ctx.ShouldBind(lm)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(lm.Events), &lm.LineEvents)
	if err != nil {
		return err
	}
	return nil
}

type TravelReq struct {
	Country  string `form:"country"`
	Location string `form:"location"`
}

func (tr *TravelReq) Init(ctx *gin.Context) error {
	err := ctx.ShouldBind(tr)
	if err != nil {
		return err
	}
	return nil
}

type TouristAttractionList struct {
	Place    string
	URL      string
	Address  string
	Distance string
}

func (a TouristAttractionList) GetStruct() interface{} {
	return TouristAttractionList{}
}

func (a TouristAttractionList) GetStructPtr() interface{} {
	return &TouristAttractionList{}
}
