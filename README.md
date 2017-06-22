[![GoDoc](https://godoc.org/github.com/bellavista/router?status.svg)](https://godoc.org/github.com/bellavista/router)
[![Version](https://badge.fury.io/gh/bellavista%2Frouter.svg)](https://badge.fury.io/gh/bellavista%2Frouter)
[![Build Status](https://travis-ci.org/bellavista/router.svg?branch=master)](https://travis-ci.org/bellavista/router)
[![Coverage Status](https://coveralls.io/repos/github/bellavista/router/badge.svg?branch=master)](https://coveralls.io/github/bellavista/router?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/bellavista/router)](https://goreportcard.com/report/github.com/bellavista/router)


# Bella Vista Router

**Pure Go's stdlib, idiomatic, ultra-fast, simple, net/http compatible, context compatible, http router/mux**

Your existent http.Handler works with NO CHANGES with this package. And now they can get route parameters!
The parameters are stored on the new Context object inside http.Request since Go 1.7. 


## Requirements

- Go 1.7+


## Getting started

```go
import (
    "github.com/bellavista/router"
    "log"
    "net/http"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello " + router.GetString("name")))
}

func main() {
    // Create route
    r := New("/")
    r.Add("/hello/:name", http.HandlerFunc(sayHello))
    
    s := &http.Server{
        Addr:           ":8080",
        Handler:        router.Route(r),
    }
    
    log.Fatal(s.ListenAndServe())
}
```

