package http

import (
	"yandex_GophKeeper_client/internal/app/requiredInterfaces"
)

// API paths.
const (
	registrationPath         string = "api/register"
	loginPath                string = "api/login"
	saveBankCardPath         string = "api/bankcard"
	getBankCardPath          string = "api/bankcard"
	saveLoginAndPasswordPath string = "api/loginandpassword"
	getLoginAndPasswordPath  string = "api/loginandpassword"
	saveTextPath             string = "api/text"
	getTextPath              string = "api/text"
)

// Cookie names.
const (
	JwtCookieName string = "AuthJWT"
)

// Requester prepares and sends request to the backend server.
type Requester struct {
	//Conf       config.AppConfig
	HTTPClient requiredInterfaces.HTTPClient
	JWT        string
	ApiAddress string
}

func NewRequester(apiAddress string, httpClient requiredInterfaces.HTTPClient, jwt string) *Requester {
	return &Requester{
		//Conf:       conf,
		HTTPClient: httpClient,
		JWT:        jwt,
		ApiAddress: apiAddress,
	}
}
