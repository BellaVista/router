package router

import (
	"net/http"
	"path"
)

// Router implements the needed methods for the Dispatcher
// to be able to match and execute requests.
type Router interface {
	// Add takes a route path and a handler to store for further matching
	Add(path string, handler http.Handler)

	// Wrap takes a Middleware to wrap all handlers in order (from inside out) at router level.
	Wrap(Middleware)

	// Match checks if a request matches this router.
	// If so, adds the route parameters to the request context and returns the corresponding handler.
	// If route doesn't matches, the response is nil
	Match(*http.Request) http.Handler
}

// New creates a new Router with the provided prefix
func New(prefix string) Router {
	// Create router
	return &router{
		prefix:     prefix,
		tree:       rootNode("/", nil),
		middleware: make([]Middleware, 0),
	}
}

// router implements Router interface
type router struct {
	// Routes prefix for this router
	prefix string

	// Routes tree
	tree *node

	// Middlewares collection
	middleware []Middleware
}

func (r *router) Add(route string, h http.Handler) {
	r.tree.add(path.Join(r.prefix, route), h)
}

func (r *router) Wrap(m Middleware) {
	r.middleware = append(r.middleware, m)
}

func (r *router) Match(req *http.Request) http.Handler {
	h := r.tree.match(req)
	for _, m := range r.middleware {
		h = m(h)
	}

	return h
}
