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

func TestTheApiRepository_Register(t *testing.T) {
	t.Run("success register", func(t *testing.T) {
		r := New(config.Config{})

		fnIOCopy = func(dst io.Writer, src io.Reader) (written int64, err error) {
			mockContent := "sample content"
			io.Copy(dst, strings.NewReader(mockContent))
			return 1, nil
		}

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

		_, err := r.Register(context.TODO(), model.CreateRegisterRequest{})
		if err != nil {
			t.Error("unit test failed, please check error")
		}
	})

	t.Run("failed register", func(t *testing.T) {
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

		_, err := m.Register(context.TODO(), model.CreateRegisterRequest{})
		assert.NotNil(t, err)
	})
}

func TestTheApiRepository_GetProfile(t *testing.T) {
	t.Run("success get profile", func(t *testing.T) {
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

		_, err := r.GetProfile(context.TODO())
		if err != nil {
			t.Error("unit test failed, please check error")
		}
	})

	t.Run("failed get profile", func(t *testing.T) {
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

		_, err := m.GetProfile(context.TODO())
		assert.NotNil(t, err)
	})
}
