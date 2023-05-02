package loop

import (
	"book-newread/calendarscripting"
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var bookloop_ch chan bool

func BookLoop(ctx context.Context, ch chan<- error) {
	fmt.Println("BookLoop start")
bookloop:
	for {
		select {
		case <-ctx.Done():
			break bookloop
		case <-bookloop_ch:
			BookloopStart()
		}
	}
	close(bookloop_ch)
	fmt.Println("BookLoop end")
	ch <- nil
}

// 新刊をチェックするループリスト
func BookloopStart() {
	fmt.Println("start new book data count")
	dataStatusSet(BOOK_SELECT, "Reload")
	t := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var listdata BListData
			var listdata_mux sync.Mutex
			listdata.NowTime = time.Now()
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
				tmp = calendarscripting.FilterComicList(tmp, maxbackday)
				listdata_mux.Lock()
				listdata.LiteNobel = tmp
				listdata_mux.Unlock()
			}()
			go func() {
				defer wgg.Done()
				tmp := calendarscripting.GetComicList(listdata.Year, listdata.Month, calendarscripting.COMIC)
				tmp = calendarscripting.FilterComicList(tmp, maxbackday)
				listdata_mux.Lock()
				listdata.Comic = tmp
				listdata_mux.Unlock()

			}()
			wgg.Wait()
			addbookloop(i, listdata)
		}(i)
	}

	wg.Wait()
	endtime := time.Now()
	dataStatusSet(BOOK_SELECT, "OK")
	timeStatusSet(BOOK_SELECT, endtime)
	fmt.Println("read new book data end", (endtime.Sub(t)).Seconds(), "s")
}

func addbookloop(i int, list BListData) {
	listtemp_mux.Lock()
	listtemp.listdata[i] = list
	listtemp_mux.Unlock()
}

func ReadBookListData() []BListData {
	listtemp_mux.Lock()
	tmp := listtemp.listdata
	listtemp_mux.Unlock()
	return tmp

}
