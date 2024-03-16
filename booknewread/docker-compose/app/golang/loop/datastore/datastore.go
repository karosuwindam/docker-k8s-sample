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
	data     []novelchack.List
	newData  []BListData
	status   Status
	count    int //URLの処理結果
	maxcount int //URLの最大処理
	mu       sync.Mutex
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

func Read(v interface{}) error {
	listtemp.mu.Lock()
	defer listtemp.mu.Unlock()
	switch v.(type) {
	case *Status:
		*v.(*Status) = listtemp.status
	case *[]novelchack.List:
		*v.(*[]novelchack.List) = listtemp.data
	case *[]BListData:
		*v.(*[]BListData) = listtemp.newData
	}
	return nil
}

func ReadCount() int {
	listtemp.mu.Lock()
	defer listtemp.mu.Unlock()
	return listtemp.count
}

func ReadPerCount() float64 {
	listtemp.mu.Lock()
	defer listtemp.mu.Unlock()
	c := float64(listtemp.count)
	m := float64(listtemp.maxcount)
	return c / m

}

func ClearCount() {
	listtemp.mu.Lock()
	defer listtemp.mu.Unlock()
	listtemp.count = 0
}

func AddCount() {
	listtemp.mu.Lock()
	defer listtemp.mu.Unlock()
	listtemp.count++
}

func SetMaxCount(i int) {
	listtemp.mu.Lock()
	defer listtemp.mu.Unlock()
	listtemp.maxcount = i
}
