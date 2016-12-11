package routes

import (
	"encoding/json"
	"io"

	"harkd/core"
	"harkd/errors"
)

type requestDecoder interface {
	Decode(io.Reader, core.Validator) error
}

func jsonRequestDecoder() requestDecoder {
	return jsonRequestDecoderImpl{}
}

type jsonRequestDecoderImpl struct{}

// Requests to hark nest the actual request entity inside of a 'payload' key.
//
// The payload is a json.RawMessage so that we can can defer decoding it until
// we can specify the type requested by the caller.
type wrappedRequestPayload struct {
	Payload json.RawMessage `json:"payload"`
	Meta    map[string]interface{}
}

func (jrd jsonRequestDecoderImpl) Decode(r io.Reader, into core.Validator) error {
	// Decode the full request
	var wrapped wrappedRequestPayload

	dec := json.NewDecoder(r)
	err := dec.Decode(&wrapped)

	if err != nil {
		return errors.ErrBadRequestEntity(err)
	}

	// Now decode the actual entity payload into the interface provided
	// by the caller.
	if err = json.Unmarshal([]byte(wrapped.Payload), into); err != nil {
		return err
	}

	// Now that we've decoded the entity, we'll validate it
	return into.Validate()
}
