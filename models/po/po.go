package po

import "time"

// Currency ...
type Currency struct {
	AddedTime time.Time `gorm:"column:added_time"`
	Value     float64   `gorm:"column:value"`
	Server    string    `gorm:"column:server"`
	Title     string    `gorm:"column:title"`
	URL       string    `gorm:"column:url"`
	Abnormal  int       `gorm:"column:abnormal"`
}

// LineUser ...
type LineUser struct {
	UserID string `gorm:"column:user_id"`
	//AddedTime time.Time `gorm:"column:added_time"`
}

// Website ...
type Website struct {
	Url       string    `gorm:"column:url"`
	Tag       string    `gorm:"column:tag"`
	AddedTime time.Time `gorm:"column:added_time"`
}

type Maple8591Product struct {
	Title     string `gorm:"column:title"`
	Server    string `gorm:"server"`
	Amount    string `gorm:"column:amount"`
	Number    string `gorm:"column:number"`
	Pageviews string `gorm:"column:pageviews"`
	URL       string `gorm:"column:url"`
}

func (mp Maple8591Product) TableName() string {
	return "maple_product"
}

type KktixActivity struct {
	Title             string
	URL               string
	Introduction      string
	Category          string
	CreateTime        string
	TicketStatus      string
	ParticipateNumber string
	ActivityTime      string
}

type TTActivity struct {
	Title        string
	URL          string
	ActivityTime string
	Viewers      string // 觀看人數
}

type TravelList struct {
	ID       int    `gorm:"primary_key;column:id;type:int(11);not null"`
	Title    string `gorm:"column:title;type:varchar(45)"`
	Category string `gorm:"column:category;type:varchar(45)"`
	URL      string `gorm:"unique;column:url;type:varchar(45)"`
	Status   int    `gorm:"column:status;type:int(11)"` // 0:沒去過,1:去過了,2:沒興趣
}

func (t TravelList) TableName() string {
	return "travel_list"
}
