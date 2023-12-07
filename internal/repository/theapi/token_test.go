package theapi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stephanvebrian/go-simple-api-jwt/internal/config"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestTheApiRepository_GetToken(t *testing.T) {
	t.Run("success get token", func(t *testing.T) {
		r := New(config.Config{})

		r.fnRequestContext = func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
			return http.NewRequest(method, url, body)
		}
		r.fnClientDo = func(client *http.Client, req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"token": {} }`)),
			}, nil
		}
		r.fnParseBody = func(r io.Reader) ([]byte, error) {
			return io.ReadAll(r)
		}

		_, err := r.GetToken(context.TODO(), model.GetTokenRequest{})
		if err != nil {
			t.Error("unit test failed, please check error")
		}
	})

	t.Run("failed get token", func(t *testing.T) {
		m := New(config.Config{})

		m.fnRequestContext = func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
			return http.NewRequest(method, url, body)
		}
		m.fnClientDo = func(client *http.Client, req *http.Request) (*http.Response, error) {
			return &http.Response{}, errors.New("error occurred")
		}
		m.fnParseBody = func(r io.Reader) ([]byte, error) {
			return io.ReadAll(r)
		}

		_, err := m.GetToken(context.Background(), model.GetTokenRequest{})
		assert.NotNil(t, err)
	})
}

func TestTheApiRepository_GetRefreshToken(t *testing.T) {
	t.Run("success get refresh token", func(t *testing.T) {
		r := New(config.Config{})

		r.fnRequestContext = func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
			return http.NewRequest(method, url, body)
		}
		r.fnClientDo = func(client *http.Client, req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"token": {} }`)),
			}, nil
		}
		r.fnParseBody = func(r io.Reader) ([]byte, error) {
			return io.ReadAll(r)
		}

		_, err := r.GetRefreshToken(context.TODO(), model.GetRefreshTokenRequest{})
		if err != nil {
			t.Error("unit test failed, please check error")
		}
	})

	t.Run("failed get refresh token", func(t *testing.T) {
		m := New(config.Config{})

		m.fnRequestContext = func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
			return http.NewRequest(method, url, body)
		}
		m.fnClientDo = func(client *http.Client, req *http.Request) (*http.Response, error) {
			return &http.Response{}, errors.New("error occurred")
		}
		m.fnParseBody = func(r io.Reader) ([]byte, error) {
			return io.ReadAll(r)
		}

		_, err := m.GetRefreshToken(context.Background(), model.GetRefreshTokenRequest{})
		assert.NotNil(t, err)
	})
}
