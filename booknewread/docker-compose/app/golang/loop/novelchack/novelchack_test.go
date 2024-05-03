package novelchack_test

import (
	"book-newread/config"
	"book-newread/loop/novelchack"
	"fmt"
	"testing"
)

func TestNovelChack(t *testing.T) {
	config.Init()
	novelchack.Init()

	url := "https://kakuyomu.jp/works/4852201425154978969"
	// url := "https://kakuyomu.jp/works/16816452218254294002"
	if tmp, err := novelchack.ChackURLData(url); err != nil {
		t.Error(err)
	} else {
		fmt.Println(tmp)
	}

	url = "https://kakuyomu.jp/works/16817330662159451369"
	if tmp, err := novelchack.ChackURLData(url); err != nil {
		t.Error(err)
	} else {
		fmt.Println(tmp)
	}

	url = "https://ncode.syosetu.com/n1976ey/"
	// url := "https://kakuyomu.jp/works/16816452218254294002"
	if tmp, err := novelchack.ChackURLData(url); err != nil {
		t.Error(err)
	} else {
		fmt.Println(tmp)
	}

	url = "https://novel18.syosetu.com/n6719in/"
	if tmp, err := novelchack.ChackURLData(url); err != nil {
		t.Error(err)
	} else {
		fmt.Println(tmp)
	}

}
