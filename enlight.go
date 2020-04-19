package enlight

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"path"
	"path/filepath"
	"sync"
	"time"

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
	Listener         net.Listener
	TLSListener      net.Listener
	premiddleware    []MiddlewareFunc
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

// Use adds middleware to the chain which is run after router
func (e *Enlight) Use(middleware ...MiddlewareFunc) {
	e.middleware = append(e.middleware, middleware...)
}

// CONNECT registers a new CONNECT route for a path with matching handler
func (e *Enlight) CONNECT(path string, h HandleFunc) {
	e.Router.Handle(fasthttp.MethodConnect, path, h)
}

// DELETE registers a new DELETE route for a path with matching handler
func (e *Enlight) DELETE(path string, handle HandleFunc) {
	e.Router.Handle(fasthttp.MethodDelete, path, handle)
}

// GET registers a new GET route for a path withh matching handler
func (e *Enlight) GET(path string, handle HandleFunc) {
	e.Router.Handle(fasthttp.MethodGet, path, handle)
}

// HEAD registers a new HEAD route for a path withh matching handler
func (e *Enlight) HEAD(path string, handle HandleFunc) {
	e.Router.Handle(fasthttp.MethodHead, path, handle)
}

// OPTIONS registers a new OPTIONS route for a path withh matching handler
func (e *Enlight) OPTIONS(path string, handle HandleFunc) {
	e.Router.Handle(fasthttp.MethodOptions, path, handle)
}

// PATCH registers a new PATCH route for a path with matching handler
func (e *Enlight) PATCH(path string, handle HandleFunc) {
	e.Router.Handle(fasthttp.MethodPatch, path, handle)
}

// POST registers a new POST route for a path with matching handler
func (e *Enlight) POST(path string, handle HandleFunc) {
	e.Router.Handle(fasthttp.MethodPost, path, handle)
}

// PUT registers a new PUT route for a path with matching handler
func (e *Enlight) PUT(path string, handle HandleFunc) {
	e.Router.Handle(fasthttp.MethodPut, path, handle)
}

// TRACE registers a new TRACE route for a path with matching handler
func (e *Enlight) TRACE(path string, handle HandleFunc) {
	e.Router.Handle(fasthttp.MethodTrace, path, handle)
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
func (e *Enlight) Any(path string, handle HandleFunc) {
	for _, m := range methods {
		e.Router.Handle(m, path, handle)
	}
}

func (e *Enlight) Match(methods []string, path string, handle HandleFunc) {
	for _, m := range methods {
		e.Router.Handle(m, path, handle)
	}
}

// Static serves static files
func (e *Enlight) Static(prefix, root string) {
	if root == "" {
		root = "."
	}
	e.static(prefix, root, e.GET)
}

func (common) static(prefix, root string, get func(string, HandleFunc)) {
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
	// this is were we call the "preMiddleware"
	// the pre middleware allows us to cast dynamic routes to our router
	// we must register them and deregister them afterwards

	h := NotFoundHandler

	if e.premiddleware == nil {
		e.Router.Find(c)
		h = c.handler
		h = applyMiddleware(h, e.middleware...)
	} else {

	}

	if err := h(c); err != nil {
		fmt.Println(err)
		e.HTTPErrorHandler(err, c)
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
	if e.Listener == nil {
		e.Listener, err = newListener(address)
		if err != nil {
			return err
		}
	}
	fmt.Printf("⇨ http server started on %s\n", e.Listener.Addr())
	return fasthttp.Serve(e.Listener, e.ServeHTTP)
	//return s.Serve(e.Listener)
}

// Close immediately stops the listeners.
func (e *Enlight) Close() error {
	if e.TLSListener != nil {
		if err := e.TLSListener.Close(); err != nil {
			return err
		}
	}

	return e.Listener.Close()
}

// Shutdown stops the server gracefully.
func (e *Enlight) Shutdown() error {
	fmt.Print("⇨ Server is shutting down...")
	return e.Server.Shutdown()
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	if c, err = ln.AcceptTCP(); err != nil {
		return
	} else if err = c.(*net.TCPConn).SetKeepAlive(true); err != nil {
		return
	}
	// Ignore error from setting the KeepAlivePeriod as some systems, such as
	// OpenBSD, do not support setting TCP_USER_TIMEOUT on IPPROTO_TCP
	_ = c.(*net.TCPConn).SetKeepAlivePeriod(3 * time.Minute)
	return
}

func newListener(address string) (*tcpKeepAliveListener, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return &tcpKeepAliveListener{l.(*net.TCPListener)}, nil
}
