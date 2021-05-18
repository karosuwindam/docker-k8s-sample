package novel_chack

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type BookBark struct {
	Title string `json:title`
	Url   string `json:url`
}

func ReadBookBark(path string) []BookBark {
	output := []BookBark{}
	fileInfos, _ := ioutil.ReadFile(path)
	stringReader := strings.NewReader(string(fileInfos))
	doc, err := goquery.NewDocumentFromReader(stringReader)
	// doc, err := goquery.NewDocument(path)
	if err != nil {
		fmt.Println(err.Error())
		return output
	}
	doc.Find("dt").Each(func(i int, s *goquery.Selection) {
		var tmp BookBark
		tmp.Title = s.Find("a").Text()
		tmp.Url, _ = s.Find("a").Attr("href")
		if tmp.Url != "" {
			if strings.Index(tmp.Url, "http") == 0 {
				output = append(output, tmp)
			}
		}
	})
	return output
}
