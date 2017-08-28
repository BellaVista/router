package router

import (
	"net/http"
	"testing"
)

func paramHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello test!"))
}

func TestGetParam(t *testing.T) {
	r := New("/")
	r.Add("/:param", http.HandlerFunc(paramHandler))

	req, _ := http.NewRequest("GET", "http://example.com/value", nil)
	h := r.Match(req)
	if h == nil {
		t.Errorf("%s should have matched our routes", "http://example.com/value")
	} else if GetParam(req, "param") != "value" {
		t.Errorf("Param :param should be set to 'value'. Got %s", GetParam(req, "param"))
	}
}

func TestGetWrongParam(t *testing.T) {
	r := New("/")
	r.Add("/:param", http.HandlerFunc(paramHandler))

	req, _ := http.NewRequest("GET", "http://example.com/value", nil)
	h := r.Match(req)
	if h == nil {
		t.Fatalf("%s should have matched our routes", "http://example.com/value")
	} else if GetParam(req, "invalid") != "" {
		t.Errorf("Param :invalid should be set to ''. Got %v", GetParam(req, "invalid"))
	}

	if GetParam(req, "invalid") != "" {
		t.Errorf("GetParam for :invalid should have been ''. Got %s", GetParam(req, "invalid"))
	}
}
