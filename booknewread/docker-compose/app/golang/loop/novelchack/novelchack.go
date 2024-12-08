package novelchack

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"book-newread/config"

	"github.com/PuerkitoBio/goquery"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type List struct {
	Title      string    `json:tilte`
	Count      int       `json:count`
	Lastdate   time.Time `json:lastdate`
	Url        string    `json:url`
	LastStoryT string    `json:laststoryt`
	LastUrl    string    `json:lasturl`
}

type NarouAPIOUT struct {
	Allcount        int    `json:allcount`
	Title           string `json:title`
	Ncode           string `json:ncode`
	Userid          int    `json:userid`
	Writer          string `json:writer`
	Story           string `json:story`
	Biggenre        int    `json:biggenre`
	Genre           int    `json:genre`
	Gensaku         string `json:gensaku`
	Keyword         string `json:keyword`
	General_firstup string `json:general_firstup`
	General_lastup  string `json:general_lastup`
	Novel_type      int    `json:novel_type`
	End             int    `json:end`
	General_all_no  int    `json:general_all_no`
	Length          int    `json:length`
	Time            int    `json:time`
	Isstop          int    `json:isstop`
	Isr15           int    `json:isr15`
	Isbl            int    `json:isbl`
	Isgl            int    `json:isgl`
	Iszankoku       int    `json:iszankoku`
	Istensei        int    `json:istensei`
	Istenni         int    `json:istenni`
	Global_point    int    `json:global_point`
	Daily_point     int    `json:daily_point`
	Weekly_point    int    `json:weekly_point`
	Monthly_point   int    `json:monthly_point`
	Quarter_point   int    `json:quarter_point`
	Yearly_point    int    `json:yearly_point`
	Fav_novel_cnt   int    `json:fav_novel_cnt`
	Impression_cnt  int    `json:impression_cnt`
	Review_cnt      int    `json:review_cnt`
	All_point       int    `json:all_point`
	All_hyoka_cnt   int    `json:all_hyoka_cnt`
	Sasie_cnt       int    `json:sasie_cnt`
	Kaiwaritu       int    `json:kaiwaritu`
	Novelupdated_at int    `json:novelupdated_at`
	Updated_at      int    `json:updated_at`
}

type NaroudData struct {
	baseurl string
	apikey  string
	body    []NarouAPIOUT
}

type documentdata struct {
	url  string
	data *goquery.Document
}

const (
	BASE_URL_NAROUS   = "https://ncode.syosetu.com"
	BASE_URL_KAKUYOMU = "https://kakuyomu.jp"
	BASE_URL_NOCKUS   = "https://novel18.syosetu.com"
	BASE_URL_ALPHAS   = "https://www.alphapolis.co.jp"
	KAKUYOMU_SIDEBER  = "episode_sidebar"
)

const (
	API_URL_NAROU = "https://api.syosetu.com/novelapi/api/"
	API_URL_NOCKU = "https://api.syosetu.com/novel18api/api/"
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
	OTHOR_WEB
)

var narou_ch chan string    //なろう
var kakuyomu_ch chan string //カクヨム
var nnocku_ch chan string   //ノクターン
var alpha_ch chan string    //アルファポリス

func Init() error {
	narou_ch = make(chan string, config.Nobel.MaxNarouAPI)
	kakuyomu_ch = make(chan string, config.Nobel.MaxKakuyomuAPI)
	nnocku_ch = make(chan string, config.Nobel.MaxNarou18API)
	alpha_ch = make(chan string, config.Nobel.MaxAlphaAPI)
	return nil

}

// URKLから対象のホームページを判断
func chackUrlPage(url string) nobelWebType {
	ckDnsTmp := []string{
		DNS_NAROU,
		DNS_KAKUYOMU,
		DNS_NOCKU,
		DNS_ALPHA,
	}
	output := OTHOR_WEB
	for i, dns := range ckDnsTmp {
		if strings.Index(url, dns) > 0 {
			output = nobelWebType(i)
			break
		}
	}

	return output
}

var ErrAtherUrl = errors.New("other url data")

// URLから小説の種類を判別してデータを取得
func ChackURLData(ctx context.Context, url string, name string) (List, error) {
	var output List
	var outerr error = nil
	url = strings.Replace(url, "http://", "https://", 1)

	switch chackUrlPage(url) {
	case NAROU_WEB:
		ctx, span := config.TracerS(ctx, "ChackURLData", "Narou")
		defer span.End()
		slog.DebugContext(ctx, "Check Narou URL", "url", url)
		//なろうのAPIを使用して解析
		if d, err := narouAPIRead(ctx, API_URL_NAROU, url); err == nil {
			output, outerr = d.narouChangeList(ctx)
			if output.Title == "" {
				output.Title = name
				output.Url = url
			}
		} else {
			return output, err
		}
	case KAKUYOMU_WEB:

		ctx, span := config.TracerS(ctx, "ChackURLData", "Kakuyomu")
		defer span.End()
		slog.DebugContext(ctx, "Check Kakuyomu URL", "url", url)
		//カクヨムのページからスクレイピング
		for i := 0; i < 3; i++ {
			if data, err := getKakuyomu(ctx, url); err != nil {
				slog.ErrorContext(ctx, "getKakuyomu", "error", err)
				outerr = err
			} else {
				outerr = nil
				output = chackKakuyomu(ctx, data)
				break
			}
			time.Sleep(time.Microsecond * 200)
		}
	case NNOCKU_WEB:
		ctx, span := config.TracerS(ctx, "ChackURLData", "Narou18")
		defer span.End()
		slog.DebugContext(ctx, "Check Narou18 URL", "url", url)
		//なろうR18のAPIを使用して解析
		if d, err := narouAPIRead(ctx, API_URL_NOCKU, url); err == nil {
			output, outerr = d.narouChangeList(ctx)
			if output.Title == "" {
				output.Title = name
				output.Url = url
			}
		} else {
			slog.ErrorContext(ctx, "narouAPIRead", "error", err)
			return output, err
		}
	case ALPHA_WEB:

		ctx, span := config.TracerS(ctx, "ChackURLData", "Alpha")
		defer span.End()
		slog.DebugContext(ctx, "Check Alpha URL", "url", url)
		//アルファポリスのページからスクレイピング
		for i := 0; i < 3; i++ {
			if data, err := getAlpha(ctx, url); err != nil {
				slog.ErrorContext(ctx, "getAlpha", "error", err)
				outerr = err
			} else {
				outerr = nil
				output = chackAlpha(ctx, data)
				break
			}
			time.Sleep(time.Microsecond * 200)
		}
	case OTHOR_WEB:
		outerr = ErrAtherUrl
	}
	return output, outerr
}

// なろう系のAPIを使用して解析
func narouAPIRead(ctx context.Context, api, url string) (NaroudData, error) {
	var output NaroudData
	// 起動制限
	if api == API_URL_NAROU {
		output.baseurl = BASE_URL_NAROUS
		narou_ch <- url
		defer func(ch chan string) {
			<-ch
		}(narou_ch)
	} else if api == API_URL_NOCKU {
		output.baseurl = BASE_URL_NOCKUS
		nnocku_ch <- url
		defer func(ch chan string) {
			<-ch
		}(nnocku_ch)
	} else {
		return output, errors.New("api is not input data")
	}
	if n := strings.Index(url, output.baseurl); n >= 0 {
		tmp_url := url[n+len(output.baseurl):]
		if tmp_url[0:1] == "/" {
			tmp_url = tmp_url[1:]
		}
		if i := strings.Index(tmp_url, "/"); i >= 1 {
			output.apikey = tmp_url[:i]
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, api+"?ncode="+output.apikey+"&out=json", nil)
		if err != nil {
			return output, err
		}
		// client := new(http.Client)
		client := http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		}
		resp, err := client.Do(req)

		if err != nil {
			return output, err
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return output, err
		}
		output.body = []NarouAPIOUT{}
		json.Unmarshal(b, &output.body)
	} else {
		return output, errors.New("error for url")
	}
	return output, nil
}

func (nd *NaroudData) narouChangeList(ctx context.Context) (List, error) {
	ctx, span := config.TracerS(ctx, "narouChangeList", nd.baseurl+"/"+nd.apikey+"/")
	defer span.End()
	out := List{}
	if len(nd.body) == 2 {
		if nd.body[1].Title != "" {
			out.Title = nd.body[1].Title
		}
		out.Url = nd.baseurl + "/" + nd.apikey + "/"
		tmpLTime := nd.body[1].General_lastup
		t, _ := time.Parse("2006-01-02 15:04:05 JST", tmpLTime+" JST")
		out.Lastdate = t.UTC().Add(-9 * time.Hour)
		count := nd.body[1].General_all_no
		out.LastUrl = out.Url + strconv.Itoa(count) + "/"
		out.LastStoryT = strconv.Itoa(count) + "話"
		if nd.body[1].End == 0 {
			out.LastStoryT = "全" + out.LastStoryT
		}

		span.SetAttributes(attribute.String("url", out.Url))
		span.SetAttributes(attribute.String("title", out.Title))
		span.SetAttributes(attribute.String("LastUrl", out.LastUrl))
		span.SetAttributes(attribute.String("LastStoryT", out.LastStoryT))
	} else {
		slog.WarnContext(ctx, "Not Found Data URL:"+nd.baseurl+"/"+nd.apikey+"/", "url", nd.baseurl, "apikey", nd.apikey)
		span.SetStatus(codes.Error, "Not Found Data")
		// out.Title =
	}
	return out, nil
}

// アルファポリスの取得
func getAlpha(ctx context.Context, urldata string) (documentdata, error) {
	var output documentdata
	output.url = urldata

	kakuyomu_ch <- urldata
	defer func(ch chan string) {
		<-ch
	}(kakuyomu_ch)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urldata, nil)
	// client := new(http.Client)
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

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

// カクヨムの取得
func getKakuyomu(ctx context.Context, urldata string) (documentdata, error) {
	var output documentdata
	output.url = urldata

	kakuyomu_ch <- urldata
	defer func(ch chan string) {
		<-ch
	}(kakuyomu_ch)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urldata, nil)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	// client := new(http.Client)
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

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

// カクヨムのチェック
func chackKakuyomu(ctx context.Context, data documentdata) List {
	ctx, span := config.TracerS(ctx, "checkKakuyomu", data.url)
	defer span.End()
	var output List
	output.Url = data.url
	span.SetAttributes(attribute.String("url", output.Url))
	doc := data.data
	//div.NewBox_box__45ont.NewBox_padding-px-4l__Kx_xT.NewBox_padding-pt-7l__Czm59
	output.Title = doc.Find("div#workHeader-inner").Find("h1#workTitle").Text()
	if output.Title == "" {
		// tmpTitle, _ := doc.Find("div.NewBox_box__45ont.NewBox_padding-px-4l__Kx_xT.NewBox_padding-pt-7l__Czm59").Find("a").Attr("title")
		tmpTitle, _ := doc.Find("div.NewBox_box__45ont.NewBox_padding-px-4l__Kx_xT").Find("a").Attr("title")
		if tmpTitle != "" {
			output.Title = tmpTitle
		}
	}
	span.SetAttributes(attribute.String("title", output.Title))

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
	//第一話の取り出しとサイドバー処理
	//Layout_layout__5aFuw Layout_items-normal__4mOqD Layout_justify-normal__zqNe7 Layout_direction-row__boh0Z Layout_wrap-wrap__yY3zM Layout_gap-2s__xUCm0
	doc.Find("div.Layout_layout__5aFuw.Layout_items-normal__4mOqD.Layout_justify-normal__zqNe7.Layout_direction-row__boh0Z.Layout_wrap-wrap__yY3zM.Layout_gap-2s__xUCm0").Each(func(i int, s *goquery.Selection) {
		tmpurl, _ := s.Find("a").Attr("href")
		if tmpurl != "" {
			output.LastUrl = BASE_URL_KAKUYOMU + tmpurl
		}
		tmpdoc, err := getKakuyomu(ctx, BASE_URL_KAKUYOMU+tmpurl+"/"+KAKUYOMU_SIDEBER)
		if err == nil {
			eTempUrl := ""
			eEpText := ""
			eEpTime := ""
			tmpdoc.data.Find("li").Each(func(i int, ss *goquery.Selection) {
				eTempUrl, _ = ss.Find("a").Attr("href")
				eEpText = ss.Find("a").Find("span").Text()
				eEpTime, _ = ss.Find("a").Find("time").Attr("datetime")
			})
			t, _ := time.Parse("2006-01-02T15:04:05Z", eEpTime)
			if output.Lastdate.Before(t.Local()) {
				output.LastStoryT = eEpText
				output.LastUrl = BASE_URL_KAKUYOMU + eTempUrl
				output.Lastdate = t.Local()
			}
		}
	})

	if (output.LastStoryT == "") && (output.LastUrl == "") {
		infoTmp := doc.Find("div.Typography_fontSize-m__mskXq.Typography_color-gray__ObCRz.Typography_lineHeight-3s__OOxkK.Base_block__H4wj4")
		timedata, _ := infoTmp.Find("time").Attr("datetime")
		t, _ := time.Parse("2006-01-02T15:04:05Z", timedata)
		output.Lastdate = t.Local()

		countTmp := infoTmp.Find("ul.Meta_meta__7tVPt.Meta_disc__uPSnA.Meta_lightGray__mzmje.Meta_lineHeightXsmall__66NnD").Find("div.Meta_metaItem__8eZTP")
		countTmp.Each(func(i int, s *goquery.Selection) {
			if strings.Index(s.Text(), "話") > 0 {
				if output.LastStoryT == "" {
					output.LastStoryT = s.Text()
				}
			}
		})
		output.LastUrl = output.Url

	}
	return output

}

// アルファポリスのチェック
func chackAlpha(ctx context.Context, data documentdata) List {
	ctx, span := config.TracerS(ctx, "chackAlpha", data.url)
	defer span.End()
	var output List
	output.Url = data.url
	doc := data.data
	output.Title = doc.Find("h1.title").Text()
	span.SetAttributes(attribute.String("url", output.Url))
	span.SetAttributes(attribute.String("title", output.Title))

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
			if len(tmp) > len(BASE_URL_ALPHAS) {
				if BASE_URL_ALPHAS == tmp[:len(BASE_URL_ALPHAS)] {
					output.LastUrl = tmp
				}

			}
		}
		output.Count = i + 1
	})
	span.SetAttributes(attribute.String("LastUrl", output.LastUrl))
	span.SetAttributes(attribute.String("LastStoryT", output.LastStoryT))
	return output
}
