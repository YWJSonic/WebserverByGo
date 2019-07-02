package guest

import (
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/playerinfo"
)

// GuestInfo ...
type GuestInfo struct {
	Account string
}

// PartyAccount ...
func (g *GuestInfo) PartyAccount() string {
	return g.Account
}

// GameAccount ...
func (g *GuestInfo) GameAccount(encodeStr string) string {
	return foundation.NewGameAccount(encodeStr, string(g.Account))
}

// AccountType ...
func (g *GuestInfo) AccountType() int64 {
	return playerinfo.Guest
}
