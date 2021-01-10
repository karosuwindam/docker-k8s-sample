package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const (
	URL_SQL_API        = 0
	URL_SQL_TYPE       = 1
	URL_SQL_TABLE_NAME = 2
	URL_SQL_DATA_1     = 3
	URL_SQL_DATA_MAX   = 4
)

//編集用CGI
func edit(w http.ResponseWriter, r *http.Request) {
	urltmp := urlAnalysis(r.URL.Path)
	tmp := map[string]string{}
	var tmpdata []string
	if len(urltmp) < URL_SQL_DATA_MAX {
		w.WriteHeader(400)
		fmt.Fprintf(w, "400 Bad Request")
		return
	}
	if r.Method == "POST" {
		switch urltmp[URL_SQL_TABLE_NAME] {
		case "isbn":
			tmpdata = bookdatatable.Column_name
		default:
			w.WriteHeader(400)
			fmt.Fprintf(w, "400 Bad Request")
			return
		}
		err := r.ParseForm()
		if err != nil {
			// エラー処理
			fmt.Fprint(w, "ng")
			return
		} else {
			for _, keyword := range tmpdata {
				tmp[keyword] = r.FormValue(keyword)
			}
			if urltmp[URL_SQL_TABLE_NAME] == "isbn" {
				fmt.Println(urltmp[URL_SQL_DATA_1], tmp["isbn"], tmp["title"], tmp["writer"], tmp["brand"], tmp["ext"], tmp["synopsis"], tmp["image"])
				bookdatatable.Edit(urltmp[URL_SQL_DATA_1], tmp["isbn"], tmp["title"], tmp["writer"], tmp["brand"], tmp["ext"], tmp["synopsis"], tmp["image"])
			}
			tmp["inputdata"] = ConvertData(ReadHtml("tmp_html/"+urltmp[URL_SQL_TABLE_NAME]+"/show.html"), tmp)
			fmt.Fprint(w, ConvertData(ReadHtml("tmp_html/"+"show_.html"), tmp))
		}

	} else {

		if urltmp[URL_SQL_TABLE_NAME] == "isbn" {
			tmpdata := bookdatatable.ScansqlId(urltmp[URL_SQL_DATA_1])
			tmp["id"] = strconv.Itoa(tmpdata.Id)
			tmp["isbn"] = tmpdata.Isbn
			tmp["title"] = tmpdata.Title
			tmp["writer"] = tmpdata.Writer
			tmp["brand"] = tmpdata.Brand
			tmp["ext"] = tmpdata.Ext
			tmp["synopsis"] = tmpdata.Synopsis
			tmp["image"] = tmpdata.Image
		}
		tmp["inputdata"] = ConvertData(ReadHtml("tmp_html/"+urltmp[URL_SQL_TABLE_NAME]+"/new.html"), tmp)
		fmt.Fprint(w, ConvertData(ReadHtml("tmp_html/"+"edit_.html"), tmp))

	}

}

//削除用CGI
func destory(w http.ResponseWriter, r *http.Request) {
	urltmp := urlAnalysis(r.URL.Path)

	if len(urltmp) < URL_SQL_DATA_MAX {
		w.WriteHeader(400)
		fmt.Fprintf(w, "400 Bad Request")
		return
	}
	if r.Method == "POST" {
		switch urltmp[URL_SQL_TABLE_NAME] {
		case "isbn":
			bookdatatable.Delte(urltmp[URL_SQL_DATA_1])
		default:
			w.WriteHeader(400)
			fmt.Fprintf(w, "400 Bad Request")
			return
		}
	} else {
		w.WriteHeader(400)
		fmt.Fprintf(w, "400 Bad Request")
		return
	}

}

func sqlapiserver(w http.ResponseWriter, r *http.Request) {
	urldata := urlAnalysis(r.URL.Path)
	log.Printf("request:%v\n", r.URL.Path)
	if len(urldata) > 1 {
		switch urldata[URL_SQL_TYPE] {
		case "destory":
			destory(w, r)
		case "edit":
			edit(w, r)
		default:

		}

	} else {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Err API request")
	}
}
