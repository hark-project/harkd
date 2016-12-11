package routes

import (
	"net/http"

	"harkd/errors"
)

// responseWriter is an interface which handles all the logic of writing a response,
// including wrapping the response payload and determining the HTTP status code.
type responseWriter interface {
	WriteResponse(http.ResponseWriter, interface{})
	WriteResponseWithStatus(http.ResponseWriter, int, interface{})
}

func newResponseWriter() responseWriter {
	errorHandler := errors.NewErrorHandlerService()
	responseEncoder := jsonResponseEncoder()
	return responseWriterImpl{errorHandler, responseEncoder}
}

type responseWriterImpl struct {
	errors.ErrorHandlerService
	responseEncoder
}

func (rw responseWriterImpl) WriteResponse(w http.ResponseWriter, val interface{}) {
	var statusCode = 200
	// if the value is an error, get the error handler service to determine the
	// status code.
	if resErr, ok := val.(error); ok {
		statusCode = rw.GetHTTPStatusCode(resErr)
	}

	rw.WriteResponseWithStatus(w, statusCode, val)
}

func (rw responseWriterImpl) WriteResponseWithStatus(w http.ResponseWriter, statusCode int, val interface{}) {
	w.WriteHeader(statusCode)
	rw.Encode(w, val)
}
