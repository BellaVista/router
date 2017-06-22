package router

import (
	"net/http"
)

// Dispatcher is constructed by Route() and works as a replacement
// for http.Handler to be used on any http.Server
type Dispatcher interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// Route constructs a Dispatcher that implements http.Handler and will contain
// all routes defined in the Router objects passed as parameters.
func Route(routes ...Router) Dispatcher {
	d := &dispatcher{
		routes: make([]Router, len(routes)),
	}

	for i, r := range routes {
		d.routes[i] = r
	}

	return d
}

type dispatcher struct {
	routes []Router
}

func (d *dispatcher) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, r := range d.routes {
		if h := r.Match(req); h != nil {
			// Found
			h.ServeHTTP(w, req)
			return
		}
	}

	// 404 Not Found
}
