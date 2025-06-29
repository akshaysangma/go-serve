package middleware

import "net/http"

type contextKey string

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			m := xs[i]
			next = m(next)
		}
		return next
	}
}
