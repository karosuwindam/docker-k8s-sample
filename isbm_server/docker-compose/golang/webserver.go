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
	"strings"
)

type WebSetup struct {
	Ip       string `json:ip`
	Port     string `json:port`
	RootPath string `json:rootpath`
}

type WebSetupData struct {
	Data WebSetup
	flag bool
}

const (
	IPDATA       = ""
	PORTDATA     = "8080"
	ROOTPATHDATA = "./html"
	CONFIG_PATH  = "./config"
	CONFIG_FILE  = "websetup.json"
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
	}
	if f, err := os.Stat(t.Data.RootPath); os.IsNotExist(err) || !f.IsDir() {
		errtext := t.Data.RootPath + "フォルダが見つかりません。"
		return errors.New(errtext)
	} else {
		t.flag = true
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

func (t *WebSetupData) webstart() {
	if !t.flag {
		fmt.Println("Don't start web setup")
		return
	}
	fmt.Println(t.Data.Ip + ":" + t.Data.Port + "server start")
	http.HandleFunc("/v1/", apiserver)
	http.HandleFunc("/sql/", sqlapiserver)
	http.HandleFunc("/", t.viewhtml)
	http.ListenAndServe(t.Data.Ip+":"+t.Data.Port, nil)
}
