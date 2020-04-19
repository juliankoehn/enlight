package main

import (
	"fmt"

	"github.com/juliankoehn/enlight"
)

type TestStruct struct {
	Message string `json:"message"`
}

func getHandler(c enlight.Context) error {
	fmt.Println("getHandler")
	return c.String(200, "Hello World")
}

func getAutoHandler(c enlight.Context) error {
	m := TestStruct{"Hello from AutoHandler"}
	fmt.Println("getAutoHandler")
	return c.JSON(200, m)
}

func main() {
	e := enlight.New()
	e.GET("/", getHandler)
	e.GET("/auto", getAutoHandler)

	e.Start(":8085")
}
