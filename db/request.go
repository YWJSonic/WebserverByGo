package db

import (
	"fmt"
	"time"

	"../code"
	"../log"
	"../messagehandle/errorlog"
	"github.com/go-sql-driver/mysql"
)

// GetSetting get db setting
func GetSetting() {
	request, err := CallRead("SettingGet_Read")
	if err.ErrorCode != code.OK {
		panic(err.Msg)
	}

	fmt.Println(request)
}

// GetAccountInfo Check Account existence and get
func GetAccountInfo(Account string) ([]map[string]interface{}, errorlog.ErrorMsg) {
	result, err := CallReadOutMap("AccountGet_Read", Account)
	return result, err
}

// NewAccount Create new Account
func NewAccount(age ...interface{}) errorlog.ErrorMsg {
	_, err := CallWrite("AccountNew_Write", age...)
	return err
}

// UpdateAccount update
func UpdateAccount(age ...interface{}) errorlog.ErrorMsg {
	_, err := CallWrite("AccountSet_Update", age...)
	return err
}

// NewGameAccount ...
func NewGameAccount(args ...interface{}) (int64, errorlog.ErrorMsg) {
	QuertStr := "INSERT INTO gameaccount VALUE (NULL,"
	err := errorlog.New()
	if len(args) > 0 {
		for range args {
			QuertStr += "?,"
		}
		QuertStr = QuertStr[:len(QuertStr)-1]
	}
	QuertStr += ");"

	request, errMSg := exec(QuertStr, args...)
	if errMSg.ErrorCode != code.OK {
		fmt.Println(errMSg)
		panic(errMSg)
	}
	playerID, reserr := request.LastInsertId()
	if reserr != nil {
		fmt.Println("NewGameAccount", reserr)
	}
	// err := errorlog.New()
	return playerID, err
}

// GetPlayerInfoByGameAccount ...
func GetPlayerInfoByGameAccount(GameAccount string) ([]map[string]interface{}, errorlog.ErrorMsg) {
	result, err := CallReadOutMap("GameAccountGet_Read", GameAccount)
	return result, err
}

// GetPlayerInfoByPlayerID ...
func GetPlayerInfoByPlayerID(PlayerID int64) (interface{}, errorlog.ErrorMsg) {
	result, err := CallReadOutMap("GameAccountGet_Read", PlayerID)

	return result, err
}

func NewLogTable(TableName string) {
	TableName = time.Now().AddDate(0, 0, 1).Format("20060102")
	query := fmt.Sprintf("CREATE TABLE `%s` (`index` BIGINT NOT NULL AUTO_INCREMENT,`Account` VARCHAR(128) NOT NULL,`PlayerID` BIGINT NOT NULL,`Time` BIGINT NOT NULL,`ActivityEvent` INT NOT NULL,`IValue1` BIGINT NOT NULL DEFAULT 0,`IValue2` BIGINT NOT NULL DEFAULT 0,`IValue3` BIGINT NOT NULL DEFAULT 0,`SValue1` VARCHAR(128) NOT NULL DEFAULT '',`SValue2` VARCHAR(128) NOT NULL DEFAULT '',`SValue3` VARCHAR(128) NOT NULL DEFAULT '',`Msg` TEXT NOT NULL,PRIMARY KEY (`index`));", TableName)
	fmt.Println(query)
	_, err := logDBSQL.DB.Exec(query)

	if err != nil {
		mysqlerr := err.(*mysql.MySQLError)
		if mysqlerr.Number == 1050 { // Table alerady exists
			return
		}
		panic(err)
	}
}

func SetLog(loginfo log.Log) errorlog.ErrorMsg {
	_, err := CallWrite("AccountSet_Update",
		loginfo.Account,
		loginfo.PlayerID,
		loginfo.Time,
		loginfo.ActivityEvent,
		loginfo.IValue1,
		loginfo.IValue2,
		loginfo.IValue3,
		loginfo.SValue1,
		loginfo.SValue2,
		loginfo.SValue3,
		loginfo.Msg)
	return err
}
