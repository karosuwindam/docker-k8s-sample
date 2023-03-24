package loop

import "sync"

type RESET_FLAG int

const (
	RESET_NOBEL = RESET_FLAG(1)
	RESET_BOOK  = RESET_FLAG(1 << 1)
	RESET_DATA  = RESET_NOBEL | RESET_BOOK
)

type reset struct {
	flag RESET_FLAG
	mu   sync.Mutex
}

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
