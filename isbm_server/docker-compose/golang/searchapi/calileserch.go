package searchapi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type CalileNameType struct {
	Isbm     string `json:isbn`
	Title    string `json:title`
	Writer   string `json:writer`
	Ext      string `json:ext`
	Url      string `json:url`
	Synopsis string `json:synopsis`
	Brand    string `json:brand`
	Image    string `json:image`
}

func isbm13to10(isbm string) string {
	if len(isbm) == 13 {
		tmp := isbm[3:12]
		count := 0
		for i := 0; i < len(tmp); i++ {
			j, _ := strconv.Atoi(tmp[i : i+1])
			count += j * (10 - i)
			// tmp[i]
		}
		count = 11 - count%11
		if count == 11 {
			tmp += "0"
		} else if count == 10 {
			tmp += "X"
		} else {
			tmp += strconv.Itoa(count)
		}
		return tmp
	} else {
		return ""
	}
}

func GetPageCalilURL(isbm string) CalileNameType {
	isbm10 := isbm13to10(isbm)
	var output CalileNameType
	if isbm10 == "" {
		return output
	} else {
		output.Isbm = isbm
		url := "https://calil.jp/book/" + isbm10
		doc, err := goquery.NewDocument(url)
		if err != nil {
			fmt.Println(err.Error())
			return output
		}
		output.Url = url
		doc.Find("h1.title").Each(func(i int, s *goquery.Selection) {
			output.Title = strings.TrimSpace(s.Text())
		})
		doc.Find("div.author").Each(func(i int, s *goquery.Selection) {
			s.Find("a").Each(func(j int, ss *goquery.Selection) {
				if output.Writer == "" {
					output.Writer = strings.TrimSpace(ss.Text())
				} else {
					if output.Ext == "" {
						output.Ext = strings.TrimSpace(ss.Text())
					} else {
						output.Ext += "," + strings.TrimSpace(ss.Text())
					}
				}

			})
		})
		doc.Find("a.cover").Each(func(i int, s *goquery.Selection) {
			imageurl, _ := s.Find("img").Attr("src")
			if imageurl != "" {
				output.Image = imageurl
				return
			}
		})
		doc.Find("div.detail").Each(func(i int, s *goquery.Selection) {
			doc.Find("span").Each(func(i int, ss *goquery.Selection) {
				tmp, _ := ss.Attr("itemprop")
				if tmp == "publisher" {
					output.Brand = ss.Text()
				}
			})
		})
		doc.Find("div.openbd_description").Each(func(i int, s *goquery.Selection) {
			output.Synopsis = s.Find("p").Text()
		})
		return output
	}
}
