package main

import (
	"log"
	"strconv"
	"time"
)

type ListData struct {
	Month     string     `json:Month`
	Year      string     `json:Year`
	Comic     []BookList `json:Comic`
	LiteNobel []BookList `json:LiteNobel`
}

var GrobalListData map[int]ListData

func main() {
	tmp := map[int]ListData{}
	GrobalListData = tmp
	ch1 := make(chan bool)
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
	// fmt.Println(GetComicList("", "", LITENOVEL))

	var web WebSetupData
	err := web.websetup()
	if err == nil {
		web.webstart()
	} else {
		log.Fatal(err)
	}
}
