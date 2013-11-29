package rocketpants

import (
	"log"
	"net/http"
	"time"
)

type RequestTimer struct {
	wrapped http.Handler
}

func (rt *RequestTimer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	before := time.Now()
	rt.wrapped.ServeHTTP(w, r)
	taken := float64(time.Now().Sub(before) / time.Millisecond)
	// TODO: Add in a header...
	log.Printf("Request to %s took %0.10fms", r.URL, taken)
}

func TimeRequests(wrapped http.Handler) http.Handler {
	return &RequestTimer{wrapped}
}
