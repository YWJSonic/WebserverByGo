package transmission

import (
	"github.com/julienschmidt/httprouter"
)

// RESTfulURL ...
type RESTfulURL struct {
	RequestType string
	URL         string
	Fun         httprouter.Handle
}
