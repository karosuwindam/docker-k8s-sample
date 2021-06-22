package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ashwanthkumar/slack-go-webhook"
)

const (
	CONFIGPASS = "conf/conf.json"
	COUNTTIME  = 5
	TMPFILE    = "./tmp"
)

type SenserPut struct {
	Node  string        `json:node`
	Datas []Sennserdata `json:datas`
}

type Sennserdata struct {
	Senser string `json:senser`
	Type   string `json:type`
	Data   string `json:data`
}

type SlackSetUp struct {
	WebUrl    string
	Channel   string
	UserName  string
	CountTime int
	TimeZone  string
}

//コンテナ用、環境変数収集
func EnvConfRead() SlackSetUp {
	var tmp SlackSetUp
	if str := os.Getenv("SLACK_URL"); str != "" {
		tmp.WebUrl = str
	}
	if str := os.Getenv("SLACK_CHANNEL"); str != "" {
		tmp.Channel = str
	}
	if str := os.Getenv("SLACK_USERNAME"); str != "" {
		tmp.UserName = str
	}
	if str := os.Getenv("SLACK_COUNTTIME"); str != "" {
		tmp.CountTime, _ = strconv.Atoi(str)
	} else {
		tmp.CountTime = COUNTTIME
	}
	if str := os.Getenv("TZ"); str != "" {
		tmp.TimeZone = str
	} else {
		tmp.TimeZone = "Asia/Tokyo"
	}
	return tmp
}

//Slackポスト
func (t *SlackSetUp) PostSlack(msg string) {
	field1 := slack.Field{Title: "部屋のセンサー値", Value: msg}

	attachment := slack.Attachment{}
	attachment.AddField(field1)
	color := "good"
	attachment.Color = &color
	payload := slack.Payload{
		Username:    t.UserName,
		Channel:     t.Channel,
		Attachments: []slack.Attachment{attachment},
	}
	err := slack.Send(t.WebUrl, "", payload)
	if err != nil {
		log.Println(err)
	}
}

//Slack送信加工用
func output_data(datas []SenserPut) string {
	output := ""
	for _, data := range datas {
		if output != "" {
			output += "\n"
		}
		output += data.Node
		for _, tdata := range data.Datas {
			output += "\n " + tdata.Senser + " " + tdata.Type + ":" + tdata.Data
		}
	}
	return output
}

// URLからセンサーデータ取得プログラム
func GetSennserdata(datas []UrlData) []SenserPut {
	output := []SenserPut{}
	for _, data := range datas {
		tmp := getjson(data.Url)
		if tmp != "" {
			var json_tmp []Sennserdata
			json.Unmarshal([]byte(tmp), &json_tmp)
			if json_tmp != nil {
				var tmpdata SenserPut
				tmpdata.Node = data.Node
				tmpdata.Datas = json_tmp
				output = append(output, tmpdata)
			}
		}
	}
	return output
}

//センサーJSONデータ取得部分
func getjson(url string) string {
	client := &http.Client{}
	client.Timeout = time.Second * 15
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}
