package main

import (
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"time"
)

var (
	ErrorTooManyReq = errors.New("too many request.")
)

func main() {
	http.Handle("/", ClientBasedLimit(http.HandlerFunc(InvestigateRequestHandler)))
	log.Fatal(http.ListenAndServe(":3000", nil))
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

	var users = map[string]Client{}

	// this goroutine prevents users map to grow infinitely.
	go func() {
		// every three seconds, a cleaning operation is done. If user's lastSeen is >60 seconds ago, its ipAddr's will be removed.
		for _ = range time.Tick(time.Second * 3) {
			for _, u := range users {
				if u.lastSeen.Before(time.Now().Add(-time.Second * 60)) {
					delete(users, u.ipAddr)
				}
			}
		}

	}()

	h := func(w http.ResponseWriter, r *http.Request) {

		c, ok := users[r.RemoteAddr]
		if !ok {
			// if there is no such a user, create one and allow the traffic immediately.
			users[r.RemoteAddr] = Client{
				lastSeen: time.Now(),
				ipAddr:   r.RemoteAddr,
				limit:    rate.NewLimiter(4, 8),
			}
			handler.ServeHTTP(w, r)
			return
		}

		if !c.limit.Allow() {
			http.Error(w, ErrorTooManyReq.Error(), http.StatusTooManyRequests)
			return
		}

		// update the lastSeen for every successful call.
		c.lastSeen = time.Now()
		handler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(h)
}

type Client struct {
	lastSeen time.Time
	ipAddr   string
	limit    *rate.Limiter
}
