package errornames

import "errors"

var ErrFailure = errors.New("failure") // want `error variable "ErrFailure" should be named err`

type errorAlias = error

func names() {
	closeErr := errors.New("close") // want `error variable "closeErr" should be named err`
	_, readErr := result()          // want `error variable "readErr" should be named err`
	var writeErr error              // want `error variable "writeErr" should be named err`
	var aliasErr errorAlias         // want `error variable "aliasErr" should be named err`
	var err error
	err = closeErr
	_, err = result()
	_, closeErr = result() // Do not report uses or reassignments again.
	_, _ = result()
	println(readErr, writeErr, aliasErr, err)
}

func result() (int, error) { return 0, nil }

func parameters(closeErr error) (resultErr error) { // want `error variable "closeErr" should be named err` `error variable "resultErr" should be named err`
	return closeErr
}

func singleLetter() {
	e := errors.New("failure") // want `error variable "e" should be named err`
	println(e)
	println(e)
	println(e)
	println(e)
	println(e)
	println(e)
}

// Names alone must not trigger the rule, nor should fields or concrete types
// implementing error be confused with variables of the built-in error type.
type record struct{ closeErr error }
type concreteError struct{}

func (concreteError) Error() string { return "failure" }

func allowed(err error) error {
	closeErr := "not an error"
	concreteErr := concreteError{}
	var pointerErr *concreteError
	var _ error
	println(closeErr, concreteErr.Error(), pointerErr)
	return err
}
