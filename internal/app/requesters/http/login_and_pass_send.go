package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"yandex_GophKeeper_client/internal/app/entities"
	"yandex_GophKeeper_client/pkg/gophKeeperErrors"
)

// SendLoginAndPassword sends login and password to the backend.
// If http status code != 201 - this func returns a gophKeeperErrors.ErrWithHTTPCode.
func (h *Requester) SendLoginAndPassword(login string, password string) error {
	data := entities.LoginAndPassword{
		Login:    login,
		Password: password,
	}

	//prepare request
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cant marshal login and password, err: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, h.ApiAddress+"/"+saveLoginAndPasswordPath, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("cant create request, err: %w", err)
	}
	req.AddCookie(&http.Cookie{
		Name:  JwtCookieName,
		Value: h.JWT,
	})
	req.Header.Set("Content-Type", "application/json")

	//send request
	resp, err := h.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("cant send request, err: %w", err)
	}
	defer resp.Body.Close()

	//read request
	if resp.StatusCode != http.StatusCreated {
		return gophKeeperErrors.NewErrWithHTTPCode(resp.StatusCode, "Server`s response`s status is not CREATED")
	}
	return nil
}
