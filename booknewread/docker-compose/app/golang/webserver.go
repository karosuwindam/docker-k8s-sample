package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type WebSetup struct {
	Ip       string `json:ip`
	Port     string `json:port`
	RootPath string `json:rootpath`
	Template string `json:template`
}

type WebSetupData struct {
	Data WebSetup
	flag bool
}

const (
	IPDATA           = ""
	PORTDATA         = "8080"
	ROOTPATHDATA     = "./"
	ROOTTEMPLATEDATA = "./"
	CONFIG_PATH      = "./config"
	CONFIG_FILE      = "websetup.json"
)
const (
	basicAuthUser     = "admin"
	basicAuthPassword = "admin"
)

func (t *WebSetupData) websetup() error {
	config_json := CONFIG_PATH + "/" + CONFIG_FILE
	raw, err := ioutil.ReadFile(config_json)
	var buf bytes.Buffer
	if err != nil {
		t.Data.Ip = IPDATA
		t.Data.Port = PORTDATA
		t.Data.RootPath = ROOTPATHDATA
		t.Data.Template = ROOTTEMPLATEDATA
		if f, err := os.Stat(CONFIG_PATH); os.IsNotExist(err) || !f.IsDir() {
			_ = os.Mkdir(CONFIG_PATH, 0777)
		}
		fp, err := os.Create(config_json)
		if err != nil {
			return err
		}
		// jsonエンコード
		outputJson, err := json.Marshal(&t.Data)
		if err != nil {
			return err
		}
		json.Indent(&buf, outputJson, "", "  ")
		fp.Write(buf.Bytes())
		fp.Close()
	} else {
		var fc WebSetup
		json.Unmarshal(raw, &fc)
		t.Data = fc
		if t.Data.Port == "" {
			t.Data.Port = PORTDATA
		}
		if t.Data.RootPath == "" {
			t.Data.RootPath = ROOTPATHDATA
		}
		if t.Data.Template == "" {
			t.Data.Template = ROOTTEMPLATEDATA
		}
	}
	if f, err := os.Stat(t.Data.RootPath); os.IsNotExist(err) || !f.IsDir() {
		errtext := t.Data.RootPath + "フォルダが見つかりません。"
		return errors.New(errtext)
	} else {
		if f, err := os.Stat(t.Data.Template); os.IsNotExist(err) || !f.IsDir() {
			errtext := t.Data.Template + "フォルダが見つかりません。"
			return errors.New(errtext)
		} else {
			t.flag = true
		}
	}
	return nil
}

//静的HTMLのページを返す
func (t *WebSetupData) viewhtml(w http.ResponseWriter, r *http.Request) {
	textdata := []string{".html", ".htm", ".css", ".js"}
	upath := r.URL.Path
	tmp := map[string]string{}
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	if upath == "/" {
		upath += "index.html"
	}
	if !Exists(t.Data.RootPath + upath) {
		w.WriteHeader(404)
		log.Printf("ERROR request:%v\n", r.URL.Path)
		return
	} else {
		for _, data := range textdata {
			if len(upath) > len(data) {
				if upath[len(upath)-len(data):] == data {
					fmt.Fprint(w, ConvertData(ReadHtml(t.Data.RootPath+upath), tmp))
					return
				}
			}
		}
		buffer := ReadOther(t.Data.RootPath + upath)
		// bodyに書き込み
		w.Write(buffer)
	}
	return
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
func (t *WebSetupData) json(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, getjson("https://books.rakuten.co.jp/event/book/comic/calendar/2021/05/js/booklist.json"))
}

func (t *WebSetupData) getlocaljson(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		jsondata, err := json.Marshal(GrobalListData[0])
		if err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			fmt.Fprintf(w, "%s", jsondata)
		}
	} else if r.Method == "POST" {
		mf := r.MultipartForm

		// 通常のリクエスト
		if mf != nil {
			for k, v := range mf.Value {
				fmt.Printf("%v : %v", k, v)
			}
		}
		page := r.FormValue("page")
		tmp_page, _ := strconv.Atoi(page)
		if tmp_page >= len(GrobalListData) {
			tmp_page = 0
		}
		jsondata, err := json.Marshal(GrobalListData[tmp_page])
		if err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			fmt.Fprintf(w, "%s", jsondata)
		}
	}
}
func (t *WebSetupData) getnowdata(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		jsondata, err := json.Marshal(Listdata)
		if err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			fmt.Fprintf(w, "%s", jsondata)
		}
	}
}

func (t *WebSetupData) webstart() {
	if !t.flag {
		fmt.Println("Don't start web setup")
		return
	}
	fmt.Println(t.Data.Ip + ":" + t.Data.Port + "server start")
	http.HandleFunc("/", t.viewhtml)
	http.HandleFunc("/json", t.json)
	http.HandleFunc("/jsonb", t.getlocaljson)
	http.HandleFunc("/jsonnobel", t.getnowdata)
	http.ListenAndServe(t.Data.Ip+":"+t.Data.Port, nil)
}
