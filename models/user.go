package models

import "time"

type LineUser struct {
	UserID    string    `gorm:"column:user_id"`
	AddedTime time.Time `gorm:"column:added_time"`
}
