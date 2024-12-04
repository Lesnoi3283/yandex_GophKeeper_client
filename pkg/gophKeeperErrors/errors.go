package gophKeeperErrors

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
