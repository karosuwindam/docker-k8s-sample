package main

import (
	"booknewread/dirread"
	"booknewread/loop"
	"booknewread/novel_chack"
	"booknewread/webserver"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	LOOP_TIME_SEC = 60 * 60
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello World")
}

func RootConfg() []webserver.WebConfig {
	output := []webserver.WebConfig{}
	tmp := webserver.WebConfig{Pass: "/", Handler: hello}
	output = append(output, tmp)
	return output
}

func Config(cfg *webserver.SetupServer) error {
	webserver.Config(cfg, RootConfg())
	return nil
}

func Run(ctx context.Context) error {
	cfg, err := webserver.NewSetup()
	if err != nil {
		return err
	}
	if err := Config(cfg); err != nil {
		return err
	}
	s, err := cfg.NewServer()
	if err != nil {
		return err
	}

	return s.Run(ctx)
}
func EndRun() {}

func bookmarkFalderRead() (dirread.Dirtype, error) {
	var fpass dirread.Dirtype

	fpass.Setup("bookmark")
	fpass.Read("/")
	if len(fpass.Data) == 0 {
		return fpass, errors.New("err bookmark not file")
	}
	return fpass, nil
}

func bookmarkread(fpass *dirread.Dirtype) []string {
	data := []novel_chack.BookBark{}
	if len(fpass.Data) != 0 {
		data = []novel_chack.BookBark{}
		for _, fd := range fpass.Data {
			pass := fd.RootPath + fd.Name
			ncr := novel_chack.ReadBookBark(pass)
			for _, tmpd := range ncr {
				data = append(data, tmpd)
			}
			data = novel_chack.BookBarkSout(data)
		}
	}
	output := []string{}
	for _, bookmark := range data {
		output = append(output, bookmark.Url)
	}
	return output

}

func main() {
	fpass, err := bookmarkFalderRead()
	if err != nil {
		fmt.Println(err)
		return
	}
	bookmarklists := bookmarkread(&fpass)

	go func(list []string) { //ループによるURLチェックスタート
		fmt.Println("start novel data count")
		for {
			starttime := time.Now()
			loop.NobelLoop(list)
			for {
				if time.Now().Sub(starttime).Seconds() > LOOP_TIME_SEC {
					break
				} else {
					time.Sleep(time.Microsecond * 100)
				}
			}
			fmt.Println("reload novel data")
		}
	}(bookmarklists)
	go func() {
		fmt.Println("start new book data count")
		starttime := time.Now()
		for {
			if time.Now().Sub(starttime).Seconds() > LOOP_TIME_SEC {
				break
			} else {
				time.Sleep(time.Microsecond * 100)
			}
		}
		fmt.Println("reload new book data")

	}()

	//動作確認--start
	time.Sleep(time.Second * 60)
	loop.Read()
	fmt.Println(len(bookmarklists), loop.Count())
	for _, listdata := range loop.ListData {
		fmt.Println(listdata)
	}
	time.Sleep(time.Second * 60)
	//動作確認--end
	return
	fmt.Println("start")
	ctx := context.Background()
	if err := Run(ctx); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	EndRun()
	fmt.Println("end")
}
