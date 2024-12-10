package http

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"yandex_GophKeeper_client/internal/app/requiredInterfaces"
	"yandex_GophKeeper_client/internal/app/requiredInterfaces/mocks"
)

func TestHandler_RegisterUser(t *testing.T) {
	type fields struct {
		HTTPClient func(c *gomock.Controller, lt *testing.T) requiredInterfaces.HTTPClient
	}
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantJwt string
		wantErr bool
	}{
		{
			name: "Normal",
			fields: fields{
				HTTPClient: func(c *gomock.Controller, lt *testing.T) requiredInterfaces.HTTPClient {
					client := mocks.NewMockHTTPClient(c)
					client.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						assert.Equal(lt, http.MethodPost, req.Method, "request method must be POST")
						assert.Equal(lt, "application/json", req.Header.Get("Content-Type"), "content type must be application/json")
						bodyBytes, err := io.ReadAll(req.Body)
						assert.NoError(lt, err, "cant read request body")
						body := string(bodyBytes)
						assert.Equal(lt, `{"login":"testlogin@example.com","password":"testpassword"}`, body, "wrong request body")

						responseWriter := httptest.NewRecorder()
						http.SetCookie(responseWriter, &http.Cookie{
							Name:  JwtCookieName,
							Value: "test.jwt.token",
						})
						responseWriter.WriteHeader(http.StatusCreated)
						return responseWriter.Result(), nil
					})
					return client
				},
			},
			args: args{
				login:    "testlogin@example.com",
				password: "testpassword",
			},
			wantJwt: "test.jwt.token",
			wantErr: false,
		},
		{
			name: "Internal Server Error (no cookie)",
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
				login:    "testlogin@example.com",
				password: "testpassword",
			},
			wantJwt: "",
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
				login:    "testlogin@example.com",
				password: "testpassword",
			},
			wantJwt: "",
			wantErr: true,
		},
		{
			name: "Status created, but no cookie",
			fields: fields{
				HTTPClient: func(c *gomock.Controller, lt *testing.T) requiredInterfaces.HTTPClient {
					client := mocks.NewMockHTTPClient(c)
					client.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						responseWriter := httptest.NewRecorder()
						responseWriter.WriteHeader(http.StatusCreated)
						return responseWriter.Result(), nil
					})
					return client
				},
			},
			args: args{
				login:    "testlogin@example.com",
				password: "testpassword",
			},
			wantJwt: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			h := &Requester{
				HTTPClient: tt.fields.HTTPClient(c, t),
			}
			gotJwt, err := h.RegisterUser(tt.args.login, tt.args.password)
			if tt.wantErr {
				assert.Error(t, err, "no error, but have to")
			} else {
				assert.NoError(t, err, "some error have happened")
			}

			assert.Equal(t, tt.wantJwt, gotJwt, "wrong jwt value")
		})
	}
}
