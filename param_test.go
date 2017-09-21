package router

import (
	"net/http"
	"testing"
)

func paramHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello test!"))
}

func TestParam(t *testing.T) {
	r := New("/")
	r.Add("/:param", http.HandlerFunc(paramHandler))

	req, _ := http.NewRequest("GET", "http://example.com/value", nil)
	h := r.Match(req)
	if h == nil {
		t.Errorf("%s should have matched our routes", "http://example.com/value")
	} else if Param(req, "param") != "value" {
		t.Errorf("Param :param should be set to 'value'. Got %s", Param(req, "param"))
	}
}

func TestGetWrongParam(t *testing.T) {
	r := New("/")
	r.Add("/:param", http.HandlerFunc(paramHandler))

	req, _ := http.NewRequest("GET", "http://example.com/value", nil)
	h := r.Match(req)
	if h == nil {
		t.Fatalf("%s should have matched our routes", "http://example.com/value")
	} else if Param(req, "invalid") != "" {
		t.Errorf("Param :invalid should be set to ''. Got %v", Param(req, "invalid"))
	}

	if Param(req, "invalid") != "" {
		t.Errorf("Param for :invalid should have been ''. Got %s", Param(req, "invalid"))
	}
}
