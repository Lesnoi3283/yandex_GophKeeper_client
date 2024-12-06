package http

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"yandex_GophKeeper_client/config"
	"yandex_GophKeeper_client/internal/app/entities"
	"yandex_GophKeeper_client/internal/app/requiredInterfaces"
	"yandex_GophKeeper_client/internal/app/requiredInterfaces/mocks"
)

func TestHandler_GetBankCard(t *testing.T) {
	type fields struct {
		Conf       config.AppConfig
		HTTPClient func(c *gomock.Controller) requiredInterfaces.HTTPClient
	}
	type args struct {
		lastFourDigits string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantCard   entities.BankCard
		wantErr    bool
		httpStatus int
		httpBody   string
	}{
		{
			name: "Normal response",
			fields: fields{
				HTTPClient: func(c *gomock.Controller) requiredInterfaces.HTTPClient {
					client := mocks.NewMockHTTPClient(c)
					client.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						assert.Equal(t, http.MethodGet, req.Method, "request method must be GET")
						assert.Equal(t, "text/plain", req.Header.Get("Content-Type"), "content type must be text/plain")

						bodyBytes, err := io.ReadAll(req.Body)
						assert.NoError(t, err, "cant read request body")
						assert.Equal(t, "1234", string(bodyBytes), "wrong request body")

						responseWriter := httptest.NewRecorder()
						responseWriter.WriteHeader(http.StatusOK)
						responseWriter.Write([]byte(`{"PAN":"1234567890123456","expires_at":"12/24","owner_lastname":"Ivanov","owner_firstname":"Ivan"}`))
						return responseWriter.Result(), nil
					})
					return client
				},
			},
			args: args{
				lastFourDigits: "1234",
			},
			wantCard: entities.BankCard{
				PAN:            "1234567890123456",
				ExpiresAt:      "12/24",
				OwnerLastname:  "Ivan",
				OwnerFirstname: "Ivanov",
			},
			wantErr: false,
		},
		{
			name: "Internal server error response",
			fields: fields{
				HTTPClient: func(c *gomock.Controller) requiredInterfaces.HTTPClient {
					client := mocks.NewMockHTTPClient(c)
					client.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						responseWriter := httptest.NewRecorder()
						responseWriter.WriteHeader(http.StatusInternalServerError)
						return responseWriter.Result(), nil
					})
					return client
				},
			},
			args: args{
				lastFourDigits: "1234",
			},
			wantCard: entities.BankCard{},
			wantErr:  true,
		},
		{
			name: "Empty response body",
			fields: fields{
				HTTPClient: func(c *gomock.Controller) requiredInterfaces.HTTPClient {
					client := mocks.NewMockHTTPClient(c)
					client.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						responseWriter := httptest.NewRecorder()
						responseWriter.WriteHeader(http.StatusOK)
						return responseWriter.Result(), nil
					})
					return client
				},
			},
			args: args{
				lastFourDigits: "1234",
			},
			wantCard: entities.BankCard{},
			wantErr:  true,
		},
		{
			name: "Bad json in response",
			fields: fields{
				HTTPClient: func(c *gomock.Controller) requiredInterfaces.HTTPClient {
					client := mocks.NewMockHTTPClient(c)
					client.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						responseWriter := httptest.NewRecorder()
						responseWriter.WriteHeader(http.StatusOK)
						responseWriter.Write([]byte(`{bad json}`))
						return responseWriter.Result(), nil
					})
					return client
				},
			},
			args: args{
				lastFourDigits: "1234",
			},
			wantCard: entities.BankCard{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			h := &Handler{
				HTTPClient: tt.fields.HTTPClient(c),
			}
			gotCard, err := h.GetBankCard(tt.args.lastFourDigits)

			if tt.wantErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "unexpected error occurred")
			}

			assert.Equal(t, tt.wantCard, gotCard, "unexpected card data")
		})
	}
}
