package router

import (
	"net/http"
)

// Dispatcher is constructed by Route() and works as a replacement
// for http.Handler to be used on any http.Server
type Dispatcher interface {
	// ServeHTTP implements http.Handler
	ServeHTTP(w http.ResponseWriter, r *http.Request)

	// Add inserts a Router to the end of the Dispatcher's queue
	Add(r Router)

	// Wrap takes a Middleware to wrap all handlers in order (from inside out) at dispatcher level.
	Wrap(Middleware)
}

// Build constructs a Dispatcher that implements http.Handler and will contain
// all routes defined in the Router objects passed as parameters.
func Build(routes ...Router) Dispatcher {
	d := &dispatcher{
		routes:     make([]Router, len(routes)),
		middleware: make([]Middleware, 0),
	}

	for i, r := range routes {
		d.routes[i] = r
	}

	return d
}

type dispatcher struct {
	routes     []Router
	middleware []Middleware
}

// ServeHTTP implements http.Handler interface.
// Takes care of middleware execution and stops the request flow if at any point the Context is cancelled.
func (d *dispatcher) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Match
	for _, r := range d.routes {
		// Found
		if h := r.Match(req); h != nil {
			// Add middleware
			for _, m := range d.middleware {
				h = m(h)
			}

			// Dispatch
			h.ServeHTTP(w, req)

			// Return at route match
			return
		}
	}

	// 404 Not Found
	return
}

func (d *dispatcher) Add(r Router) {
	d.routes = append(d.routes, r)
}

func (d *dispatcher) Wrap(m Middleware) {
	d.middleware = append(d.middleware, m)
}
