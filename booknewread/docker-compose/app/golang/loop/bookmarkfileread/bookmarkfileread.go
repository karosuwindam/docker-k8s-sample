package bookmarkfileread

import (
	"book-newread/config"
	"io/ioutil"
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var filepass string

type BookBark struct {
	Title string `json:title`
	Url   string `json:url`
}

func Init() error {
	filepass = config.Loop.BookmarkF
	if _, err := os.Stat(filepass); os.IsNotExist(err) {
		return err
	}
	return nil
}

func ReadBookmark() map[string]string {
	// out := []string{}
	tmpBook := make(map[string]string)
	if files := readfilepass(); len(files) != 0 {
		for _, f := range files {
			if tmps := readBookBark(f); len(tmps) != 0 {
				for _, tmp := range tmps {
					tmpBook[tmp.Url] = tmp.Title
				}
			}
		}
	}
	// for url, s := range tmpBook {
	// 	out = append(out, url)
	// }
	return tmpBook
}

func readfilepass() []string {
	output := []string{}
	if files, err := ioutil.ReadDir(filepass); err == nil {
		for _, f := range files {
			tmp := ""
			if filepass[len(filepass)-1:] == "/" {
				tmp = filepass + f.Name()
			} else {
				tmp = filepass + "/" + f.Name()
			}
			output = append(output, tmp)
		}
	}

	return output
}

func readBookBark(path string) []BookBark {
	output := []BookBark{}
	fileInfos, _ := ioutil.ReadFile(path)
	stringReader := strings.NewReader(string(fileInfos))
	doc, err := goquery.NewDocumentFromReader(stringReader)
	// doc, err := goquery.NewDocument(path)
	if err != nil {
		slog.Error("readBookBark", "error", err)
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
	return output
}
