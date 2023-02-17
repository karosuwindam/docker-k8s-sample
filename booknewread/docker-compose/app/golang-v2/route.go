package main

import (
	"booknewread/textread"
	"booknewread/webserver"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	ROUTE = "./html"
)

type Htmldata struct {
	Loopdata *LoopData
}

//静的HTMLのページを返す
func viewhtml(w http.ResponseWriter, r *http.Request) {
	textdata := []string{".html", ".htm", ".css", ".js"}
	upath := r.URL.Path
	tmp := map[string]string{}
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	log.Println(r.Method + ":" + r.URL.Path)
	if upath == "/" {
		upath += "index.html"
	}
	if !textread.Exists(ROUTE + upath) {
		w.WriteHeader(404)
		log.Printf("ERROR request:%v %v\n", r.URL.Path, ROUTE+upath)
		return
	} else {
		for _, data := range textdata {
			if len(upath) > len(data) {
				if upath[len(upath)-len(data):] == data {
					fmt.Fprint(w, textread.ConvertData(textread.ReadHtml(ROUTE+upath), tmp))
					return
				}
			}
		}
		buffer := textread.ReadOther(ROUTE + upath)
		// bodyに書き込み
		w.Write(buffer)
	}
	return
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Path: %v", r.URL.Path[1:])
}

func setupbaseRoute() (Htmldata, error) {
	var err error
	output := Htmldata{}
	output.Loopdata, err = loopSetup()
	if f, err := os.Stat(ROUTE); os.IsNotExist(err) || !f.IsDir() {
		errtext := ROUTE + "フォルダが見つかりません。"
		return output, errors.New(errtext)
	}
	if err != nil {
		return output, err
	}
	return output, nil
}

func (t *Htmldata) setupRoute(cfg *webserver.SetupServer) {
	cfg.Add("/", viewhtml)
	cfg.Add("/status", t.Loopdata.status)
	cfg.Add("/restart", t.Loopdata.restart)
	cfg.Add("/jsonb", t.Loopdata.getlocaljson)
	cfg.Add("/jsonnobel", t.Loopdata.getnowdata)
}
