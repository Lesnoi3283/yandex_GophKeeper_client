package http

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"yandex_GophKeeper_client/config"
	"yandex_GophKeeper_client/internal/app/requiredInterfaces"
	"yandex_GophKeeper_client/internal/app/requiredInterfaces/mocks"
)

func TestHandler_SendBankCard(t *testing.T) {
	type fields struct {
		Conf       config.AppConfig
		HTTPClient func(c *gomock.Controller) requiredInterfaces.HTTPClient
	}
	type args struct {
		PAN            string
		OwnerFirstName string
		OwnerLastName  string
		ExpiresAt      string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
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
						assert.Equal(t, http.MethodPost, req.Method, "request method must be POST")
						assert.Equal(t, "application/json", req.Header.Get("Content-Type"), "content type must be application/json")

						bodyBytes, err := io.ReadAll(req.Body)
						assert.NoError(t, err, "cant read request body")
						expectedBody := `{"PAN":"1234567890123456","expires_at":"12/24","owner_lastname":"Ivanov","owner_firstname":"Ivan"}`
						assert.Equal(t, expectedBody, string(bodyBytes), "wrong request body")

						responseWriter := httptest.NewRecorder()
						responseWriter.WriteHeader(http.StatusCreated)
						return responseWriter.Result(), nil
					})
					return client
				},
			},
			args: args{
				PAN:            "1234567890123456",
				OwnerFirstName: "Ivan",
				OwnerLastName:  "Ivanov",
				ExpiresAt:      "12/24",
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
				PAN:            "1234567890123456",
				OwnerFirstName: "Ivan",
				OwnerLastName:  "Ivanov",
				ExpiresAt:      "12/24",
			},
			wantErr: true,
		},
		{
			name: "HTTP Client error",
			fields: fields{
				HTTPClient: func(c *gomock.Controller) requiredInterfaces.HTTPClient {
					client := mocks.NewMockHTTPClient(c)
					client.EXPECT().Do(gomock.Any()).Return(nil, fmt.Errorf("test error"))
					return client
				},
			},
			args: args{
				PAN:            "1234567890123456",
				OwnerFirstName: "Ivan",
				OwnerLastName:  "Ivanov",
				ExpiresAt:      "12/24",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			h := &Handler{
				HTTPClient: tt.fields.HTTPClient(c),
			}
			err := h.SendBankCard(tt.args.PAN, tt.args.OwnerFirstName, tt.args.OwnerLastName, tt.args.ExpiresAt)

			if tt.wantErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "unexpected error occurred")
			}
		})
	}
}
