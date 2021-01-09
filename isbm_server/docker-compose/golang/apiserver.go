package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"./searchapi"
)

func urlAnalysis(url string) []string {
	tmp := []string{}
	for _, str := range strings.Split(url[1:], "/") {
		tmp = append(tmp, str)
	}
	return tmp
}
func sqlserch(isbn string) []SqlBookData {
	tmp := bookdatatable.Serch_sql(isbn)
	if len(tmp) > 0 {
		return tmp
	}
	return nil
}
func newserchisbn(isbn string, no int) []SqlBookData {
	var adata searchapi.AmazonNameType
	ch1 := make(chan bool)
	if len(isbn) == 13 {
		go func() {
			adata = searchapi.GetPageAmazonURL(isbn)
			ch1 <- true
		}()
	} else {
		go func() {
			ch1 <- true
		}()
	}
	data := searchapi.GetOpenBdData(isbn)
	if data.Title == "" || data.Image == "" {
		fmt.Println("OpenBD not serch data")
		<-ch1
		if data.Title == "" {
			data.Title = adata.Title
		}
		if data.Image == "" {
			data.Image = adata.Image
		}

	}
	if len(isbn) == 13 {
		if no > 0 {
			bookdatatable.Edit(strconv.Itoa(no), data.Isbn, data.Title, data.Writer, data.Brand, data.Ext, data.Synopsis, data.Image)

		} else {
			bookdatatable.Add(data.Isbn, data.Title, data.Writer, data.Brand, data.Ext, data.Synopsis, data.Image)
		}
	}
	tmp := bookdatatable.Serch_sql(data.Isbn)
	return tmp
}
func view(table string) string {
	var err error
	var tmp []byte
	if table == "" {
		tmp, err = json.Marshal(bookdatatable.Scansql())

	} else {
		tmp, err = json.Marshal(bookdatatable.ScansqlId(table))
	}
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(tmp)
}

func apiSql(url []string, r *http.Request) (string, error) {
	var output string
	if len(url) > 3 {
		switch url[2] {
		case "view":
			output = view(url[3])
		default:
			return "", errors.New("input url err")
		}
	} else {
		return "", errors.New("input url err")
	}
	return output, nil
}
func apiserver(w http.ResponseWriter, r *http.Request) {
	urldata := urlAnalysis(r.URL.Path)
	log.Printf("request:%v\n", r.URL.Path)
	if len(urldata) > 1 {
		switch urldata[1] {
		case "":
			w.WriteHeader(400)
			fmt.Fprintf(w, "Err API request")
		case "sql":
			start := time.Now()
			jsondata, err := apiSql(urldata, r)
			end := time.Now()
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintf(w, "Err API request")
			}
			fmt.Fprintf(w, "{\"Data\":%s,\"Time\":%f}", jsondata, (end.Sub(start)).Seconds())
		case "isbn":
			if len(urldata) > 2 {
				start := time.Now()
				data := sqlserch(urldata[2])
				if data == nil {
					data = newserchisbn(urldata[2], 0)
				} else {
					if len(data) == 1 {
						if data[0].Title == "" {
							data = newserchisbn(urldata[2], data[0].Id)
						}
					}
				}
				jsondata, err := json.Marshal(data)
				end := time.Now()
				if err == nil {
					fmt.Fprintf(w, "{\"Data\":%s,\"Time\":%f}", jsondata, (end.Sub(start)).Seconds())
				} else {
					w.WriteHeader(400)
					fmt.Fprintf(w, "Err API request")
				}
			} else {
				w.WriteHeader(400)
				fmt.Fprintf(w, "Err API request")
			}
		default:
			fmt.Fprintf(w, "%s", r.URL.Path)
		}
	} else {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Err API request")
	}
	return

}
