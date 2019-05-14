package db

import (
	"fmt"

	"../code"
	"../foundation"
	"../messagehandle/errorlog"
	"github.com/go-sql-driver/mysql"
)

// GetSetting get db setting
func GetSetting() {
	_, err := CallRead("SettingGet_Read")
	if err.ErrorCode != code.OK {
		panic(err.Msg)
	}
}

// GetAccountInfo Check Account existence and get
func GetAccountInfo(Account string) ([]map[string]interface{}, errorlog.ErrorMsg) {
	result, err := CallReadOutMap(gameBDSQL.DB, "AccountGet_Read", Account)
	return result, err
}

// NewAccount new goruting set new Account
func NewAccount(args ...interface{}) { //errorlog.ErrorMsg {
	qu := sqlQuary{
		Quary: makeProcedureQueryStr("AccountNew_Write", len(args)),
		Args:  args,
	}
	go func() {
		WriteGameChan <- qu
	}()
}

// UpdateAccount update
func UpdateAccount(args ...interface{}) errorlog.ErrorMsg {
	_, err := CallWrite(gameBDSQL.DB, makeProcedureQueryStr("AccountSet_Update", len(args)), args...)
	return err
}

// NewGameAccount ...
func NewGameAccount(args ...interface{}) (int64, errorlog.ErrorMsg) {
	QuertStr := "INSERT INTO gameaccount VALUE (NULL,"
	if len(args) > 0 {
		for range args {
			QuertStr += "?,"
		}
		QuertStr = QuertStr[:len(QuertStr)-1]
	}
	QuertStr += ");"

	request, err := exec(gameBDSQL.DB, QuertStr, args...)
	if err.ErrorCode != code.OK {
		err.ErrorCode = code.FailedPrecondition
		err.Msg = "NewGameAccountError"
		errorlog.ErrorLogPrintln("DB", err.Msg, QuertStr)
		return -1, err
	}
	playerID, errMsg := request.LastInsertId()
	if errMsg != nil {
		errorlog.ErrorLogPrintln("DB NewGameAccount", errMsg)
	}
	// err := errorlog.New()
	return playerID, err
}

// GetPlayerInfoByGameAccount ...
func GetPlayerInfoByGameAccount(GameAccount string) ([]map[string]interface{}, errorlog.ErrorMsg) {
	result, err := CallReadOutMap(gameBDSQL.DB, "GameAccountGet_Read", GameAccount)
	return result, err
}

// GetPlayerInfoByPlayerID ...
func GetPlayerInfoByPlayerID(PlayerID int64) (interface{}, errorlog.ErrorMsg) {
	result, err := CallReadOutMap(gameBDSQL.DB, "GameAccountGet_Read", PlayerID)

	return result, err
}

// UpdatePlayerInfo ...
func UpdatePlayerInfo(args ...interface{}) {
	qu := sqlQuary{
		Quary: makeProcedureQueryStr("GameAccountSet_Update", len(args)),
		Args:  args,
	}
	go func() {
		WriteGameChan <- qu
	}()
}

/////////////////		Log DB		/////////////////

// NewLogTable Create new LogTable if table alerady exists return FailedPrecondition Error
func NewLogTable(TableName string) errorlog.ErrorMsg {
	query := fmt.Sprintf("CREATE TABLE `%s` (`index` BIGINT NOT NULL AUTO_INCREMENT,`Account` VARCHAR(128) NOT NULL,`PlayerID` BIGINT NOT NULL,`Time` BIGINT NOT NULL,`ActivityEvent` INT NOT NULL,`IValue1` BIGINT NOT NULL DEFAULT 0,`IValue2` BIGINT NOT NULL DEFAULT 0,`IValue3` BIGINT NOT NULL DEFAULT 0,`SValue1` VARCHAR(128) NOT NULL DEFAULT '',`SValue2` VARCHAR(128) NOT NULL DEFAULT '',`SValue3` VARCHAR(128) NOT NULL DEFAULT '',`Msg` TEXT NOT NULL,PRIMARY KEY (`index`));", TableName)
	_, errMsg := logDBSQL.DB.Exec(query)
	err := errorlog.New()

	if errMsg != nil {
		mysqlerr := errMsg.(*mysql.MySQLError)
		if mysqlerr.Number == 1050 { // Table alerady exists
			return err
		}
		err.ErrorCode = code.FailedPrecondition
		err.Msg = "NewLogTableError"
		errorlog.ErrorLogPrintln("DB", err, query)
		return err
	}

	return err
}

// SetLog new goruting set log
func SetLog(Account string, PlayerID, Time int64, ActivityEvent int, IValue1, IValue2, IValue3 int64, SValue1, SValue2, SValue3, Msg string) {
	TableName := foundation.ServerNow().Format("20060102")
	query := fmt.Sprintf("INSERT INTO `%s` VALUE(NULL,\"%s\",%d,%d, %d, %d,%d,%d,\"%s\",\"%s\",\"%s\",\"%s\");", TableName, Account, PlayerID, Time, ActivityEvent, IValue1, IValue2, IValue3, SValue1, SValue2, SValue3, Msg)

	go func() { QueryLogChan <- query }()
}

// ExecLog Exec Use to INSTER, UPDATE, DELETE
// func ExecLog(query string, args ...interface{}) {

// 	_, err := logDBSQL.DB.Exec(query, args...)
// 	if err != nil {
// 		errorlog.ErrorLogPrintln("DB", err, query)
// 	}
// }

/////////////////		Pay DB		////////////////

// SetExchange new goruting set exchange log
func SetExchange(args ...interface{}) {
	qu := sqlQuary{
		Quary: makeProcedureQueryStr("ExchangeNew_Write", len(args)),
		Args:  args,
	}
	go func() { WritePayChan <- qu }()
}

// ExecPay ...
// func ExecPay(query string, args ...interface{}) {

// 	_, err := payDBSQL.DB.Exec(query, args...)
// 	if err != nil {
// 		errorlog.ErrorLogPrintln("DB", err, query)
// 	}
// }
