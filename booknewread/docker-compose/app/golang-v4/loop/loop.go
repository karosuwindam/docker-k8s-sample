package loop

import (
	"book-newread/config"
	"book-newread/loop/datastore"
	"book-newread/loop/novelchack"
	"context"
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
	shutdown = make(chan bool)
	return nil
}

func Run(ctx context.Context) error {

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-shutdown:
			break loop
		case <-time.After(time.Duration(config.Loop.LoopTIme) * time.Second):
		}
	}

	return nil
}

func Stop() error {
	shutdown <- true
	return nil
}
