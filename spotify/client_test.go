package spotify

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRequest(t *testing.T) {
	ctx := context.Background()
	t.Run("NoDeps", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:3000", nil)
		_, err := makeRequest(ctx, req)

		assert.Equal(t, err, errors.New("couldnt find deps"))
	})

	t.Run("ErrorResponse", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:3000", nil)
		deps := getTestDependencies(ctx, 400, "{}")
		_, err := makeRequest(deps, req)

		assert.NotEqual(t, nil, err)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:3000", nil)
		deps := getTestDependencies(ctx, 401, "{}")
		_, err := makeRequest(deps, req)

		assert.Equal(t, ErrTokenExpired(""), err)

	})

	t.Run("BadRequest", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:3000", nil)
		deps := getTestDependencies(ctx, 400, "{}")
		_, err := makeRequest(deps, req)

		assert.Equal(t, ErrBadRequest("response code: 400"), err)

	})
}
