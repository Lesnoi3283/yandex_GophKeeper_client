package handlers

import (
	"yandex_GophKeeper_client/config"
	"yandex_GophKeeper_client/internal/app/requiredInterfaces"
)

// API paths.
const (
	registration_path            string = "/api/register"
	login_path                   string = "/api/login"
	save_bank_card_path          string = "/api/bankcard"
	get_bank_card_path           string = "/api/bankcard"
	save_login_and_password_path string = "/api/loginandpassword"
	get_login_and_password_path  string = "/api/loginandpassword"
	save_text_path               string = "/api/text"
	get_text_path                string = "/api/text"
)

// Cookie names.
const (
	JWT_cookie_name string = "AuthJWT"
)

type Handler struct {
	Conf       config.AppConfig
	HTTPClient requiredInterfaces.HTTPClient
	JWT        string
	//logger *zap.SugaredLogger //todo: maybe unnecessary. Delete before merge request if it`s still unused.
}
