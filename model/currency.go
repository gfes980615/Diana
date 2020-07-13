package model

import "time"

type Currency struct {
	AddedTime time.Time `gorm:"column:added_time"`
	Value     float64   `gorm:"column:value"`
	Server    string    `gorm:"column:server"`
}

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
