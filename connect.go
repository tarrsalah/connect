package connect

import (
	"errors"
	"net/http"
)

var errInvalidMiddleware = errors.New("invalid middleware")

// Middleware wraps an HTTP handler with extra before/after behavior.
type Middleware interface {
	Wrap(http.Handler) http.Handler
}

// MiddlewareFunc type is an adapter to allow the use of ordinary
// functions as HTTP Middlewares.
type MiddlewareFunc func(http.Handler) http.Handler

// Wrap calls m(h)
func (m MiddlewareFunc) Wrap(h http.Handler) http.Handler {
	return m(h)
}

// Chain compose a list of middlewares.
func Chain(middlewares ...interface{}) Middleware {
	var h http.Handler

	fn := func(handler http.Handler) http.Handler {
		h = handler
		for i := len(middlewares) - 1; i >= 0; i-- {
			m := middlewares[i]
			switch m.(type) {
			case func(http.Handler) http.Handler:
				h = MiddlewareFunc(m.(func(http.Handler) http.Handler)).Wrap(h)

			case Middleware:
				h = m.(Middleware).Wrap(h)
			default:
				panic(errInvalidMiddleware)
			}
		}
		return h
	}

	return MiddlewareFunc(fn)
}
