package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/karlseguin/ccache"
)

//go:generate mockgen -package=cache -source=cache.go -destination=cache_mock_test.go
type Repo struct {
	cache *ccache.Cache
}

const (
	// ext token
	extTokenKey     = "ext_token"
	extTokenExpTime = time.Duration(25) * time.Minute // 25m

	// refresh token
	extRefreshTokenKey      = "ext_refresh_token"
	extRefreshgTokenExpTime = time.Duration(50) * time.Minute // 50m
)

func New() *Repo {
	cache := ccache.New(ccache.Configure())

	return &Repo{
		cache: cache,
	}
}

func (r Repo) SaveTokenByUsername(ctx context.Context, username string, extToken string) {
	cacheKey := fmt.Sprintf("%s:%s", extTokenKey, username)
	r.cache.Set(cacheKey, extToken, extTokenExpTime)
}

func (r Repo) GetTokenByUsername(ctx context.Context, username string) string {
	cacheKey := fmt.Sprintf("%s:%s", extTokenKey, username)
	item := r.cache.Get(cacheKey)
	if item == nil {
		return ""
	}

	strToken, ok := item.Value().(string)
	if !ok {
		return ""
	}

	return strToken
}

func (r Repo) SaveRefreshTokenByUsername(ctx context.Context, username string, extRefreshToken string) {
	cacheKey := fmt.Sprintf("%s:%s", extRefreshTokenKey, username)
	r.cache.Set(cacheKey, extRefreshToken, extRefreshgTokenExpTime)
}

func (r Repo) GetRefreshTokenByUsername(ctx context.Context, username string) string {
	cacheKey := fmt.Sprintf("%s:%s", extRefreshTokenKey, username)
	item := r.cache.Get(cacheKey)
	if item == nil {
		return ""
	}

	strToken, ok := item.Value().(string)
	if !ok {
		return ""
	}

	return strToken
}
