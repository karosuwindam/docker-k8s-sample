package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

type ApiData struct {
	Message string
	slack   SlackApi
}

type SlackSentMessage struct {
	Message  string
	Filename string
	File     *os.File
}

const (
	TMP_FOLDER = "tmp/"
)

//SLACK設定ファイルの初期化
func (t *ApiData) InitSlack(slack SlackApi) {
	t.slack = slack
}

//url解析
func (t *ApiData) urlAnalysis(url string) []string {
	tmp := []string{}

	for _, str := range strings.Split(url[1:], "/") {
		tmp = append(tmp, str)
	}
	return tmp
}

//query解析
func (t *ApiData) queryAnalysis(querydata string) map[string]string {
	tmp := map[string]string{}
	for _, str := range strings.Split(querydata, "&") {
		str_s := strings.Split(str, "=")
		if len(str_s) == 2 {
			tmp[str_s[0]] = str_s[1]
		}
	}
	return tmp
}

//slackへのメッセージポスト
func (t *ApiData) slackPostMessage(api SlackSentMessage) {
	t.slack.postSlackMessage(api.Message)
}

//slackへのファイルポスト
func (t *ApiData) slackPostFile(api SlackSentMessage) {
	t.slack.postSlackFileData(api.Filename, api.File)
}

//apiの実行
func (t *ApiData) selectRun(urldata []string, querydata map[string]string, api SlackSentMessage) {
	if len(urldata) < 3 {
		t.Message = "url input data error."
		return
	}
	if urldata[0] != "api" {
		t.Message = "url input data not \"api\""
		return
	}
	if urldata[1] != "v1" {
		t.Message = "url input data not \"v1\""
		return
	}
	switch urldata[2] {
	case "postmessage":
		if api.Message != "" {
			t.slackPostMessage(api)
		} else {
			t.Message = "output message not input"
			return
		}
		break
	case "uploadfile":
		if api.File != nil {
			t.slackPostFile(api)
		} else {
			return
		}

		break
	}
	t.Message = "OK"
	return
}

//Web API解析
func (t *ApiData) ApiSelect(w http.ResponseWriter, r *http.Request) {
	urldata := t.urlAnalysis((r.URL.Path))
	querydata := t.queryAnalysis(r.URL.RawQuery)
	var tmp SlackSentMessage
	if str := r.FormValue("message"); str != "" {
		tmp.Message = str
	}
	if str := r.FormValue("filename"); str != "" {
		tmp.Filename = ""
	}
	// if false {
	// 	var err error
	// 	tmp.File, err = os.Open("./test.png")
	// 	if err != nil {
	// 		return
	// 	}
	// 	defer tmp.File.Close()
	// }
	if file, fileHeader, e := r.FormFile("file"); e == nil {
		// file, fileHeader, e := r.FormFile("file")
		defer file.Close()
		tmp.Filename = fileHeader.Filename
		fp, err := os.Create(TMP_FOLDER + tmp.Filename)
		if err != nil {

		}
		var data []byte = make([]byte, 1024)
		var tmplength int64 = 0
		for {
			n, e := file.Read(data)
			if n == 0 {
				break
			}
			if e != nil {
				return
			}
			fp.WriteAt(data, tmplength)
			tmplength += int64(n)
		}
		fp.Close()

		tmp.File, err = os.Open(TMP_FOLDER + tmp.Filename)
		if err != nil {
			return
		}
		// defer tmp.File.Close()
		defer t.rmclosefile(&tmp)

	} else {
		t.Message = e.Error()
		tmp.File = nil
	}
	if r.Method == "POST" {
		t.selectRun(urldata, querydata, tmp)
		fmt.Println(urldata, querydata)
	}
}

func (t *ApiData) rmclosefile(tmp *SlackSentMessage) {
	err := tmp.File.Close()
	if err != nil {
		return
	}
	if err := os.Remove(TMP_FOLDER + tmp.Filename); err != nil {
		fmt.Println(err)
	}

}
