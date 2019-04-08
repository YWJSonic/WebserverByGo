package Account

type AccountInfo struct {
	ThirdPartyAccount string
	GameAccount       string
	PassWord          string
	LoginTime         int64 // Microsecond
	LogoutTime        int64 // Microsecond
}
