package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"yandex_GophKeeper_client/internal/app/entities"
	"yandex_GophKeeper_client/pkg/gophKeeperErrors"
)

// GetBankCard sends last 4 digits of a bank card to the backend and returns password or error.
// If http status code != 200 - returns a gophKeeperErrors.ErrWithHTTPCode.
func (h *Requester) GetBankCard(lastFourDigits string) (entities.BankCard, error) {
	//prepare request
	req, err := http.NewRequest(http.MethodGet, h.ApiAddress+"/"+getBankCardPath, bytes.NewBufferString(lastFourDigits))
	if err != nil {
		return entities.BankCard{}, fmt.Errorf("cant create request, err: %w", err)
	}
	req.AddCookie(&http.Cookie{
		Name:  JwtCookieName,
		Value: h.JWT,
	})
	req.Header.Set("Content-Type", "text/plain")

	//send request
	resp, err := h.HTTPClient.Do(req)
	if err != nil {
		return entities.BankCard{}, fmt.Errorf("cant send request, err: %w", err)
	}
	defer resp.Body.Close()

	//read request
	if resp.StatusCode != http.StatusOK {
		return entities.BankCard{}, gophKeeperErrors.NewErrWithHTTPCode(resp.StatusCode, "Server`s response`s status is not OK")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return entities.BankCard{}, fmt.Errorf("cant read response body, err: %w", err)
	}
	if len(bodyBytes) == 0 {
		return entities.BankCard{}, fmt.Errorf("empty response body")
	}

	bankCard := entities.BankCard{}
	err = json.Unmarshal(bodyBytes, &bankCard)
	if err != nil {
		return entities.BankCard{}, fmt.Errorf("cant unmarshal response body, err: %w", err)
	}
	return bankCard, nil
}
