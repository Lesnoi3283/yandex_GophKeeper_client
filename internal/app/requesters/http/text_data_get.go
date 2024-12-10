package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"yandex_GophKeeper_client/pkg/gophKeeperErrors"
)

// GetText text name to the backend and returns text data or error.
// If http status code != 200 - returns a gophKeeperErrors.ErrWithHTTPCode.
func (h *Requester) GetText(textName string) (string, error) {
	//prepare request
	req, err := http.NewRequest(http.MethodGet, h.ApiAddress+"/"+getTextPath, bytes.NewBufferString(textName))
	if err != nil {
		return "", fmt.Errorf("cant create request, err: %w", err)
	}
	req.AddCookie(&http.Cookie{
		Name:  JwtCookieName,
		Value: h.JWT,
	})
	req.Header.Set("Content-Type", "text/plain")

	//send request
	resp, err := h.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("cant send request, err: %w", err)
	}
	defer resp.Body.Close()

	//read request
	if resp.StatusCode != http.StatusOK {
		return "", gophKeeperErrors.NewErrWithHTTPCode(resp.StatusCode, "Server`s response`s status is not OK")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cant read response body, err: %w", err)
	}

	return string(bodyBytes), nil
}
