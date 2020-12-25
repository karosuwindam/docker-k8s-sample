package main

import (
	"fmt"
	"net/http"
	"strings"
)

type ApiData struct {
	Message string
}

func (t *ApiData) urlAnalysis(url string) []string {
	tmp := []string{}

	for _, str := range strings.Split(url[1:], "/") {
		tmp = append(tmp, str)
	}
	return tmp
}
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

func (t *ApiData) slackPostMessage() {

}
func (t *ApiData) slackPostFile() {

}

func (t *ApiData) selectRun(urldata []string, querydata map[string]string) {
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
		t.slackPostMessage()
		break
	case "uploadfile":
		t.slackPostFile()
		break
	}
	t.Message = "OK"
	return
}

func (t *ApiData) ApiSelect(w http.ResponseWriter, r *http.Request) {
	urldata := t.urlAnalysis((r.URL.Path))
	querydata := t.queryAnalysis(r.URL.RawQuery)
	t.selectRun(urldata, querydata)
	fmt.Println(urldata, querydata)
}
