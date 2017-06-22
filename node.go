package router

import (
	"context"
	"net/http"
	"path"
	"strings"
)

// node represents each path part in a route and constructs a tree
type node struct {
	path        string
	handler     http.Handler
	parent      *node
	children    []*node
	childParams []*node
}

// rootNode is a helper function to initialize the root "/" node for any tree.
func rootNode(route string, handler http.Handler) *node {
	n := &node{
		path:        "/",
		children:    make([]*node, 0),
		childParams: make([]*node, 0),
	}

	n.add(route, handler)

	return n
}

// add construct the children tree for the current node matching the route provided.
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
				path:        p,
				children:    make([]*node, 0),
				childParams: make([]*node, 0),
				parent:      nn,
			}

			// Add children
			ch.add(strings.Join(parts[i+1:], "/"), handler)

			// Save on route or param tree
			if ch.path[0] == ':' {
				nn.childParams = append(nn.childParams, ch)
			} else {
				nn.children = append(nn.children, ch)
			}

			return
		}
	}

	// At last, set the handler
	nn.handler = handler
}

// walk moves through nodes for a given path until no further match is found.
func (n *node) walk(route string) (*node, string) {
	route = cleanupRoute(route)
	parts := strings.Split(route, "/")

	// Loop through route parts
	for i, p := range parts {
		// Skip empty parts, usually the first
		if p == "" {
			continue
		}

		// Stop at catch-all route
		if p == "*" {
			return n, route
		}

		// Look for a match on existing tree
		for _, ch := range n.children {
			if ch.path == p {
				return ch.walk(strings.Join(parts[i+1:], "/"))
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

	// Check if route is root route
	if cleanupRoute(r.URL.EscapedPath()) == "/" {
		return n.handler
	}

	// Split parts
	parts := strings.Split(cleanupRoute(r.URL.EscapedPath()), "/")

	// Look for children match, skip first "/"
	return n.matchChild(parts[1:], r)
}

// matchChild does the recursive work of matching the tree parts and try to find the correct path for a route.
func (n *node) matchChild(parts []string, r *http.Request) http.Handler {
	// Invalid route parts
	if parts == nil || len(parts) == 0 {
		return nil
	}

	// First look for exact match
	for _, ch := range n.children {
		if parts[0] == ch.path {
			// Last one
			if len(parts) == 1 {
				return ch.handler
			}

			// Go deeper
			h := ch.matchChild(parts[1:], r)
			if h != nil {
				return h
			}
		}
	}

	// Now look for params match
	for _, ch := range n.childParams {
		*r = *r.WithContext(
			context.WithValue(
				r.Context(),
				Param(strings.TrimPrefix(ch.path, ":")),
				parts[0]))

		// Last one
		if len(parts) == 1 {
			return ch.handler
		} else {
			// Go deeper
			h := ch.matchChild(parts[1:], r)
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

	// No match found
	return nil
}

// cleanupRoute ensures proper routes definition formatting
func cleanupRoute(route string) string {
	// Ensure route starts with /
	if route == "" || route[0] != '/' {
		route = path.Join("/", route)
	} else {
		route = path.Clean(route)
	}

	// Remove trailing "/"
	for len(route) > 1 && route[len(route)-1] == '/' {
		route = strings.TrimSuffix(route, "/")
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
