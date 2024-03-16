package datastore

import (
	"book-newread/loop/novelchack"
	"sync"
	"testing"
	"time"
)

// データストアのテスト
func TestDataStore(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	tnow := time.Now()
	tmpStatus := Status{
		BookNowTIme:     tnow,
		BookStatus:      "OK",
		BookMarkNowTime: tnow,
		BookMarkStatus:  "OK",
	}
	tmpBlistdata := []BListData{
		{
			NowTime:   tnow,
			Month:     "1",
			Year:      "2020",
			Comic:     nil,
			LiteNobel: nil,
		},
		{
			NowTime:   tnow,
			Month:     "2",
			Year:      "2020",
			Comic:     nil,
			LiteNobel: nil,
		},
	}
	tmpListdata := novelchack.List{
		Title:      "test",
		Url:        "http://test.com",
		LastStoryT: "test",
	}
	if err := Write(tmpStatus); err != nil {
		t.Fatal(err)
	}
	tStatus := Status{}
	if err := Read(&tStatus); err != nil {
		t.Fatal(err)
	}
	if tStatus.BookStatus != "OK" {
		t.Fatal("Status Error")
	}
	if tStatus.BookMarkStatus != "OK" {
		t.Fatal("Status Error")
	}
	if tStatus.BookNowTIme != tnow {
		t.Fatal("Status Error")
	}
	if tStatus.BookMarkNowTime != tnow {
		t.Fatal("Status Error")
	}
	if err := Write(tmpBlistdata); err != nil {
		t.Fatal(err)
	}
	tBlistdata := []BListData{}
	if err := Read(&tBlistdata); err != nil {
		t.Fatal(err)
	}
	if len(tBlistdata) != 2 {
		t.Fatal("Blistdata Error")
	}
	tmpListdata.LastStoryT = "test2"
	if err := Write(tmpListdata); err != nil {
		t.Fatal(err)
	}
	tListdata := []novelchack.List{}
	if err := Read(&tListdata); err != nil {
		t.Fatal(err)
	}
	if len(tListdata) != 1 {
		t.Fatal("Listdata Error")
	}
	if tListdata[0].LastStoryT != "test2" {
		t.Fatal("Listdata Error")
	}
	if err := Write(tmpListdata); err != nil {
		t.Fatal(err)
	}
	if err := Read(&tListdata); err != nil {
		t.Fatal(err)
	}
	if len(tListdata) != 1 {
		t.Fatal("Listdata Error")
	}
	tmpListdata = novelchack.List{
		Title: "test2",
		Url:   "http://test2.com",
	}
	if err := Write(tmpListdata); err != nil {
		t.Fatal(err)
	}
	if err := Read(&tListdata); err != nil {
		t.Fatal(err)
	}
	if len(tListdata) != 2 {
		t.Fatal("Listdata Error")
	}
}

func TestCount(t *testing.T) {
	var wg sync.WaitGroup
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	max := 10
	SetMaxCount(max)
	AddCount()
	if ReadCount() != 1 {
		t.Fatal("Not Add count")
	}
	if ReadPerCount() != 0.1 {
		t.Fatal("Not Set Max")
	}

	for i := 0; i < 9; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			AddCount()
		}()
	}
	wg.Wait()
	if ReadPerCount() != 1 {
		t.Fatal("Not Add ")
	}
}
