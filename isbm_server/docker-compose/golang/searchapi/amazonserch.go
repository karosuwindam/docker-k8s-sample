package searchapi

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type AmazonNameType struct {
	Isbn   string `json:isbn`
	Title  string `json:title`
	Writer string `json:writer`
	Image  string `json:image`
	Url    string `json:url`
}

func GetPageAmazonURL(isbn string) AmazonNameType {
	var output AmazonNameType
	var outurl string
	url := "https://www.amazon.co.jp/s?k=" + isbn + "&i=stripbooks"
	url2 := "https://www.amazon.co.jp/s?k=" + isbn
	doc, err := goquery.NewDocument(url)
	output.Url = url
	if err != nil {
		fmt.Println(err.Error())
		return output
	}
	for i := 0; i < 100; i++ {
		doc.Find("div.sg-col-4-of-12.sg-col-8-of-16").Each(func(i int, s *goquery.Selection) {
			if outurl != "" {
				return
			}
			s.Find("a").Each(func(i int, ss *goquery.Selection) {
				tmp, _ := ss.Attr("href")
				if (tmp != "") && (strings.Index(tmp, "gp/help") < 1) && (strings.Index(tmp, "footer_logo") < 1) {
					outurl = "https://www.amazon.co.jp" + tmp
					return
				}
			})
		})
		if outurl == "" {

			doc.Find("div.a-section.a-spacing-medium").Each(func(i int, s *goquery.Selection) {
				if outurl != "" {
					return
				}
				s.Find("a").Each(func(i int, ss *goquery.Selection) {
					tmp, _ := ss.Attr("href")
					if (tmp != "") && (strings.Index(tmp, "gp/help") < 1) && (strings.Index(tmp, "offer-listing") < 1) {
						outurl = "https://www.amazon.co.jp" + tmp
						return
					}
				})
			})
		}
		if outurl != "" {
			doc.Find("div.sg-col-4-of-24.sg-col-4-of-12.sg-col-4-of-36.sg-col-4-of-28.sg-col-4-of-16.sg-col.sg-col-4-of-20.sg-col-4-of-32").Each(func(i int, s *goquery.Selection) {
				imageurl, _ := s.Find("img").Attr("src")
				if imageurl != "" {
					output.Image = imageurl
					return
				}

			})

			doc.Find("span.a-size-medium").Each(func(i int, s *goquery.Selection) {
				output.Title = s.Text()
			})
			doc.Find("span.a-size-base-plus").Each(func(i int, s *goquery.Selection) {
				if output.Title == "" {
					output.Title = s.Text()
				}
			})
			doc.Find("div.sg-col-4-of-12.sg-col-8-of-16.sg-col-12-of-32.sg-col-12-of-20.sg-col-12-of-36.sg-col.sg-col-12-of-24.sg-col-12-of-28").Each(func(i int, s *goquery.Selection) {
				tmp := s.Find(".a-size-base.a-link-normal").Text()
				output.Writer = strings.TrimSpace(tmp)
			})
			break
		}
		time.Sleep(time.Millisecond * 100) //100ms
		if i%2 == 0 {
			doc, _ = goquery.NewDocument(url)
			output.Url = url
		} else {
			doc, _ = goquery.NewDocument(url2)
			output.Url = url2
		}
	}

	return output

}
