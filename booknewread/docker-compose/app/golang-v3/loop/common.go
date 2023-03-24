package loop

import (
	"sync"
	"time"
)

type RESET_FLAG int

const (
	RESET_NOBEL = RESET_FLAG(1)
	RESET_BOOK  = RESET_FLAG(1 << 1)
	RESET_DATA  = RESET_NOBEL | RESET_BOOK
)

type SELECT_STATUS int

const (
	NOBEL_SELECT = SELECT_STATUS(0)
	BOOK_SELECT  = SELECT_STATUS(1)
)

type reset struct {
	flag RESET_FLAG
	mu   sync.Mutex
}

type Status struct {
	BookNowTIme     time.Time `json:booknowtime`
	BookStatus      string    `json:bookstatus`
	BookMarkNowTime time.Time `json:bookmarknowtime`
	BookMarkStatus  string    `json:bookmarkstatus`
}

var Statusdata Status

var resetflag reset

func Reset_ON(flag RESET_FLAG) {
	resetflag.mu.Lock()
	resetflag.flag = resetflag.flag | flag
	resetflag.mu.Unlock()
}

func Reset_OFF(flag RESET_FLAG) {
	tmp := RESET_DATA - flag
	resetflag.mu.Lock()
	resetflag.flag = resetflag.flag & tmp
	resetflag.mu.Unlock()
}

func ResetRead(flag RESET_FLAG) bool {
	resetflag.mu.Lock()
	tmp := resetflag.flag
	resetflag.mu.Unlock()
	if flag&tmp > 0 {
		return true
	}
	return false

}

func dataStatusSet(num SELECT_STATUS, str string) {
	listtemp.mu.Lock()
	switch num {
	case BOOK_SELECT:
		listtemp.status.BookStatus = str
	case NOBEL_SELECT:
		listtemp.status.BookMarkStatus = str
	}
	listtemp.mu.Unlock()
}

func timeStatusSet(num SELECT_STATUS, tdata time.Time) {
	listtemp.mu.Lock()
	switch num {
	case BOOK_SELECT:
		listtemp.status.BookNowTIme = tdata
	case NOBEL_SELECT:
		listtemp.status.BookMarkNowTime = tdata
	}
	listtemp.mu.Unlock()
}
