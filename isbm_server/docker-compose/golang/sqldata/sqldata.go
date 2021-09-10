package sqldata

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

//SQLの設定ファイル
type SqlConfig struct {
	dataSourceName string
	dbname         string
	dbtype         string
	setupflag      bool
	dbpath         *sql.DB
	tableName      string
	column_name    []string
	column_type    []string
}

const timeout_ms = 500

//SqlSetup
// dbtype :sqlite3 or mysql
// len(v) = 1 or 6
func (t *SqlConfig) SqlSetup(dbtype string, v ...interface{}) error {
	if t.setupflag {
		fmt.Println("セットアップ済み")
		return nil
	}
	switch dbtype {
	case "sqlite3":
		if len(v) == 1 {
			fmt.Println(v[0])
			t.sqlSetupSqlite(dbtype, v[0].(string))
			return nil
		}
		break
	case "mysql":
		if len(v) == 6 {
			dbUser := v[0].(string)
			dbPass := v[1].(string)
			conectType := v[2].(string)
			ipAddr := v[3].(string)
			portN := v[4].(string)
			dbName := v[5].(string)
			err := t.sqlSetupMySql(dbtype, dbUser, dbPass, conectType, ipAddr, portN, dbName)
			return err
		}
		break
	default:
		break
	}
	return errors.New("input data Err")
}

//sqlSetupSqlite
func (t *SqlConfig) sqlSetupSqlite(dbtype, dbpath string) {
	t.dbtype = dbtype
	t.dataSourceName = dbpath
}

//sqlSetupMySql
func (t *SqlConfig) sqlSetupMySql(dbtype, dbUser, dbPass, conectType, ipAddr, portN, dbName string) error {
	t.dbtype = dbtype
	t.dataSourceName = dbUser + ":" + dbPass + "@" + conectType + "(" + ipAddr + ":" + portN + ")/"
	t.dbname = dbName
	t.setupflag = true
	db, err := t.openDBMySql("")
	if err != nil {
		return err
	}
	fmt.Println("データベースの接続成功")
	db, err = t.openDBMySql(t.dbname)
	if err != nil {
		err = t.createDBMySql()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		db, err = t.openDBMySql(t.dbname)
	}
	t.dbpath = db
	return err
}

//createDBMySql はデータベースを作成する
func (t *SqlConfig) createDBMySql() error {
	if !t.setupflag {
		return errors.New("not run setup")
	}
	db, err := sql.Open("mysql", t.dataSourceName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	str := "CREATE DATABASE " + t.dbname + ";"
	_, err = db.Exec(str)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully created database..")
	}
	db.Close()
	return err
}

//openDBMySql()
//データベースを開く
func (t *SqlConfig) openDBMySql(table string) (*sql.DB, error) {
	if !t.setupflag {
		return nil, errors.New("not run setup")
	}
	fmt.Println("OpenDB : " + t.dataSourceName + table + "?parseTime=true")
	db, err := sql.Open("mysql", t.dataSourceName+table+"?parseTime=true")
	if err != nil {
		fmt.Println(err.Error())
	}
	for i := 1; i < 4; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		if table != "" {
			break
		}
		fmt.Printf("リトライ %v\n", i)
		time.Sleep(time.Millisecond * timeout_ms)

	}
	if err != nil {
		if table == "" {
			fmt.Println("データベースの接続失敗")
		} else {
			fmt.Println(table + "のテーブルがありません")
		}
		db.Close()
	}
	return db, err
}
func (t *SqlConfig) TableSetup(table string, columnName, columnType []string) error {
	var err error
	t.tableName = table
	if len(columnName) != len(columnType) {
		return errors.New("The number of rows does not match")
	}
	t.column_name = columnName
	t.column_type = columnType
	switch t.dbtype {
	case "sqlite3":
		err = t.creatDbSqlite()
	case "mysql":
		err = t.createTableDBMysql()
		break
	default:
		errout := "sql type don't :" + t.dbtype
		return errors.New(errout)
	}

	if err == nil {
		t.setupflag = true
	}

	return err
}

//CreateTableDB()
func (t *SqlConfig) createTableDBMysql() error {
	if !t.setupflag {
		return errors.New("not run setup")
	}
	cmd := "create table "
	cmd += t.tableName + " "
	cmd += "("
	cmd += "id" + " " + "int " + "NOT NULL AUTO_INCREMENT" + ","
	for i := 0; i < len(t.column_type); i++ {
		cmd += t.column_name[i] + " " + t.column_type[i] + ","
	}
	cmd += "created_at DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00',"
	// cmd += "updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,"
	cmd += "updated_at DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00',"
	cmd += "PRIMARY KEY (id)"
	cmd += ")"
	fmt.Println(cmd)
	stmt, err := t.dbpath.Prepare(cmd)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("データベース作成")
	}
	return nil

}
func (t *SqlConfig) creatDbSqlite() error {
	var err error
	t.dbpath, err = sql.Open(t.dbtype, t.dataSourceName)
	if err != nil {
		return err
	}
	cmd := `CREATE TABLE IF NOT EXISTS ` + t.tableName + `(
`
	cmd += `id integer PRIMARY KEY, `
	for i := 0; i < len(t.column_name); i++ {
		cmd += t.column_name[i] + " " + t.column_type[i] + ","
	}
	cmd += `created_at datatime NOT NULL,
		updated_at datatime NOT NULL
		)`
	_, err = t.dbpath.Exec(cmd)
	if err != nil {
		return nil
	} else {
		fmt.Printf("Create table %v\n", t.tableName)
	}
	return err

}

//CloseDB()
func (t *SqlConfig) CloseDB() {
	t.setupflag = false
	t.dbpath.Close()
}

//AddDB()
//
func (t *SqlConfig) AddDB(v ...interface{}) int {
	id := 0
	if !t.setupflag {
		return id
	}
	time_now := time.Now()
	cmd := "insert into " + t.tableName + "("
	for i := 0; i < len(t.column_name); i++ {
		if i == 0 {
			cmd += t.column_name[i]
		} else {
			cmd += "," + t.column_name[i]
		}
	}
	cmd += "," + "created_at"
	cmd += "," + "updated_at"
	cmd += ") values("
	for i := 0; i < len(t.column_name); i++ {
		if i == 0 {
			cmd += "?"
		} else {
			cmd += ",?"
		}
	}
	cmd += ",'" + time_now.Format("2006-01-02 15:04:05.999999999") + "'"
	cmd += ",'" + time_now.Format("2006-01-02 15:04:05.999999999") + "'"
	cmd += ")"
	fmt.Println(cmd)
	stmt, err := t.dbpath.Prepare(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return -1
	}
	defer stmt.Close()
	_, err = stmt.Exec(v...)
	if err != nil {
		fmt.Println(err.Error())
		return -1
	} else {
		fmt.Println("追加成功")
	}
	cmd = "select max(id) from " + t.tableName + ";"
	row, err := t.dbpath.Query(cmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer row.Close()
	row.Next()
	err = row.Scan(&id)
	if err != nil {
		fmt.Println(err.Error())
	}
	return id
}

//ReadAllDB()
func (t *SqlConfig) ReadAllDB() (*sql.Rows, error) {
	cmd := "select * from " + t.tableName
	row, err := t.dbpath.Query(cmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	return row, err
}

//SerchTimeDB()
func (t *SqlConfig) SerchTimeDB(num int) (*sql.Rows, error) {
	nowtime := time.Now()
	cmd := "select * from " + t.tableName
	switch num {
	case 1: //week
		cmd += " " + "where " + "updated_at "
		cmd += "between '" + nowtime.Add(-24*time.Hour*7).Format("2006-01-02") + "' and '"
		cmd += nowtime.Add(24*time.Hour).Format("2006-01-02") + "'"
	case 2: //month
		cmd += " " + "where " + "updated_at "
		cmd += "between '" + nowtime.Add(-24*time.Hour*30).Format("2006-01-02") + "' and '"
		cmd += nowtime.Add(24*time.Hour).Format("2006-01-02") + "'"
	default: //today
		cmd += " " + "where " + "updated_at "
		cmd += "between '" + nowtime.Format("2006-01-02") + "' and '"
		cmd += nowtime.Add(24*time.Hour).Format("2006-01-02") + "'"
	}
	row, err := t.dbpath.Query(cmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	return row, err
}

//ReadIdDB()
func (t *SqlConfig) ReadIdDB(Id string) (*sql.Rows, error) {
	cmd := "select * from " + t.tableName
	cmd += " " + "where id=" + Id
	row, err := t.dbpath.Query(cmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	return row, err
}

func (t *SqlConfig) SeachReadDB(word string, serchKey []string) (*sql.Rows, error) {
	cmd := "select * from " + t.tableName
	cmd += " " + "where "
	for i := 0; i < len(serchKey); i++ {
		if i == 0 {
			cmd += serchKey[i] + " " + "like '%" + word + "%'"
		} else {
			cmd += " or " + serchKey[i] + " " + "like '%" + word + "%'"
		}
	}
	row, err := t.dbpath.Query(cmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	return row, err
}

//EditDB()
//
func (t *SqlConfig) EditDB(No string, v ...interface{}) {
	if !t.setupflag {
		return
	}
	time_now := time.Now()
	cmd := "update " + t.tableName + " set "
	for i := 0; i < len(t.column_name); i++ {
		if i == 0 {
			cmd += t.column_name[i] + "=?"
		} else {
			cmd += "," + t.column_name[i] + "=?"
		}
	}
	cmd += ",updated_at='" + time_now.Format("2006-01-02 15:04:05.999999999") + "'"
	cmd += " where id=" + No
	fmt.Println(cmd)
	stmt, err := t.dbpath.Prepare(cmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(v...)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("編集成功")
	}
}

//DeleteDB()
//
func (t *SqlConfig) DeleteDB(No int) {
	if !t.setupflag {
		return
	}
	cmd := "delete from " + t.tableName + " where id=?"
	stmtDelete, err := t.dbpath.Prepare(cmd)
	if err != nil {
		panic(err.Error())
	}
	defer stmtDelete.Close()

	result, err := stmtDelete.Exec(No)
	if err != nil {
		panic(err.Error())
	}

	_, err = result.RowsAffected()
	if err != nil {
		panic(err.Error())
	}

}
