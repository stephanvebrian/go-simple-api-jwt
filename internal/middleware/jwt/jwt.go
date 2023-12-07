package jwt

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/config"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/model"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/repository/cache"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/repository/theapi"
)

var secretKey = []byte("lUuuULKmdd7bhYkS5PwQFpnsYtYb6M")

// User represents a user in the system
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims represents the JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type JWTMiddleware struct {
	cfg          config.Config
	cacheRepo    *cache.Repo
	taRepository theapi.Repository
}

type JWTMiddlewareInterface interface {
	AuthMiddleware(next http.Handler) http.Handler
}

func New(cfg config.Config, cacheRepository *cache.Repo, taRepository theapi.Repository) JWTMiddleware {
	return JWTMiddleware{
		cfg:          cfg,
		cacheRepo:    cacheRepository,
		taRepository: taRepository,
	}
}

func (m *JWTMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msg("[auth] accessing jwt middleware")

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization token is missing", http.StatusUnauthorized)
			return
		}

		// Check for the "TSTMY" prefix
		const prefix = "TSTMY"
		if !strings.HasPrefix(tokenString, prefix+" ") {
			http.Error(w, "Invalid token prefix", http.StatusUnauthorized)
			return
		}

		tokenString = tokenString[len(prefix)+1:]

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %+v", err), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		extToken := m.cacheRepo.GetTokenByUsername(context.Background(), claims.Username)
		if extToken == "" {
			extRefreshToken := m.cacheRepo.GetRefreshTokenByUsername(context.TODO(), claims.Username)
			_, _ = m.taRepository.GetRefreshToken(context.TODO(), model.GetRefreshTokenRequest{
				Refresh: extRefreshToken,
			})
			// TODO: res-set extToken from GetRefreshToken
		}

		// Token is valid, proceed with the next handler
		ctx := context.WithValue(r.Context(), model.ContextUsernameKey, claims.Username)
		ctx = context.WithValue(ctx, model.ContextExtApiTokenKey, extToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GenerateToken(username string) (string, error) {
	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
