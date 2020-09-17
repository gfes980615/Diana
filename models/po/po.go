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
