package service

import "github.com/gfes980615/Diana/glob"

const (
	baseRegexp = `<li><a href="/([a-z]+)/">(.{4,6})</a></li>`
	subRegexp  = `<li><h3><a href="(.{0,20})" title="(.{0,10})" target="_blank">`
	GBK        = "gbk"
	UTF8       = "utf-8"
)

var (
	bot          = glob.Bot
	baseURL      = "https://www.1juzi.com/"
	mapleBaseURL = "https://tw.beanfun.com/maplestory/"
)
