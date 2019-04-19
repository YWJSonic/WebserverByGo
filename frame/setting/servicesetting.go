package setting

import "fmt"

// ServerOption Server Setting
type ServerOption struct {
	IP       string
	PORT     string
	Maintain bool
}

// URL ...
func (Sopt *ServerOption) URL() string {
	return fmt.Sprintf("%s:%s", Sopt.IP, Sopt.PORT)
}
