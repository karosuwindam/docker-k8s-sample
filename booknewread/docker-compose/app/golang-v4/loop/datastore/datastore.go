package datastore

import (
	"book-newread/loop/calendarscripting"
	"book-newread/loop/novelchack"
	"sync"
	"time"
)

// 新刊チェック用の保存型
type BListData struct {
	NowTime   time.Time                    `json:Nowtime`
	Month     string                       `json:Month`
	Year      string                       `json:Year`
	Comic     []calendarscripting.BookList `json:Comic`
	LiteNobel []calendarscripting.BookList `json:LiteNobel`
}

// 基本データストア
type Listdata struct {
	data    []novelchack.List
	newData []BListData
	status  Status
	count   int
	mu      sync.Mutex
}

type Status struct { //状態確認について
	BookNowTIme     time.Time `json:booknowtime`
	BookStatus      string    `json:bookstatus`
	BookMarkNowTime time.Time `json:bookmarknowtime`
	BookMarkStatus  string    `json:bookmarkstatus`
}

var listtemp Listdata

// 初期化作業
func Init() error {
	listtemp.data = []novelchack.List{}
	listtemp.newData = []BListData{}
	listtemp.status = Status{}
	return nil
}

// データストアに追加作業
func Write(v interface{}) error {
	listtemp.mu.Lock()
	defer listtemp.mu.Unlock()
	switch v.(type) {
	case Status:
		listtemp.status = v.(Status)
	case []novelchack.List:
		listtemp.data = v.([]novelchack.List)
	case []BListData:
		listtemp.newData = v.([]BListData)
	case novelchack.List:
		addNovelChack(v.(novelchack.List))
		listtemp.count++
	}
	return nil
}

func addNovelChack(nl novelchack.List) {
	if nl.Url != "" || nl.Title != "" {
		flag := true
		for i, tmp := range listtemp.data {
			if tmp.Url == nl.Url {
				listtemp.data[i] = nl
				flag = false
				break
			}
		}
		if flag {
			listtemp.data = append(listtemp.data, nl)
		}

	}
}

func Read(v interface{}) {
	listtemp.mu.Lock()
	defer listtemp.mu.Unlock()
	switch v.(type) {
	case *Listdata:
		tmp := listtemp.newData
		v = &tmp
	case *novelchack.List:
		tmp := listtemp.data
		v = &tmp
	case *Status:
		tmp := listtemp.status
		v = &tmp

	}
}

func ReadCount() int {
	listtemp.mu.Lock()
	defer listtemp.mu.Unlock()
	return listtemp.count
}
