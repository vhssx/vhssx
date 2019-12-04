package shortener

import (
	"net/http"
	"strings"
)

// Domain Configures
// Fallthrough or not.
// Domain -+> Prefix Checking -+> Route Mapping -+> Destination
// Domain -+> Route Mapping -+> Destination

// Get record by route.
type RouteMapper = map[string]*Record

type Domain struct {
	// Do further research and support later.
	Prefixes map[string]RouteMapper

	Routes RouteMapper
	// By default, yes!
	// Do fallthrough to ordinary static server.
	Fallthrough bool `json:"fallthrough"`
}

func (m *Domain) GetRecord(path string) *Record {
	if m.Routes[path] != nil {
		return m.Routes[path]
	}
	path = strings.ToLower(path)
	record := m.Routes[path]
	if record != nil && record.CaseInsensitive {
		// Case-insensitively matched.
		return record
	}
	return nil
}

func (m *Domain) ServeHttp(w http.ResponseWriter, req *http.Request) bool {
	record := m.GetRecord(req.URL.Path)
	if record == nil {
		if !m.Fallthrough {
			// handle the 404
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Not supported redirection yet, enable fallthrough option please!"))
			return true
		}
		// Do Fallthrough
		return false
	}
	// Following the record.Code to send the code.
	// FIX-ME Following the record.Title to send custom content.
	http.Redirect(w, req, record.Target, http.StatusMovedPermanently)
	return true
}
