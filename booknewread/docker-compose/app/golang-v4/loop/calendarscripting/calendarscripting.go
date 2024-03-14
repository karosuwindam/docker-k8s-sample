package calendarscripting

import (
	"fmt"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type BookList struct {
	Type   string `json:type`
	Months string `json:months`
	Days   string `json:days`
	Img    string `json:img`
	Title  string `json:title`
	Writer string `json:writer`
	Bround string `json:bround`
	Ext    string `json:ext`
}

const (
	BOOKURL      string = "https://calendar.gameiroiro.com/"
	COMICURL     string = "manga.php"
	LITENOVELURL string = "litenovel.php"
	MAGAZINEURL  string = "magazine.php"
)
const (
	COMIC     = 0 //漫画
	LITENOVEL = 1 //ライトノベル
	MAGAZINE  = 2 //雑誌
)

func FilterComicList(data []BookList, lastday int) []BookList {
	output := []BookList{}
	for _, tmp := range data {
		output = append(output, tmp)
	}
	return output
}

var readmutex sync.Mutex

func GetComicList(year, month string, booktype int) []BookList {
	var output []BookList
	mounth_tmp := month

	url := BOOKURL
	switch booktype {
	case COMIC:
		url += COMICURL
	case LITENOVEL:
		url += LITENOVELURL
	case MAGAZINE:
		url += MAGAZINEURL
	}
	if (year != "") && (month != "") {
		url += "?year=" + year + "&month=" + month
	}
	readmutex.Lock()
	doc, err := goquery.NewDocument(url)
	readmutex.Unlock()
	if err != nil {
		fmt.Println(err.Error())
		return output
	}

	if mounth_tmp == "" {
		mounth_tmp, _ = doc.Find("input.month").Attr("value")
	}
	doc.Find("div#content-inner").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(j int, ss *goquery.Selection) {
			days := strings.TrimSpace(ss.Find("td.day-td").Text())
			ss.Find("div.div-wrap").Each(func(k int, sss *goquery.Selection) {
				var tmp BookList
				tmp.Days = days
				tmp.Months = mounth_tmp
				sss.Find("div.product-image-left").Each(func(l int, image *goquery.Selection) {
					tmp.Img, _ = image.Find("img").Attr("src")
				})
				sss.Find("div.product-description-right").Each(func(k int, data *goquery.Selection) {
					tmp.Type = strings.TrimSpace(data.Find("p.p-genre").Text())
					tmp.Title = strings.TrimSpace(data.Find("a").Text())
					data.Find("p.p-company").Each(func(i int, data2 *goquery.Selection) {
						if tmp.Bround == "" {
							tmp.Bround = strings.TrimSpace(data2.Text())
						} else if tmp.Writer == "" {
							tmp.Writer = strings.TrimSpace(data2.Text())
						} else if tmp.Ext == "" {
							tmp.Ext = strings.TrimSpace(data2.Text())
						} else {
							tmp.Ext += "," + strings.TrimSpace(data2.Text())
						}

					})
				})
				output = append(output, tmp)
			})

		})
	})
	return output
}
