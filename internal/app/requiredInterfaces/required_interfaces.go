package requiredInterfaces

import "net/http"

//go:generate mockgen -source=required_interfaces.go -destination=./mocks/mocks.go -package=mocks

//type EncryptionWriter interface {
//}

// HTTPClient have to send an HTTP request and return a response.
// It was created to make mocking requests possible.
type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}
