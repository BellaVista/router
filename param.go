package router

import (
	"net/http"
)

type routeParamsKey struct{}

// GetParams returns a map[string]string containing all route parameters
func GetParams(req *http.Request) map[string]string {
	params := req.Context().Value(routeParamsKey{})
	if _, ok := params.(map[string]string); ok {
		return params.(map[string]string)
	}

	return nil
}

// GetParam is a convenience function to retrieve a route param from the current request.
// It just wraps req.Context().Value(Param(key))
func GetParam(req *http.Request, key string) string {
	params := GetParams(req)
	if params != nil {
		if v, ok := params[key]; ok {
			return v
		}
	}

	return ""
}
