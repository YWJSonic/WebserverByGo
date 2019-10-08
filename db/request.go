package db

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
	"gitlab.com/WeberverByGoGame5/code"
	"gitlab.com/WeberverByGoGame5/foundation"
	"gitlab.com/WeberverByGoGame5/messagehandle/errorlog"
)

// GetSetting get db setting
func GetSetting() ([]map[string]interface{}, errorlog.ErrorMsg) {
	result, err := CallReadOutMap(gameBDSQL.DB, "SettingGet_Read")
	return result, err
}

// GetSettingKey get db setting
func GetSettingKey(key string) ([]map[string]interface{}, errorlog.ErrorMsg) {
	result, err := CallReadOutMap(gameBDSQL.DB, "SettingKeyGet_Read", key)
	return result, err
}

// NewSetting ...
func NewSetting(args ...interface{}) {
	qu := sqlQuary{
		Quary: makeProcedureQueryStr("SettingNew_Write", len(args)),
		Args:  args,
	}

	if UseChanQueue {
		go func() { WriteGameChan <- qu }()
	} else {
		CallWrite(gameBDSQL.DB, qu.Quary, qu.Args...)
	}
}

// UpdateSetting ...
func UpdateSetting(args ...interface{}) errorlog.ErrorMsg {
	_, err := CallWrite(gameBDSQL.DB, makeProcedureQueryStr("SettomgSet_Update", len(args)), args...)
	return err
}

// GetAttachKind get db attach kind
func GetAttachKind(playerid int64, kind int64) ([]map[string]interface{}, errorlog.ErrorMsg) {
	result, err := CallReadOutMap(gameBDSQL.DB, "AttachKindGet_Read", playerid, kind)
	return result, err
}

// NewAttach ...
func NewAttach(args ...interface{}) {
	qu := sqlQuary{
		Quary: makeProcedureQueryStr("AttachNew_Write", len(args)),
		Args:  args,
	}

	if UseChanQueue {
		go func() { WriteGameChan <- qu }()
	} else {
		CallWrite(gameBDSQL.DB, qu.Quary, qu.Args...)
	}
}

// UpdateAttach ...
func UpdateAttach(args ...interface{}) errorlog.ErrorMsg {
	_, err := CallWrite(gameBDSQL.DB, makeProcedureQueryStr("AttachGame5Set_Update", len(args)), args...)
	return err
}

// GetAccountInfo Check Account existence and get
func GetAccountInfo(account string) ([]map[string]interface{}, errorlog.ErrorMsg) {
	result, err := CallReadOutMap(gameBDSQL.DB, "AccountGet_Read", account)
	return result, err
}

// NewAccount new goruting set new Account
func NewAccount(args ...interface{}) { //errorlog.ErrorMsg {
	qu := sqlQuary{
		Quary: makeProcedureQueryStr("AccountNew_Write", len(args)),
		Args:  args,
	}

	if UseChanQueue {
		go func() { WriteGameChan <- qu }()
	} else {
		CallWrite(gameBDSQL.DB, qu.Quary, qu.Args...)
	}
}

// UpdateAccount update
func UpdateAccount(args ...interface{}) errorlog.ErrorMsg {
	_, err := CallWrite(gameBDSQL.DB, makeProcedureQueryStr("AccountSet_Update", len(args)), args...)
	return err
}

// NewGameAccount gameaccount, money, gametoken
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
func GetPlayerInfoByGameAccount(gameAccount string) ([]map[string]interface{}, errorlog.ErrorMsg) {
	result, err := CallReadOutMap(gameBDSQL.DB, "GameAccountGet_Read", gameAccount)
	return result, err
}

// GetPlayerInfoByPlayerID ...
func GetPlayerInfoByPlayerID(playerID int64) (interface{}, errorlog.ErrorMsg) {
	result, err := CallReadOutMap(gameBDSQL.DB, "GameAccountGet_Read", playerID)
	return result, err
}

// UpdatePlayerInfo ...
func UpdatePlayerInfo(args ...interface{}) {
	qu := sqlQuary{
		Quary: makeProcedureQueryStr("GameAccountSet_Update", len(args)),
		Args:  args,
	}

	if UseChanQueue {
		go func() { WriteGameChan <- qu }()
	} else {
		CallWrite(gameBDSQL.DB, qu.Quary, qu.Args...)
	}
}

// third party request

// NewULGInfoRow gametoken, playerid
func NewULGInfoRow(args ...interface{}) errorlog.ErrorMsg {
	_, err := CallWrite(gameBDSQL.DB, makeProcedureQueryStr("ULGNew_Write", len(args)), args...)
	return err
}

// UpdateULGInfoRow gametoken ,totalwin, totallost ,checkout
func UpdateULGInfoRow(args ...interface{}) errorlog.ErrorMsg {
	_, err := CallWrite(gameBDSQL.DB, makeProcedureQueryStr("ULGSet_Update", len(args)), args...)
	return err
}

// GetULGInfoRow ...
func GetULGInfoRow(gametoken string) ([]map[string]interface{}, errorlog.ErrorMsg) {
	result, err := CallReadOutMap(gameBDSQL.DB, "ULGGet_Read", gametoken)
	return result, err
}

// UpdateCheckUlgRow ...
func UpdateCheckUlgRow(gametoken string) errorlog.ErrorMsg {
	_, err := CallWrite(gameBDSQL.DB, makeProcedureQueryStr("ULGCheckout_Update", 2), gametoken, true)
	return err
}

// ULGMaintainCheckoutRow ...
func ULGMaintainCheckoutRow() ([]map[string]interface{}, errorlog.ErrorMsg) {
	result, err := CallReadOutMap(gameBDSQL.DB, "ULGMaintainCheckoutGet_Read")
	return result, err
}

// ULGMaintainCheckOutUpdate ...
func ULGMaintainCheckOutUpdate() errorlog.ErrorMsg {
	_, err := CallReadOutMap(gameBDSQL.DB, "ULGMaintainCheckOutSet_Update")
	return err

}

/////////////////		Log DB		/////////////////

// NewLogTable Create new LogTable if table alerady exists return FailedPrecondition Error
func NewLogTable(tableName string) {
	query := fmt.Sprintf("CREATE TABLE `%s` (`index` BIGINT NOT NULL AUTO_INCREMENT,`Account` VARCHAR(128) NOT NULL,`PlayerID` BIGINT NOT NULL,`Time` BIGINT NOT NULL,`ActivityEvent` INT NOT NULL,`IValue1` BIGINT NOT NULL DEFAULT 0,`IValue2` BIGINT NOT NULL DEFAULT 0,`IValue3` BIGINT NOT NULL DEFAULT 0,`SValue1` VARCHAR(128) NOT NULL DEFAULT '',`SValue2` VARCHAR(128) NOT NULL DEFAULT '',`SValue3` VARCHAR(128) NOT NULL DEFAULT '',`Msg` TEXT NOT NULL,PRIMARY KEY (`index`));", tableName)
	_, errMsg := logDBSQL.DB.Exec(query)
	err := errorlog.New()

	errorlog.LogPrintln("DB NewLogTable", tableName)
	if errMsg != nil {
		mysqlerr := errMsg.(*mysql.MySQLError)
		if mysqlerr.Number == 1050 { // Table alerady exists
			return
		}
		err.ErrorCode = code.FailedPrecondition
		err.Msg = "NewLogTableError"
		errorlog.ErrorLogPrintln("DB", err, query)
	}
}

// SetLog new goruting set log
func SetLog(account string, playerID, time int64, activityEvent int, iValue1, iValue2, iValue3 int64, sValue1, sValue2, sValue3, msg string) {
	tableName := foundation.ServerNow().Format("20060102")
	query := fmt.Sprintf("INSERT INTO `%s` VALUE(NULL,\"%s\",%d,%d, %d, %d,%d,%d,\"%s\",\"%s\",\"%s\",\"%s\");", tableName, account, playerID, time, activityEvent, iValue1, iValue2, iValue3, sValue1, sValue2, sValue3, msg)

	if UseChanQueue {
		go func() { QueryLogChan <- query }()
	} else {
		CallWrite(logDBSQL.DB, query)
	}
}

/////////////////		Pay DB		////////////////

// SetExchange new goruting set exchange log
func SetExchange(args ...interface{}) {
	qu := sqlQuary{
		Quary: makeProcedureQueryStr("ExchangeNew_Write", len(args)),
		Args:  args,
	}

	if UseChanQueue {
		go func() { WritePayChan <- qu }()
	} else {
		CallWrite(payDBSQL.DB, qu.Quary, qu.Args...)
	}
}
