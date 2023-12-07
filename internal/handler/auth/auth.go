package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/config"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/middleware"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/model"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/repository/cache"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/repository/theapi"
)

type AuthHandler struct {
	cfg             config.Config
	taRepository    theapi.IRepository
	cacheRepository *cache.Repo
}

func New(router *mux.Router, mw middleware.Middleware, cfg config.Config, taRepository theapi.IRepository, cacheRepository *cache.Repo) AuthHandler {
	h := AuthHandler{
		cfg:             cfg,
		taRepository:    taRepository,
		cacheRepository: cacheRepository,
	}

	router.HandleFunc("/register", h.RegisterHandler).Methods(http.MethodPost)
	router.HandleFunc("/login", h.LoginHandler).Methods(http.MethodPost)

	return h
}

func (m AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Get the uploaded file
	file, _, err := r.FormFile("profile_image")
	if err != nil {
		http.Error(w, "Failed to get profile_image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	resp, err := m.taRepository.Register(ctx, model.CreateRegisterRequest{
		Username:     r.FormValue("username"),
		Password:     r.FormValue("password"),
		FirstName:    r.FormValue("first_name"),
		LastName:     r.FormValue("last_name"),
		Telephone:    r.FormValue("telephone"),
		Address:      r.FormValue("address"),
		City:         r.FormValue("city"),
		Province:     r.FormValue("province"),
		Country:      r.FormValue("country"),
		ProfileImage: file,
	})
	if err != nil {
		log.Error().Err(err).Msg("[handler] register handler")
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

	// TODO: call SaveTokenByUsername to save user token

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

func (m AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Failed to unmarshal JSON body", http.StatusBadRequest)
		return
	}

	resp, err := m.taRepository.GetToken(ctx, model.GetTokenRequest{
		Username: requestBody.Username,
		Password: requestBody.Password,
	})
	if err != nil {
		log.Error().Err(err).Msg("[handler] gettoken handler")
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

	// TODO: save token into cache

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
