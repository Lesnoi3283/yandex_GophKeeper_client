package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"yandex_GophKeeper_client/internal/app/entities"
	"yandex_GophKeeper_client/pkg/gophKeeperErrors"
)

// SendText sends text data to the backend.
// If http status code != 201 - this func returns a gophKeeperErrors.ErrWithHTTPCode.
func (h *Requester) SendText(textName string, text string) error {
	data := entities.TextData{
		TextName: textName,
		Text:     text,
	}

	//prepare request
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cant marshal login and password, err: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, h.ApiAddress+"/"+saveTextPath, bytes.NewBuffer(jsonData))
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
