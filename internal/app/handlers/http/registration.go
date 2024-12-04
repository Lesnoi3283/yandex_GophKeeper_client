package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"yandex_GophKeeper_client/internal/app/entities"
	"yandex_GophKeeper_client/pkg/gophKeeperErrors"
)

// RegisterUser sends registration request to the backend and returns JWT string or error.
// If http status code != 201 - this func returns a gophKeeperErrors.ErrWithHTTPCode.
func (h *Handler) RegisterUser(login string, password string) (jwt string, err error) {
	user := entities.User{
		Login:    login,
		Password: password,
	}

	//create request
	jsonUser, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("can`t marshal user: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, h.Conf.APIAddress+registration_path, bytes.NewReader(jsonUser))
	if err != nil {
		return "", fmt.Errorf("can`t create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	//send request
	resp, err := h.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("can`t send request: %v", err)
	}
	defer resp.Body.Close()

	//read request
	if resp.StatusCode != http.StatusCreated {
		return "", gophKeeperErrors.NewErrWithHTTPCode(resp.StatusCode, fmt.Sprintf("Server`s response has status code `%v`", resp.StatusCode))
	}

	//get jwt from cookies
	for _, cookie := range resp.Cookies() {
		if cookie.Name == JWT_cookie_name {
			if cookie.Value != "" {
				return cookie.Value, nil
			} else {
				return "", fmt.Errorf("JWT cookie is empty")
			}
		}
	}
	return "", fmt.Errorf("no cookies with name `%v`", JWT_cookie_name)
}
