package loop

import (
	"book-newread/config"
	"book-newread/loop/bookmarkfileread"
	"book-newread/loop/datastore"
	"book-newread/loop/novelchack"
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var shutdown chan bool
var loogflag bool

func Init() error {
	if err := novelchack.Init(); err != nil {
		return err
	}
	if err := datastore.Init(); err != nil {
		return err
	}
	if err := bookmarkfileread.Init(); err != nil {
		return err
	}
	shutdown = make(chan bool, 1)
	return nil
}

func Run(ctx context.Context) error {
	var wg sync.WaitGroup
	fmt.Println("start loop")
	wg.Add(2)
	go func() {
		defer wg.Done()
		readNarouData()
	}()
	go func() {
		defer wg.Done()
		readNewBookData()
	}()
	loogflag = true
	defer func() {
		loogflag = false
	}()
	wg.Wait()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-shutdown:
			break loop
		case <-time.After(time.Duration(config.Loop.LoopTIme) * time.Second):
			wg.Add(2)
			go func() { //新刊を取得する
				defer wg.Done()
				readNewBookData()
			}()
			go func() { //小説のデータを取得する
				defer wg.Done()
				readNarouData()
			}()
			wg.Wait()
		}
	}
	fmt.Println("end loop")
	return nil
}

func Stop() error {
	if loogflag {
		shutdown <- true
	}
	time.Sleep(1 * time.Second)
	return nil
}

// 新刊情報の取得
func readNewBookData() {

}

// Web小説のデータを取得する
func readNarouData() {
	fmt.Println("start novel data count")
	statusUpdate(NOBEL_SELECT, "Reload")
	now := time.Now()
	limit := 10
	slots := make(chan struct{}, limit)
	var wg sync.WaitGroup
	// bookmarkのファイル読み取り
	urls := bookmarkfileread.ReadBookmark()
	if len(urls) != 0 {
		datastore.ClearCount()
		wg.Add(len(urls))
		datastore.SetMaxCount(len(urls))
		for _, url := range urls {
			slots <- struct{}{}

			go func(url string) {
				defer func() {
					<-slots
					wg.Done()
				}()
				//urlによる解析処理

				if tmp, err := novelchack.ChackURLData(url); err == nil {
					if err = datastore.Write(tmp); err != nil {
						fmt.Println(err)
					}
				} else {
					fmt.Println(err)
				}
				datastore.AddCount()
				statusUpdate(NOBEL_SELECT, "Reload:"+strconv.Itoa(int(datastore.ReadPerCount())*100)+"%")
			}(url)
		}
		wg.Done()
	}
	endtime := time.Now()
	fmt.Println("read novel data end", (endtime.Sub(now)).Seconds(), "s")
	statusUpdate(NOBEL_SELECT, "ok")
}
