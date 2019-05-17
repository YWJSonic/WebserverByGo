package foundation

import (
	"github.com/julienschmidt/httprouter"
)

// ConnType ...
const (
	Client  = "cli"
	Backend = "back"
)

// RESTfulURL ...
type RESTfulURL struct {
	RequestType string
	URL         string
	Fun         httprouter.Handle
	ConnType    string
}
