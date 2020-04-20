package enlight

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
)

type (
	// Context represents the context of the current HTTP request. It holds request and
	// response objects, path, path parameters, data and registered handler.
	Context interface {
		Request() *fasthttp.RequestCtx
		Response() *fasthttp.Response

		// Param returns path parameter by name.
		Param(name string) string

		// QueryParams returns the query parameters as `*fasthttp.Args`.
		QueryParams() *fasthttp.Args

		// HTML sends an HTTP response with status code.
		HTML(code int, html string) error

		// HTMLBlob sends an HTTP blob response with status code.
		HTMLBlob(code int, b []byte) error

		// Blob sends a blob response with status code and content type.
		Blob(code int, contentType string, b []byte) error

		// String sends a string response with status code.
		String(code int, s string) error

		// JSON sends a JSON response with status code.
		JSON(code int, i interface{}) error

		// File sends a response with the content of the file.
		File(file string) error

		// NoContent sends a response with no body and a status code.
		NoContent(code int) error

		// Redirect redirects the request to a provided URL with status code.
		Redirect(code int, url string) error

		// Error invokes the registered HTTP error handler. Generally used by middleware.
		Error(err error)

		// WantsJSON checks if contentType or Accept header contains "application/json"
		WantsJSON() bool

		// Peek gets value of key from header or ""
		Peek(key string) string

		// Enlight returns the `Enlight` instance
		Enlight() *Enlight

		Handler() HandleFunc
	}

	context struct {
		RequestCtx *fasthttp.RequestCtx
		request    *http.Request
		path       string
		params     Params
		pnames     []string
		pvalues    []string
		query      *fasthttp.Args
		handler    HandleFunc
		enlight    *Enlight
	}
)

const (
	defaultMemory = 32 << 20 // 32 MB
	indexPage     = "index.html"
)

// WantsJSON checks if contentType or Accept header contains "application/json"
func (c *context) WantsJSON() bool {
	contentType := c.Peek("Content-Type")
	if contentType == "application/json" {
		return true
	}
	acceptable := c.RequestCtx.Request.Header.Peek("Accept")
	return strings.Contains(string(acceptable), "application/json")
}

// Peek gets value of key from header or ""
func (c *context) Peek(key string) string {
	return string(c.RequestCtx.Request.Header.Peek(key))
}

func (c *context) Handler() HandleFunc {
	return c.handler
}

func (c *context) Response() *fasthttp.Response {
	return &c.RequestCtx.Response
}

func (c *context) Request() *fasthttp.RequestCtx {
	return c.RequestCtx
}

func (c *context) Param(name string) string {
	return c.params.ByName(name)
}

func (c *context) ParamNames() []string {
	return c.pnames
}

func (c *context) QueryParams() *fasthttp.Args {
	if c.query == nil {
		c.query = c.RequestCtx.QueryArgs()
	}
	return c.query
}

func (c *context) HTML(code int, html string) (err error) {
	return c.HTMLBlob(code, []byte(html))
}

func (c *context) HTMLBlob(code int, b []byte) (err error) {
	return c.Blob(code, MIMETextHTMLCharsetUTF8, b)
}

func (c *context) Blob(code int, contentType string, b []byte) (err error) {
	c.RequestCtx.SetContentType(contentType)
	c.RequestCtx.SetStatusCode(code)
	_, err = c.RequestCtx.Write(b)

	return
}

func (c *context) NoContent(code int) error {
	c.RequestCtx.Response.SetStatusCode(code)
	return nil
}

func (c *context) Redirect(code int, url string) error {
	if code < 300 || code > 308 {
		return ErrInvalidRedirectCode
	}
	c.RequestCtx.Response.Header.Set(HeaderLocation, url)
	c.RequestCtx.Redirect(url, code)
	return nil
}

func (c *context) String(code int, s string) (err error) {
	return c.Blob(code, MIMETextPlainCharsetUTF8, []byte(s))
}

func (c *context) json(code int, i interface{}) error {
	enc := json.NewEncoder(c.RequestCtx)

	c.RequestCtx.SetContentType(MIMEApplicationJSONCharsetUTF8)
	c.RequestCtx.Response.SetStatusCode(code)
	//c.response.Status = code
	return enc.Encode(i)
}

func (c *context) JSON(code int, i interface{}) (err error) {
	return c.json(code, i)
}

func (c *context) File(file string) (err error) {
	c.RequestCtx.SendFile(file)
	return
}

func (c *context) Error(err error) {
	c.enlight.HTTPErrorHandler(err, c)
}

func (c *context) Enlight() *Enlight {
	return c.enlight
}

// func (c *context) Reset(r *http.Request, w http.ResponseWriter) {
func (c *context) Reset(ctx *fasthttp.RequestCtx) {
	//c.request = r
	// c.response.reset(ctx)
	c.RequestCtx = ctx
	c.query = nil
	c.handler = NotFoundHandler
	c.path = ""
	c.pnames = nil
	c.params = nil
}
