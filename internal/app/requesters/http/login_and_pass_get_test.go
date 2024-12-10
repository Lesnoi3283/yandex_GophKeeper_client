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

func TestHandler_GetLoginAndPassword(t *testing.T) {
	type fields struct {
		Conf       config.AppConfig
		HTTPClient func(c *gomock.Controller, lt *testing.T) requiredInterfaces.HTTPClient
	}
	type args struct {
		login string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantPwd    string
		wantErr    bool
		httpStatus int
		httpBody   string
	}{
		{
			name: "Normal response",
			fields: fields{
				HTTPClient: func(c *gomock.Controller, lt *testing.T) requiredInterfaces.HTTPClient {
					client := mocks.NewMockHTTPClient(c)
					client.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						assert.Equal(lt, http.MethodGet, req.Method, "request method must be GET")
						assert.Equal(lt, "text/plain", req.Header.Get("Content-Type"), "content type must be text/plain")

						bodyBytes, err := io.ReadAll(req.Body)
						assert.NoError(lt, err, "cant read request body")
						assert.Equal(lt, "testlogin@example.com", string(bodyBytes), "wrong request body")

						responseWriter := httptest.NewRecorder()
						responseWriter.WriteHeader(http.StatusOK)
						responseWriter.Write([]byte("testpassword"))
						return responseWriter.Result(), nil
					})
					return client
				},
			},
			args: args{
				login: "testlogin@example.com",
			},
			wantPwd: "testpassword",
			wantErr: false,
		},
		{
			name: "Internal server error response",
			fields: fields{
				HTTPClient: func(c *gomock.Controller, lt *testing.T) requiredInterfaces.HTTPClient {
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
				login: "testlogin@example.com",
			},
			wantPwd: "",
			wantErr: true,
		},
		{
			name: "Empty response body",
			fields: fields{
				HTTPClient: func(c *gomock.Controller, lt *testing.T) requiredInterfaces.HTTPClient {
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
				login: "testlogin@example.com",
			},
			wantPwd: "",
			wantErr: true,
		},
		{
			name: "HTTP Client error",
			fields: fields{
				HTTPClient: func(c *gomock.Controller, lt *testing.T) requiredInterfaces.HTTPClient {
					client := mocks.NewMockHTTPClient(c)
					client.EXPECT().Do(gomock.Any()).Return(nil, fmt.Errorf("test error"))
					return client
				},
			},
			args: args{
				login: "testlogin@example.com",
			},
			wantPwd: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			h := &Requester{
				HTTPClient: tt.fields.HTTPClient(c, t),
			}
			gotPwd, err := h.GetLoginAndPassword(tt.args.login)

			if tt.wantErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "unexpected error occurred")
			}

			assert.Equal(t, tt.wantPwd, gotPwd, "unexpected password")
		})
	}
}
