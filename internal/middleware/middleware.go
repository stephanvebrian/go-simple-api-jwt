package middleware

import (
	"github.com/stephanvebrian/go-simple-api-jwt/internal/config"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/middleware/jwt"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/repository/cache"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/repository/theapi"
)

type Middleware struct {
	JWTMiddleware jwt.JWTMiddleware
}

func New(cfg config.Config, cacheRepository *cache.Repo, taRepository theapi.Repository) Middleware {
	return Middleware{
		JWTMiddleware: jwt.New(cfg, cacheRepository, taRepository),
	}
}
