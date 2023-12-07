package theapi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stephanvebrian/go-simple-api-jwt/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetArticles(t *testing.T) {
	t.Run("success get articles", func(t *testing.T) {
		m := New(config.Config{})

		m.fnRequestContext = func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
			return http.NewRequest(method, url, body)
		}
		m.fnClientDo = func(client *http.Client, req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"articles": [{}, {}]}`)),
			}, nil
		}
		m.fnParseBody = func(r io.Reader) ([]byte, error) {
			return io.ReadAll(r)
		}

		_, err := m.GetArticles(context.Background(), 0, 0)
		if err != nil {
			t.Error("unit test failed, please check error")
		}
	})

	t.Run("failed get articles", func(t *testing.T) {
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

		_, err := m.GetArticles(context.Background(), 0, 0)
		assert.NotNil(t, err)
	})
}

func TestGetArticleById(t *testing.T) {
	t.Run("success get article by id", func(t *testing.T) {
		m := New(config.Config{})

		m.fnRequestContext = func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
			return http.NewRequest(method, url, body)
		}
		m.fnClientDo = func(client *http.Client, req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"article": {}}`)),
			}, nil
		}
		m.fnParseBody = func(r io.Reader) ([]byte, error) {
			return io.ReadAll(r)
		}

		_, err := m.GetArticleById(context.Background(), 123)
		if err != nil {
			t.Error("unit test failed, please check error")
		}
	})

	t.Run("failed to get article by id", func(t *testing.T) {
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

		_, err := m.GetArticleById(context.Background(), 123)
		assert.NotNil(t, err)
	})
}
