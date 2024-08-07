package loop

import (
	"book-newread/config"
	"book-newread/loop/bookmarkfileread"
	"book-newread/loop/calendarscripting"
	"book-newread/loop/datastore"
	"book-newread/loop/novelchack"
	"context"
	"errors"
	"log"
	"math"
	"strconv"
	"sync"
	"time"
)

var shutdown chan bool
var done chan bool
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
	done = make(chan bool, 1)
	return nil
}

func Run(ctx context.Context) error {
	var wg sync.WaitGroup
	log.Println("info:", "start loop")
	resetflag = true
	loogflag = true
	wg.Add(1)
	go func() {
		wg.Done()
		readNarouData(context.Background())
	}()
	readNewBookData(context.Background())
	wg.Wait()
	resetflag = false

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-shutdown:
			done <- true
			break loop
		case <-resetCH:
			resetflag = true
			wg.Add(2)
			go func() { //新刊を取得する
				defer wg.Done()
				readNewBookData(context.Background())
			}()
			go func() { //小説のデータを取得する
				defer wg.Done()
				readNarouData(context.Background())
			}()
			wg.Wait()
			resetflag = false
		case <-time.After(time.Duration(config.Loop.LoopTIme) * time.Second):
			resetflag = true
			wg.Add(2)
			go func() { //新刊を取得する
				defer wg.Done()
				readNewBookData(context.Background())
			}()
			go func() { //小説のデータを取得する
				defer wg.Done()
				readNarouData(context.Background())
			}()
			wg.Wait()
			resetflag = false
		}
	}
	log.Println("info:", "end loop")
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

func Stop(ctx context.Context) error {
	if loogflag {
		loogflag = false
		shutdown <- true
	}
	select {
	case <-done:
		break
	case <-ctx.Done():
		return errors.New("contex done")
	case <-time.After(time.Second):
		return errors.New("time over 1 sec")
	}
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
func readNewBookData(ctx context.Context) {
	log.Println("info:", "start new book data count")
	ctx, traceSpan := config.TracerS(ctx, "readNewBookData", "loop NewBookData")
	defer traceSpan.End()
	statusUpdate(BOOK_SELECT, "Reload")
	t := time.Now()
	var wg sync.WaitGroup
	var output []datastore.BListData = make([]datastore.BListData, 3)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int, ctx context.Context) {
			ctx, spango := config.TracerS(ctx, "readNewBookData_min", "count "+strconv.Itoa(i))
			defer spango.End()
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
			go func(ctx context.Context) {
				ctx, spango := config.TracerS(ctx, "readNewBookData_Nobel", "count "+strconv.Itoa(i))
				defer spango.End()
				defer wgg.Done()
				tmp := calendarscripting.GetComicList(listdata.Year, listdata.Month, calendarscripting.LITENOVEL, ctx)
				tmp = calendarscripting.FilterComicList(tmp)
				listdata_mux.Lock()
				listdata.LiteNobel = tmp
				listdata_mux.Unlock()
			}(ctx)
			go func(ctx context.Context) {
				ctx, spango := config.TracerS(ctx, "readNewBookData_Comic", "count "+strconv.Itoa(i))
				defer spango.End()
				defer wgg.Done()
				tmp := calendarscripting.GetComicList(listdata.Year, listdata.Month, calendarscripting.COMIC, ctx)
				tmp = calendarscripting.FilterComicList(tmp)
				listdata_mux.Lock()
				listdata.Comic = tmp
				listdata_mux.Unlock()
			}(ctx)
			wgg.Wait()
			output[i] = listdata
		}(i, ctx)
	}
	wg.Wait()
	if err := datastore.Write(output); err != nil {
		log.Println("error:", err)
	}

	endtime := time.Now()
	log.Println("info:", "read new book data end", "spantime", (endtime.Sub(t)).Seconds(), "s")
	statusUpdate(BOOK_SELECT, "ok")

}

// Web小説のデータを取得する
func readNarouData(ctx context.Context) {
	log.Println("info:", "start novel data count")
	ctx, traceSpan := config.TracerS(ctx, "readNarouData", "loop Nobel")
	defer traceSpan.End()
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
				ctx, span := config.TracerS(ctx, "readNarouData_min", url)
				defer func() {
					span.End()
					<-slots
					wg.Done()
				}()
				if !loogflag { //ループ終了は何もしない
					return
				}
				//urlによる解析処理
				if tmp, err := novelchack.ChackURLData(ctx, url); err == nil {
					if err = datastore.Write(tmp); err != nil {
						log.Println("error:", err)
					}
				} else if err != novelchack.ErrAtherUrl {
					log.Println("error:", err)
				}
				datastore.AddCount()
				per := datastore.ReadPerCount()
				per = math.Floor(per*1000) / 1000
				// float64をstringに変換
				pers := strconv.FormatFloat(per*100, 'f', -1, 64)
				statusUpdate(NOBEL_SELECT, "Reload:"+pers+"%")
			}(url)
		}
		wg.Wait()
	}
	endtime := time.Now()
	log.Println("info:", "read novel data end", "spantime", (endtime.Sub(now)).Seconds(), "s")
	statusUpdate(NOBEL_SELECT, "ok")
}
