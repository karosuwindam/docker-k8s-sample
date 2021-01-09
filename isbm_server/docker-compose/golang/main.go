package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	// _ "github.com/mattn/go-sqlite3"
)

type DbsetupData struct {
	Dbtype     string `json:type`
	Dbuser     string `json:user`
	Dbpassword string `json:password`
	Ipaddr     string `json:host`
	Port       string `json:port`
	Databs     string `json:database`
	DbPath     string `json:dbpath`
}

const ( //環境変数による設定読み込み
	envSqlType     = "APP_SQL_TYPE"
	envSqlUser     = "APP_SQL_USER"
	envSqlPass     = "APP_SQL_PASS"
	envSqlHost     = "APP_SQL_HOST"
	envSqlPort     = "APP_SQL_PORT"
	envSqlDatabase = "APP_SQL_DATABASE"
	envSqlDbPath   = "APP_SQL_PATH"
)
const ( //初期化時の設定値
	dbtype     = "mysql"
	dbuser     = "bookserver"
	dbpassword = "bookserver"
	conecttype = "tcp"
	ipaddr     = "mysql_host"
	port       = "3306"
	databs     = "isbn_bookbase"
	Dbpath     = "./test.db"
)

var bookdatatable BookData

func dbsetup() DbsetupData {
	var tmp = DbsetupData{}
	var buf bytes.Buffer

	config_json := "config/dbsetup.json"
	raw, err := ioutil.ReadFile(config_json)
	if err != nil {
		tmp.Dbtype = dbtype
		tmp.Dbuser = dbuser
		tmp.Dbpassword = dbpassword
		tmp.Ipaddr = ipaddr
		tmp.Port = port
		tmp.Databs = databs
		tmp.DbPath = Dbpath
		fp, err := os.Create(config_json)
		if err != nil {
			panic(err)
		}
		// jsonエンコード
		outputJson, err := json.Marshal(&tmp)
		if err != nil {
			panic(err)
		}
		json.Indent(&buf, outputJson, "", "  ")
		fp.Write(buf.Bytes())
		fp.Close()

	} else {
		var fc DbsetupData
		json.Unmarshal(raw, &fc)
		if fc.Dbtype == "" {
			fc.Dbtype = dbtype
		}
		if fc.Databs == "" {
			fc.Dbuser = dbuser
		}
		if fc.Dbpassword == "" {
			fc.Dbpassword = dbpassword
		}
		if fc.Ipaddr == "" {
			fc.Ipaddr = ipaddr
		}
		if fc.Port == "" {
			fc.Port = port
		}
		if fc.Databs == "" {
			fc.Databs = databs
		}
		if fc.DbPath == "" {
			fc.DbPath = Dbpath
		}
		tmp = fc

	}
	if str := os.Getenv(envSqlType); str != "" {
		tmp.Dbtype = str
	}
	if str := os.Getenv(envSqlHost); str != "" {
		tmp.Ipaddr = str
	}
	if str := os.Getenv(envSqlDatabase); str != "" {
		tmp.Databs = str
	}
	if str := os.Getenv(envSqlUser); str != "" {
		tmp.Dbuser = str
	}
	if str := os.Getenv(envSqlPass); str != "" {
		tmp.Dbpassword = str
	}
	if str := os.Getenv(envSqlPort); str != "" {
		tmp.Port = str
	}
	if str := os.Getenv(envSqlDbPath); str != "" {
		tmp.DbPath = str
	}

	return tmp
}

func main() {
	var web WebSetupData
	dbdata := dbsetup()
	var err error
	bookdatatable.sqltype = dbdata.Dbtype
	switch bookdatatable.sqltype {
	case "mysql":
		err = bookdatatable.bookdatas_sql.SqlSetup(dbdata.Dbtype, dbdata.Dbuser, dbdata.Dbpassword, conecttype, dbdata.Ipaddr, dbdata.Port, dbdata.Databs)
	case "sqlite3":
		err = bookdatatable.bookdatas_sql.SqlSetup(dbdata.Dbtype, dbdata.DbPath)
	default:
		str := "SQL Type" + bookdatatable.sqltype + "is not"
		err = errors.New(str)
	}

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer bookdatatable.bookdatas_sql.CloseDB()

	bookdatatable.Sqlsetup("bookdatabase")

	err = web.websetup()
	if err == nil {
		web.webstart()
	} else {
		log.Fatal(err)
	}
}
