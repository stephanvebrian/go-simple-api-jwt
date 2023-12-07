package jwt

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stephanvebrian/go-simple-api-jwt/internal/config"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/repository/cache"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/repository/theapi"
)

func TestJWTMiddleware_AuthMiddleware(t *testing.T) {
	jwtMiddleware := New(config.Config{}, cache.New(), theapi.New(config.Config{}))

	// Mock HTTP handler for the next middleware or endpoint
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.Context().Value("username")
		if username == nil {
			t.Error("Expected username in context, but it is missing")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Implement your expected behavior based on the username, if needed
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		authorization  string
		expectedStatus int
	}{
		{
			name: "ValidToken",
			authorization: func() string {
				testtoken, _ := GenerateToken("testest")
				return fmt.Sprintf("TSTMY %s", testtoken)
			}(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "MissingToken",
			authorization:  "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "InvalidTokenPrefix",
			authorization:  "InvalidPrefix validtoken",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "InvalidToken",
			authorization:  "TSTMYinvalidtoken",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", tt.authorization)

			rr := httptest.NewRecorder()

			authHandler := jwtMiddleware.AuthMiddleware(mockHandler)
			authHandler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v, want %v", status, tt.expectedStatus)
			}
		})
	}
}
