package thirdparty

// thirdparty id
const (
	None = iota
	Self
	Ulg
)

// IPratyAccount thirdparty api interface
type IPratyAccount interface {
	PartyAccount() string
	GameAccount() string
}
