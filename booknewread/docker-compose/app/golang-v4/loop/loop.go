package loop

import (
	"book-newread/config"
	"book-newread/loop/bookmarkfileread"
	"book-newread/loop/calendarscripting"
	"book-newread/loop/datastore"
	"book-newread/loop/novelchack"
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var shutdown chan bool
var resetflag bool
var resetCH chan bool
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
	resetCH = make(chan bool, 1)
	return nil
}

func Run(ctx context.Context) error {
	var wg sync.WaitGroup
	fmt.Println("start loop")
	resetflag = true
	wg.Add(2)
	go func() {
		defer wg.Done()
		readNarouData()
	}()
	go func() {
		defer wg.Done()
		readNewBookData()
		loogflag = true
	}()
	defer func() {
		loogflag = false
	}()
	wg.Wait()
	resetflag = false

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-shutdown:
			break loop
		case <-resetCH:
			resetflag = true
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
			resetflag = false
		case <-time.After(time.Duration(config.Loop.LoopTIme) * time.Second):
			resetflag = true
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
			resetflag = false
		}
	}
	fmt.Println("end loop")
	return nil
}

// ループがスタートするまでの確認待ち
func RunWait() error {
	if loogflag {
		return nil
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
loop:
	for {
		select {
		case <-ctx.Done():
			return errors.New("10 sec over")
		case <-time.After(500 * time.Microsecond):
			if loogflag {
				break loop

			}
		}
	}
	return nil
}

func Stop() error {
	if loogflag {
		shutdown <- true
	}
	time.Sleep(1 * time.Second)
	return nil
}

// ループを強制的に実行
func Reset() {
	if resetflag {
		return
	}
	if len(resetCH) == 0 {
		resetCH <- true
	}
}

// 新刊情報の取得
func readNewBookData() {
	fmt.Println("start new book data count")
	statusUpdate(BOOK_SELECT, "Reload")
	t := time.Now()
	var wg sync.WaitGroup
	var output []datastore.BListData = make([]datastore.BListData, 3)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var listdata datastore.BListData
			var listdata_mux sync.Mutex
			if (int(t.Month()) + i) <= 12 {
				listdata.Month = strconv.Itoa((int(t.Month()) + i))
				listdata.Year = strconv.Itoa(int(t.Year()))
			} else {
				listdata.Month = strconv.Itoa((int(t.Month()) + i) % 12)
				listdata.Year = strconv.Itoa(int(t.Year()) + 1)
			}
			var wgg sync.WaitGroup
			wgg.Add(2)
			go func() {
				defer wgg.Done()
				tmp := calendarscripting.GetComicList(listdata.Year, listdata.Month, calendarscripting.LITENOVEL)
				tmp = calendarscripting.FilterComicList(tmp)
				listdata_mux.Lock()
				listdata.LiteNobel = tmp
				listdata_mux.Unlock()
			}()
			go func() {
				defer wgg.Done()
				tmp := calendarscripting.GetComicList(listdata.Year, listdata.Month, calendarscripting.COMIC)
				tmp = calendarscripting.FilterComicList(tmp)
				listdata_mux.Lock()
				listdata.Comic = tmp
				listdata_mux.Unlock()
			}()
			wgg.Wait()
			output[i] = listdata
		}(i)
	}
	wg.Wait()
	if err := datastore.Write(output); err != nil {
		fmt.Println(err)
	}

	endtime := time.Now()
	fmt.Println("read new book data end", (endtime.Sub(t)).Seconds(), "s")
	statusUpdate(BOOK_SELECT, "ok")

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
		datastore.SetMaxCount(len(urls))
		wg.Add(len(urls))
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
				} else if err != novelchack.ErrAtherUrl {
					fmt.Println(err)
				}
				datastore.AddCount()
				per := datastore.ReadPerCount()
				// float64をstringに変換
				pers := strconv.FormatFloat(per*100, 'f', -1, 64)
				statusUpdate(NOBEL_SELECT, "Reload:"+pers+"%")
			}(url)
		}
		wg.Wait()
	}
	endtime := time.Now()
	fmt.Println("read novel data end", (endtime.Sub(now)).Seconds(), "s")
	statusUpdate(NOBEL_SELECT, "ok")
}
