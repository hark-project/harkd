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

// ErrUserLookup is an error used when hark fails to look up the local
// user in order to provide a hark context.
type ErrUserLookup struct {
	Err error
}

func (e ErrUserLookup) Error() string {
	return fmt.Sprintf("Error looking up current user: %s", e.Err)
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

// ErrEntityConflictError creates an error for 409 responses
func ErrEntityConflictError(msg string) error {
	return harkConflictError{404002, msg}
}

// ErrSerializationError creates an error for 500 responses
func ErrSerializationError(msg string, err error) error {
	fullMsg := fmt.Sprintf("Failed %s: %q", msg, err)
	return harkInternalServerError{500001, fullMsg}
}

// ErrStatePersistError creates an error for 500 responses
func ErrStatePersistError(err error) error {
	return harkInternalServerError{500002, "failed to persist state: " + err.Error()}
}

// ErrStateLockError creates an error for 500 responses
func ErrStateLockError(err error) error {
	return harkInternalServerError{500003, "failed to lock state for writing: " + err.Error()}
}
