package theapi

import (
	"context"
	"io"
	"net/http"

	"github.com/stephanvebrian/go-simple-api-jwt/internal/config"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/model"
)

//go:generate mockery --name IRepository --inpackage --case=underscore
type IRepository interface {
	GetArticles(ctx context.Context, limit, offset int64) (model.GetArticlesResponse, error)
	GetArticleById(ctx context.Context, id int64) (model.GetArticleDetailResponse, error)
	GetToken(ctx context.Context, request model.GetTokenRequest) (model.GetTokenResponse, error)
	GetRefreshToken(ctx context.Context, request model.GetRefreshTokenRequest) (model.GetRefreshTokenResponse, error)
	Register(ctx context.Context, request model.CreateRegisterRequest) (model.CreateRegisterResponse, error)
	GetProfile(ctx context.Context) (model.GetProfileResponse, error)
}

type Repository struct {
	cfg config.Config

	fnParseBody      func(r io.Reader) ([]byte, error)
	fnRequestContext func(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error)
	fnClientDo       func(client *http.Client, req *http.Request) (*http.Response, error)
}

func New(cfg config.Config) Repository {
	return Repository{
		cfg: cfg,

		fnParseBody:      io.ReadAll,
		fnRequestContext: http.NewRequestWithContext,
		fnClientDo: func(client *http.Client, req *http.Request) (*http.Response, error) {
			return client.Do(req)
		},
	}
}
