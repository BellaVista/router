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

	// AddBefore inserts a http.Handler to the middleware queue, before the request matching handler.
	AddBefore(handler http.Handler)

	// AddAfter inserts a http.Handler to the middleware queue, after the request matching handler.
	AddAfter(handler http.Handler)

	// Before returns all the middleware handlers that have to run before the request matching handler
	Before() []http.Handler

	// After returns all the middleware handlers that have to run after the request matching handler
	After() []http.Handler

	// Match checks if a request matches this router.
	// If so, adds the route parameters to the request context and returns the corresponding handler
	// If route doesn't matches, the response is nil
	Match(*http.Request) http.Handler
}

// New creates a new Router with the provided prefix
func New(prefix string) Router {
	// Create router
	return &router{
		prefix: prefix,
		tree:   rootNode("/", nil),
		before: make([]http.Handler, 0),
		after:  make([]http.Handler, 0),
	}
}

// router implements Router interface
type router struct {
	// Routes prefix for this router
	prefix string

	// Routes tree
	tree *node

	// Middleware
	before []http.Handler
	after  []http.Handler
}

func (r *router) Add(route string, h http.Handler) {
	r.tree.add(path.Join(r.prefix, route), h)
}

func (r *router) AddBefore(h http.Handler) {
	r.before = append(r.before, h)
}

func (r *router) AddAfter(h http.Handler) {
	r.after = append(r.after, h)
}

func (r *router) Before() []http.Handler {
	return r.before
}

func (r *router) After() []http.Handler {
	return r.after
}

func (r *router) Match(req *http.Request) http.Handler {
	return r.tree.match(req)
}
