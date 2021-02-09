package main

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type ResponseWrapper struct {
	http.ResponseWriter
	Status int
}

func (wrap *ResponseWrapper) WriteHeader(status int) {
	wrap.ResponseWriter.WriteHeader(status)
	wrap.Status = status
}

type Middleware func(http.Handler) http.Handler

func LoggingMiddleware(logger *log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Println(
						"err", err,
						"trace", string(debug.Stack()),
					)
				}
			}()

			wrapper := &ResponseWrapper{
				ResponseWriter: w,
			}
			start := time.Now()
			next.ServeHTTP(wrapper, r)
			logger.Println(r.Method, r.URL.EscapedPath(), "status", wrapper.Status, "response time", time.Since(start))
		}
		return http.HandlerFunc(fn)
	}
}
