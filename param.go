package router

import (
	"net/http"
)

type routeParamsKey struct{}

// Params returns a map[string]string containing all route parameters
func Params(req *http.Request) map[string]string {
	params := req.Context().Value(routeParamsKey{})
	if _, ok := params.(map[string]string); ok {
		return params.(map[string]string)
	}

	return nil
}

// Param is a convenience function to retrieve a route param from the current request.
func Param(req *http.Request, key string) string {
	params := Params(req)
	if params != nil {
		if v, ok := params[key]; ok {
			return v
		}
	}

	return ""
}
