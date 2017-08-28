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
	} else if GetParam(req, "param").(string) != "value" {
		t.Errorf("Param :param should be set to 'value'. Got %s", GetParam(req, "param").(string))
	}
}

func TestGetWrongParam(t *testing.T) {
	r := New("/")
	r.Add("/:param", http.HandlerFunc(paramHandler))

	req, _ := http.NewRequest("GET", "http://example.com/value", nil)
	h := r.Match(req)
	if h == nil {
		t.Fatalf("%s should have matched our routes", "http://example.com/value")
	} else if GetParam(req, "invalid") != nil {
		t.Errorf("Param :invalid should be set to 'nil'. Got %v", GetParam(req, "invalid"))
	}

	if GetString(req, "invalid") != "" {
		t.Errorf("GetString for :invalid should have been ''. Got %s", GetString(req, "invalid"))
	}
}

func TestGetString(t *testing.T) {
	r := New("/")
	r.Add("/:param", http.HandlerFunc(paramHandler))

	req, _ := http.NewRequest("GET", "http://example.com/value", nil)
	h := r.Match(req)
	if h == nil {
		t.Errorf("%s should have matched our routes", "http://example.com/value")
	} else if GetString(req, "param") != "value" {
		t.Errorf("Param :param should be set to 'value'. Got %s", GetString(req, "param"))
	}
}

func TestGetInt(t *testing.T) {
	r := New("/")
	r.Add("/:int", http.HandlerFunc(paramHandler))

	req, _ := http.NewRequest("GET", "http://example.com/42", nil)
	h := r.Match(req)
	if h == nil {
		t.Errorf("%s should have matched our routes", "http://example.com/42")
	} else {
		i, err := GetInt(req, "int")
		if err != nil {
			t.Fatal(err.Error())
		}
		if i != 42 {
			t.Errorf("Param :int should be set to 42. Got %d", i)
		}
	}
}

func TestGetInt64(t *testing.T) {
	r := New("/")
	r.Add("/:int", http.HandlerFunc(paramHandler))

	req, _ := http.NewRequest("GET", "http://example.com/42", nil)
	h := r.Match(req)
	if h == nil {
		t.Errorf("%s should have matched our routes", "http://example.com/42")
	} else {
		i, err := GetInt64(req, "int")
		if err != nil {
			t.Fatal(err.Error())
		}
		if i != 42 {
			t.Errorf("Param :int should be set to 42. Got %d", i)
		}
	}
}

func TestGetUint64(t *testing.T) {
	r := New("/")
	r.Add("/:int", http.HandlerFunc(paramHandler))

	req, _ := http.NewRequest("GET", "http://example.com/42", nil)
	h := r.Match(req)
	if h == nil {
		t.Errorf("%s should have matched our routes", "http://example.com/42")
	} else {
		i, err := GetUint64(req, "int")
		if err != nil {
			t.Fatal(err.Error())
		}
		if i != 42 {
			t.Errorf("Param :int should be set to 42. Got %d", i)
		}
	}
}

func TestGetFloat64(t *testing.T) {
	r := New("/")
	r.Add("/:float", http.HandlerFunc(paramHandler))

	req, _ := http.NewRequest("GET", "http://example.com/42.321", nil)
	h := r.Match(req)
	if h == nil {
		t.Errorf("%s should have matched our routes", "http://example.com/42.321")
	} else {
		f, err := GetFloat64(req, "float")
		if err != nil {
			t.Fatal(err.Error())
		}
		if f != 42.321 {
			t.Errorf("Param :float should be set to 42.321. Got %f", f)
		}
	}
}

func TestGetBool(t *testing.T) {
	r := New("/")
	r.Add("/:bool", http.HandlerFunc(paramHandler))

	req, _ := http.NewRequest("GET", "http://example.com/true", nil)
	h := r.Match(req)
	if h == nil {
		t.Errorf("%s should have matched our routes", "http://example.com/true")
	} else {
		b, err := GetBool(req, "bool")
		if err != nil {
			t.Fatal(err.Error())
		}
		if b != true {
			t.Errorf("Param :bool should be set to true. Got %t", b)
		}
	}
}
