package router

import (
	"context"
	"net/http"
	"path"
	"path/filepath"
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
	// Root and matches
	if route == n.path || n.path == "*" {
		n.handler = handler
		return
	}

	// Remove starting and trailing "/"
	for len(route) > 1 && route[0] == '/' {
		route = route[1:]
	}
	for len(route) > 1 && route[len(route)-1] == '/' {
		route = route[:len(route)-1]
	}

	// Lookup as far as possible
	nn, remain := n.walk(strings.Split(route, "/"))

	// Add pending parts if any and stop adding after catch-all
	if len(remain) > 0 {
		if nn.path != "*" {
			// Create child
			ch := &node{
				path:     remain[0],
				children: make([]*node, 0),
				parent:   nn,
			}

			// Go deeper
			if len(remain) > 1 {
				ch.add(strings.Join(remain[1:], "/"), handler)
			} else {
				ch.handler = handler
			}

			// Save route
			nn.children = append(nn.children, ch)
			return
		}
	}
}

// walk moves through nodes for a given path until no further match is found.
func (n *node) walk(parts []string) (*node, []string) {
	// End at empty
	if len(parts) == 0 {
		return n, parts
	}

	// Search for first match in children
	for _, ch := range n.children {
		if parts[0] == ch.path {
			return ch.walk(parts[1:])
		}
	}

	// If no match, return current node and path
	return n, parts
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

	// Create parameters storage
	params := make(map[string]string)

	// Cleanup path
	r.URL.Path = filepath.Clean(r.URL.Path)

	// Get handler
	h := n.matchChild(r.URL.Path[1:], r, params)

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
					// Set last param and return
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

// buildPath returns a string representing the entire path
// from the root to the current node.
func (n *node) buildPath() string {
	if n.parent == nil {
		return n.path
	}

	return path.Join(n.parent.buildPath(), n.path)
}
