package main

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/stephanvebrian/go-simple-api-jwt/internal/config"
	articleHandler "github.com/stephanvebrian/go-simple-api-jwt/internal/handler/article"
	authHandler "github.com/stephanvebrian/go-simple-api-jwt/internal/handler/auth"
	userHandler "github.com/stephanvebrian/go-simple-api-jwt/internal/handler/user"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/middleware"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/repository/cache"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/repository/theapi"
)

func startServer(cfg config.Config) error {
	router := mux.NewRouter()

	cacheRepository := cache.New()
	taRepository := theapi.New(cfg)

	cacheRepository.SaveTokenByUsername(context.TODO(), "admin@gmail.com", "HelloThisIsToken")

	mw := middleware.New(cfg, cacheRepository, taRepository)

	authHandler.New(router, mw, cfg, taRepository, cacheRepository)
	userHandler.New(router, mw, cfg, taRepository)
	articleHandler.New(router, mw, cfg, taRepository)

	return http.ListenAndServe(cfg.Port, router)
}
