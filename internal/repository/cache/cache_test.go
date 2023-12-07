package cache

import (
	"context"
	"fmt"
	"testing"

	"github.com/karlseguin/ccache"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("initialize", func(t *testing.T) {
		got := New()
		assert.NotNil(t, got)
	})
}

func TestRepo_SaveTokenByUsername(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cache := ccache.New(ccache.Configure())

		r := Repo{
			cache: cache,
		}

		r.SaveTokenByUsername(context.TODO(), "example-username", "example-token")

		item := cache.Get(fmt.Sprintf("%s:%s", extTokenKey, "example-username"))
		assert.NotNil(t, item)
	})
}

func TestRepo_GetTokenByUsername(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cache := ccache.New(ccache.Configure())

		r := Repo{
			cache: cache,
		}

		cache.Set(fmt.Sprintf("%s:%s", extTokenKey, "example-username"), "token-1", extTokenExpTime)

		token := r.GetTokenByUsername(context.TODO(), "example-username")
		assert.EqualValues(t, "token-1", token)
	})
}

func TestRepo_SaveRefreshTokenByUsername(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cache := ccache.New(ccache.Configure())

		r := Repo{
			cache: cache,
		}

		r.SaveRefreshTokenByUsername(context.TODO(), "example-username", "example-token")

		item := cache.Get(fmt.Sprintf("%s:%s", extRefreshTokenKey, "example-username"))
		assert.NotNil(t, item)
	})
}

func TestRepo_GetRefreshTokenByUsername(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cache := ccache.New(ccache.Configure())

		r := Repo{
			cache: cache,
		}

		cache.Set(fmt.Sprintf("%s:%s", extRefreshTokenKey, "example-username"), "token-1", extTokenExpTime)

		token := r.GetRefreshTokenByUsername(context.TODO(), "example-username")
		assert.EqualValues(t, "token-1", token)
	})
}
