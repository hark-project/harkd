package core

// A Validator is any type which can be validated.
type Validator interface {
	Validate() error
}
