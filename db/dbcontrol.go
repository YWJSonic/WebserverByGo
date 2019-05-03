package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"../code"
	"../data"
	"../messagehandle/errorlog"
)

type sqlCLi struct {
	DB *sql.DB
}

var gameBDSQL *sqlCLi
var logDBSQL *sqlCLi

func SetDBCOnn() {
	connectGameDB()
	connectLogDB()
}

// Connect New connect
func connectGameDB() (db *sql.DB, err error) {
	if gameBDSQL == nil {
		gameBDSQL = new(sqlCLi)
		//root:123456@/gamedb
		sqlstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&timeout=30s", data.DBUser, data.DBPassword, data.DBIP, data.DBPORT, "gamedb")
		fmt.Println(sqlstr, data.DBDataSourceName)
		db, err := sql.Open("mysql", sqlstr) // data.DBDataSourceName)

		connMaxLifetime := 4
		maxIdleConns := 3
		maxOpenConns := 6

		fmt.Printf("connMaxLifetime:%d\n", connMaxLifetime)
		db.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Hour)

		fmt.Printf("maxIdleConns:%d\n", maxIdleConns)
		db.SetMaxIdleConns(maxIdleConns)

		fmt.Printf("maxOpenConns:%d\n", maxOpenConns)
		db.SetMaxOpenConns(maxOpenConns)
		if err != nil {
			return nil, err
		}

		gameBDSQL.DB = db
	}

	return gameBDSQL.DB, nil

	// if gameBDSQL == nil {
	// 	gameBDSQL = new(MySQLCli)
	// 	var err error
	// 	gameBDSQL.DB, err = sql.Open("mysql", data.DBDataSourceName)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// gameBDSQL.DB.SetMaxIdleConns(3)
	// }

	// return gameBDSQL.DB, nil
}

// Connect New connect
func connectLogDB() (db *sql.DB, err error) {
	if logDBSQL == nil {
		logDBSQL = new(sqlCLi)
		//root:123456@/gamedb
		sqlstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&timeout=30s", data.DBUser, data.DBPassword, data.DBIP, data.DBPORT, "logdb")
		fmt.Println(sqlstr, data.DBDataSourceName)
		db, err := sql.Open("mysql", sqlstr) // data.DBDataSourceName)

		connMaxLifetime := 4
		maxIdleConns := 3
		maxOpenConns := 6

		fmt.Printf("connMaxLifetime:%d\n", connMaxLifetime)
		db.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Hour)

		fmt.Printf("maxIdleConns:%d\n", maxIdleConns)
		db.SetMaxIdleConns(maxIdleConns)

		fmt.Printf("maxOpenConns:%d\n", maxOpenConns)
		db.SetMaxOpenConns(maxOpenConns)
		if err != nil {
			return nil, err
		}

		logDBSQL.DB = db
	}

	return gameBDSQL.DB, nil
}

// CallRead call stored procedure
func CallRead(name string, args ...interface{}) ([]interface{}, errorlog.ErrorMsg) {
	QueryStr := makeQueryStr(name, len(args))
	request, err := query(QueryStr, args...)
	return request, err

}

// CallReadOutMap call stored procedure
func CallReadOutMap(name string, args ...interface{}) ([]map[string]interface{}, errorlog.ErrorMsg) {
	QueryStr := makeQueryStr(name, len(args))
	request, err := queryOutMap(QueryStr, args...)
	return request, err

}

// CallWrite ...
func CallWrite(name string, args ...interface{}) (sql.Result, errorlog.ErrorMsg) {
	QueryStr := makeQueryStr(name, len(args))
	request, err := exec(QueryStr, args...)
	return request, err
}

// Query Use to SELECT return array, first is Keys
func query(query string, args ...interface{}) ([]interface{}, errorlog.ErrorMsg) {
	err := errorlog.New()

	res, errMsg := gameBDSQL.DB.Query(query, args...)
	if errMsg != nil {
		err.ErrorCode = code.FailedPrecondition
		err.Msg = "DBExecFail"
		return nil, err
	}

	request := makeScanArray(res)
	defer res.Close()

	return request, err
}

// QueryOutMap Use to SELECT return array map
func queryOutMap(query string, args ...interface{}) ([]map[string]interface{}, errorlog.ErrorMsg) {
	err := errorlog.New()

	res, errMsg := gameBDSQL.DB.Query(query, args...)
	if errMsg != nil {
		err.ErrorCode = code.FailedPrecondition
		err.Msg = "DBExecFail"
		return nil, err
	}

	request := makeScanMap(res)
	defer res.Close()

	return request, err
}

// Exec Use to INSTER, UPDATE, DELETE
func exec(query string, args ...interface{}) (sql.Result, errorlog.ErrorMsg) {
	err := errorlog.New()

	res, errMsg := gameBDSQL.DB.Exec(query, args...)
	if errMsg != nil {
		err.ErrorCode = code.FailedPrecondition
		err.Msg = "DBExecFail"
		return nil, err
	}
	res.LastInsertId()
	return res, err

}

func makeQueryStr(name string, args int) string {
	query := "CALL " + name + "("

	if args > 0 {
		for i := 0; i < args; i++ {
			query += "?,"
		}
		query = query[:len(query)-1]
	}
	query += ");"
	return query
}

func makeScanArray(rows *sql.Rows) []interface{} {
	var Result []interface{}
	Keys, err := rows.Columns()
	Result = append(Result, Keys)

	for rows.Next() {
		Row := make([]interface{}, len(Keys))
		for i := range Keys {
			Row[i] = new(sql.RawBytes)
		}

		err = rows.Scan(Row...)
		if err != nil {
			log.Fatalln(err)
		}

		Result = append(Result, Row)

	}

	return Result
}

func makeScanMap(rows *sql.Rows) []map[string]interface{} {
	Keys, err := rows.Columns()
	types, err := rows.ColumnTypes()
	scanArgs := make([]interface{}, len(Keys))
	values := make([][]byte, len(Keys))
	Result := []map[string]interface{}{}

	for rows.Next() {
		for i := range Keys {
			scanArgs[i] = &values[i]
		}

		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Fatalln(err)
		}

		tomap := map[string]interface{}{}
		for i, key := range Keys {
			if types[i].DatabaseTypeName() == "BIGINT" {
				val, errMsg := strconv.ParseInt(string(values[i]), 10, 64)
				if errMsg != nil {
					panic("makeScanMap Error BIGINT")
				}
				tomap[key] = val
			} else if types[i].DatabaseTypeName() == "INT" {
				val, errMsg := strconv.ParseInt(string(values[i]), 10, 0)
				if errMsg != nil {
					panic("makeScanMap Error INT")
				}
				tomap[key] = val
			} else if types[i].DatabaseTypeName() == "TINYINT" {
				val, errMsg := strconv.ParseInt(string(values[i]), 10, 0)
				if errMsg != nil {
					panic("makeScanMap Error TINYINT")
				}
				tomap[key] = val == 1
			} else if types[i].DatabaseTypeName() == "VARCHAR" {
				tomap[key] = string(values[i])
			} else if types[i].DatabaseTypeName() == "TEXT" {
				tomap[key] = string(values[i])
			} else {
				tomap[key] = values[i]
			}
		}

		Result = append(Result, tomap)
		fmt.Println(Result)
	}

	return Result
}
