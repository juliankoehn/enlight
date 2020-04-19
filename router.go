package enlight

import (
	"net/http"
	"strings"
	"sync"
)

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions
type Router struct {
	trees      map[string]*node
	paramsPool sync.Pool
	maxParams  uint16

	RedirectTrailingSlash bool

	// If enabled, adds the matched route path onto the http.Request context
	// before invoking the handler.
	// The matched route path is only added to handlers of routes that where
	// registered when this option was enabled.
	SaveMatchedRoutePath bool
	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is anwered with "Method Not Allowed"
	// and HTTP status code 405
	//
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// Cached value of global (*) allowed methods
	globalAllowed string

	// Configurable http.Handler which is called when no matching route is found
	// If it is not set, http.NotFound is used.
	NotFound http.Handler

	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNowAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	MethodNotAllowed http.Handler
}

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

// HandleFunc is a function that can be registered to a route to handle HTTP
// requests. Like http.HandlerFunc.
type HandleFunc func(Context) error

// NewRouter returns a new initialized Router.
//
func NewRouter() *Router {
	return &Router{
		HandleMethodNotAllowed: true,
		RedirectTrailingSlash:  true,
	}
}

// Handle registers a new request handle with the given path and method.
func (r *Router) Handle(method, path string, handle HandleFunc) {

	if method == "" {
		panic("method must not be empty")
	}

	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}
	if handle == nil {
		panic("handle must not be nil")
	}

	varsCount := uint16(0)

	if r.SaveMatchedRoutePath {
		varsCount++
		handle = r.saveMatchedRoutePath(path, handle)
	}

	if r.trees == nil {
		r.trees = make(map[string]*node)
	}

	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root

		r.globalAllowed = r.allowed("*", "")
	}

	root.addRoute(path, handle)

	// Update maxParams
	if paramsCount := countParams(path); paramsCount+varsCount > r.maxParams {
		r.maxParams = paramsCount + varsCount
	}

	// Lazy-init paramsPool alloc func
	if r.paramsPool.New == nil && r.maxParams > 0 {
		r.paramsPool.New = func() interface{} {
			ps := make(Params, 0, r.maxParams)
			return &ps
		}
	}
}

func (r *Router) allowed(path, reqMethod string) (allow string) {
	allowed := make([]string, 0, 9)

	if path == "*" { // server-wide
		// empty method is used for internal calls to refresh the cache
		if reqMethod == "" {
			for method := range r.trees {
				if method == http.MethodOptions {
					continue
				}
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		} else {
			return r.globalAllowed
		}
	} else {
		// specific path
		for method := range r.trees {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == http.MethodOptions {
				continue
			}

			handle, _, _ := r.trees[method].getValue(path)
			if handle != nil {
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		}
	}

	if len(allowed) > 0 {
		// add request method to list of allowed methods
		allowed = append(allowed, http.MethodOptions)

		// Sort allowed methods.
		// sort.strings(allowed) unfortunately causes unnecessary allocations
		// due to allowed being moved to the heap and interface conversion
		for i, l := 1, len(allowed); i < l; i++ {
			for j := i; j > 0 && allowed[j] < allowed[j-1]; j-- {
				allowed[j], allowed[j-1] = allowed[j-1], allowed[j]
			}
		}

		return strings.Join(allowed, ", ")
	}
	return
}

func (r *Router) saveMatchedRoutePath(path string, handle HandleFunc) HandleFunc {
	return func(c Context) error {
		return handle(c)
	}
}

// Find lookup a handler registered for method and path.
func (r *Router) Find(c Context) {
	ctx := c.(*context)

	method := string(ctx.RequestCtx.Method())
	path := string(ctx.RequestCtx.Path())

	ctx.path = path

	if root := r.trees[method]; root != nil {
		if handle, param, tsr := root.getValue(path); handle != nil {
			if param != nil {
				ctx.params = param
			}
			ctx.handler = handle
			return
		} else if method != "CONNECT" && path != "/" {
			if param != nil {
				ctx.params = param
			}
			code := 301
			if method != "GET" {
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308
				code = 307
			}

			if tsr && r.RedirectTrailingSlash {
				var uri string

				if len(path) > 1 && path[len(path)-1] == '/' {
					if len(path) > 1 && path[len(path)-1] == '/' {
						uri = path[:len(path)-1]
					} else {
						uri = path + "/"
					}
				}

				ctx.Redirect(code, uri)
				return
			}
		}
	}
}
