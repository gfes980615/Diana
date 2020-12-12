package utils

import (
	"math"
	"runtime"
	"unicode"
)

func RemoveExtraChar(title string) string {
	var s []int32
	for _, t := range title {
		if unicode.Is(unicode.Han, t) || unicode.IsDigit(t) || unicode.IsLetter(t) {
			s = append(s, t)
		}
	}
	return string(s)
}

func EarthDistance(lat1, lng1, lat2, lng2 float64) float64 {
	radius := 6371000.0 //6378137.0
	rad := math.Pi / 180.0
	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * radius / 1000
}

func TraceMemStats() uint64 {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	//log.Printf("Alloc:%d(bytes) HeapIdle:%d(bytes) HeapReleased:%d(bytes)", ms.Alloc, ms.HeapIdle, ms.HeapReleased)
	return ms.Alloc
}

func EmptyString(str string) bool {
	if len(str) == 0 {
		return true
	}
	return false
}
