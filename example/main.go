package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/juliankoehn/enlight"
	"github.com/juliankoehn/enlight/middleware"
	"github.com/valyala/fasthttp"
)

type (
	App struct {
		Enlight *enlight.Enlight
	}
	TestStruct struct {
		Message string `json:"message"`
	}
	Stats struct {
		Uptime       time.Time      `json:"uptime"`
		RequestCount uint64         `json:"requestCount"`
		Statuses     map[string]int `json:"statuses"`
		mutex        sync.RWMutex
	}
)

func NewStats() *Stats {
	return &Stats{
		Uptime:   time.Now(),
		Statuses: map[string]int{},
	}
}

// AddDynamicRoutes middleware to add dynamic routes before "routing happens"
func (a *App) AddDynamicRoutes(next enlight.HandleFunc) enlight.HandleFunc {
	return func(c enlight.Context) error {
		a.Enlight.GET("/dynamic", DynamicRoute)
		return next(c)
	}
}

// CleanupDynamicRoutes removes dynamic routes after serve
func (a *App) CleanupDynamicRoutes(next enlight.HandleFunc) enlight.HandleFunc {
	return func(c enlight.Context) error {
		a.Enlight.Drop("GET", "/dynamic")
		return next(c)
	}
}

func DynamicRoute(c enlight.Context) error {
	return c.String(200, "Hello from dynamic Route")
}

func PanicRoute(c enlight.Context) error {
	panic("Panic!")
}

// ServerHeader middleware adds a `Server` header to the response.
func ServerHeader(next enlight.HandleFunc) enlight.HandleFunc {
	return func(c enlight.Context) error {
		c.Response().Header.Set(enlight.HeaderServer, "Enlight/3.0")
		return next(c)
	}
}

// ProcessStats is the middleware function.
func (s *Stats) ProcessStats(next enlight.HandleFunc) enlight.HandleFunc {
	return func(c enlight.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.RequestCount++
		status := strconv.Itoa(c.Response().StatusCode())
		s.Statuses[status]++
		return nil
	}
}

// RouteBasedMiddleware demonstrates route based middleware
func RouteBasedMiddleware(next enlight.HandleFunc) enlight.HandleFunc {
	return func(c enlight.Context) error {
		c.Response().Header.Set(enlight.HeaderServer, "RouteBased")
		return next(c)
	}
}

// Handle is the endpoint to get stats
func (s *Stats) Handle(c enlight.Context) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return c.JSON(fasthttp.StatusOK, s)
}

func getAutoHandler(c enlight.Context) error {
	m := TestStruct{"Hello from AutoHandler"}
	fmt.Println("getAutoHandler")
	return c.JSON(200, m)
}

func showHTML(c enlight.Context) error {
	return c.HTML(200, "<h1>Hello World</h1>")
}

func userHandler(c enlight.Context) error {
	hello := string(c.QueryParams().Peek("name"))
	name := c.Param("name")
	return c.HTML(200, fmt.Sprintf("<h1>Hello %s & %s</h1>", name, hello))
}

func serve() (err error) {
	app := &App{}
	e := enlight.New()
	app.Enlight = e

	e.Use(middleware.Recover())

	e.GET("/", showHTML)
	e.GET("/auto", getAutoHandler, RouteBasedMiddleware)

	// testing Static
	e.Static("/public", "")

	s := NewStats()
	e.Use(s.ProcessStats)
	e.GET("/stats", s.Handle)

	e.GET("/panic", PanicRoute)

	e.GET("/user/:name", userHandler)

	// Server header
	e.Use(ServerHeader)

	e.Before(app.AddDynamicRoutes)

	e.After(app.CleanupDynamicRoutes)

	if err = e.Start(":8085"); err != nil {
		fmt.Printf("listen:%+s\n", err)
		return err
	}

	return
}

func main() {
	if err := serve(); err != nil {
		fmt.Printf("failed to serve:+%v\n", err)
	}
}
