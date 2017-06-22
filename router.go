package router

import (
	"net/http"
	"path"
)

// Router implements the needed methods for the Dispatcher
// to be able to match and execute requests.
type Router interface {
	// Add takes a route path and a handler to store for further matching
	Add(path string, hanlder http.Handler)

	// Match checks if a request matches this router.
	// If so, adds the route parameters to the request and returns the corresponding handler
	// If route doesn't matches, the response is nil
	Match(*http.Request) http.Handler
}

// New creates a new Router with the provided prefix
func New(prefix string) Router {
	// Create router
	return &router{
		prefix: prefix,
		tree:   rootNode("/", nil),
	}
}

// router implements Router interface
type router struct {
	// Routes prefix for this router
	prefix string

	// Routes tree
	tree *node
}

func (r *router) Add(route string, handler http.Handler) {
	r.tree.add(path.Join(r.prefix, route), handler)
}

func (r *router) Match(req *http.Request) http.Handler {
	return r.tree.match(req)
}
