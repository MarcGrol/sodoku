package main

import (
	"expvar"
	"log"
	"net/http"
	"time"
)

var (
	myCounters = expvar.NewMap("counters")
)

func init() {
	// initialize hits-counter in map to zero
	myCounters.Add("hits", 0)
}

// middleware example
func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

// middleware example
func countingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		myCounters.Add("hits", 1)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
