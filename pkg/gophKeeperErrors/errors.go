package gophKeeperErrors

import "errors"

var errWrongLoginOrPassword = errors.New("wrong login or password (or user not exists)")

// NewErrWrongLoginOrPassword returns errWrongLoginOrPassword error.
// errWrongLoginOrPassword is just a errors.New() err.
func NewErrWrongLoginOrPassword() error {
	return errWrongLoginOrPassword
}

var errUserAlreadyExists = errors.New("user already exists")

// NewErrUserAlreadyExists returns errUserAlreadyExists error.
// errUserAlreadyExists is just a errors.New() err.
func NewErrUserAlreadyExists() error {
	return errUserAlreadyExists
}

// ErrWithHTTPCode contains HTTP status and error text.
type ErrWithHTTPCode struct {
	StatusCode int
	Text       string
}

func (e *ErrWithHTTPCode) Error() string {
	return e.Text
}

func (e *ErrWithHTTPCode) Code() int {
	return e.StatusCode
}

// NewErrWithHTTPCode returns new ErrWithHTTPCode.
func NewErrWithHTTPCode(statusCode int, text string) *ErrWithHTTPCode {
	return &ErrWithHTTPCode{StatusCode: statusCode, Text: text}
}
