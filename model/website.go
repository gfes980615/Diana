package model

import "time"

type Website struct {
	Url       string    `gorm:"column:url"`
	Tag       string    `gorm:"column:tag"`
	AddedTime time.Time `gorm:"column:added_time"`
}
