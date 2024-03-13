package loop

import (
	"book-newread/config"
	"book-newread/loop/bookmarkfileread"
	"book-newread/loop/datastore"
	"book-newread/loop/novelchack"
	"context"
	"fmt"
	"sync"
	"time"
)

var shutdown chan bool

func Init() error {
	if err := novelchack.Init(); err != nil {
		return err
	}
	if err := datastore.Init(); err != nil {
		return err
	}
	if err := bookmarkfileread.Init(); err != nil {
		return err
	}
	shutdown = make(chan bool, 1)
	return nil
}

func Run(ctx context.Context) error {
	var wg sync.WaitGroup
	readNarouData()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-shutdown:
			break loop
		case <-time.After(time.Duration(config.Loop.LoopTIme) * time.Second):
			wg.Add(2)
			go func() {
				defer wg.Done()
			}()
			go func() {
				defer wg.Done()
				readNarouData()
			}()
			wg.Wait()
		}
	}

	return nil
}

func Stop() error {
	shutdown <- true
	return nil
}

func readNarouData() {
	limit := 10
	slots := make(chan struct{}, limit)
	var wg sync.WaitGroup
	// bookmarkのファイル読み取り
	urls := bookmarkfileread.ReadBookmark()
	if len(urls) != 0 {
		wg.Add(len(urls))
		for _, url := range urls {
			slots <- struct{}{}

			go func(url string) {
				defer func() {
					<-slots
					wg.Done()
				}()
				//urlによる解析処理

				if tmp, err := novelchack.ChackURLData(url); err == nil {
					fmt.Println(tmp)
				}

			}(url)
		}
		wg.Done()
	}

}
