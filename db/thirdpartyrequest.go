package db

import "../messagehandle/errorlog"

// NewULGInfo gametoken, playerid
func NewULGInfoRow(args ...interface{}) errorlog.ErrorMsg {
	_, err := CallWrite(gameBDSQL.DB, makeProcedureQueryStr("ULGNew_Write", len(args)), args...)
	return err
}

// UpdateULGInfo gametoken ,totalwin, totallost ,checkout
func UpdateULGInfoRow(args ...interface{}) errorlog.ErrorMsg {
	_, err := CallWrite(gameBDSQL.DB, makeProcedureQueryStr("ULGSet_Update", len(args)), args...)
	return err
}

// GetULGInfo ...
func GetULGInfoRow(gametoken string) ([]map[string]interface{}, errorlog.ErrorMsg) {
	result, err := CallReadOutMap(gameBDSQL.DB, "ULGGet_Read", gametoken)
	return result, err
}
