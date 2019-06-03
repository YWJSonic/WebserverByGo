package guest

import (
	"gitlab.com/WeberverByGo/foundation"
	"gitlab.com/WeberverByGo/player"
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
func (g *GuestInfo) GameAccount() string {
	return foundation.NewGameAccount(string(g.Account))
}

func (g *GuestInfo) AccountType() int64 {
	return player.Guest
}
