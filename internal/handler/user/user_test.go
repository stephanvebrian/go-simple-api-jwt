package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/config"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/model"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/repository/theapi"
	"github.com/stretchr/testify/mock"
)

func TestProfileHandler(t *testing.T) {
	mockTheApiRepository := &theapi.MockIRepository{}
	mockTheApiRepository.On("GetProfile", mock.Anything).Return(model.GetProfileResponse{}, nil)

	userHandler := UserHandler{
		cfg:          config.Config{},
		taRepository: mockTheApiRepository,
	}

	router := mux.NewRouter()

	router.HandleFunc("/profile", userHandler.ProfileHandler).Methods(http.MethodGet)

	req, err := http.NewRequest("GET", "/profile", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var res map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}

	if message, ok := res["message"]; !ok {
		t.Error("Response does not contain 'message' key")
	} else {
		// Assert that the value of "message" is "SUCCESS"
		if message != "SUCCESS" {
			t.Errorf("Expected 'message' value to be 'SUCCESS', got '%s'", message)
		}
	}
}
