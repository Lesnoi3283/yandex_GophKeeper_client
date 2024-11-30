package gophKeeperErrors

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

func NewErrWithHTTPCode(statusCode int, text string) *ErrWithHTTPCode {
	return &ErrWithHTTPCode{StatusCode: statusCode, Text: text}
}
