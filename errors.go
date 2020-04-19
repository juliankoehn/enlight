package enlight

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/valyala/fasthttp"
)

// Errors
var (
	ErrNotFound            = NewHTTPError(http.StatusNotFound)
	ErrInvalidRedirectCode = errors.New("invalid redirect status code")
)

// HTTPError represents an error that occured while handling a request.
type HTTPError struct {
	Code     int         `json:"-"`
	Message  interface{} `json:"message"`
	Internal error       `json:"-"` // Stores the error returned by an external dependency
}

func NewHTTPError(code int, message ...interface{}) *HTTPError {
	he := &HTTPError{
		Code:    code,
		Message: http.StatusText(code),
	}
	if len(message) > 0 {
		he.Message = message[0]
	}
	return he
}

// Error makes it compatible with `error` interface
func (he *HTTPError) Error() string {
	if he.Internal == nil {
		return fmt.Sprintf("code=%d, message=%v", he.Code, he.Message)
	}
	return fmt.Sprintf("code=%d, message=%v, internal=%v", he.Code, he.Message, he.Internal)
}

// SetInternal sets error to HTTPError.Internal
func (he *HTTPError) SetInternal(err error) *HTTPError {
	he.Internal = err
	return he
}

// HTTPErrorHandler is a centralized HTTP error handler.
type HTTPErrorHandler func(error, Context)

// DefaultHTTPErrorHandler is the default HTTP error handler. It sends a JSON response
// with status code.
func (e *Enlight) DefaultHTTPErrorHandler(err error, c Context) {
	he, ok := err.(*HTTPError)
	if ok {
		if he.Internal != nil {
			if herr, ok := he.Internal.(*HTTPError); ok {
				he = herr
			}
		}
	} else {
		he = &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}

	code := he.Code
	message := he.Message

	if m, ok := message.(string); ok {
		message = Map{"message": m}
	}

	if string(c.Request().Method()) == fasthttp.MethodHead {
		err = c.NoContent(he.Code)
	} else {
		err = c.JSON(code, message)
	}
	if err != nil {
		fmt.Println(err)
		// TODO: log the error
	}

}
