package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"yandex_GophKeeper_client/internal/app/entities"
)

// SendLoginAndPassword sends login and password to the backend and returns error.
// If http status code != 201 - this func returns a gophKeeperErrors.ErrWithHTTPCode.
func (h *Handler) SendLoginAndPassword(login string, password string) error {
	user := entities.User{
		Login:    login,
		Password: password,
	}

}
