package novelchack_test

import (
	"book-newread/novelchack"
	"fmt"
	"testing"
)

func TestChackUrlData(t *testing.T) {
	novelchack.Setup()
	url := "https://kakuyomu.jp/works/4852201425154978969"
	// url := "https://kakuyomu.jp/works/16816452218254294002"
	if tmp, err := novelchack.ChackUrlData(novelchack.KAKUYOMU_WEB, url); err != nil {
		t.Error(err)
	} else {
		fmt.Println(tmp)
	}

}
