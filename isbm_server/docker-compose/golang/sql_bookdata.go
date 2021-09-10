package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"./sqldata"
)

//
type SqlBookData struct {
	Id         int       `json:id`
	Isbn       string    `json:isbn`
	Title      string    `json:title`
	Writer     string    `json:writer`
	Brand      string    `json:brand`
	Ext        string    `json:ext`
	Synopsis   string    `json:synopsis`
	Image      string    `json:image`
	Created_at time.Time `json:created_at`
	Updated_at time.Time `json:updated_at`
}

type BookData struct {
	bookdatas_sql sqldata.SqlConfig
	sqltype       string
	Column_name   []string
}

//テーブルセットアップ
func (t *BookData) Sqlsetup(TName string) {
	tableName := TName
	column_type := []string{}
	column_name := []string{"isbn", "title", "writer", "brand", "ext", "synopsis", "image"}
	if t.sqltype == "sqlite3" {
		column_type = []string{"varchar", "varchar", "varchar", "varchar", "varchar", "varchar", "varchar"}
	} else {
		column_type = []string{"varchar(14)", "varchar(255)", "varchar(255)", "varchar(255)", "varchar(255)", "varchar(2048)", "varchar(255)"}
	}
	t.Column_name = column_name
	t.bookdatas_sql.TableSetup(tableName, column_name, column_type)
}

//検索用
func (t *BookData) scandata(row *sql.Rows) (SqlBookData, error) {
	var tmp SqlBookData
	var cat, uat string
	var err error
	layout := "2006-01-02 15:04:05"
	if t.sqltype == "sqlite3" {
		//"isbn", "title","writer","brand","ext","synopsis","image"
		err = row.Scan(&tmp.Id, &tmp.Isbn, &tmp.Title, &tmp.Writer, &tmp.Brand, &tmp.Ext, &tmp.Synopsis, &tmp.Image, &cat, &uat)
		if err == nil {
			tmp.Created_at, _ = time.Parse(layout, cat)
			tmp.Updated_at, _ = time.Parse(layout, uat)
		}
	} else {
		err = row.Scan(&tmp.Id, &tmp.Isbn, &tmp.Title, &tmp.Writer, &tmp.Brand, &tmp.Ext, &tmp.Synopsis, &tmp.Image, &tmp.Created_at, &tmp.Updated_at)
	}
	return tmp, err
}

//すべて表示
func (t *BookData) Scansql() []SqlBookData {
	output := []SqlBookData{}
	row, _ := t.bookdatas_sql.ReadAllDB()
	defer row.Close()
	for row.Next() {
		tmp, err := t.scandata(row)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		output = append(output, tmp)
	}
	return output
}

//キーワード検索
func (t *BookData) Serch_sql(keyword string) []SqlBookData {
	var row *sql.Rows
	timekeyword := []string{"today", "toweek", "tomonth"}
	keyflag := -1
	output := []SqlBookData{}
	for i := 0; i < len(timekeyword); i++ {
		if keyword == timekeyword[i] {
			keyflag = i
			break
		}
	}
	if keyflag >= 0 {
		row, _ = t.bookdatas_sql.SerchTimeDB(keyflag)
	} else {
		row, _ = t.bookdatas_sql.SeachReadDB(keyword, t.Column_name)
	}
	defer row.Close()
	for row.Next() {
		tmp, err := t.scandata(row)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		output = append(output, tmp)
	}
	return output
}

//ID検索
func (t *BookData) ScansqlId(Id string) SqlBookData {
	row, _ := t.bookdatas_sql.ReadIdDB(Id)
	defer row.Close()
	row.Next()
	tmp, err := t.scandata(row)
	if err != nil {
		fmt.Println(err.Error())
	}
	return tmp
}

//ADD
func (t *BookData) Add(isbn, title, writer, brand, ext, Synopsis, image string) int {
	return t.bookdatas_sql.AddDB(isbn, title, writer, brand, ext, Synopsis, image)
}

//Edit
func (t *BookData) Edit(id, isbn, title, writer, brand, ext, Synopsis, image string) {
	t.bookdatas_sql.EditDB(id, isbn, title, writer, brand, ext, Synopsis, image)
}

//Delete
func (t *BookData) Delte(id string) {
	num, err := strconv.Atoi(id)
	if err != nil {
		println(err.Error())
		return
	}
	t.bookdatas_sql.DeleteDB(num)
}
