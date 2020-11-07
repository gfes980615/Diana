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
	ID                int    `gorm:"column:id;primary_key;not null;autoIncrement"`
	Title             string `gorm:"column:title;type:text"`
	URL               string `gorm:"column:url;type:text"`
	Introduction      string `gorm:"column:introduction;type:text"`
	Category          string `gorm:"column:category;type:text"`
	CreateTime        string `gorm:"column:create_time;type:text"`
	TicketStatus      string `gorm:"column:ticket_status;type:text"`
	ParticipateNumber string `gorm:"column:participate_number;type:text"`
	ActivityTime      string `gorm:"column:activity_time;type:text"`
}

func (k KktixActivity) TableName() string {
	return "kktix_activity"
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

type TouristAttractionList struct {
	Place        string  `gorm:"place"`
	URL          string  `gorm:"url"`
	ActivityTime string  `gorm:"activity_time"`
	Country      string  `gorm:"country"`
	Location     string  `gorm:"location"`
	Address      string  `gorm:"address"`
	Latitude     float64 `gorm:"latitude"`
	Longitude    float64 `gorm:"longitude"`
}

func (t TouristAttractionList) TableName() string {
	return "tourist_attraction_list"
}
