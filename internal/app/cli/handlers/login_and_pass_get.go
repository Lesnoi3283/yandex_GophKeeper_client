package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"yandex_GophKeeper_client/internal/app/entities"
	"yandex_GophKeeper_client/pkg/gophKeeperErrors"
)

// GetLoginAndPassword sends login to the backend and returns password or error.
// If http status code != 200 - returns a gophKeeperErrors.ErrWithHTTPCode.
func (h *Handler) GetLoginAndPassword(login string) (string, error) {
	//prepare request
	req, err := http.NewRequest(http.MethodGet, get_login_and_password_path, bytes.NewBufferString(login))
	if err != nil {
		return "", fmt.Errorf("cant create request, err: %w", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	//send request
	resp, err := http.DefaultClient.Do(req)
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

	password := string(bodyBytes)
	if password == "" {
		return "", fmt.Errorf("password is empty, but server`s response`s status is OK")
	}
	return password, nil
}
