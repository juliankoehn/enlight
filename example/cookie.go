package main

import (
	"fmt"

	"github.com/juliankoehn/enlight"
)

func setCookie(c enlight.Context) error {
	name := c.Param("name")
	value := c.Param("value")

	c.SetCookie(name, value)

	return c.String(200, fmt.Sprintf("cookie added: %s = %s", name, value))
}

func retrieveCookie(c enlight.Context) error {
	name := c.Param("name")

	value := c.Cookie(name)

	return c.String(200, value)
}
func deleteCookie(c enlight.Context) error {
	name := c.Param("name")

	c.RemoveCookie(name)

	return c.String(200, fmt.Sprintf("cookie %s removed", name))
}
