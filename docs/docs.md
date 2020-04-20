# Basic Routing

The mos tbasic route binding requires a URL and a Handler

```go
    e.GET("/", func(c enlight.Context) error {
        return c.String(200, "Hello World")
    })
```

The return value will be sent back to the client as a response.

## Available Router Methods
```go
    e.CONNECT(path, func, ...middleware)
    e.DELETE(path, func, ...middleware)
    e.GET(path, func, ...middleware)
    e.HEAD(path, func, ...middleware)
    e.OPTIONS(path, func, ...middleware)
    e.PATCH(path, func, ...middleware)
    e.POST(path, func, ...middleware)
    e.PUT(path, func, ...middleware)
    e.TRACE(path, func, ...middleware)
```

To register a route that responds to all Methods use `e.Any()`

```go
    e.Any(path, func)
```

To register a route that responds to multiple Methods use `e.Match()`

```go
    e.Match([...]string{"GET", "POST", "PUT"}, path, func, ...middleware)
```

Additional you can Drop routes in runtime, this is usefull if you have some kind of dynamic API Service

```go
    e.Drop(method, path)
```

# Route Parameters

## Required Parameters

For dynamic routes, you can define route parameters like so:
```go
    e.GET("/posts/:id", func(c enlight.Context) error {
        id := c.Param("id")

        return c.String(200, id)
    })
```

In the example above, :id is a route parameter.
Its value is then retrieved via the Param function.