package calendarscripting

import (
	"book-newread/config"
	"context"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
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

func FilterComicList(data []BookList) []BookList {
	output := []BookList{}
	for _, tmp := range data {
		output = append(output, tmp)
	}
	return output
}

var readmutex sync.Mutex

func GetComicList(year, month string, booktype int, ctx context.Context) []BookList {
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
	var err_ch chan error
	var doc_ch chan *goquery.Document
	err_ch = make(chan error, 1)
	doc_ch = make(chan *goquery.Document, 1)
	go func() {
		readmutex.Lock()
		defer readmutex.Unlock()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			err_ch <- err
			return
		}
		client := http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		}
		resp, err := client.Do(req)
		if err != nil {
			err_ch <- err
			return
		}
		defer resp.Body.Close()
		doc_tmp, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {
			err_ch <- err
			return
		}
		doc_ch <- doc_tmp
	}()
	var doc *goquery.Document
	select {
	case err := <-err_ch:
		log.Println("error:", err)
		return output
	case doc = <-doc_ch:
	}
	_, span := config.TracerS(ctx, "check data", url)
	defer span.End()
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
	span.SetAttributes(attribute.Int("data count", len(output)))
	return output
}
