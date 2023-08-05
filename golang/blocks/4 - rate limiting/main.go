package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
)

func main() {
	http.Handle("/", GlobalLimit(http.HandlerFunc(InvestigateRequestHandler)))
	http.ListenAndServe(":3000", nil)
}

func InvestigateRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "rate limited end point")
}

// GlobalLimit is a middleware to limit requests globally. Not too useful.
func GlobalLimit(handler http.Handler) http.Handler {

	// r for rate, b for bucket size.
	limiter := rate.NewLimiter(4, 8)

	h := func(w http.ResponseWriter, r *http.Request) {

		// Allow will remove one token from the bucket and return true when request is possible.
		// If there is no token in bucket, Allow will return false.
		if !limiter.Allow() {
			http.Error(w, "exceed the limit", http.StatusTooManyRequests)
			return
		}

		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(h)
}

// ClientBasedLimit is a middleware to limit requests based on ip addresses.
func ClientBasedLimit(handler http.Handler) http.Handler {

	h := func(w http.ResponseWriter, r *http.Request) {

		// TODO: implement ip based rate limiting with a map.

		panic("not implemented yet.")
	}

	return http.HandlerFunc(h)

}
