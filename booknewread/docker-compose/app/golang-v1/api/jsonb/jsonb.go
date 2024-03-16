package jsonb

import (
	"book-newread/config"
	"book-newread/loop"
	"book-newread/webserver"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var apiname string = "jsonb"

func getlocaljson(w http.ResponseWriter, r *http.Request) {
	form_data := ""
	statusdata := loop.ReadStatus()
	timedata := time.Now().Sub(statusdata.BookNowTIme).Seconds()
	if timedata > 300 {
		loop.Reset_ON(loop.RESET_DATA)
	}
	r.ParseForm()
	for cnt, strs := range r.Form {
		form_data += " " + cnt + ":"
		for i, str := range strs {
			if i == 0 {
				form_data += str
			} else {
				form_data += "," + str
			}
		}
	}
	fmt.Println(r.Method + ":" + r.URL.Path + " " + form_data)
	booklistdata := loop.ReadBookListData()
	if r.Method == "GET" {
		jsondata, err := json.Marshal(booklistdata[0])
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
		if tmp_page >= len(booklistdata) {
			tmp_page = 0
		}
		jsondata, err := json.Marshal(booklistdata[tmp_page])
		if err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			fmt.Fprintf(w, "%s", jsondata)
		}
	}
}

var route []webserver.WebConfig = []webserver.WebConfig{
	{Pass: "/" + apiname, Handler: getlocaljson},
}

// Setup
func Setup(cfg *config.Config) ([]webserver.WebConfig, error) {
	return route, nil
}
