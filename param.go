package router

import (
	"net/http"
	"strconv"
)

// Param type is used to add key values to http.Request's Context.
type Param string

// GetParam is a convenience function to retrieve a route param from the current request.
// It just wraps req.Context().Value(Param(key))
func GetParam(req *http.Request, key string) interface{} {
	return req.Context().Value(Param(key))
}

// GetString wraps GetParam and returns the value as a string type.
func GetString(req *http.Request, key string) string {
	v := GetParam(req, key)
	if v != nil {
		return v.(string)
	}

	return ""
}

// GetInt wraps GetParam and returns the value as a int type.
// An error is returned if the int value can't be parsed
func GetInt(req *http.Request, key string) (int, error) {
	v := GetParam(req, key)
	return strconv.Atoi(v.(string))
}

// GetInt64 wraps GetParam and returns the value as a int64 type.
// An error is returned if the int value can't be parsed
func GetInt64(req *http.Request, key string) (int64, error) {
	v := GetParam(req, key)
	return strconv.ParseInt(v.(string), 10, 64)
}

// GetUint64 wraps GetParam and returns the value as a uint64 type.
// An error is returned if the int value can't be parsed
func GetUint64(req *http.Request, key string) (uint64, error) {
	v := GetParam(req, key)
	return strconv.ParseUint(v.(string), 10, 64)
}

// GetFloat64 wraps GetParam and returns the value as a float64 type.
// An error is returned if the int value can't be parsed
func GetFloat64(req *http.Request, key string) (float64, error) {
	v := GetParam(req, key)
	return strconv.ParseFloat(v.(string), 64)
}

// GetBool wraps GetParam and returns the value as a bool type.
// An error is returned if the int value can't be parsed
func GetBool(req *http.Request, key string) (bool, error) {
	v := GetParam(req, key)
	return strconv.ParseBool(v.(string))
}
