package enlight

import (
	ctx "context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

type (
	user struct {
		ID   int    `json:"id" xml:"id" form:"id" query:"id" param:"id"`
		Name string `json:"name" xml:"name" form:"name" query:"name" param:"name"`
	}
)

const (
	userJSON = `{"id":1,"name":"Jon Snow"}`
)

func serve(handler fasthttp.RequestHandler, req *http.Request) (*http.Response, error) {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx ctx.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}

	return client.Do(req)
}

func TestEnlightStart(t *testing.T) {
	e := New()
	go func() {
		assert.NoError(t, e.Start(":0"))
	}()

	time.Sleep(200 * time.Millisecond)
}

func TestEnlightHandler(t *testing.T) {
	e := New()

	e.GET("/", func(c Context) error {
		return c.String(200, "Hello")
	})

	r, err := http.NewRequest("GET", "http://test/", nil)
	if err != nil {
		t.Error(err)
	}

	res, err := serve(e.ServeHTTP, r)
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, string(body), "Hello")
}
