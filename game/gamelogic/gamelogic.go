package gamelogic

import (
	"encoding/json"
	"fmt"

	"../../code"
	"../../foundation"
	"../../messagehandle/errorlog"
)

type LogicResult struct {
	TotalWinMoney   int64
	PerTypeWinItem  [][]int
	PerTypeWinMoney [][]int
}

// GetGameResult ...
func GetGameResult(thirdparty string, gametypeID, betmoney int64) (map[string]interface{}, errorlog.ErrorMsg) {
	err := errorlog.New()
	POSTvalues := map[string][]string{"POST": []string{fmt.Sprintf(`{"betmoney": %d, "gametypeID":%d,"thirdparty":%s}`, betmoney, gametypeID, thirdparty)}}
	reuslt := foundation.HTTPPostRequest("http://192.168.1.15:8100/slot2/gameresult", POSTvalues)

	var gamereuslt map[string]interface{}
	if errMsg := json.Unmarshal(reuslt, &gamereuslt); errMsg != nil {
		errorlog.ErrorLogPrintln("GameResult", errMsg, "[", reuslt, "]")
		err.ErrorCode = code.NoThisPlayer
		err.Msg = "PlayerFormatError"
		return nil, err
	}

	return gamereuslt, err

}
