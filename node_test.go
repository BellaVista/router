package router

import (
	"net/http"
	"testing"
)

type emptyHandler struct{}

func (e emptyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func TestAdd(t *testing.T) {
	root := rootNode("/", emptyHandler{})
	if root.path != "/" {
		t.Errorf("root.path should be '/'. Got %s", root.path)
	}

	root.add("/some/route/with/five/parts", emptyHandler{})
	if len(root.children) != 1 {
		for _, ch := range root.children {
			t.Errorf("Error data: %s", ch.path)
		}
		t.Fatalf("root.children should have 1 items. Got %d", len(root.children))
	}

	root.add("/test/action", emptyHandler{})
	if len(root.children) != 2 {
		for _, ch := range root.children {
			t.Errorf("Error data: %s", ch.path)
		}
		t.Fatalf("root.children should have 2 items. Got %d", len(root.children))
	}
}

func TestBuildPath(t *testing.T) {
	root := rootNode("/", emptyHandler{})
	if root.buildPath() != "/" {
		t.Errorf("root.buildPath() should be '/'. Got %s", root.buildPath())
	}

	root = rootNode("/test", emptyHandler{})
	if root.path != "/" {
		t.Errorf("root.path should be '/'. Got %s", root.path)
	}
	if root.buildPath() != "/" {
		t.Errorf("root.buildPath() should be '/'. Got %s", root.buildPath())
	}
	if len(root.children) != 1 {
		t.Fatalf("root.children should have 1 items. Got %d", len(root.children))
	}
	if root.children[0].buildPath() != "/test" {
		t.Errorf("root.children[0].buildPath() should be '/test'. Got %s", root.children[0].buildPath())
	}
	if root.children[0].path != "test" {
		t.Errorf("root.children[0].path should be '/test'. Got %s", root.children[0].path)
	}

	root.add("/test/action", emptyHandler{})
	if len(root.children) != 1 {
		t.Fatalf("root.children should have 1 items after adding another route with same prefix. Got %d", len(root.children))
	}
}
