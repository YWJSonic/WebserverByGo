package foundation

import (
	"github.com/julienschmidt/httprouter"
)

// RESTfulURL ...
type RESTfulURL struct {
	RequestType string
	URL         string
	Fun         httprouter.Handle
}
