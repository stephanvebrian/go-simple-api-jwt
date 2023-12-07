package theapi

import (
	"testing"

	"github.com/stephanvebrian/go-simple-api-jwt/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("initialize", func(t *testing.T) {
		got := New(config.Config{})
		assert.NotNil(t, got)
	})
}
