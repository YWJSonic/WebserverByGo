package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"../code"
	"../crontab"
	"../data"
	"../foundation"
	"../messagehandle/errorlog"
)

type sqlCLi struct {
	DB *sql.DB
}
type sqlQuary struct {
	Quary string
	Args  []interface{}
}

var gameBDSQL *sqlCLi
var logDBSQL *sqlCLi
var payDBSQL *sqlCLi

// QueryLogChan channel for write log
var QueryLogChan chan string

// WriteGameChan channel for write game db
var WriteGameChan chan sqlQuary

// WritePayChan channel for write pay db
var WritePayChan chan sqlQuary

// SetDBConn init value
func SetDBConn() {
	QueryLogChan = make(chan string)
	WriteGameChan = make(chan sqlQuary)
	WritePayChan = make(chan sqlQuary)
	connectGameDB()
	connectLogDB()
	connectPayDB()

	// server start check today log table.
	NewLogTable(foundation.ServerNow().Format("20060102"))

	// set Schedule check next day log table.
	crontab.NewCronBaseJob("0 35 15 * * *", &crontab.LogCrontab{
		Params: func() string { return foundation.ServerNow().AddDate(0, 0, 1).Format("20060102") },
		FUN:    NewLogTable,
	})
}

// SQLSelect channel loop
func SQLSelect() {
	select {
	case dbLogQuery := <-QueryLogChan:
		CallWrite(logDBSQL.DB, dbLogQuery)
	case dbgameQuery := <-WriteGameChan:
		CallWrite(gameBDSQL.DB, dbgameQuery.Quary, dbgameQuery.Args...)
	case dbpayQuary := <-WritePayChan:
		CallWrite(payDBSQL.DB, dbpayQuary.Quary, dbpayQuary.Args...)
	}
}

// Connect New connect
func connectGameDB() (db *sql.DB, err error) {
	if gameBDSQL == nil {
		gameBDSQL = new(sqlCLi)
		sqlstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&timeout=30s", data.DBUser, data.DBPassword, data.DBIP, data.DBPORT, "gamedb")
		errorlog.LogPrintln("DB Connect:", sqlstr)
		db, err := sql.Open("mysql", sqlstr)

		connMaxLifetime := 59 * time.Second
		maxIdleConns := 5
		maxOpenConns := 15

		errorlog.LogPrintf("connMaxLifetime:%d\n", connMaxLifetime/time.Second)
		db.SetConnMaxLifetime(time.Duration(connMaxLifetime))

		errorlog.LogPrintf("maxIdleConns:%d\n", maxIdleConns)
		db.SetMaxIdleConns(maxIdleConns)

		errorlog.LogPrintf("maxOpenConns:%d\n", maxOpenConns)
		db.SetMaxOpenConns(maxOpenConns)
		if err != nil {
			return nil, err
		}

		gameBDSQL.DB = db
	}

	return gameBDSQL.DB, nil

}

// Connect New connect
func connectLogDB() (db *sql.DB, err error) {
	if logDBSQL == nil {
		logDBSQL = new(sqlCLi)
		sqlstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&timeout=30s", data.DBUser, data.DBPassword, data.DBIP, data.DBPORT, "logdb")
		errorlog.LogPrintln("DB Connect:", sqlstr)
		db, err := sql.Open("mysql", sqlstr)

		connMaxLifetime := 59 * time.Second
		maxIdleConns := 5
		maxOpenConns := 15

		errorlog.LogPrintf("connMaxLifetime:%d second\n", connMaxLifetime/time.Second)
		db.SetConnMaxLifetime(time.Duration(connMaxLifetime))

		errorlog.LogPrintf("maxIdleConns:%d\n", maxIdleConns)
		db.SetMaxIdleConns(maxIdleConns)

		errorlog.LogPrintf("maxOpenConns:%d\n", maxOpenConns)
		db.SetMaxOpenConns(maxOpenConns)
		if err != nil {
			return nil, err
		}

		logDBSQL.DB = db
	}

	return logDBSQL.DB, nil
}

// Connect New connect
func connectPayDB() (db *sql.DB, err error) {
	if payDBSQL == nil {
		payDBSQL = new(sqlCLi)
		sqlstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&timeout=30s", data.DBUser, data.DBPassword, data.DBIP, data.DBPORT, "paydb")
		errorlog.LogPrintln("DB Connect:", sqlstr)
		db, err := sql.Open("mysql", sqlstr)

		connMaxLifetime := 59 * time.Second
		maxIdleConns := 3
		maxOpenConns := 6

		errorlog.LogPrintf("connMaxLifetime:%d second\n", connMaxLifetime/time.Second)
		db.SetConnMaxLifetime(time.Duration(connMaxLifetime))

		errorlog.LogPrintf("maxIdleConns:%d\n", maxIdleConns)
		db.SetMaxIdleConns(maxIdleConns)

		errorlog.LogPrintf("maxOpenConns:%d\n", maxOpenConns)
		db.SetMaxOpenConns(maxOpenConns)
		if err != nil {
			return nil, err
		}

		payDBSQL.DB = db
	}

	return payDBSQL.DB, nil
}

// CallRead call stored procedure
func CallRead(name string, args ...interface{}) ([]interface{}, errorlog.ErrorMsg) {
	QueryStr := makeProcedureQueryStr(name, len(args))
	request, err := query(QueryStr, args...)
	return request, err
}

// CallReadOutMap call stored procedure
func CallReadOutMap(DB *sql.DB, name string, args ...interface{}) ([]map[string]interface{}, errorlog.ErrorMsg) {
	QueryStr := makeProcedureQueryStr(name, len(args))
	request, err := queryOutMap(DB, QueryStr, args...)
	return request, err

}

// CallWrite ...
func CallWrite(DB *sql.DB, name string, args ...interface{}) (sql.Result, errorlog.ErrorMsg) {
	request, err := exec(DB, name, args...)
	return request, err
}

// Query Use to SELECT return array, first is Keys
func query(query string, args ...interface{}) ([]interface{}, errorlog.ErrorMsg) {
	err := errorlog.New()

	res, errMsg := gameBDSQL.DB.Query(query, args...)
	if errMsg != nil {
		err.ErrorCode = code.FailedPrecondition
		err.Msg = "DBExecFail"
		errorlog.ErrorLogPrintln("DB", errMsg, query)
		return nil, err
	}

	request := makeScanArray(res)
	defer res.Close()

	return request, err
}

// QueryOutMap Use to SELECT return array map
func queryOutMap(DB *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, errorlog.ErrorMsg) {
	err := errorlog.New()

	res, errMsg := DB.Query(query, args...)
	if errMsg != nil {
		err.ErrorCode = code.FailedPrecondition
		err.Msg = "DBExecFail"
		errorlog.ErrorLogPrintln("DB", errMsg, query)
		return nil, err
	}

	request := makeScanMap(res)
	defer res.Close()

	return request, err
}

// Exec Use to INSTER, UPDATE, DELETE
func exec(DB *sql.DB, query string, args ...interface{}) (sql.Result, errorlog.ErrorMsg) {
	err := errorlog.New()

	res, errMsg := DB.Exec(query, args...)
	if errMsg != nil {
		err.ErrorCode = code.FailedPrecondition
		err.Msg = "DBExecFail"
		errorlog.ErrorLogPrintln("DB", errMsg, query)
		return nil, err
	}
	return res, err
}

func makeProcedureQueryStr(name string, args int) string {
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
	}

	return Result
}
