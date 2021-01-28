package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// var Services []Service

type Middleware func(http.Handler) http.Handler

func InitMiddleware() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Before")
			defer fmt.Println("After")
			h.ServeHTTP(w, r)
		})
	}
}

func LoggingMiddleware(logger *log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logger.Println("Logger middleware start here")
			defer logger.Println("Logger in the end of the request")
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func initAction(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello")
}

// func getStatus(w http.ResponseWriter, r *http.Request) {

// }

// func postStatus()

func main() {
	logFile, err := os.Create("tmp.log")
	if err != nil {
		fmt.Println("Failed to open a file")
	}
	mw := io.MultiWriter(os.Stdout, logFile)

	logger := log.New(mw, "Logger bruh: ", log.Ldate|log.Lshortfile)
	logMdlw := LoggingMiddleware(logger)
	mdlwr := InitMiddleware()
	router := http.NewServeMux()
	router.HandleFunc("/", initAction)

	finalMux := logMdlw(mdlwr(router))
	http.ListenAndServe(":9091", finalMux)
}
