package errors

import "fmt"

// DefaultErrorCode is the default code used for an error.
// This applies if we handle an error which is not for hark code without wrapping it.
const DefaultErrorCode = 500000

// Coder is an error which can provide an error code.
type Coder interface {
	error
	Code() int
}

type harkBadRequestError struct {
	code int
	msg  string
}

func (hbr harkBadRequestError) Error() string {
	return hbr.msg
}

func (hbr harkBadRequestError) Code() int {
	return hbr.code
}

type harkNotFoundError struct {
	code int
	msg  string
}

func (hnf harkNotFoundError) Error() string {
	return hnf.msg
}

func (hnf harkNotFoundError) Code() int {
	return hnf.code
}

type harkConflictError struct {
	code int
	msg  string
}

func (hc harkConflictError) Error() string {
	return hc.msg
}
func (hc harkConflictError) Code() int {
	return hc.code
}

type harkInternalServerError struct {
	code int
	msg  string
}

func (his harkInternalServerError) Error() string {
	return his.msg
}

func (his harkInternalServerError) Code() int {
	return his.code
}

// ErrBadRequestEntity creates an error for 400 responses
func ErrBadRequestEntity(err error) error {
	return harkBadRequestError{400001, fmt.Sprintf("Could not decode request entity: %q", err.Error())}
}

// ErrEntityInvalid creates an error for 400 responses
func ErrEntityInvalid(msg string) error {
	return harkBadRequestError{400002, fmt.Sprintf("Request entity invalid: %q", msg)}
}

// ErrMachineNotFound creates an error for 404 responses
func ErrMachineNotFound(machineID string) error {
	return harkNotFoundError{404001, fmt.Sprintf("Machine not found: %q", machineID)}
}

// ErrEntityConflict creates an error for 409 responses
func ErrEntityConflict(msg string) error {
	return harkConflictError{404002, msg}
}

// ErrSerialization creates an error for 500 responses
func ErrSerialization(msg string, err error) error {
	fullMsg := fmt.Sprintf("Failed %s: %q", msg, err)
	return harkInternalServerError{500001, fullMsg}
}

// ErrStatePersist creates an error for 500 responses
func ErrStatePersist(err error) error {
	return harkInternalServerError{500002, "failed to persist state: " + err.Error()}
}

// ErrStateLock creates an error for 500 responses
func ErrStateLock(err error) error {
	return harkInternalServerError{500003, "failed to lock state for writing: " + err.Error()}
}

// ErrUserLookup creates an error for 500 responses
func ErrUserLookup(err error) error {
	return harkInternalServerError{500004, fmt.Sprintf("failed to look up user: %q", err.Error())}
}
