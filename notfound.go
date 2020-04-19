package enlight

import "fmt"

// NotFoundHandler is the default 404 handler
func NotFoundHandler(c Context) error {
	fmt.Println("Not FOund Error")
	return nil
}
