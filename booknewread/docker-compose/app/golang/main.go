package main

import (
	"log"
	"sort"
	"strconv"
	"time"

	"./novel_chack"
)

type ListData struct {
	NowTime   time.Time  `json:Nowtime`
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
	ch1 := make(chan bool, 2)

	ch2 := make(chan bool, 2)

	go func() {
		count := 0
		log.Println("start novel data count")
		for {
			listtmp := []novel_chack.List{}
			pass := "bookmarks_2021_05_18.html"
			data := novel_chack.ReadBookBark(pass)

			for _, str := range data {
				tmp := novel_chack.ChackUrldata(str.Url)
				if tmp.Title != "" {
					listtmp = append(listtmp, tmp)
					// fmt.Println(tmp)
				}
			}
			sort.Slice(listtmp, func(i, j int) bool { return listtmp[i].Lastdate.Unix() > listtmp[j].Lastdate.Unix() })
			Listdata = listtmp
			if count == 0 {
				ch2 <- true
			}
			time.Sleep(time.Hour)
			log.Println("reload novel data")
			count++
		}
	}()
	go func() {
		count := 0
		log.Println("start new book data count")
		for {
			t := time.Now()
			for i := 0; i < 3; i++ {
				var listdata ListData
				listdata.NowTime = time.Now()
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
			if count == 0 {
				ch1 <- true
			}
			time.Sleep(time.Hour * 12)
			log.Println("reload new book data")
			count++
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
