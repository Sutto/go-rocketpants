package rocketpants

import (
  "net/http"
  "reflect"
)


type ApiHandler interface {
	ServeHTTP(*ResponseWriter, *http.Request)
}

type ApiHandlerFunc func(*ResponseWriter, *http.Request)
// Convert ApiHandlerFunc to also be a valid ApiHandler.
func (f ApiHandlerFunc) ServeHTTP(w *ResponseWriter, r *http.Request) {
	f(w, r)
}

type Handler struct {
  Endpoint ApiHandler
}

// Handles the conversion and 'flushing' the api response to complete it.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := WrapResponse(w)
	h.Endpoint.ServeHTTP(response, r)
	// TODO: Hooks for the middleware.
	if !response.finished {
		response.Finish()
	}
}

func NewHandler(handler interface{}) http.Handler {
  reflected := reflect.TypeOf(handler)
  if(reflected.Kind() == reflect.Func) {
	 return Handler{ApiHandlerFunc(handler.(func(*ResponseWriter, *http.Request)))}
  } else {
    return Handler{handler.(ApiHandler)}
  }
}