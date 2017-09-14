[![GoDoc](https://godoc.org/github.com/bellavista/router?status.svg)](https://godoc.org/github.com/bellavista/router)
[![Version](https://badge.fury.io/gh/bellavista%2Frouter.svg)](https://badge.fury.io/gh/bellavista%2Frouter)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](LICENSE)
[![Build Status](https://travis-ci.org/BellaVista/router.svg?branch=master)](https://travis-ci.org/BellaVista/router)
[![Coverage Status](https://coveralls.io/repos/github/BellaVista/router/badge.svg?branch=master)](https://coveralls.io/github/BellaVista/router?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/bellavista/router)](https://goreportcard.com/report/github.com/bellavista/router)


# Bella Vista Router

**Pure Go's stdlib, idiomatic, ultra-fast, simple, net/http compatible, context compatible, http router/mux**

Your existent http.Handler works with NO CHANGES with this package. And now they can get route parameters!
The parameters are stored on the new Context object inside http.Request since Go 1.7. 


## Requirements

- Go 1.7+


## Getting started

```go
package main

import (
    "github.com/bellavista/router"
    "net/http"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello " + router.GetParam(r, "name")))
}

func main() {
    // Create route
    r := router.New("/")
    r.Add("/hello/:name", http.HandlerFunc(sayHello))
    
    s := &http.Server{
        Addr:           ":8080",
        Handler:        router.Build(r),
    }
    
    s.ListenAndServe()
}
```


## Handlers

Your handlers are plain simply Go's stdlib [http.Handler](https://golang.org/pkg/net/http/#Handler) objects. 


```go

import "net/http"

type MyHandler struct {}

func (mh MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello world!")
}

```


And, if you're used to use the [http.HandleFunc](https://golang.org/pkg/net/http/#HandleFunc) approach, 
you're still covered by the `net/http` package with [http.HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc). 

```go

import "net/http"

func MyHandlerFunc(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello world!")
}

h := http.HandlerFunc(MyHandlerFunc)

```

No surprises, no changes in the function signature when you want to receive route parameters, nothing. 


## Routes definition

Routes are handled by a router object that can group several routes under a single prefix. 
Then, multiple routers can join into a single dispatcher that acts as a replacement for `http.Server.Handler`.


```go

import (
    "github.com/bellavista/router"
    "net/http"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello " + router.GetParam(r, "name")))
}

func main() {
    // Create first router for routes starting at `/v1`
    r1 := router.New("/v1")
    r1.Add("/hello/:name", http.HandlerFunc(sayHello))
    
    // Create second router for routes starting at `/v2`
    r2 := router.New("/v2")
    r2.Add("/hello/:name", http.HandlerFunc(sayHello))
    
    // Create http.Server and dispatch both routers
    s := &http.Server{
        Addr:           ":8080",
        Handler:        router.Build(r1, r2), // router.Build creates the dispatcher object 
    }
    
    s.ListenAndServe()
}

```


## Route parameters

Since Go 1.7, the [context](https://golang.org/pkg/context) package is included on the stdlib, and with it, 
the http.Request object now has a Context included on its definition that allow us to pass context related values (just like our parameters) 
across the life of our request.
Bella Vista Router uses this new feature to keep it compatible with existent (and future) net/http handlers. 

Your routes can hold parameters by defining a route part starting with `:`.
So, if you want to receive a parameter called `id` at the end of your `/user` route, you can define and consume as follows


```go

import (
    "github.com/bellavista/router"
    "net/http"
)

func getUser(w http.ResponseWriter, r *http.Request) {
    id := router.GetParam(r, "id")
    
    // Do something with that id
    // ...
}

func main() {
    // Create route
    r := router.New("/")
    r.Add("/user/:id", http.HandlerFunc(getUser))
    
    s := &http.Server{
        Addr:           ":8080",
        Handler:        router.Build(r),
    }
    
    s.ListenAndServe()
}

```


## Middleware

Middleware type is a function that takes an http.Handler object and returns another http.Handler object to be executed. 
It has the following signature: 

```go
type Middleware func(http.Handler) http.Handler
``` 

It can be used to wrap handlers at handler level, router level and/or dispatcher level so different request flows can be architected:  

```go

import (
    "context"
    "github.com/bellavista/router"
    "net/http"
    
    "fake/mockup/service/users"
)

// Middleware to set content type
func mContentType(next http.Handler) http.Handler {
    return func(w http.ResponseWriter, r *http.Request) {
        // Set header
        r.Header().Set("Content-Type", "text/plain")
            
        // Continue flow
        next.ServeHTTP(w, r)
    }
}

// Handler that says hello
func hSayHello(w http.ResponseWriter, r *http.Request) {
    r.Write([]byte("Hello!"))
}

func main() {
    // Create route
    r := router.New("/")
    
    // Use middleware at router level
    r.Wrap(http.HandlerFunc(mContentType))
    
    // Route to handler
    r.Add("/hello", http.HandlerFunc(hSayHello))
    
    // Get dispatcher
    d := router.Build(r)
    
    s := &http.Server{
        Addr:           ":8080",
        Handler:        d,
    }
    
    s.ListenAndServe()
}

```
