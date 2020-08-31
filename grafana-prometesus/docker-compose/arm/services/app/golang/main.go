package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type AM2320Data struct {
	Id        int64     `json:"id"`
	Tmp       float32   `json:"tmp"`
	Hum       float32   `json:"hum"`
	CreatedAt time.Time `json:"createdAt"`
}

const (

	// ユーザー名
	DBUser = "root"

	// パスワード
	DBPass = "root"

	// プロトコル
	DBProtocol = "tcp(mysql_host1:3306)"

	// DB名
	DBName = "test_database"
)

func connectGorm() *gorm.DB {
	connectTemplate := "%s:%s@%s/%s"
	connect := fmt.Sprintf(connectTemplate, DBUser, DBPass, DBProtocol, DBName)
	db, err := gorm.Open("mysql", connect)

	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return db
}
func insert(data AM2320Data, db *gorm.DB) {
	db.NewRecord(data)
	db.Create(&data)
}

var data_hum, data_tmp float64

func metrics(w http.ResponseWriter, r *http.Request) {
	output := "senser_data{type=\"tmp\"} " + strconv.FormatFloat(data_tmp, 'f', 1, 64)
	output += "\n" + "senser_data{type=\"hum\"} " + strconv.FormatFloat(data_hum, 'f', 1, 64)
	fmt.Fprintf(w, "%s", output)

}
func data(w http.ResponseWriter, r *http.Request) {
	output := "<html><body><a href=\"/metrics\">metrics</a></body></html>"
	fmt.Fprintf(w, "%s", output)

}

func main() {
	db := connectGorm()
	if db == nil {
		for i := 0; i < 2; i++ {
			time.Sleep(10 * time.Second)
			db := connectGorm()
			if db != nil {
				break
			}
		}
	}
	if db == nil {
		fmt.Print("don't open db")
		return
	}
	db.AutoMigrate(&AM2320Data{})
	db.Close()
	go func() {
		for {
			hum, tmp := ReadAM2320()
			data_hum = float64(hum)
			data_tmp = float64(tmp)
			sample1 := AM2320Data{Tmp: tmp, Hum: hum}
			// fmt.Printf("%v,%v\n", sample1.Tmp, sample1.Hum)
			db = connectGorm()
			if db == nil {
				fmt.Print("don't open db")
			} else {
				insert(sample1, db)
				db.Close()
			}
			time.Sleep(15 * time.Second)
		}
	}()

	port := "4000"
	fmt.Println("start server " + port)
	http.HandleFunc("/metrics", metrics)
	http.HandleFunc("/", data)
	http.ListenAndServe(":"+port, nil)

}
