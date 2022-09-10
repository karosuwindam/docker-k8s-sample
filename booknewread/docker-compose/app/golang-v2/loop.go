package main

import (
	"booknewread/calendarscriping"
	"booknewread/dirread"
	"booknewread/novel_chack"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"time"
)

type ListData struct {
	NowTime   time.Time                   `json:Nowtime`
	Month     string                      `json:Month`
	Year      string                      `json:Year`
	Comic     []calendarscriping.BookList `json:Comic`
	LiteNobel []calendarscriping.BookList `json:LiteNobel`
}

type ReloadFlag struct {
	BookFlag     bool
	BookMarkFlag bool
}

type Status struct {
	BookNowTIme     time.Time `json:booknowtime`
	BookStatus      string    `json:bookstatus`
	BookMarkNowTime time.Time `json:bookmarknowtime`
	BookMarkStatus  string    `json:bookmarkstatus`
}

//環境変数で操作できる数値周
type SetupEnvData struct {
	MaxAccess  int //一つのサイトに対してアクセスできる限界数
	MaxBackDay int //今日を基準に表示する前日数
}

//デフォルト
const (
	DEF_ACCESS  = 2
	DEF_BACKDAY = 5
	BDatafolder = "bookmark"
)

type LoopData struct {
	GrobalListData map[int]ListData
	Listdata       []novel_chack.List
	GrobalStatus   Status
	Reloadflag     ReloadFlag
	EnvData        SetupEnvData
	filedata       dirread.Dirtype
	stopflag       bool
	ch1            bool
}

func SetupEnv() SetupEnvData {
	var output SetupEnvData
	if true {
		output.MaxAccess = 3
	} else {
		output.MaxAccess = DEF_ACCESS
	}
	if true {
		output.MaxBackDay = 2
	} else {
		output.MaxBackDay = DEF_BACKDAY
	}
	return output
}

func loopSetup() (*LoopData, error) {
	cfg := &LoopData{stopflag: false}
	cfg.GrobalListData = map[int]ListData{}
	cfg.GrobalStatus = Status{
		BookNowTIme:     time.Time{},
		BookStatus:      "Reload",
		BookMarkNowTime: time.Time{},
		BookMarkStatus:  "Reload",
	}
	// tmp := make(chan bool)
	cfg.ch1 = false
	cfg.EnvData = SetupEnv()
	cfg.filedata.Setup(BDatafolder)
	cfg.filedata.Read("/")
	if len(cfg.filedata.Data) == 0 {
		return nil, errors.New("err bookmark not file")
	}

	return cfg, nil
}

func loopWait(cfg *LoopData) {
	for {
		if cfg.ch1 {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func ckNobelloop(cfg *LoopData) {
	readflag := false

	log.Println("start novel data count")
	ch := novel_chack.Setup(cfg.EnvData.MaxAccess)
	data := []novel_chack.BookBark{}
	if len(data) == 0 {
		log.Println("No bookmark")
		return
	}
	for {
		starttime := time.Now()
		cfg.GrobalStatus.BookMarkStatus = "Reload"

		listtmp := []novel_chack.List{}
		cfg.filedata.Read("/")
		if len(cfg.filedata.Data) != 0 {
			data = []novel_chack.BookBark{}
			for _, fd := range cfg.filedata.Data {
				pass := fd.RootPath + fd.Name
				ncr := novel_chack.ReadBookBark(pass)
				for _, tmpd := range ncr {
					data = append(data, tmpd)
				}
				data = novel_chack.BookBarkSout(data)
			}
		}

		//マルチ処理
		count := 0
		cpus := runtime.NumCPU()
		ch_th := make(chan int, cpus)
		tempdata := make([]novel_chack.List, len(data))
		for i, str := range data {
			// log.Println(str.Url)
			ch_th <- 1
			go func(num int, url string) {
				tempdata[num] = ch.ChackUrldata(url)
				count += <-ch_th
			}(i, str.Url)
			cfg.GrobalStatus.BookMarkStatus = "Reload:" + strconv.Itoa(count*100/len(data)) + "%"
			// time.Sleep(time.Millisecond * 100)
		}
		for {
			cfg.GrobalStatus.BookMarkStatus = "Reload:" + strconv.Itoa(count*100/len(data)) + "%"
			if len(ch_th) == 0 {
				for _, str := range tempdata {
					if str.Title != "" {
						// fmt.Println(str)
						listtmp = append(listtmp, str)
					}
				}
				break
			}
			time.Sleep(time.Millisecond * 100)
		}

		sort.Slice(listtmp, func(i, j int) bool { return listtmp[i].Lastdate.Unix() > listtmp[j].Lastdate.Unix() })
		endtime := time.Now()
		cfg.Listdata = listtmp
		if !readflag {
			cfg.ch1 = true
			readflag = true
		}
		cfg.GrobalStatus.BookMarkStatus = "OK"
		cfg.GrobalStatus.BookMarkNowTime = time.Now()
		cfg.Reloadflag.BookMarkFlag = false
		log.Println("read novel data end", (endtime.Sub(starttime)).Seconds(), "s")
		for i := 0; i < 60*60; i++ {
			if cfg.Reloadflag.BookMarkFlag || cfg.stopflag {
				break
			}
			time.Sleep(time.Second)
		}
		if cfg.stopflag {
			break
		}
		log.Println("reload novel data")
	}
}

func ckBooklloop(cfg *LoopData) {
	log.Println("start new book data count")
	for {
		starttime := time.Now()
		cfg.GrobalStatus.BookStatus = "Reload"
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

			listdata.LiteNobel = calendarscriping.GetComicList(listdata.Year, listdata.Month, calendarscriping.LITENOVEL)
			listdata.Comic = calendarscriping.GetComicList(listdata.Year, listdata.Month, calendarscriping.COMIC)
			listdata.LiteNobel = calendarscriping.FilterComicList(listdata.LiteNobel, cfg.EnvData.MaxBackDay)
			listdata.Comic = calendarscriping.FilterComicList(listdata.Comic, cfg.EnvData.MaxBackDay)
			cfg.GrobalListData[i] = listdata
		}
		endtime := time.Now()
		cfg.GrobalStatus.BookStatus = "OK"
		cfg.GrobalStatus.BookNowTIme = time.Now()
		cfg.Reloadflag.BookFlag = false
		log.Println("read new book data end", (endtime.Sub(starttime)).Seconds(), "s")
		for i := 0; i < 60*60*12; i++ {
			if cfg.Reloadflag.BookFlag {
				break
			}
			time.Sleep(time.Second)
		}
		log.Println("reload new book data")
	}
}

func loopStop(cfg *LoopData) {
	cfg.stopflag = true
}

func (t *LoopData) getlocaljson(w http.ResponseWriter, r *http.Request) {
	form_data := ""
	timedata := time.Now().Sub(t.GrobalStatus.BookNowTIme).Seconds()
	if timedata > 300 {
		t.Reloadflag.BookMarkFlag = true
		t.Reloadflag.BookFlag = true
	}
	r.ParseForm()
	for cnt, strs := range r.Form {
		form_data += " " + cnt + ":"
		for i, str := range strs {
			if i == 0 {
				form_data += str
			} else {
				form_data += "," + str
			}
		}
	}
	log.Println(r.Method + ":" + r.URL.Path + " " + form_data)

	if r.Method == "GET" {
		jsondata, err := json.Marshal(t.GrobalListData[0])
		if err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			fmt.Fprintf(w, "%s", jsondata)
		}
	} else if r.Method == "POST" {
		mf := r.MultipartForm

		// 通常のリクエスト
		if mf != nil {
			for k, v := range mf.Value {
				fmt.Printf("%v : %v", k, v)
			}
		}
		page := r.FormValue("page")
		tmp_page, _ := strconv.Atoi(page)
		if tmp_page >= len(t.GrobalListData) {
			tmp_page = 0
		}
		jsondata, err := json.Marshal(t.GrobalListData[tmp_page])
		if err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			fmt.Fprintf(w, "%s", jsondata)
		}
	}

}

func (t *LoopData) getnowdata(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method + ":" + r.URL.Path)

	if r.Method == "GET" {
		jsondata, err := json.Marshal(t.Listdata)
		if err != nil {
			fmt.Fprint(w, err.Error())
		} else {
			fmt.Fprintf(w, "%s", jsondata)
		}
	}

}

func (t *LoopData) restart(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method + ":" + r.URL.Path)

	if r.Method == "POST" {
		t.Reloadflag.BookMarkFlag = true
		t.Reloadflag.BookFlag = true
		fmt.Fprintf(w, "OK")
	} else {
		fmt.Fprintf(w, "NG")

	}

}
func (t *LoopData) status(w http.ResponseWriter, r *http.Request) {
	// jsondata, err := json.Marshal(GrobalStatus)
	jsondata, err := json.Marshal(t.GrobalStatus)
	if err != nil {
		fmt.Fprint(w, err.Error())
	} else {
		fmt.Fprintf(w, "%s", jsondata)
	}

}
