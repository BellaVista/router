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

	// AddBefore inserts new pre-dispatch middleware at Dispatcher level to be used by all routers.
	AddBefore(u http.Handler)

	// AddAfter inserts new post-dispatch middleware
	AddAfter(u http.Handler)
}

// Build constructs a Dispatcher that implements http.Handler and will contain
// all routes defined in the Router objects passed as parameters.
func Build(routes ...Router) Dispatcher {
	d := &dispatcher{
		before: make([]http.Handler, 0),
		routes: make([]Router, len(routes)),
		after:  make([]http.Handler, 0),
	}

	for i, r := range routes {
		d.routes[i] = r
	}

	return d
}

type dispatcher struct {
	before []http.Handler
	routes []Router
	after  []http.Handler
}

// ServeHTTP implements http.Handler interface.
// Takes care of middleware execution and stops the request flow if at any point the Context is cancelled.
func (d *dispatcher) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Match
	for _, r := range d.routes {
		// Found
		if h := r.Match(req); h != nil {
			// Pre-dispatch middleware
			d.runMiddleware(w, req, d.before)

			// Dispatch
			if req.Context().Err() == nil {
				// Pre-route middleware
				d.runMiddleware(w, req, r.Before())

				// Check cancel
				if req.Context().Err() != nil {
					return
				}

				// Handler dispatch
				h.ServeHTTP(w, req)

				// Post-route middleware
				d.runMiddleware(w, req, r.After())
			}

			// Post-dispatch middleware
			d.runMiddleware(w, req, d.after)

			// Return at route match
			return
		}
	}

	// 404 Not Found
	return
}

func (d *dispatcher) runMiddleware(w http.ResponseWriter, req *http.Request, ms []http.Handler) {
	for _, m := range ms {
		// Stop at cancelled request context
		if req.Context().Err() != nil {
			return
		}

		// Run middleware
		m.ServeHTTP(w, req)
	}
}

func (d *dispatcher) Add(r Router) {
	d.routes = append(d.routes, r)
}

func (d *dispatcher) AddBefore(m http.Handler) {
	d.before = append(d.before, m)
}

func (d *dispatcher) AddAfter(m http.Handler) {
	d.after = append(d.after, m)
}
