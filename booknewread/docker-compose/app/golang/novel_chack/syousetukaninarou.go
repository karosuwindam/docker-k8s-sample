package novel_chack

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type List struct {
	Title      string    `json:tilte`
	Count      int       `json:count`
	Lastdate   time.Time `json:lastdate`
	Url        string    `json:url`
	LastStoryT string    `json:laststoryt`
	LastUrl    string    `json:lasturl`
}

const (
	BASE_URL_NAROU  = "http://ncode.syosetu.com"
	BASE_URL_NAROUS = "https://ncode.syosetu.com"
)

func chackurl(url string) bool {
	flag := false
	if len(BASE_URL_NAROU) <= len(url) {
		if url[:len(BASE_URL_NAROU)] == BASE_URL_NAROU {
			flag = true
		}
	} else if len(BASE_URL_NAROUS) <= len(url) {
		if url[:len(BASE_URL_NAROUS)] == BASE_URL_NAROUS {
			flag = true
		}
	}

	return flag
}

func GetSyousetu(url string) List {
	var output List
	if chackurl(url) {
		output.Url = url
		doc, err := goquery.NewDocument(url)
		if err != nil {
			fmt.Println(err.Error())
			return output
		}
		output.Title = doc.Find("p.novel_title").Text()
		doc.Find("dl.novel_sublist2").Each(func(i int, s *goquery.Selection) {
			output.LastStoryT = strings.TrimSpace(s.Find("dd.subtitle").Text())
			times := strings.TrimSpace(s.Find("dt.long_update").Text())
			if times != "" {
				if strings.Index(times, "（改）") > 0 {
					times = times[:strings.Index(times, "（改）")]
				}
				t, _ := time.Parse("2006/01/02 15:04:05 MST", times+":00 JST")
				output.Lastdate = t

			}
			tmp, _ := s.Find("dd.subtitle").Find("a").Attr("href")
			if tmp != "" {
				output.LastUrl = BASE_URL_NAROU + tmp
			}
			output.Count = i + 1
		})
	}
	return output
}
