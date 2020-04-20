package middleware

import (
	"github.com/juliankoehn/enlight"
)

type (
	// Skipper defines a function to skip middleware. Returning true skips processing
	// the middleware.
	Skipper func(enlight.Context) bool

	// BeforeFunc defines a function which is executed just before the middleware.
	BeforeFunc func(enlight.Context)
)

// DefaultSkipper returns false which processes the middleware.
func DefaultSkipper(enlight.Context) bool {
	return false
}
