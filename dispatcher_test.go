package router

import (
	"net/http"
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
