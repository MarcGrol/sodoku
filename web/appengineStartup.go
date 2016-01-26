/*   // +build appengine */

package web

import (
	"log"
	"net/http"
	"time"

	"github.com/justinas/alice"
)

func logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}

func init() {
	h := &SodokuHandler{
		Timeout:      10,
		MinSolutions: 1,
	}
	chain := alice.New(logger).Then(h)

	http.Handle("/sodoku", chain)
}
