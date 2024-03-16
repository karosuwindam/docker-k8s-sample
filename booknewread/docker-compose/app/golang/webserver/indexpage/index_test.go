package indexpage_test

import (
	"book-newread/config"
	"book-newread/webserver/indexpage"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestIndex(t *testing.T) {
	os.Setenv("WEB_FOLDER", "./html")
	config.Init()
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexpage.Init("/"))

	srv := &http.Server{
		Addr:    config.Web.Hostname + ":" + config.Web.Port,
		Handler: mux,
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer srv.Shutdown(ctx)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Fatal(err)
		}
	}()
	if str := getIndex("index.html"); str != "" {
		fmt.Println(str)
	} else {
		t.Fatal("Not Found file")
	}
	if str := getIndex("js/test.js"); str != "404 Not Found" {
		fmt.Println(str)
	} else {
		t.Fatal("Not Found file")
	}
	if str := getIndex("css/test.css"); str != "404 Not Found" {
		fmt.Println(str)
	} else {
		t.Fatal("Not Found file")
	}
	if str := getIndex("test.css"); str != "404 Not Found" {
		t.Fatal("Not 404 Code")
	}
	cancel()

}

func getIndex(url string) string {
	res, err := http.Get("http://localhost:8080/" + url)
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	return string(body)
}
