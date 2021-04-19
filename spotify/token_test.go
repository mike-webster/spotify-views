package spotify

import (
	"context"
	"testing"

	"github.com/mike-webster/spotify-views/keys"
	"github.com/stretchr/testify/assert"
)

func TestExchangeOauthCode(t *testing.T) {
	t.Run("TestGetCodeSwapRequest", func(t *testing.T) {
		ctx := context.Background()
		t.Run("TestGetCodeSwapParamsError", func(t *testing.T) {
			_, err := getCodeSwapParams(ctx)
			assert.NotNil(t, err)
		})

		t.Run("happy path", func(t *testing.T) {
			ctx = context.WithValue(ctx, keys.ContextSpotifyReturnURL, "/test")
			ctx = context.WithValue(ctx, keys.ContextSpotifyClientID, "test")
			ctx = context.WithValue(ctx, keys.ContextSpotifyClientSecret, "test")
			res, err := getCodeSwapParams(ctx)
			assert.Nil(t, err)
			assert.NotNil(t, res)
		})
	})

	t.Run("TestGetCodeSwapParams", func(t *testing.T) {
		ctx := context.Background()
		defRet := &map[string]string{}
		t.Run("NoReturnUrlError", func(t *testing.T) {
			ret, err := getCodeSwapParams(ctx)
			assert.Equal(t, defRet, ret)
			assert.Equal(t, "no return url provided", err.Error())
		})

		ctx = context.WithValue(ctx, keys.ContextSpotifyReturnURL, &map[string]int{"1": 1})
		t.Run("BadReturnURL", func(t *testing.T) {
			ret, err := getCodeSwapParams(ctx)
			assert.Equal(t, defRet, ret)
			assert.Equal(t, "return url couldn't be parsed", err.Error())
		})

		ctx = context.WithValue(ctx, keys.ContextSpotifyReturnURL, "test")
		t.Run("NoClientIDError", func(t *testing.T) {
			ret, err := getCodeSwapParams(ctx)
			assert.Equal(t, defRet, ret)
			assert.Equal(t, "no client id provided", err.Error())
		})

		ctx = context.WithValue(ctx, keys.ContextSpotifyClientID, &map[string]int{"1": 1})
		t.Run("BadClientIDURL", func(t *testing.T) {
			ret, err := getCodeSwapParams(ctx)
			assert.Equal(t, defRet, ret)
			assert.Equal(t, "client id couldn't be parsed", err.Error())
		})

		ctx = context.WithValue(ctx, keys.ContextSpotifyClientID, "test")
		t.Run("NoClientSecretError", func(t *testing.T) {
			ret, err := getCodeSwapParams(ctx)
			assert.Equal(t, defRet, ret)
			assert.Equal(t, "no client secret provided", err.Error())
		})

		ctx = context.WithValue(ctx, keys.ContextSpotifyClientSecret, &map[string]int{"1": 1})
		t.Run("BadClientSecretURL", func(t *testing.T) {
			ret, err := getCodeSwapParams(ctx)
			assert.Equal(t, defRet, ret)
			assert.Equal(t, "client secret couldn't be parsed", err.Error())
		})

		ctx = context.WithValue(ctx, keys.ContextSpotifyClientSecret, "test")
		t.Run("HappyPath", func(t *testing.T) {
			ret, err := getCodeSwapParams(ctx)
			assert.NotEqual(t, defRet, ret)
			assert.Nil(t, err)
		})
	})

	t.Run("TestParseTokenFromCodeSwapResponse", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			bytes := []byte(getTokenFromSwapPayload)

			as, err := parseTokensFromCodeSwapResponse(&bytes)
			assert.Nil(t, err)
			assert.NotNil(t, as)
		})

		t.Run("bad body", func(t *testing.T) {
			bytes := []byte("fdakslfjda;klfjad;kjadl;")
			_, err := parseTokensFromCodeSwapResponse(&bytes)
			assert.NotNil(t, err)
		})
	})
}

func TestRefreshMe(t *testing.T) {
	t.Run("TestGetRefreshRequest", func(t *testing.T) {
		ctx := context.Background()
		t.Run("NoToken", func(t *testing.T) {
			ret, err := getRefreshRequest(ctx, "")
			assert.Nil(t, ret)
			assert.Equal(t, "no refresh token provided", err.Error())
		})
		t.Run("TestGetRefreshParamsError", func(t *testing.T) {
			ret, err := getRefreshRequest(ctx, "test")
			assert.Nil(t, ret)
			assert.NotNil(t, err)
		})

		ctx = context.WithValue(ctx, keys.ContextSpotifyRefreshToken, "test")
		ctx = context.WithValue(ctx, keys.ContextSpotifyClientID, "test")
		ctx = context.WithValue(ctx, keys.ContextSpotifyClientSecret, "test")

		t.Run("happy path", func(t *testing.T) {
			req, err := getRefreshRequest(ctx, "test")
			assert.Nil(t, err)
			assert.NotNil(t, req)
		})
	})

	t.Run("TestGetRefreshParams", func(t *testing.T) {
		ctx := context.Background()
		t.Run("MissingClientId", func(t *testing.T) {
			ret, err := getRefreshParams(ctx)
			assert.Equal(t, &map[string]string{}, ret)
			assert.Equal(t, "no client id provided", err.Error())
		})
		ctx = context.WithValue(ctx, keys.ContextSpotifyClientID, "Test")
		t.Run("MissingClientSecret", func(t *testing.T) {
			ret, err := getRefreshParams(ctx)
			assert.Equal(t, &map[string]string{}, ret)
			assert.Equal(t, "no client secret provided", err.Error())
		})
		ctx = context.WithValue(ctx, keys.ContextSpotifyClientSecret, "Test")
		t.Run("HappyPath", func(t *testing.T) {
			ret, err := getRefreshParams(ctx)
			assert.Nil(t, err)
			assert.NotEqual(t, &map[string]string{}, ret)
		})
	})

	t.Run("TestParseTokenFromRefreshResponse", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			bytes := []byte(getRefreshPayload)

			as, err := parseTokenFromRefreshResponse(&bytes)
			assert.Nil(t, err)
			assert.NotNil(t, as)
		})

		t.Run("bad body", func(t *testing.T) {
			bytes := []byte("fdakslfjda;klfjad;kjadl;")
			_, err := parseTokenFromRefreshResponse(&bytes)
			assert.NotNil(t, err)
		})
	})

	t.Run("MainMethod", func(t *testing.T) {

	})
}

var (
	getTokenFromSwapPayload = `{
		"access_token": "testtok",
		"token_type": "testtype",
		"scope": "testscope",
		"expires_in": 42143,
		"refresh_token": "testref"
	}`

	getRefreshPayload = `{
		"access_token": "test"
	}`
)
