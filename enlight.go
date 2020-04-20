package enlight

import (
	"fmt"
	"io"
	"net/url"
	"path"
	"path/filepath"
	"sync"

	"github.com/valyala/fasthttp"
)

// Enlight is a http.Handler which can be
// used to dispatch requests to different handler
// functions via configurable routes
type Enlight struct {
	common
	Debug            bool
	Router           *Router
	Server           *fasthttp.Server
	TLSServer        *fasthttp.Server
	premiddleware    []MiddlewareFunc
	aftermiddleware  []MiddlewareFunc
	middleware       []MiddlewareFunc
	HTTPErrorHandler HTTPErrorHandler
	pool             sync.Pool
	Renderer         Renderer
}

// Common struct for Echo & Group.
type common struct{}

// Renderer is the interface that wraps the Render function.
type Renderer interface {
	Render(io.Writer, string, interface{}, Context) error
}

// Map defines a generic map of type `map[string]interface{}`.
type Map map[string]interface{}

// New returns a new initialized Enlight instance
func New() (e *Enlight) {
	e = &Enlight{
		Server:           new(fasthttp.Server),
		TLSServer:        new(fasthttp.Server),
		Router:           NewRouter(),
		HTTPErrorHandler: e.DefaultHTTPErrorHandler,
		Debug:            false,
	}
	e.Server.Handler = e.ServeHTTP
	e.pool.New = func() interface{} {
		return e.NewContext()
	}
	return
}

// NewContext returns a Context instance.
func (e *Enlight) NewContext() Context {
	return &context{
		enlight: e,
		handler: NotFoundHandler,
	}
}

// Before adds middleware to the chain which is run before router.
func (e *Enlight) Before(middleware ...MiddlewareFunc) {
	e.premiddleware = append(e.premiddleware, middleware...)
}

// After adds middleware to the chain which is run after router.
func (e *Enlight) After(middleware ...MiddlewareFunc) {
	e.aftermiddleware = append(e.aftermiddleware, middleware...)
}

// Use adds middleware to the chain which is run after router
func (e *Enlight) Use(middleware ...MiddlewareFunc) {
	e.middleware = append(e.middleware, middleware...)
}

// CONNECT registers a new CONNECT route for a path with matching handler
func (e *Enlight) CONNECT(path string, handle HandleFunc, m ...MiddlewareFunc) {
	e.Add(fasthttp.MethodConnect, path, handle, m...)
}

// DELETE registers a new DELETE route for a path with matching handler
func (e *Enlight) DELETE(path string, handle HandleFunc, m ...MiddlewareFunc) {
	e.Add(fasthttp.MethodDelete, path, handle, m...)
}

// GET registers a new GET route for a path withh matching handler
func (e *Enlight) GET(path string, handle HandleFunc, m ...MiddlewareFunc) {
	e.Add(fasthttp.MethodGet, path, handle, m...)
}

// HEAD registers a new HEAD route for a path withh matching handler
func (e *Enlight) HEAD(path string, handle HandleFunc, m ...MiddlewareFunc) {
	e.Add(fasthttp.MethodHead, path, handle, m...)
}

// OPTIONS registers a new OPTIONS route for a path withh matching handler
func (e *Enlight) OPTIONS(path string, handle HandleFunc, m ...MiddlewareFunc) {
	e.Add(fasthttp.MethodOptions, path, handle, m...)
}

// PATCH registers a new PATCH route for a path with matching handler
func (e *Enlight) PATCH(path string, handle HandleFunc, m ...MiddlewareFunc) {
	e.Add(fasthttp.MethodPatch, path, handle, m...)
}

// POST registers a new POST route for a path with matching handler
func (e *Enlight) POST(path string, handle HandleFunc, m ...MiddlewareFunc) {
	e.Add(fasthttp.MethodPost, path, handle, m...)
}

// PUT registers a new PUT route for a path with matching handler
func (e *Enlight) PUT(path string, handle HandleFunc, m ...MiddlewareFunc) {
	e.Add(fasthttp.MethodPut, path, handle, m...)
}

// TRACE registers a new TRACE route for a path with matching handler
func (e *Enlight) TRACE(path string, handle HandleFunc, m ...MiddlewareFunc) {
	e.Add(fasthttp.MethodTrace, path, handle, m...)
}

var (
	methods = [...]string{
		fasthttp.MethodConnect,
		fasthttp.MethodDelete,
		fasthttp.MethodGet,
		fasthttp.MethodHead,
		fasthttp.MethodOptions,
		fasthttp.MethodPatch,
		fasthttp.MethodPost,
		PROPFIND,
		fasthttp.MethodPut,
		fasthttp.MethodTrace,
		REPORT,
	}
)

// Any registers a new route for all HTTP methods and path with matching handler
func (e *Enlight) Any(path string, handle HandleFunc, middleware ...MiddlewareFunc) {
	for _, m := range methods {
		e.Add(m, path, handle, middleware...)
	}
}

// Match registers a new route for all given HTTP methods and path with matching handler
func (e *Enlight) Match(methods []string, path string, handle HandleFunc, middleware ...MiddlewareFunc) {
	for _, m := range methods {
		e.Add(m, path, handle, middleware...)
	}
}

// Drop removes a route from router-tree
func (e *Enlight) Drop(method, path string) {
	e.Router.Drop(method, path)
}

// Static serves static files
func (e *Enlight) Static(prefix, root string) {
	if root == "" {
		root = "."
	}
	e.static(prefix, root, e.GET)
}

// Add registers a new route for an HTTP method and path with matching handler
// in the router with optional route-level middleware.
func (e *Enlight) Add(method, path string, handle HandleFunc, middleware ...MiddlewareFunc) {
	e.add(method, path, handle, middleware...)
}
func (e *Enlight) add(method, path string, handle HandleFunc, middleware ...MiddlewareFunc) {
	e.Router.Handle(method, path, func(c Context) error {
		h := handle
		// Chain middleware
		for i := len(middleware) - 1; i >= 0; i-- {
			h = middleware[i](h)
		}
		return h(c)
	}, false)
}

func (common) static(prefix, root string, get func(string, HandleFunc, ...MiddlewareFunc)) {
	h := func(c Context) error {

		p, err := url.PathUnescape(c.Param("filepath"))
		if err != nil {
			return err
		}
		name := filepath.Join(root, path.Clean("/"+p)) // "/"+ for security
		return c.File(name)
	}
	if prefix == "/" {
		get(prefix+"*filepath", h)
		return
	}
	get(prefix+"/*filepath", h)
	return
}

// ServeHTTP implements `http.Handler` interface, which serves HTTP requests.
//func (e *Enlight) ServeHTTP(w http.ResponseWriter, r *http.Request) {
func (e *Enlight) ServeHTTP(ctx *fasthttp.RequestCtx) {
	// Acquire context
	c := e.pool.Get().(*context)
	c.Reset(ctx)

	h := NotFoundHandler

	if e.premiddleware == nil {
		e.Router.Find(c)
		h = c.Handler()
		h = applyMiddleware(h, e.middleware...)
	} else {
		h = func(c Context) error {
			e.Router.Find(c)
			h := c.Handler()
			h = applyMiddleware(h, e.middleware...)
			return h(c)
		}
		h = applyMiddleware(h, e.premiddleware...)
	}

	if err := h(c); err != nil {
		fmt.Println(err)
		e.HTTPErrorHandler(err, c)
	}

	// the last handleFunc is usually returning nil
	//

	if e.aftermiddleware != nil {
		after := func(c Context) error {
			// middlewares are calling next(c)
			// it's not always clear if it's the last in chain
			return nil
		}
		after = applyMiddleware(after, e.aftermiddleware...)

		if err := after(c); err != nil {
			fmt.Println(err)
			e.HTTPErrorHandler(err, c)
		}
	}

	// Clearing ref to fasthttp
	c.RequestCtx = nil
	e.pool.Put(c)
}

// Start starts an HTTP server.
func (e *Enlight) Start(address string) error {
	return e.StartServer(address)
}

// StartServer starts a custom http server.
func (e *Enlight) StartServer(address string) (err error) {
	e.Server = &fasthttp.Server{
		Name:    "Enlight",
		Handler: e.ServeHTTP,
	}

	fmt.Printf("⇨ http server started on %s\n", address)
	return e.Server.ListenAndServe(address)
}

// Shutdown stops the server gracefully.
func (e *Enlight) Shutdown() error {
	fmt.Print("⇨ Server is shutting down...")
	return e.Server.Shutdown()
}
