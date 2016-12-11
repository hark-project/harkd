package routes

import (
	"encoding/json"
	"io"

	"harkd/errors"
)

type responseWrapper interface {
	Wrap(interface{}) interface{}
}

type payloadResponseWrapper struct {
}

type wrappedResponsePayload struct {
	Payload   interface{} `json:"payload,omitempty"`
	Error     *string     `json:"error,omitempty"`
	ErrorCode *int        `json:"errorCode,omitempty"`
}

func (pr payloadResponseWrapper) Wrap(res interface{}) interface{} {
	// Meta is not implemented - create an empty map for now
	// If the response is an error, we use it as the error in the response.
	// Otherwise, we use it as the payload.
	if resErr, ok := res.(error); ok {
		code := getErrorCode(resErr)
		errMsg := resErr.Error()
		return wrappedResponsePayload{nil, &errMsg, &code}
	}
	return wrappedResponsePayload{res, nil, nil}
}

// getErrorCode determines the error code given an error.
//
// It checks to see if the error is an instance of error.Coder. If it's not,
// it uses the DefaultErrorCode from the errors package.
func getErrorCode(err error) int {
	if resCoderErr, ok := err.(errors.Coder); ok {
		return resCoderErr.Code()
	}
	return errors.DefaultErrorCode
}

type responseEncoder interface {
	Encode(w io.Writer, val interface{})
}

func jsonResponseEncoder() responseEncoder {
	return jsonResponseEncoderImpl{payloadResponseWrapper{}}
}

type jsonResponseEncoderImpl struct {
	responseWrapper
}

func (jsr jsonResponseEncoderImpl) Encode(w io.Writer, val interface{}) {
	wrapped := jsr.Wrap(val)
	enc := json.NewEncoder(w)
	enc.Encode(wrapped)
}
