package novel_chack

import (
	"fmt"
	"testing"
)

func TestChackKuyoumu(t *testing.T) {
	Setup()
	url := "https://kakuyomu.jp/works/4852201425154978969"
	if tmp, err := ChackUrldata(url); err != nil {
		t.Error(err)
	} else {
		fmt.Println(tmp)
	}

	url = "https://ncode.syosetu.com/n9350eg/?p=6"
	if tmp, err := ChackUrldata(url); err != nil {
		t.Error(err)
	} else {
		fmt.Println(tmp)
	}
	url = "https://novel18.syosetu.com/n6719in/"
	if tmp, err := ChackUrldata(url); err != nil {
		t.Error(err)
	} else {
		fmt.Println(tmp)
	}

}
