package spotify

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRequest(t *testing.T) {
	ctx := context.Background()
	t.Run("NoDeps", func(t *testing.T) {
		_, err := makeRequest(ctx, nil)

		assert.Equal(t, err, errors.New("couldnt find deps"))
	})

	t.Run("ErrorResponse", func(t *testing.T) {
		deps := getTestDependencies(ctx, 400, "{}")
		_, err := makeRequest(deps, &http.Request{})

		assert.NotEqual(t, nil, err)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		deps := getTestDependencies(ctx, 401, "{}")
		_, err := makeRequest(deps, &http.Request{})

		assert.Equal(t, ErrTokenExpired(""), err)

	})

	t.Run("BadRequest", func(t *testing.T) {
		deps := getTestDependencies(ctx, 400, "{}")
		_, err := makeRequest(deps, &http.Request{})

		assert.Equal(t, ErrBadRequest("response code: 400"), err)

	})
}
