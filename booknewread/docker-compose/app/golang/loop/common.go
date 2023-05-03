package loop

import (
	"book-newread/config"
	"book-newread/loop/dirread"
	"book-newread/novel_chack"
	"book-newread/novelchack"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type RESET_FLAG int

const (
	RESET_NOBEL = RESET_FLAG(1)            //なろうチェック理再スタート
	RESET_BOOK  = RESET_FLAG(1 << 1)       //新刊の再チェックスタート
	RESET_DATA  = RESET_NOBEL | RESET_BOOK //両方再チェックスタート
)

type SELECT_STATUS int

const (
	NOBEL_SELECT = SELECT_STATUS(0) //web小説チェックの状態
	BOOK_SELECT  = SELECT_STATUS(1) //新刊チェックの状態
)

type reset struct {
	flag RESET_FLAG
}

var resetflag reset
var resetflag_mux sync.Mutex

// リセットフラグのON
func Reset_ON(flag RESET_FLAG) {
	resetflag_mux.Lock()
	resetflag.flag = resetflag.flag | flag
	resetflag_mux.Unlock()
}

// リセットフラグのOFF
func Reset_OFF(flag RESET_FLAG) {
	tmp := RESET_DATA - flag
	resetflag_mux.Lock()
	resetflag.flag = resetflag.flag & tmp
	resetflag_mux.Unlock()
}

// リセット状態の読み取り
func ResetRead(flag RESET_FLAG) bool {
	resetflag_mux.Lock()
	tmp := resetflag.flag
	resetflag_mux.Unlock()
	if flag&tmp > 0 {
		return true
	}
	return false

}

type Listdata struct { //出力について
	data     []novelchack.List
	listdata []BListData
	status   Status
	count    int
}

type Status struct { //状態確認について
	BookNowTIme     time.Time `json:booknowtime`
	BookStatus      string    `json:bookstatus`
	BookMarkNowTime time.Time `json:bookmarknowtime`
	BookMarkStatus  string    `json:bookmarkstatus`
}

var looptime int          //待ち時間
var bookmarkfalder string //ブックマークの保存フォルダ

var listtemp Listdata
var listtemp_mux sync.Mutex

// ステータスの更新設定
func dataStatusSet(num SELECT_STATUS, str string) {
	listtemp_mux.Lock()
	switch num {
	case BOOK_SELECT:
		listtemp.status.BookStatus = str
	case NOBEL_SELECT:
		listtemp.status.BookMarkStatus = str
	}
	listtemp_mux.Unlock()
}

// ステータスの時間更新
func timeStatusSet(num SELECT_STATUS, tdata time.Time) {
	listtemp_mux.Lock()
	switch num {
	case BOOK_SELECT:
		listtemp.status.BookNowTIme = tdata
	case NOBEL_SELECT:
		listtemp.status.BookMarkNowTime = tdata
	}
	listtemp_mux.Unlock()
}

// スタート時間の測定
func timeStatusRead(num SELECT_STATUS) time.Time {
	var output time.Time
	listtemp_mux.Lock()
	switch num {
	case BOOK_SELECT:
		output = listtemp.status.BookNowTIme
	case NOBEL_SELECT:
		output = listtemp.status.BookMarkNowTime
	}
	listtemp_mux.Unlock()
	return output

}

// 特定フォルダからブックマークのURLを取り出し
func bookmarkread(list []string, fpass *dirread.Dirtype) []string {
	data := []novel_chack.BookBark{}
	for _, url := range list {
		data = append(data, novel_chack.ReadBookBark(url)...)
	}
	if len(fpass.Data) != 0 {
		for _, fd := range fpass.Data {
			fmt.Println("read file", fd.Name)
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

var bookmarklists []string

func CommonLoop(ctx context.Context, ch chan<- error) {
	fmt.Println("common loop start")
commonloop:
	for {
		select {
		case <-ctx.Done():
			break commonloop
		case <-time.After(time.Microsecond * 100):
			if (time.Now().Sub(timeStatusRead(BOOK_SELECT)).Seconds() > float64(looptime)) || ResetRead(RESET_BOOK) {
				if len(bookloop_ch) == 0 {
					Reset_OFF(RESET_BOOK)
					bookloop_ch <- true
				}
			}
			if (time.Now().Sub(timeStatusRead(NOBEL_SELECT)).Seconds() > float64(looptime)) || ResetRead(RESET_NOBEL) {
				if len(nobelloop_ch) == 0 {

					if f, err := bookmarkFalderRead(); err == nil {
						bookmarklists = bookmarkread(bookmarklists, &f)

					}
					Reset_OFF(RESET_NOBEL)
					nobelloop_ch <- bookmarklists
				}
			}
		}
	}
	ch <- nil
	fmt.Println("common loop end")
}

func bookmarkFalderRead() (dirread.Dirtype, error) {
	var fpass dirread.Dirtype

	fpass.Setup(bookmarkfalder)
	fpass.Read("/")
	if len(fpass.Data) == 0 {
		return fpass, errors.New("err bookmark not file")
	}
	return fpass, nil
}

// ループ機能を有効にする設定
func Setup(cfg *config.Config) error {
	maxaccess = cfg.Loop.MaxAccess
	maxbackday = cfg.Loop.MaxBackDay
	looptime = cfg.Loop.LoopTIme
	bookmarkfalder = cfg.Loop.BookmarkF

	booklistdata = make([]BListData, 3)
	listtemp.data = []novelchack.List{}
	listtemp.listdata = make([]BListData, 3)
	listtemp.status = Status{
		BookNowTIme:     time.Now(),
		BookStatus:      "Reload",
		BookMarkNowTime: time.Now(),
		BookMarkStatus:  "Reload",
	}
	Reset_ON(RESET_DATA)
	bookloop_ch = make(chan bool, 1)
	nobelloop_ch = make(chan []string, 1)

	if f, err := bookmarkFalderRead(); err != nil {
		return err
	} else {
		bookmarklists = bookmarkread([]string{}, &f)

	}
	return nil
}

func ReadStatus() Status {
	listtemp_mux.Lock()
	tmp := listtemp.status
	listtemp_mux.Unlock()
	return tmp
}
