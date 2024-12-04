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
func (h *Handler) SendText(textName string, text string) error {
	data := entities.TextData{
		TextName: textName,
		Text:     text,
	}

	//prepare request
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cant marshal login and password, err: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, save_text_path, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("cant create request, err: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	//send request
	resp, err := http.DefaultClient.Do(req)
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
