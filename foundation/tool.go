package foundation

import (
	"fmt"
	"time"

	"gitlab.com/WeberverByGo/code"
	"gitlab.com/WeberverByGo/data"
	"gitlab.com/WeberverByGo/messagehandle/errorlog"
	"gitlab.com/WeberverByGo/mycache"
)

// NewAccount convert all plant account to server account
func NewAccount(plant, account string) string {
	return fmt.Sprintf("%s:%s", plant, account)
}

// NewGameAccount new game account
func NewGameAccount(account string) string {
	return MD5Code(data.AccountEncodeStr + account)
}

// NewToken ...
func NewToken(gameAccount string) string {
	return MD5Code(fmt.Sprintf("%s%d", gameAccount, ServerNowTime()))
}

// CheckToken Check Token func
func CheckToken(gameAccount, token string) errorlog.ErrorMsg {
	err := errorlog.New()
	ServerToken := mycache.GetToken(gameAccount)
	if ServerToken != token {
		err.ErrorCode = code.Unauthenticated
		err.Msg = "TokenError"
	}
	return err
}

// CheckGameType Check Game Type
func CheckGameType(gameTypeID string) errorlog.ErrorMsg {
	err := errorlog.New()
	if gameTypeID != data.GameTypeID {
		err.ErrorCode = code.GameTypeError
		err.Msg = "GameTypeError"
	}
	return err
}

// ServerNowTime Get now Unix time
func ServerNowTime() int64 {
	return time.Now().Unix()
}

// ServerNow Get now time
func ServerNow() time.Time {
	return time.Now()
}

// IsInclude ...
func IsInclude(target int, src []int) bool {
	for _, value := range src {
		if value == target {
			return true
		}
	}
	return false
}
