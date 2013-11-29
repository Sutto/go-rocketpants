package rocketpants

type ApiError struct {
	Name       string
	StatusCode int
	Message    string
	Metadata   map[string]interface{}
}

type RenderableApiError map[string]interface{}

var errorRegistry = make(map[string](*ApiError))

func init() {
	RegisterError("unknown", 500, "An unknown server-side error occured.")
}

func ErrorByName(name string) *ApiError {
	return errorRegistry[name]
}

func RegisterError(name string, statusCode int, message string) *ApiError {
	apiError := &ApiError{Name: name, StatusCode: statusCode, Message: message}
	errorRegistry[apiError.Name] = apiError
	return apiError
}

func (w *ResponseWriter) RenderError(errorName string) {
	error := ErrorByName(errorName)
	// Fallback to an unknown error name.
	if error == nil {
		error = ErrorByName("unknown")
	}
	renderable := map[string](interface{}){
		"error":             error.Name,
		"error_description": error.Message,
	}
	// Merge in the values from the metadata.
	for key, value := range error.Metadata {
		renderable[key] = value
	}
	// Now, render the error code out.
	w.RenderMapWithCode(error.StatusCode, renderable)
}
