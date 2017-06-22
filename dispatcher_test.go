package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func dhandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello test!"))
}

func TestDispatcher(t *testing.T) {
	r1 := New("/v1/")
	r1.Add("/test", http.HandlerFunc(dhandler))
	r1.Add("/test/1/2/3/4/5/6", http.HandlerFunc(dhandler))

	r2 := New("/v2")
	r2.Add("/test", http.HandlerFunc(dhandler))
	r2.Add("/test/1/2/3/4/5/6", http.HandlerFunc(dhandler))

	d := Route(r1, r2)

	if len(d.(*dispatcher).routes) != 2 {
		t.Errorf("Route should have added 2 routes to dispatcher. Got %d", len(d.(*dispatcher).routes))
	}
}

func TestDispatcherParam(t *testing.T) {
	r := New("/")
	r.Add("/hello/:name", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Hello " + GetString(req, "name")))
	}))

	d := Route(r)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost/hello/joe", nil)

	d.ServeHTTP(res, req)

	if req.Context().Value(Param("name")).(string) != "joe" {
		t.Error("Request should have the :name context param set to 'joe' after dispatch")
	}
}
