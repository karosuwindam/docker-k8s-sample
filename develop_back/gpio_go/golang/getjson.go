package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	CONFIGPASS = "conf/conf.json"
	COUNTTIME  = 5
	TMPFILE    = "./tmp"
)

type Sennserdata struct {
	Senser string `json:senser`
	Type   string `json:type`
	Data   string `json:data`
}

type Urldata struct {
	Name string `json:name`
	Url  string `json:url`
}

func getdata() []Sennserdata {
	urldata, _ := ConfJsonUrlRead()
	outdata := ""
	var output []Sennserdata
	for i := 0; i < len(urldata); i++ {
		tmp := getjson(urldata[i].Url)
		if tmp != "" {
			var json_tmp []Sennserdata
			json.Unmarshal([]byte(tmp), &json_tmp)
			if outdata != "" {
				outdata += "\n"
			}
			outdata += urldata[i].Name + "\n"
			outdata += output_data(json_tmp)
			for _, tmp := range json_tmp {
				output = append(output, tmp)
			}
		}
	}
	return output
}

func ConfJsonUrlRead() ([]Urldata, error) {
	raw, err := ioutil.ReadFile(CONFIGPASS)
	if err == nil {
		var fc []Urldata
		json.Unmarshal(raw, &fc)
		return fc, nil
	} else {
		return nil, err
	}
}

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
func output_data(data []Sennserdata) string {
	output := ""
	for i := 0; i < len(data); i++ {
		if output != "" {
			output += "\n"
		}
		output += " " + data[i].Senser + " " + data[i].Type + ":" + data[i].Data
	}
	return output
}
