package router

import (
	"net/http"
	"net/http/httptest"
	"strconv"
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
		res.Write([]byte("Hello " + GetParam(req, "name")))
	}))

	d := Route(r)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost/hello/joe", nil)

	d.ServeHTTP(res, req)

	if GetParam(req, "name") != "joe" {
		t.Error("Request should have the :name context param set to 'joe' after dispatch")
	}
}

func TestConcurrentDispatch(t *testing.T) {
	r := New("/test")
	r.Add("/one/:param", http.HandlerFunc(dhandler))
	r.Add("/two/:param", http.HandlerFunc(dhandler))

	d := Route(r)

	for i := 0; i < 1000; i++ {
		res1 := httptest.NewRecorder()
		res2 := httptest.NewRecorder()

		one, _ := http.NewRequest("GET", "http://localhost:8080/test/one/"+strconv.Itoa(i), nil)
		two, _ := http.NewRequest("GET", "http://localhost:8080/test/two/"+strconv.Itoa(i), nil)

		go d.ServeHTTP(res1, one)
		go d.ServeHTTP(res2, two)
	}
}
