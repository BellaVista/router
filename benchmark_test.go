package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello test!"))
}

func helloName(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello " + GetString(r, "name")))
}

func helloNames(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello " + GetString(r, "first-name") + GetString(r, "middle-name") + GetString(r, "last-name")))
}

func BenchmarkRootMatch(b *testing.B) {
	// Create route
	r := New("/")
	r.Add("/", http.HandlerFunc(hello))

	req, _ := http.NewRequest("GET", "http://test.com/", nil)

	for i := 0; i < b.N; i++ {
		r.Match(req)
	}
}

func BenchmarkParamMatch(b *testing.B) {
	// Create route
	r := New("/")
	r.Add("/hello/:name", http.HandlerFunc(helloName))

	req, _ := http.NewRequest("GET", "http://test.com/hello/joe", nil)

	for i := 0; i < b.N; i++ {
		r.Match(req)
	}
}

func BenchmarkMultiParamMatch(b *testing.B) {
	// Create route
	r := New("/")
	r.Add("/hello/:first-name/:middle-name/:last-name", http.HandlerFunc(helloName))

	req, _ := http.NewRequest("GET", "http://test.com/hello/joe/x/smith", nil)

	for i := 0; i < b.N; i++ {
		r.Match(req)
	}
}

func BenchmarkRootDispatch(b *testing.B) {
	// Create route
	r := New("/")
	r.Add("/", http.HandlerFunc(hello))
	d := Route(r)

	req, _ := http.NewRequest("GET", "http://test.com", nil)
	res := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		d.ServeHTTP(res, req)
	}
}
