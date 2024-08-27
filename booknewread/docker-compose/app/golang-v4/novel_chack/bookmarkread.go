package novel_chack

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type BookBark struct {
	Title string `json:title`
	Url   string `json:url`
}

func chackdabul(data []BookBark) []BookBark {
	var output []BookBark
	if len(data) > 1 {
		back := data[0]
		output = []BookBark{data[0]}
		for _, ary := range data[1:] {
			if back.Url != ary.Url {
				output = append(output, ary)
			}
			back = ary
		}
	} else {
		output = data
	}
	return output
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
	sort.Slice(output, func(i, j int) bool { return output[i].Url > output[j].Url })
	output = chackdabul(output)
	return output
}

func BookBarkSout(ary []BookBark) []BookBark {
	tmp := ary
	sort.Slice(tmp, func(i, j int) bool { return tmp[i].Url > tmp[j].Url })
	output := chackdabul(tmp)
	return output
}
