package rocketpants

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type ResponseWriter struct {
	finished bool
	status   int
	body     []byte
	http.ResponseWriter
}

func WrapResponse(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{false, 200, nil, w}
}

func (w *ResponseWriter) Finish() {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(w.status)
	w.Write(w.body)
	w.finished = false
}

func (w *ResponseWriter) RenderMapWithCode(statusCode int, renderableContent interface{}) {
	// Encodes and renders the object, as a basic implementation.
	rendered, err := json.Marshal(renderableContent)
	if err != nil {
		panic("Unable to marshal the given objects to JSON")
	}
	w.status = statusCode
	w.body = rendered
}

func (w *ResponseWriter) ExposeWithMedata(content interface{}, metadata map[string]interface{}) {
	response := map[string](interface{}){"response": content}
	for key, value := range metadata {
		response[key] = value
	}
	w.RenderMapWithCode(200, response)
}

func (w *ResponseWriter) ExposeList(content []interface{}) {
	metadata := map[string](interface{}){"count": len(content)}
	w.ExposeWithMedata(content, metadata)
}

func (w *ResponseWriter) ExposeSingle(content interface{}) {
	metadata := make(map[string]interface{})
	w.ExposeWithMedata(content, metadata)
}

func (w *ResponseWriter) Expose(content interface{}) {
	switch reflect.TypeOf(content).Kind() {
	case reflect.Slice, reflect.Array:
		w.ExposeList(content.([]interface{}))
	default:
		w.ExposeSingle(content)
	}

}
