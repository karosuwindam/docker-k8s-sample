package novel_chack

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
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
	NAROU_CK_SLEEP    = time.Millisecond * 500 //500ms なろう待ち時間
	KAKUYOMU_CK_SLEEP = time.Millisecond * 500 //500ms カクヨム待ち時間
	NOCKU_CK_SLEEP    = time.Millisecond * 500 //500ms ノクターン待ち時間
	ALPHA_CK_SLEEP    = time.Millisecond * 500 //500ms アルファポリス待ち時間
	MAX_LOOP_COUNT    = 20
)

type Channel struct {
	Ch_Narou    sync.Mutex //なろう
	Ch_Kakuyomu sync.Mutex //カクヨム
	Ch_Nocku    sync.Mutex //ノクターン
	Ch_Alpha    sync.Mutex //アルファポリス
	flag        bool
}

var channel_data Channel
var maxloop chan bool

func Setup() {
	channel_data = Channel{flag: true}
	maxloop = make(chan bool, MAX_LOOP_COUNT)
}

func ChackUrldata(url string) (List, error) {
	var output List
	var wp sync.WaitGroup
	if !channel_data.flag {
		log.Println("not Setup")
		return output, errors.New("not Setup")
	}
	maxloop <- true
	if len(BASE_URL_NAROU)+1 < len(url) { //なろうのチェック
		url_tmp := ""
		if url[:len(BASE_URL_NAROU)] == BASE_URL_NAROU {
			url_tmp = strings.Replace(url, "http", "https", 1)
		} else if len(BASE_URL_NAROUS) <= len(url) {
			if url[:len(BASE_URL_NAROUS)] == BASE_URL_NAROUS {
				url_tmp = url
			}
		}
		if url_tmp != "" {
			channel_data.Ch_Narou.Lock()
			data, err := getDocument(url)
			wp.Add(1)
			go func() {
				defer wp.Done()
				time.Sleep(NAROU_CK_SLEEP)
				channel_data.Ch_Narou.Unlock()
			}()
			if err != nil {
				fmt.Println(err.Error())
				return output, err
			} else {
				return chackSyousetu(data), nil
			}
		}
	}
	if len(BASE_URL_KAKUYOMU)+1 < len(url) { //カクヨムのチェック
		if url[:len(BASE_URL_KAKUYOMU)] == BASE_URL_KAKUYOMU {
			data, err := getKakuyomu(url)
			if err != nil {
				fmt.Println(err.Error())
				return output, err
			} else {
				return chackKakuyomu(data), nil
			}
		}
	}
	if len(BASE_URL_NOCKU)+1 < len(url) { //ノクターンチェック
		url_tmp := ""
		if url[:len(BASE_URL_NOCKU)] == BASE_URL_NOCKU {
			url_tmp = strings.Replace(url, "http", "https", 1)
		} else if len(BASE_URL_NOCKUS) <= len(url) {
			if url[:len(BASE_URL_NOCKUS)] == BASE_URL_NOCKUS {
				url_tmp = url
			}
		}
		if url_tmp != "" {
			data, err := getNokutarn(url_tmp)
			if err != nil {
				fmt.Println(err.Error())
				return output, err
			} else {
				return chackNokutarn(data), nil
			}
		}
	}
	if len(BASE_URL_ALPHA)+1 < len(url) { //アルファポリス
		url_tmp := ""
		if url[:len(BASE_URL_ALPHA)] == BASE_URL_ALPHA {
			url_tmp = strings.Replace(url, "http", "https", 1)
		} else if len(BASE_URL_ALPHAS) <= len(url) {
			if url[:len(BASE_URL_ALPHAS)] == BASE_URL_ALPHAS {
				url_tmp = url
			}
		}
		if url_tmp != "" {
			channel_data.Ch_Alpha.Lock()
			data, err := getDocument(url_tmp)
			wp.Add(1)
			go func() {
				defer wp.Done()
				time.Sleep(ALPHA_CK_SLEEP)
				channel_data.Ch_Alpha.Unlock()
			}()
			if err != nil {
				fmt.Println(err.Error())
				return output, err
			} else {
				return chackAlpha(data), nil
			}
		}
	}
	wp.Wait()
	<-maxloop

	return output, nil

}

// なろうの取得
// アルファポリスの取得
func getDocument(url string) (documentdata, error) {
	var output documentdata
	output.url = url
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return output, err
	}
	output.data = doc
	return output, nil
}

// カクヨムの取得
func getKakuyomu(urldata string) (documentdata, error) {
	var output documentdata
	var wp sync.WaitGroup
	output.url = urldata

	channel_data.Ch_Kakuyomu.Lock()
	req, err := http.NewRequest(http.MethodGet, urldata, nil)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	client := new(http.Client)
	resp, err := client.Do(req)
	wp.Add(1)
	go func() {
		defer wp.Done()
		time.Sleep(KAKUYOMU_CK_SLEEP)
		channel_data.Ch_Kakuyomu.Unlock()
	}()

	if err != nil {
		return output, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return output, err
	}
	output.data = doc
	wp.Wait()
	return output, nil
}

// ノクターンノベルのゲット
func getNokutarn(urldata string) (documentdata, error) {
	var output documentdata
	var wp sync.WaitGroup
	output.url = urldata

	channel_data.Ch_Nocku.Lock()
	req, err := http.NewRequest(http.MethodPost, urldata, nil)
	if err != nil {
		channel_data.Ch_Nocku.Unlock()
		return output, err
	}
	req.Header.Set("Cookie", "over18=yes")
	resp, err := http.DefaultClient.Do(req)
	wp.Add(1)
	go func() {
		defer wp.Done()
		time.Sleep(NOCKU_CK_SLEEP)
		channel_data.Ch_Nocku.Unlock()
	}()

	if err != nil {
		wp.Wait()
		return output, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		wp.Wait()
		return output, err
	}
	output.data = doc
	wp.Wait()
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
