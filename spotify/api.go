package spotify

import (
	"context"

	"github.com/mike-webster/spotify-views/keys"
)

type ErrTokenExpired string

func (e ErrTokenExpired) Error() string {
	return string(e)
}

var (

	// EventNeedsRefreshToken holds the key to log when a user needs a to
	// refresh their session
	EventNeedsRefreshToken = "token_needs_refresh"
	// EventNon200Response holds the key to log when an external request
	// comes back with a non-200 response
	EventNon200Response = "non_200_response"
)

// GetUserInfo will perform a request to retrieve the user's ID and email.
func GetUserInfo(ctx context.Context) (context.Context, error) {
	info, err := getUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	c := context.WithValue(ctx, keys.ContextSpotifyResults, info)
	return c, nil
}
