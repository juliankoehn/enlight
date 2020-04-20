# Enlight HTTP Router

Golang HTTP Router with `valyala/fasthttp` 

`router.go` and `tree.go` are mostly based upon (julienschmidt/httprouter)[https://github.com/julienschmidt/httprouter]

## Feature Overview

- Optimized HTTP router which smartly prioritize routes
- Build robust and scalable RESTful APIs
- Extensible middleware framework
- Define middleware at root or route level
- Centralized HTTP error handling

## Example

```go
package main

import (
    "github.com/juliankoehn/enlight"
)

func main() {
    e := enlight.New()

    e.Use(middleware.Recover)

    // Routes
    e.GET("/", helloWorldHandler)

    e.Start(":8080")
}

func helloWorldHandler(c enlight.Context) error {
    return c.String(200, "Hello, World!")
}
```

## TODO

* Route based middleware