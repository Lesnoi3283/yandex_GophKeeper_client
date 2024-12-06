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

func TestHandler_SendText(t *testing.T) {
	type fields struct {
		HTTPClient func(c *gomock.Controller, lt *testing.T) requiredInterfaces.HTTPClient
	}
	type args struct {
		textName string
		text     string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		httpStatus int
	}{
		{
			name: "Valid response",
			fields: fields{
				HTTPClient: func(c *gomock.Controller, lt *testing.T) requiredInterfaces.HTTPClient {
					client := mocks.NewMockHTTPClient(c)
					client.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						assert.Equal(lt, http.MethodPost, req.Method, "request method must be POST")
						assert.Equal(lt, "application/json", req.Header.Get("Content-Type"), "content type must be application/json")

						bodyBytes, err := io.ReadAll(req.Body)
						assert.NoError(lt, err, "can`t read request body")
						expectedBody := `{"text_name":"example text","text":"This is a test text"}`
						assert.JSONEq(lt, expectedBody, string(bodyBytes), "wrong request body")

						responseWriter := httptest.NewRecorder()
						responseWriter.WriteHeader(http.StatusCreated)
						return responseWriter.Result(), nil
					})
					return client
				},
			},
			args: args{
				textName: "example text",
				text:     "This is a test text",
			},
			wantErr: false,
		},
		{
			name: "Bad request response",
			fields: fields{
				HTTPClient: func(c *gomock.Controller, lt *testing.T) requiredInterfaces.HTTPClient {
					client := mocks.NewMockHTTPClient(c)
					client.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
						responseWriter := httptest.NewRecorder()
						responseWriter.WriteHeader(http.StatusBadRequest)
						return responseWriter.Result(), nil
					})
					return client
				},
			},
			args: args{
				textName: "exampleText",
				text:     "This is a test text",
			},
			wantErr: true,
		},
		{
			name: "HTTP client error",
			fields: fields{
				HTTPClient: func(c *gomock.Controller, lt *testing.T) requiredInterfaces.HTTPClient {
					client := mocks.NewMockHTTPClient(c)
					client.EXPECT().Do(gomock.Any()).Return(nil, fmt.Errorf("test error"))
					return client
				},
			},
			args: args{
				textName: "exampleText",
				text:     "This is a test text",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			h := &Requester{
				HTTPClient: tt.fields.HTTPClient(c, t),
			}
			err := h.SendText(tt.args.textName, tt.args.text)

			if tt.wantErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "unexpected error occurred")
			}
		})
	}
}
