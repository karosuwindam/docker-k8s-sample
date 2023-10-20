package novelchack

import (
	"errors"
	"fmt"
	"log"
	"net/http"
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

type documentdata struct {
	url  string
	data *goquery.Document
}

const (
	BASE_URL_NAROU    = "http://ncode.syosetu.com"
	BASE_URL_NAROUS   = "https://ncode.syosetu.com"
	BASE_URL_KAKUYOMU = "https://kakuyomu.jp"
	BASE_URL_NOCKU    = "http://novel18.syosetu.com"
	BASE_URL_NOCKUS   = "https://novel18.syosetu.com"
	BASE_URL_ALPHA    = "http://www.alphapolis.co.jp"
	BASE_URL_ALPHAS   = "https://www.alphapolis.co.jp"
)

const (
	DNS_NAROU    = "ncode.syosetu.com"
	DNS_KAKUYOMU = "kakuyomu.jp"
	DNS_NOCKU    = "novel18.syosetu.com"
	DNS_ALPHA    = "alphapolis.co.jp"
)

type nobelWebType int

const (
	NAROU_WEB    nobelWebType = iota //なろう
	KAKUYOMU_WEB                     //カクヨム
	NNOCKU_WEB                       //ノクターン
	ALPHA_WEB                        //アルファポリス
	ATHOR_WEB
)

const MAX_CH int = 3
const MAX_Novel_CH int = 1

var narou_ch chan string
var kakuyomu_ch chan string
var nnocku_ch chan string
var alpha_ch chan string

func Setup() {
	narou_ch = make(chan string, MAX_Novel_CH)
	kakuyomu_ch = make(chan string, MAX_CH)
	nnocku_ch = make(chan string, MAX_Novel_CH)
	alpha_ch = make(chan string, MAX_CH)
}

// URLから小説の種類を判定
func ChackUrlType(url string) (nobelWebType, error) {
	ckDnsTmp := []string{
		DNS_NAROU,
		DNS_KAKUYOMU,
		DNS_NOCKU,
		DNS_ALPHA,
	}
	nobeltype := ATHOR_WEB
	for i, dns := range ckDnsTmp {
		if strings.Index(url, dns) > 0 {
			nobeltype = nobelWebType(i)
			break
		}
	}
	if nobeltype == ATHOR_WEB {
		return nobeltype, errors.New("Other URL Data:" + url)
	}
	return nobeltype, nil
}

// URLから小説の種類を判定してデータを取得
func ChackUrlData(nwt nobelWebType, url string) (List, error) {
	var output List
	var outerr error
	switch nwt {
	case NAROU_WEB:
		for i := 0; i < 3; i++ {
			if data, err := getDocument(url, narou_ch); err != nil {
				log.Println(err)
				outerr = err
			} else {
				outerr = nil
				output = chackSyousetu(data)
				break
			}
			time.Sleep(time.Microsecond * 100)
		}
	case KAKUYOMU_WEB:
		for i := 0; i < 3; i++ {
			if data, err := getKakuyomu(url, kakuyomu_ch); err != nil {
				log.Println(err)
				outerr = err
			} else {
				outerr = nil
				output = chackKakuyomu(data)
				break
			}
			time.Sleep(time.Microsecond * 200)
		}
	case NNOCKU_WEB:
		for i := 0; i < 3; i++ {
			if data, err := getNokutarn(url, nnocku_ch); err != nil {
				log.Println(err)
				outerr = err
			} else {
				outerr = nil
				output = chackNokutarn(data)
				break
			}
			time.Sleep(time.Microsecond * 200)
		}
	case ALPHA_WEB:
		for i := 0; i < 3; i++ {
			if data, err := getDocument(url, alpha_ch); err != nil {
				log.Println(err)
				outerr = err
			} else {
				outerr = nil
				output = chackAlpha(data)
				break
			}
			time.Sleep(time.Microsecond * 200)
		}
	default:
		return output, errors.New("not input type")
	}
	return output, outerr
}

// func Loop(ctx context.Context){
// loop:
// for{
// 	select{
// 	case <-ctx.Done():
// 		break loop
// 	case url:=<-narou_ch:
// 	case url:=<-kakuyomu_ch:
// 	case url:=<-nnocku_ch:
// 	case url:=<-alpha_ch:
// 	}
// }
// }

// なろうの取得
// アルファポリスの取得
func getDocument(url string, ch chan string) (documentdata, error) {
	var output documentdata
	output.url = url
	ch <- url
	doc, err := goquery.NewDocument(url)
	if err == nil {
		output.data = doc
	}
	<-ch
	return output, err
}

// カクヨムの取得
func getKakuyomu(urldata string, ch chan string) (documentdata, error) {
	var output documentdata
	output.url = urldata

	ch <- urldata
	req, err := http.NewRequest(http.MethodGet, urldata, nil)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	client := new(http.Client)
	resp, err := client.Do(req)

	if err != nil {
		return output, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return output, err
	}
	output.data = doc
	<-ch
	return output, nil
}

// ノクターンノベルのゲット
func getNokutarn(urldata string, ch chan string) (documentdata, error) {
	var output documentdata
	output.url = urldata

	ch <- urldata
	req, err := http.NewRequest(http.MethodPost, urldata, nil)
	if err != nil {
		return output, err
	}
	req.Header.Set("Cookie", "over18=yes")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return output, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return output, err
	}
	output.data = doc
	<-ch
	return output, nil
}

// ノクターンノベルのチェック
func chackNokutarn(data documentdata) List {
	var output List
	output.Url = data.url
	doc := data.data
	output.Title = doc.Find("p.novel_title").Text()
	doc.Find("dl.novel_sublist2").Each(func(i int, s *goquery.Selection) {
		output.LastStoryT = strings.TrimSpace(s.Find("dd.subtitle").Text())
		times := strings.TrimSpace(s.Find("dt.long_update").Text())
		if times != "" {
			if strings.Index(times, "（改）") > 0 {
				times = times[:strings.Index(times, "（改）")]
			}
			t, _ := time.Parse("2006/01/02 15:04:05 JST", times+":00 JST")
			output.Lastdate = t.UTC().Add(-9 * time.Hour)

		}
		tmp, _ := s.Find("dd.subtitle").Find("a").Attr("href")
		if tmp != "" {
			output.LastUrl = BASE_URL_NAROU + tmp
		}
		output.Count = i + 1
	})
	return output
}

// 小説になろうのチェック
func chackSyousetu(data documentdata) List {
	var output List
	output.Url = data.url
	doc := data.data
	output.Title = doc.Find("p.novel_title").Text()
	doc.Find("dl.novel_sublist2").Each(func(i int, s *goquery.Selection) {
		output.LastStoryT = strings.TrimSpace(s.Find("dd.subtitle").Text())
		times := strings.TrimSpace(s.Find("dt.long_update").Text())
		if times != "" {
			if strings.Index(times, "（改）") > 0 {
				times = times[:strings.Index(times, "（改）")]
			}
			t, _ := time.Parse("2006/01/02 15:04:05 JST", times+":00 JST")
			output.Lastdate = t.UTC().Add(-9 * time.Hour)

		}
		tmp, _ := s.Find("dd.subtitle").Find("a").Attr("href")
		if tmp != "" {
			output.LastUrl = BASE_URL_NAROU + tmp
		}
		output.Count = i + 1
	})
	return output
}

// カクヨムのチェック
func chackKakuyomu(data documentdata) List {
	var output List
	output.Url = data.url
	doc := data.data
	output.Title = doc.Find("div#workHeader-inner").Find("h1#workTitle").Text()

	doc.Find("div.widget-toc-main").Each(func(i int, s *goquery.Selection) {
		s.Find("li.widget-toc-episode").Each(func(j int, ss *goquery.Selection) {
			output.LastStoryT = ss.Find("a").Find("span").Text()
			tmpurl, _ := ss.Find("a").Attr("href")
			if tmpurl != "" {
				output.LastUrl = BASE_URL_KAKUYOMU + tmpurl
			}
			tmpdate, _ := ss.Find("a").Find("time").Attr("datetime")
			if tmpdate != "" {
				t, _ := time.Parse("2006-01-02T15:04:05Z", tmpdate)
				output.Lastdate = t.Local()
			}
		})

	})
	return output

}

// アルファポリスのチェック
func chackAlpha(data documentdata) List {
	var output List
	output.Url = data.url
	doc := data.data
	output.Title = doc.Find("h1.title").Text()
	doc.Find("div.episode ").Each(func(i int, s *goquery.Selection) {
		output.LastStoryT = strings.TrimSpace(s.Find("span.title").Text())
		times := strings.TrimSpace(s.Find("span.open-date").Text())
		if times != "" {
			t, _ := time.Parse("2006.01.02 15:04:05 JST", times+":00 JST")
			output.Lastdate = t.UTC().Add(-9 * time.Hour)

		}
		tmp, _ := s.Find("a").Attr("href")
		if tmp != "" {
			output.LastUrl = BASE_URL_ALPHAS + tmp
		}
		output.Count = i + 1
	})
	return output
}

// 無効にする
func chackurl(url string) bool {
	flag := false
	if len(BASE_URL_NAROU) <= len(url) {
		if url[:len(BASE_URL_NAROU)] == BASE_URL_NAROU {
			flag = true
		}
	}
	if len(BASE_URL_NAROUS) <= len(url) {
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
