package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/config"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/middleware"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/repository/theapi"
)

type UserHandler struct {
	cfg          config.Config
	taRepository theapi.IRepository
}

func New(router *mux.Router, mw middleware.Middleware, cfg config.Config, taRepository theapi.IRepository) UserHandler {
	h := UserHandler{
		cfg:          cfg,
		taRepository: taRepository,
	}

	authRouter := router.NewRoute().Subrouter()
	authRouter.Use(mw.JWTMiddleware.AuthMiddleware)

	authRouter.HandleFunc("/profile", h.ProfileHandler).Methods(http.MethodGet)

	return h
}

func (h UserHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := h.taRepository.GetProfile(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Repository Register Error")
		// TODO: refactor
		response := map[string]string{"status": "FAILED"}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)
		return
	}

	log.Info().Str("response", fmt.Sprintf("%+v", resp)).Msg("")

	// Write JSON response {"message": "SUCCESS"}
	response := map[string]string{"message": "SUCCESS"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
