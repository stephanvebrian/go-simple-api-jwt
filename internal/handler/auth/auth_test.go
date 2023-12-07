package auth

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/config"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/model"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/repository/cache"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/repository/theapi"
	"github.com/stretchr/testify/mock"
)

func TestRegisterHandler(t *testing.T) {
	mockTheApiRepository := &theapi.MockIRepository{}
	mockTheApiRepository.On("Register", mock.Anything, mock.Anything).Return(model.CreateRegisterResponse{}, nil)

	userHandler := AuthHandler{
		cfg:             config.Config{},
		taRepository:    mockTheApiRepository,
		cacheRepository: cache.New(),
	}

	router := mux.NewRouter()
	router.HandleFunc("/register", userHandler.RegisterHandler).Methods(http.MethodPost)

	// Create a sample request with a dummy file
	requestBody := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(requestBody)

	fields := map[string]string{
		"username":      "testuser",
		"password":      "testpassword",
		"first_name":    "John",
		"last_name":     "Doe",
		"telephone":     "123456789",
		"address":       "123 Main St",
		"city":          "Anytown",
		"province":      "State",
		"country":       "Country",
		"profile_image": "dummyfile",
	}

	for key, value := range fields {
		writer.WriteField(key, value)
	}

	// Dummy file content
	part, _ := writer.CreateFormFile("profile_image", "dummyfile")
	part.Write([]byte("dummyfilecontent"))

	writer.Close()

	// Create a recorder to capture the response
	rr := httptest.NewRecorder()

	// Simulate an HTTP POST request to "/register" with the sample request
	req, err := http.NewRequest(http.MethodPost, "/register", requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Serve the request to the router
	router.ServeHTTP(rr, req)

	// Check the status code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Decode the response body
	var res map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}

	// Assert that "message" key exists in the response
	if message, ok := res["message"]; !ok {
		t.Error("Response does not contain 'message' key")
	} else {
		// Assert that the value of "message" is "SUCCESS"
		if message != "SUCCESS" {
			t.Errorf("Expected 'message' value to be 'SUCCESS', got '%s'", message)
		}
	}
}

func TestLoginHandler(t *testing.T) {
	router := mux.NewRouter()

	mockTheApiRepository := &theapi.MockIRepository{}
	mockTheApiRepository.On("GetToken", mock.Anything, mock.Anything).Return(model.GetTokenResponse{}, nil)

	authHandler := AuthHandler{
		taRepository: mockTheApiRepository,
	}

	router.HandleFunc("/login", authHandler.LoginHandler).Methods(http.MethodPost)

	requestBody := `{"username": "testuser", "password": "testpassword"}`

	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

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
		if message != "SUCCESS" {
			t.Errorf("Expected 'message' value to be 'SUCCESS', got '%s'", message)
		}
	}
}
