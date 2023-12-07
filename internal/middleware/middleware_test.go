package middleware

import (
	"testing"

	"github.com/stephanvebrian/go-simple-api-jwt/internal/config"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/repository/cache"
	"github.com/stephanvebrian/go-simple-api-jwt/internal/repository/theapi"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("initialize", func(t *testing.T) {
		got := New(config.Config{}, cache.New(), theapi.New(config.Config{}))
		assert.NotNil(t, got)
	})
}
