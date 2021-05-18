package main

import (
	"log"
	"sort"
	"strconv"
	"time"

	"./novel_chack"
)

type ListData struct {
	Month     string     `json:Month`
	Year      string     `json:Year`
	Comic     []BookList `json:Comic`
	LiteNobel []BookList `json:LiteNobel`
}

var GrobalListData map[int]ListData

var Listdata []novel_chack.List

func main() {
	tmp := map[int]ListData{}
	GrobalListData = tmp
	ch1 := make(chan bool)

	ch2 := make(chan bool)

	go func() {
		for {
			listtmp := []novel_chack.List{}
			pass := "bookmarks_2021_05_18.html"
			data := novel_chack.ReadBookBark(pass)

			for _, str := range data {
				tmp := novel_chack.GetSyousetu(str.Url)
				if tmp.Title != "" {
					listtmp = append(listtmp, tmp)
					// fmt.Println(tmp)
				}
			}
			sort.Slice(listtmp, func(i, j int) bool { return listtmp[i].Lastdate.Unix() > listtmp[j].Lastdate.Unix() })
			Listdata = listtmp
			ch2 <- true
			time.Sleep(time.Hour)
		}
	}()
	go func() {
		for {
			t := time.Now()
			for i := 0; i < 3; i++ {
				var listdata ListData
				if (int(t.Month()) + i) <= 12 {
					listdata.Month = strconv.Itoa((int(t.Month()) + i))
					listdata.Year = strconv.Itoa(int(t.Year()))
				} else {
					listdata.Month = strconv.Itoa((int(t.Month()) + i) % 12)
					listdata.Year = strconv.Itoa(int(t.Year()) + 1)
				}

				listdata.LiteNobel = GetComicList(listdata.Year, listdata.Month, LITENOVEL)
				listdata.Comic = GetComicList(listdata.Year, listdata.Month, COMIC)
				GrobalListData[i] = listdata
			}
			ch1 <- true
			time.Sleep(time.Hour * 12)
		}
	}()
	<-ch1
	<-ch2
	// fmt.Println(GetComicList("", "", LITENOVEL))

	var web WebSetupData
	err := web.websetup()
	if err == nil {
		web.webstart()
	} else {
		log.Fatal(err)
	}
}