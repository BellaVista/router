package router

import (
	"context"
	"net/http"
	"path"
	"strings"
)

// node represents each path part in a route and constructs a tree
type node struct {
	path     string
	handler  http.Handler
	parent   *node
	children []*node
}

// rootNode is a helper function to initialize the root "/" node for any tree.
func rootNode(route string, handler http.Handler) *node {
	n := &node{
		path:     "/",
		children: make([]*node, 0),
	}

	n.add(route, handler)

	return n
}

// add constructs the children tree for the current node matching the route provided.
// It sets the http.Handler to the final element.
func (n *node) add(route string, handler http.Handler) {
	// Lookup as far as possible
	nn, remain := n.walk(route)

	// Add pending parts if any and stop adding after catch-all
	if nn.path != "*" {
		parts := strings.Split(remain, "/")
		for i, p := range parts {
			if p == "" {
				continue
			}

			// Create child
			ch := &node{
				path:     p,
				children: make([]*node, 0),
				parent:   nn,
			}

			// Add children
			ch.add(strings.Join(parts[i+1:], "/"), handler)

			// Save route
			nn.children = append(nn.children, ch)

			return
		}
	}

	// At last, set the handler
	nn.handler = handler
}

// walk moves through nodes for a given path until no further match is found.
func (n *node) walk(route string) (*node, string) {
	route = cleanupRoute(route)

	// Loop through route parts
	for i := range route {
		if route[i] != '/' && i != len(route) {
			continue
		}

		// Skip empty parts, usually the first
		if route[:i] == "" {
			continue
		}

		// Stop at catch-all route
		if route[:i] == "*" {
			return n, route
		}

		// Look for a match on existing tree
		for _, ch := range n.children {
			if ch.path == route[:i] {
				return ch.walk(route[i+1:])
			}
		}
	}

	// If no match, return current node and path
	return n, route
}

// match searches for a matching route to the current request.
// If found, it adds the route params to the request context and return the corresponding handler.
func (n *node) match(r *http.Request) http.Handler {
	// Validate root node match
	if n.path != "/" {
		return nil
	}

	if r.URL.Path == "/" || r.URL.Path == "" {
		return n.handler
	}

	// Look for children match, skip first and last '/'
	i := 0
	for r.URL.Path[len(r.URL.Path)-1-i] == '/' {
		i++
	}

	// Create parameters storage
	params := make(map[string]string)

	// Get handler
	h := n.matchChild(r.URL.Path[1:len(r.URL.Path)-i], r, params)

	// Set params if needed
	if h != nil && len(params) > 0 {
		*r = *r.WithContext(context.WithValue(
			r.Context(),
			routeParamsKey{},
			params))
	}

	return h
}

// matchChild does the recursive work of matching the tree parts and try to find the correct path for a route.
func (n *node) matchChild(part string, r *http.Request, params map[string]string) http.Handler {
	// Invalid route parts
	if part == "" {
		return nil
	}

	// Remove trailing slashes
	for len(part) > 0 && part[len(part)-1] == '/' {
		n.matchChild(part[:len(part)-1], r, params)
	}

	// Split parts
	for i := range part {
		if part[i] != '/' && len(part) != (i+1) {
			continue
		}

		// Look for matches
		for _, ch := range n.children {
			// It's a parameter?
			if ch.path[0] == ':' {
				// Are we done?
				if len(part) == (i + 1) {
					// Set last param
					params[ch.path[1:]] = part[:i+1]

					return ch.handler
				}

				// Set param
				params[ch.path[1:]] = part[:i]

				// Go deeper
				h := ch.matchChild(part[i+1:], r, params)
				if h != nil {
					return h
				}
			}

			// Last route part
			if len(part) == (i + 1) {
				if part[:i+1] == ch.path {
					return ch.handler
				}
			}

			// Match current
			if part[:i] == ch.path {
				// Go deeper
				h := ch.matchChild(part[i+1:], r, params)
				if h != nil {
					return h
				}
			}
		}

		// Check for catch-all routes.
		for _, ch := range n.children {
			if ch.path == "*" {
				return ch.handler
			}
		}

		// No match found so far
		return nil
	}

	// No match found
	return nil
}

// cleanupRoute ensures proper routes definition formatting
func cleanupRoute(route string) string {
	// Ensure route starts with /
	if route == "" || route[0] != '/' {
		route = "/" + route
	}

	// Remove trailing "/"
	for len(route) > 1 && route[len(route)-1] == '/' {
		route = route[:len(route)-1]
	}

	return route
}

// buildPath returns a string representing the entire path
// from the root to the current node.
func (n *node) buildPath() string {
	if n.parent == nil {
		return n.path
	}

	return path.Join(n.parent.buildPath(), n.path)
}
