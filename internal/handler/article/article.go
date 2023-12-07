package article

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/config"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/middleware"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/repository/theapi"
)

type ArticleHandler struct {
	cfg          config.Config
	taRepository theapi.Repository
}

func New(router *mux.Router, mw middleware.Middleware, cfg config.Config, taRepository theapi.Repository) ArticleHandler {
	h := ArticleHandler{
		cfg:          cfg,
		taRepository: taRepository,
	}

	authRouter := router.NewRoute().Subrouter()
	authRouter.Use(mw.JWTMiddleware.AuthMiddleware)

	authRouter.HandleFunc("/articles", h.ArticlesHandler).Methods(http.MethodGet)
	authRouter.HandleFunc("/articles/{id}", h.ArticleDetailHandler).Methods(http.MethodGet)

	return h
}

func (h ArticleHandler) ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	queryParams := r.URL.Query()

	limitStr := queryParams.Get("limit")
	offsetStr := queryParams.Get("offset")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
		return
	}

	resp, err := h.taRepository.GetArticles(ctx, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("[handler] articles handler")
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

	// TODO: set articles data to response

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

func (h ArticleHandler) ArticleDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract the 'id' parameter from the request URL
	idStr := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	resp, err := h.taRepository.GetArticleById(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("[handler] article handler")
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

	// TODO: set article detail data to response

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
