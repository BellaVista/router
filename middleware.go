package router

import (
	"net/http"
)

// Middleware type defines the function signature for middleware implementation
type Middleware func(http.Handler) http.Handler
