package loop

import (
	"booknewread/calendarscripting"
	"fmt"
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

var BookListData []BListData

// 環境変数で操作できる数値周
type SetupEnvData struct {
	MaxAccess  int //一つのサイトに対してアクセスできる限界数
	MaxBackDay int //今日を基準に表示する前日数
}

// デフォルト
const (
	DEF_ACCESS  = 2
	DEF_BACKDAY = 5
)

var EnvData SetupEnvData

func SetupEnv() SetupEnvData {
	var output SetupEnvData
	if true {
		output.MaxAccess = 3
	} else {
		output.MaxAccess = DEF_ACCESS
	}
	if true {
		output.MaxBackDay = 2
	} else {
		output.MaxBackDay = DEF_BACKDAY
	}
	return output
}

var bookloopmu sync.Mutex

func Bookloop() {
	fmt.Println("start new book data count")
	dataStatusSet(BOOK_SELECT, "Reload")
	t := time.Now()
	ch := make(chan bool, 3)
	for i := 0; i < 3; i++ {
		go func(i int) {
			var listdata BListData
			ch1 := make(chan bool)
			ch2 := make(chan bool)
			listdata.NowTime = time.Now()
			if (int(t.Month()) + i) <= 12 {
				listdata.Month = strconv.Itoa((int(t.Month()) + i))
				listdata.Year = strconv.Itoa(int(t.Year()))
			} else {
				listdata.Month = strconv.Itoa((int(t.Month()) + i) % 12)
				listdata.Year = strconv.Itoa(int(t.Year()) + 1)
			}
			go func() {
				bookloopmu.Lock()
				listdata.LiteNobel = calendarscripting.GetComicList(listdata.Year, listdata.Month, calendarscripting.LITENOVEL)
				bookloopmu.Unlock()
				listdata.LiteNobel = calendarscripting.FilterComicList(listdata.LiteNobel, EnvData.MaxBackDay)
				ch1 <- true
			}()
			go func() {
				bookloopmu.Lock()
				listdata.Comic = calendarscripting.GetComicList(listdata.Year, listdata.Month, calendarscripting.COMIC)
				bookloopmu.Unlock()
				listdata.Comic = calendarscripting.FilterComicList(listdata.Comic, EnvData.MaxBackDay)
				ch2 <- true
			}()
			<-ch1
			<-ch2
			addbookloop(i, listdata)
			ch <- true

		}(i)
	}
	for i := 0; i < 3; i++ {
		<-ch
	}
	endtime := time.Now()
	dataStatusSet(BOOK_SELECT, "OK")
	timeStatusSet(BOOK_SELECT, endtime)
	fmt.Println("read new book data end", (endtime.Sub(t)).Seconds(), "s")

}

func addbookloop(i int, list BListData) {
	listtemp.mu.Lock()
	listtemp.listdata[i] = list
	listtemp.mu.Unlock()
}
