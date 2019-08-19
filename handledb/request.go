package db

import (
	"database/sql"
	"fmt"
	"time"

	"gitlab.com/ServerUtility/code"
	"gitlab.com/ServerUtility/dbinfo"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/messagehandle"
	"gitlab.com/ServerUtility/mysql"
	crontab "gitlab.com/WeberverByGoGame8/handlecrontab"
	"gitlab.com/WeberverByGoGame8/serversetting"
)

var gameBDSQL *dbinfo.SqlCLi
var logDBSQL *dbinfo.SqlCLi
var payDBSQL *dbinfo.SqlCLi

// QueryLogChan channel for write log
var QueryLogChan chan string

// WriteGameChan channel for write game db
var WriteGameChan chan dbinfo.SqlQuary

// WritePayChan channel for write pay db
var WritePayChan chan dbinfo.SqlQuary

// UseChanQueue use other goruting to write/update.
var UseChanQueue = false

// Connect New connect
func connectGameDB() (db *sql.DB, err error) {
	if gameBDSQL == nil {
		gameBDSQL = new(dbinfo.SqlCLi)
		sqlstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&timeout=30s", serversetting.DBUser, serversetting.DBPassword, serversetting.DBIP, serversetting.DBPORT, "gamedb")
		messagehandle.LogPrintln("DB Connect:", sqlstr)
		db, err := sql.Open("mysql", sqlstr)

		connMaxLifetime := 59 * time.Second
		maxIdleConns := 50
		maxOpenConns := 50

		messagehandle.LogPrintf("connMaxLifetime:%d\n", connMaxLifetime/time.Second)
		db.SetConnMaxLifetime(time.Duration(connMaxLifetime))

		messagehandle.LogPrintf("maxIdleConns:%d\n", maxIdleConns)
		db.SetMaxIdleConns(maxIdleConns)

		messagehandle.LogPrintf("maxOpenConns:%d\n", maxOpenConns)
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
		logDBSQL = new(dbinfo.SqlCLi)
		sqlstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&timeout=30s", serversetting.DBUser, serversetting.DBPassword, serversetting.DBIP, serversetting.DBPORT, "logdb")
		messagehandle.LogPrintln("DB Connect:", sqlstr)
		db, err := sql.Open("mysql", sqlstr)

		connMaxLifetime := 59 * time.Second
		maxIdleConns := 50
		maxOpenConns := 50

		messagehandle.LogPrintf("connMaxLifetime:%d second\n", connMaxLifetime/time.Second)
		db.SetConnMaxLifetime(time.Duration(connMaxLifetime))

		messagehandle.LogPrintf("maxIdleConns:%d\n", maxIdleConns)
		db.SetMaxIdleConns(maxIdleConns)

		messagehandle.LogPrintf("maxOpenConns:%d\n", maxOpenConns)
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
		payDBSQL = new(dbinfo.SqlCLi)
		sqlstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&timeout=30s", serversetting.DBUser, serversetting.DBPassword, serversetting.DBIP, serversetting.DBPORT, "paydb")
		messagehandle.LogPrintln("DB Connect:", sqlstr)
		db, err := sql.Open("mysql", sqlstr)

		connMaxLifetime := 59 * time.Second
		maxIdleConns := 50
		maxOpenConns := 50

		messagehandle.LogPrintf("connMaxLifetime:%d second\n", connMaxLifetime/time.Second)
		db.SetConnMaxLifetime(time.Duration(connMaxLifetime))

		messagehandle.LogPrintf("maxIdleConns:%d\n", maxIdleConns)
		db.SetMaxIdleConns(maxIdleConns)

		messagehandle.LogPrintf("maxOpenConns:%d\n", maxOpenConns)
		db.SetMaxOpenConns(maxOpenConns)
		if err != nil {
			return nil, err
		}

		payDBSQL.DB = db
	}

	return payDBSQL.DB, nil
}

// SetDBConn init value
func SetDBConn() {
	QueryLogChan = make(chan string)
	WriteGameChan = make(chan dbinfo.SqlQuary)
	WritePayChan = make(chan dbinfo.SqlQuary)
	connectGameDB()
	connectLogDB()
	connectPayDB()

	// server start check today log table.
	NewLogTable(foundation.ServerNow().Format("20060102"))
	NewLogTable(foundation.ServerNow().AddDate(0, 0, 1).Format("20060102"))

	// set Schedule check next day log table.
	crontab.NewCronBaseJob("35 15 * * *",
		crontab.NewLogCrontab(
			func() string { return foundation.ServerNow().AddDate(0, 0, 1).Format("20060102") },
			NewLogTable))
}

// SQLSelect channel loop
func SQLSelect() {
	select {
	case dbLogQuery := <-QueryLogChan:
		dbinfo.CallWrite(logDBSQL.DB, dbLogQuery)
	case dbgameQuery := <-WriteGameChan:
		dbinfo.CallWrite(gameBDSQL.DB, dbgameQuery.Quary, dbgameQuery.Args...)
	case dbpayQuary := <-WritePayChan:
		dbinfo.CallWrite(payDBSQL.DB, dbpayQuary.Quary, dbpayQuary.Args...)
	}
}

// GetSetting get db setting
func GetSetting() ([]map[string]interface{}, messagehandle.ErrorMsg) {
	result, err := dbinfo.CallReadOutMap(gameBDSQL.DB, "SettingGet_Read")
	return result, err
}

// GetSettingKey get db setting
func GetSettingKey(key string) ([]map[string]interface{}, messagehandle.ErrorMsg) {
	result, err := dbinfo.CallReadOutMap(gameBDSQL.DB, "SettingKeyGet_Read", key)
	return result, err
}

// UpdateSetting ...
func UpdateSetting(args ...interface{}) messagehandle.ErrorMsg {
	_, err := dbinfo.CallWrite(gameBDSQL.DB, dbinfo.MakeProcedureQueryStr("SettingSet_Update", len(args)), args...)
	return err
}

// GetAttachTypeRange ...
func GetAttachTypeRange(playerid, kind, miniAttType, maxAttType int64) ([]map[string]interface{}, messagehandle.ErrorMsg) {
	result, err := dbinfo.CallReadOutMap(gameBDSQL.DB, "AttachTypeRangeGet_Read", playerid, kind, miniAttType, maxAttType)
	return result, err
}

// GetAttachType ...
func GetAttachType(playerid int64, kind int64, attType int64) ([]map[string]interface{}, messagehandle.ErrorMsg) {
	result, err := dbinfo.CallReadOutMap(gameBDSQL.DB, "AttachTypeGet_Read", playerid, kind, attType)
	return result, err
}

// GetAttachKind get db attach kind
func GetAttachKind(playerid int64, kind int64) ([]map[string]interface{}, messagehandle.ErrorMsg) {
	result, err := dbinfo.CallReadOutMap(gameBDSQL.DB, "AttachKindGet_Read", playerid, kind)
	return result, err
}

// NewAttach ...
func NewAttach(args ...interface{}) {
	qu := dbinfo.SqlQuary{
		Quary: dbinfo.MakeProcedureQueryStr("AttachNew_Write", len(args)),
		Args:  args,
	}

	if UseChanQueue {
		go func() { WriteGameChan <- qu }()
	} else {
		dbinfo.CallWrite(gameBDSQL.DB, qu.Quary, qu.Args...)
	}
}

// UpdateAttach ...
func UpdateAttach(args ...interface{}) messagehandle.ErrorMsg {
	_, err := dbinfo.CallWrite(gameBDSQL.DB, dbinfo.MakeProcedureQueryStr("AttachSet_Write", len(args)), args...)
	return err
}

// GetAccountInfo Check Account existence and get
func GetAccountInfo(account string) ([]map[string]interface{}, messagehandle.ErrorMsg) {
	result, err := dbinfo.CallReadOutMap(gameBDSQL.DB, "AccountGet_Read", account)
	return result, err
}

// NewAccount new goruting set new Account
func NewAccount(args ...interface{}) { //messagehandle.ErrorMsg {
	qu := dbinfo.SqlQuary{
		Quary: dbinfo.MakeProcedureQueryStr("AccountNew_Write", len(args)),
		Args:  args,
	}

	if UseChanQueue {
		go func() { WriteGameChan <- qu }()
	} else {
		dbinfo.CallWrite(gameBDSQL.DB, qu.Quary, qu.Args...)
	}
}

// UpdateAccount update
func UpdateAccount(args ...interface{}) messagehandle.ErrorMsg {
	_, err := dbinfo.CallWrite(gameBDSQL.DB, dbinfo.MakeProcedureQueryStr("AccountSet_Update", len(args)), args...)
	return err
}

// NewGameAccount gameaccount, money, gametoken
func NewGameAccount(args ...interface{}) (int64, messagehandle.ErrorMsg) {
	QuertStr := "INSERT INTO gameaccount VALUE (NULL,"
	if len(args) > 0 {
		for range args {
			QuertStr += "?,"
		}
		QuertStr = QuertStr[:len(QuertStr)-1]
	}
	QuertStr += ");"

	request, err := dbinfo.Exec(gameBDSQL.DB, QuertStr, args...)
	if err.ErrorCode != code.OK {
		err.ErrorCode = code.FailedPrecondition
		err.Msg = "NewGameAccountError"
		messagehandle.ErrorLogPrintln("DB", err.Msg, QuertStr)
		return -1, err
	}
	playerID, errMsg := request.LastInsertId()
	if errMsg != nil {
		messagehandle.ErrorLogPrintln("DB NewGameAccount", errMsg)
	}
	// err := messagehandle.New()
	return playerID, err
}

// GetPlayerInfoByGameAccount ...
func GetPlayerInfoByGameAccount(gameAccount string) ([]map[string]interface{}, messagehandle.ErrorMsg) {
	result, err := dbinfo.CallReadOutMap(gameBDSQL.DB, "GameAccountGet_Read", gameAccount)
	return result, err
}

// GetPlayerInfoByPlayerID ...
func GetPlayerInfoByPlayerID(playerID int64) (interface{}, messagehandle.ErrorMsg) {
	result, err := dbinfo.CallReadOutMap(gameBDSQL.DB, "GameAccountGet_Read", playerID)
	return result, err
}

// UpdatePlayerInfo ...
func UpdatePlayerInfo(args ...interface{}) {
	qu := dbinfo.SqlQuary{
		Quary: dbinfo.MakeProcedureQueryStr("GameAccountSet_Update", len(args)),
		Args:  args,
	}

	if UseChanQueue {
		go func() { WriteGameChan <- qu }()
	} else {
		dbinfo.CallWrite(gameBDSQL.DB, qu.Quary, qu.Args...)
	}
}

// third party request

// NewULGInfoRow gametoken, playerid
func NewULGInfoRow(args ...interface{}) messagehandle.ErrorMsg {
	_, err := dbinfo.CallWrite(gameBDSQL.DB, dbinfo.MakeProcedureQueryStr("ULGNew_Write", len(args)), args...)
	return err
}

// UpdateULGInfoRow gametoken ,totalwin, totallost ,checkout
func UpdateULGInfoRow(args ...interface{}) messagehandle.ErrorMsg {
	_, err := dbinfo.CallWrite(gameBDSQL.DB, dbinfo.MakeProcedureQueryStr("ULGSet_Update", len(args)), args...)
	return err
}

// GetULGInfoRow ...
func GetULGInfoRow(gametoken string) ([]map[string]interface{}, messagehandle.ErrorMsg) {
	result, err := dbinfo.CallReadOutMap(gameBDSQL.DB, "ULGGet_Read", gametoken)
	return result, err
}

// UpdateCheckUlgRow ...
func UpdateCheckUlgRow(gametoken string) messagehandle.ErrorMsg {
	_, err := dbinfo.CallWrite(gameBDSQL.DB, dbinfo.MakeProcedureQueryStr("ULGCheckout_Update", 2), gametoken, true)
	return err
}

// ULGMaintainCheckoutRow ...
func ULGMaintainCheckoutRow() ([]map[string]interface{}, messagehandle.ErrorMsg) {
	result, err := dbinfo.CallReadOutMap(gameBDSQL.DB, "ULGMaintainCheckoutGet_Read")
	return result, err
}

// ULGMaintainCheckOutUpdate ...
func ULGMaintainCheckOutUpdate() messagehandle.ErrorMsg {
	_, err := dbinfo.CallReadOutMap(gameBDSQL.DB, "ULGMaintainCheckOutSet_Update")
	return err

}

/////////////////		Log DB		/////////////////

// NewLogTable Create new LogTable if table alerady exists return FailedPrecondition Error
func NewLogTable(tableName string) {
	query := fmt.Sprintf("CREATE TABLE `%s` (`index` BIGINT NOT NULL AUTO_INCREMENT,`Account` VARCHAR(128) NOT NULL,`PlayerID` BIGINT NOT NULL,`Time` BIGINT NOT NULL,`ActivityEvent` INT NOT NULL,`IValue1` BIGINT NOT NULL DEFAULT 0,`IValue2` BIGINT NOT NULL DEFAULT 0,`IValue3` BIGINT NOT NULL DEFAULT 0,`SValue1` VARCHAR(128) NOT NULL DEFAULT '',`SValue2` VARCHAR(128) NOT NULL DEFAULT '',`SValue3` VARCHAR(128) NOT NULL DEFAULT '',`Msg` TEXT NOT NULL,PRIMARY KEY (`index`));", tableName)
	_, errMsg := logDBSQL.DB.Exec(query)
	err := messagehandle.New()

	messagehandle.LogPrintln("DB NewLogTable", tableName)
	if errMsg != nil {
		mysqlerr := errMsg.(*mysql.MySQLError)
		if mysqlerr.Number == 1050 { // Table alerady exists
			return
		}
		err.ErrorCode = code.FailedPrecondition
		err.Msg = "NewLogTableError"
		messagehandle.ErrorLogPrintln("DB", err, query)
	}
}

// SetLog new goruting set log
func SetLog(account string, playerID, time int64, activityEvent int, iValue1, iValue2, iValue3 int64, sValue1, sValue2, sValue3, msg string) {
	tableName := foundation.ServerNow().Format("20060102")
	query := fmt.Sprintf("INSERT INTO `%s` VALUE(NULL,\"%s\",%d,%d, %d, %d,%d,%d,\"%s\",\"%s\",\"%s\",\"%s\");", tableName, account, playerID, time, activityEvent, iValue1, iValue2, iValue3, sValue1, sValue2, sValue3, msg)

	if UseChanQueue {
		go func() { QueryLogChan <- query }()
	} else {
		dbinfo.CallWrite(logDBSQL.DB, query)
	}
}

/////////////////		Pay DB		////////////////

// SetExchange new goruting set exchange log
func SetExchange(args ...interface{}) {
	qu := dbinfo.SqlQuary{
		Quary: dbinfo.MakeProcedureQueryStr("ExchangeNew_Write", len(args)),
		Args:  args,
	}

	if UseChanQueue {
		go func() { WritePayChan <- qu }()
	} else {
		dbinfo.CallWrite(payDBSQL.DB, qu.Quary, qu.Args...)
	}
}
