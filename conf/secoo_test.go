package conf_test

import (
	"fmt"
	"regexp"
	"testing"
)

func TestFilterCrawlers(t *testing.T) {
	// @see https://stackoverflow.com/questions/20084513/detect-search-crawlers-via-javascript
	reg := regexp.MustCompile("(?i)google|baidu|bot|crawler|spider|crawling")
	targets := []string{
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Googlebot/2.1 (+http://www.google.com/bot.html)",
		"WhatgoOgleever",
		"WhatbaIduever",
		"Baiduspider+(+http://www.baidu.com/search/spider.htm)",
	}
	for _, target := range targets {
		ok := reg.MatchString(target)
		fmt.Println(ok, "<--", target)
	}
}
