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
	UserID    string    `gorm:"column:user_id"`
	AddedTime time.Time `gorm:"column:added_time"`
}

// Website ...
type Website struct {
	Url       string    `gorm:"column:url"`
	Tag       string    `gorm:"column:tag"`
	AddedTime time.Time `gorm:"column:added_time"`
}
