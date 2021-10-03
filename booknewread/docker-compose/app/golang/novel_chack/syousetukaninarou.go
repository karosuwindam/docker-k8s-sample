package novel_chack

import (
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
)

type Channel struct {
	Ch_Narou    chan bool
	Ch_Kakuyomu chan bool
	Ch_Nocku    chan bool
	Setup       bool
}

func Setup(count int) Channel {
	var output Channel
	if count > 2 {
		output.Ch_Narou = make(chan bool, 2)
	} else {
		output.Ch_Narou = make(chan bool, count)
	}
	output.Ch_Kakuyomu = make(chan bool, count)
	output.Ch_Nocku = make(chan bool, count)
	output.Setup = true
	return output
}

func (t *Channel) ChackUrldata(url string) List {
	var output List
	if !t.Setup {
		log.Println("not Setup")
		return output
	}
	// if len(BASE_URL_NAROU) <= len(url) {
	// 	if url[:len(BASE_URL_NAROU)] == BASE_URL_NAROU {
	// 		t.Ch_Narou <- true
	// 		output = chackSyousetu(url)
	// 		<-t.Ch_Narou
	// 	}
	// }
	// if len(BASE_URL_NAROUS) <= len(url) {
	// 	if url[:len(BASE_URL_NAROUS)] == BASE_URL_NAROUS {
	// 		t.Ch_Narou <- true
	// 		output = chackSyousetu(url)
	// 		<-t.Ch_Narou
	// 	}
	// }
	if len(BASE_URL_NAROU) <= len(url) { //なろうのチェック
		url_tmp := ""
		if url[:len(BASE_URL_NAROU)] == BASE_URL_NAROU {
			// url_tmp = url
			url_tmp = strings.Replace(url, "http", "https", 1)
		} else if len(BASE_URL_NAROUS) <= len(url) {
			if url[:len(BASE_URL_NAROUS)] == BASE_URL_NAROUS {
				url_tmp = url
			}
		}
		if url_tmp != "" {
			t.Ch_Narou <- true
			data, err := getDocument(url)
			<-t.Ch_Narou
			if err != nil {
				fmt.Println(err.Error())
				return output
			} else {
				return chackSyousetu(data)
			}
		}
	}
	if len(BASE_URL_KAKUYOMU) <= len(url) { //カクヨムのチェック
		if url[:len(BASE_URL_KAKUYOMU)] == BASE_URL_KAKUYOMU {
			t.Ch_Kakuyomu <- true
			// data, err := getDocument(url)
			data, err := getKakuyomu(url)
			<-t.Ch_Kakuyomu
			if err != nil {
				fmt.Println(err.Error())
				return output
			} else {
				return chackKakuyomu(data)
			}
		}
	}
	if len(BASE_URL_NOCKU) <= len(url) { //ノクターンチェック
		url_tmp := ""
		if url[:len(BASE_URL_NOCKU)] == BASE_URL_NOCKU {
			url_tmp = strings.Replace(url, "http", "https", 1)
		} else if len(BASE_URL_NOCKUS) <= len(url) {
			if url[:len(BASE_URL_NOCKUS)] == BASE_URL_NOCKUS {
				url_tmp = url
			}
		}
		if url_tmp != "" {
			t.Ch_Nocku <- true
			data, err := getNokutarn(url_tmp)
			<-t.Ch_Nocku
			if err != nil {
				fmt.Println(err.Error())
				return output
			} else {
				return chackNokutarn(data)
			}
		}
	}

	return output

}

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

//カクヨムの取得
func getKakuyomu(urldata string) (documentdata, error) {
	var output documentdata
	output.url = urldata
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
	return output, nil
}

//ノクターンノベルのゲット
func getNokutarn(urldata string) (documentdata, error) {
	var output documentdata
	output.url = urldata
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
	return output, nil
}

//ノクターンノベルのチェック
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
			t, _ := time.Parse("2006/01/02 15:04:05 MST", times+":00 JST")
			// fmt.Println(t.Local())
			// output.Lastdate = t.Local().Add(-9 * time.Hour)
			output.Lastdate = t.Local()

		}
		tmp, _ := s.Find("dd.subtitle").Find("a").Attr("href")
		if tmp != "" {
			output.LastUrl = BASE_URL_NAROU + tmp
		}
		output.Count = i + 1
	})
	return output

}

//小説になろうのチェック
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
			t, _ := time.Parse("2006/01/02 15:04:05 MST", times+":00 JST")
			// fmt.Println(t)
			// output.Lastdate = t.Local().Add(-9 * time.Hour)
			output.Lastdate = t.Local()

		}
		tmp, _ := s.Find("dd.subtitle").Find("a").Attr("href")
		if tmp != "" {
			output.LastUrl = BASE_URL_NAROU + tmp
		}
		output.Count = i + 1
	})
	return output
}

//カクヨムのチェック
func chackKakuyomu(data documentdata) List {
	var output List
	output.Url = data.url
	doc := data.data
	output.Title = doc.Find("div#workHeader-inner").Find("h1#workTitle").Text()
	// fmt.Println(doc.Find("div.widget-toc-main").Text())

	doc.Find("div.widget-toc-main").Each(func(i int, s *goquery.Selection) {
		s.Find("li.widget-toc-episode").Each(func(j int, ss *goquery.Selection) {
			output.LastStoryT = ss.Find("a").Find("span").Text()
			tmpurl, _ := ss.Find("a").Attr("href")
			if tmpurl != "" {
				output.LastUrl = BASE_URL_KAKUYOMU + tmpurl
			}
			// fmt.Println(ss.Find("a").Find("time").Text())
			tmpdate, _ := ss.Find("a").Find("time").Attr("datetime")
			if tmpdate != "" {
				t, _ := time.Parse("2006-01-02T15:04:05Z", tmpdate)
				// jst := time.FixedZone("Asia/Tokyo", 9*60*60)

				// fmt.Println(t.Local())
				output.Lastdate = t.Local()
			}
		})

	})
	return output

}

//無効にする
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
