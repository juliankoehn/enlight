package enlight

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/valyala/fasthttp"
)

type (
	// Context represents the context of the current HTTP request. It holds request and
	// response objects, path, path parameters, data and registered handler.
	Context interface {
		Request() *fasthttp.RequestCtx

		// QueryParams returns the query parameters as `url.Values`.
		QueryParams() url.Values

		// String sends a string response with status code.
		String(code int, s string) error

		// JSON sends a JSON response with status code.
		JSON(code int, i interface{}) error

		// NoContent sends a response with no body and a status code.
		NoContent(code int) error
	}

	context struct {
		RequestCtx *fasthttp.RequestCtx
		request    *http.Request
		path       string
		query      url.Values
		handler    HandleFunc
	}
)

const (
	defaultIndent = "  "
)

func (c *context) writeContentType(value string) {
	header := &c.RequestCtx.Response.Header
	if string(header.ContentType()) == "" {
		header.SetContentType(value)
	}
}

func (c *context) Request() *fasthttp.RequestCtx {
	return c.RequestCtx
}

func (c *context) QueryParams() url.Values {
	if c.query == nil {
		c.query = c.request.URL.Query()
	}
	return c.query
}

func (c *context) Blob(code int, contentType string, b []byte) (err error) {
	c.writeContentType(contentType)

	c.RequestCtx.SetStatusCode(code)
	_, err = c.RequestCtx.Write(b)

	return
}

func (c *context) NoContent(code int) error {
	c.RequestCtx.Response.SetStatusCode(code)
	return nil
}

func (c *context) String(code int, s string) (err error) {
	return c.Blob(code, MIMETextPlainCharsetUTF8, []byte(s))
}

func (c *context) json(code int, i interface{}) error {
	enc := json.NewEncoder(c.RequestCtx)

	c.writeContentType(MIMEApplicationJSONCharsetUTF8)
	c.RequestCtx.Response.SetStatusCode(code)
	//c.response.Status = code
	return enc.Encode(i)
}

func (c *context) JSON(code int, i interface{}) (err error) {
	return c.json(code, i)
}

// func (c *context) Reset(r *http.Request, w http.ResponseWriter) {
func (c *context) Reset(ctx *fasthttp.RequestCtx) {
	//c.request = r
	// c.response.reset(ctx)
	c.RequestCtx = ctx
	c.query = nil
	c.handler = NotFoundHandler
}
