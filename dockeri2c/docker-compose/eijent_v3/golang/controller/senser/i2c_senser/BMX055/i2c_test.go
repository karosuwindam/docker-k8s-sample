package bmx055

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRead(t *testing.T) {
	var i2cMu sync.Mutex
	Init(&i2cMu)
	Test()
	accInit()
	for i := 0; i < 20; i++ {
		fmt.Println(i)
		fmt.Println(getACCRAW())
		fmt.Println(getGyroRAW())
		fmt.Println(getMag())
		time.Sleep(time.Millisecond * 10)
	}
}
