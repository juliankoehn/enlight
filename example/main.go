package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/juliankoehn/enlight"
)

type TestStruct struct {
	Message string `json:"message"`
}

func getAutoHandler(c enlight.Context) error {
	m := TestStruct{"Hello from AutoHandler"}
	fmt.Println("getAutoHandler")
	return c.JSON(200, m)
}

func showHTML(c enlight.Context) error {
	return c.HTML(200, "<h1>Hello World</h1>")
}

func serve(ctx context.Context) (err error) {
	e := enlight.New()
	e.GET("/", showHTML)
	e.GET("/auto", getAutoHandler)

	e.Static("/public", "")

	go func() {
		if err := e.Start(":8085"); err != nil {
			fmt.Printf("listen:%+s\n", err)
		}
	}()

	<-ctx.Done()

	if err = e.Shutdown(); err != nil {
		fmt.Printf("server Shutdown Failed:%+s\n", err)
	}

	fmt.Printf("server exited properly\n")

	return
}

func main() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		fmt.Printf("system call:%+v\n", oscall)
		cancel()
	}()

	if err := serve(ctx); err != nil {
		fmt.Printf("failed to serve:+%v\n", err)
	}
}
