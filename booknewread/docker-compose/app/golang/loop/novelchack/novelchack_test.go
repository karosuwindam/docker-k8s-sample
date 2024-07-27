package novelchack_test

import (
	"book-newread/config"
	"book-newread/loop/novelchack"
	"context"
	"fmt"
	"testing"
)

func TestNovelChack(t *testing.T) {
	config.Init()
	novelchack.Init()
	ctx := context.TODO()

	url := "https://kakuyomu.jp/works/1177354054883819762"
	if tmp, err := novelchack.ChackURLData(ctx, url); err != nil {
		t.Error(err)
	} else {
		fmt.Println(tmp)
	}

	url = "https://kakuyomu.jp/works/16817330662159451369"
	if tmp, err := novelchack.ChackURLData(ctx, url); err != nil {
		t.Error(err)
	} else {
		fmt.Println(tmp)
	}

	url = "https://ncode.syosetu.com/n1976ey/"
	// url := "https://kakuyomu.jp/works/16816452218254294002"
	if tmp, err := novelchack.ChackURLData(ctx, url); err != nil {
		t.Error(err)
	} else {
		fmt.Println(tmp)
	}

	url = "https://novel18.syosetu.com/n6719in/"
	if tmp, err := novelchack.ChackURLData(ctx, url); err != nil {
		t.Error(err)
	} else {
		fmt.Println(tmp)
	}

}
