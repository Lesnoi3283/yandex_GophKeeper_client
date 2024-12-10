package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"yandex_GophKeeper_client/internal/app/entities"
	"yandex_GophKeeper_client/pkg/gophKeeperErrors"
)

// Login sends request with user data to API and returns jwt or error.
// If http status code != 200 - this func returns a gophKeeperErrors.ErrWithHTTPCode.
func (h *Requester) Login(login string, password string) (jwt string, err error) {
	user := entities.User{
		Login:    login,
		Password: password,
	}

	//create request
	jsonUser, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("can`t marshal user: %v", err)
	}
	req, err := http.NewRequest(http.MethodGet, h.ApiAddress+"/"+loginPath, bytes.NewReader(jsonUser))
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
	if resp.StatusCode != http.StatusOK {
		return "", gophKeeperErrors.NewErrWithHTTPCode(resp.StatusCode, fmt.Sprintf("Server`s response has status code `%v`", resp.StatusCode))
	}

	//get jwt from cookies
	return getJWTFromCookies(resp.Cookies())
}

func getJWTFromCookies(cookies []*http.Cookie) (string, error) {
	for _, cookie := range cookies {
		if cookie.Name == JwtCookieName {
			return cookie.Value, nil
		}
	}
	return "", fmt.Errorf("no cookies with name `%v`", JwtCookieName)
}
