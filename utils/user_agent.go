package utils

import (
	"math/rand"
	"time"
)

const (
	userAgent1 = "AppleWebKit/537.36 (KHTML, like Gecko)"
	userAgent2 = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
	userAgent3 = "Chrome/85.0.4183.102 Safari/537.36"
	userAgent4 = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36"
)

var (
	userAgent = map[int]string{
		1: userAgent1,
		2: userAgent2,
		3: userAgent3,
		4: userAgent4,
	}
)

func RandomAgent() string {
	rn := random(1, 4)
	return userAgent[rn]
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
