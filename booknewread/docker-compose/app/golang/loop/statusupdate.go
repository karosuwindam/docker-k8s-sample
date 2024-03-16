package loop

import (
	"book-newread/loop/datastore"
	"fmt"
	"sync"
	"time"
)

var smux sync.Mutex //ステータスの更新用

type SELECT_STATUS int

const (
	NOBEL_SELECT = SELECT_STATUS(0) //web小説チェックの状態
	BOOK_SELECT  = SELECT_STATUS(1) //新刊チェックの状態
)

// ステータスの更新
func statusUpdate(st SELECT_STATUS, data string) {
	smux.Lock()
	defer smux.Unlock()
	s := datastore.Status{}
	if err := datastore.Read(&s); err != nil {
		fmt.Println(err)
		return
	}
	if st == NOBEL_SELECT {
		s.BookMarkStatus = data
		s.BookMarkNowTime = time.Now()
	} else if st == BOOK_SELECT {
		s.BookStatus = data
		s.BookNowTIme = time.Now()
	}
	if err := datastore.Write(s); err != nil {
		fmt.Println(err)
	}
}
