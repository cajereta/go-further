package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rr *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, rr, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, rr)
	})
}


func (app *application) rateLimit(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(2,4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow(){
			app.rateLimitExceededResponse(w,r)
			return
		}
		next.ServeHTTP(w,r)
	})
}