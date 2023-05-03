package loop

import (
	"book-newread/calendarscripting"
	"book-newread/novelchack"
	"context"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

type BListData struct {
	NowTime   time.Time                    `json:Nowtime`
	Month     string                       `json:Month`
	Year      string                       `json:Year`
	Comic     []calendarscripting.BookList `json:Comic`
	LiteNobel []calendarscripting.BookList `json:LiteNobel`
}

var booklistdata []BListData
var booklistdata_mux sync.Mutex

var maxaccess int  //一つのサイトに対してアクセスできる限界数
var maxbackday int //今日を基準に表示する前日数

// 対象のURLから小説の更新確認
func (t *Listdata) chackurl(url string) novelchack.List {
	if t, err := novelchack.ChackUrlType(url); err != nil {
		return novelchack.List{}
	} else {
		if tmp, err := novelchack.ChackUrlData(t, url); err != nil {
			return tmp
		} else {
			return tmp
		}
	}
}

var nobelloop_ch chan []string

func NobelLoop(ctx context.Context, ch chan<- error) {
	fmt.Println("Nobel loop start")
nobelloop:
	for {
		select {
		case <-ctx.Done():
			break nobelloop
		case list := <-nobelloop_ch:
			nobelLoopStart(list)
		}
	}
	close(nobelloop_ch)
	fmt.Println("Nobel loop end")
	ch <- nil
}

func nobelLoopStart(urlLists []string) {
	fmt.Println("start novel data count")
	dataStatusSet(NOBEL_SELECT, "Reload")
	now := time.Now()
	listtemp.count = 1
	novelchack.Setup()
	var wg sync.WaitGroup
	for i, url := range urlLists {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			tmp := listtemp.chackurl(url)
			addNovelChack(tmp)
			dataStatusSet(NOBEL_SELECT, "Reload:"+strconv.Itoa(ReadListCount()*100/len(urlLists))+"%")
		}(i, url)
	}
	wg.Wait()
	endtime := time.Now()
	timeStatusSet(NOBEL_SELECT, endtime)
	dataStatusSet(NOBEL_SELECT, "OK")
	fmt.Println("read novel data end", (endtime.Sub(now)).Seconds(), "s")
}

func addNovelChack(data novelchack.List) {
	listtemp_mux.Lock()
	listtemp.count++
	listtemp_mux.Unlock()
	if data.Url == "" {
		return
	}
	flag := true
	listtemp_mux.Lock()
	for i, tmp := range listtemp.data {
		if data.Title == "" {
			flag = false
			break
		}
		if tmp.Url == data.Url {
			listtemp.data[i] = data
			flag = false
			break
		}
	}
	if flag {
		listtemp.data = append(listtemp.data, data)
	}
	listtemp_mux.Unlock()

}

func ReadListCount() int {
	listtemp_mux.Lock()
	count := listtemp.count
	listtemp_mux.Unlock()
	return count
}

func ReadNListData() []novelchack.List {
	listtemp_mux.Lock()
	tmp := listtemp.data
	listtemp_mux.Unlock()
	sort.Slice(tmp, func(i, j int) bool { return tmp[i].Lastdate.Unix() > tmp[j].Lastdate.Unix() })
	return tmp
}
