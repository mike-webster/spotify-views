package spotify

import (
	"context"
	"errors"
	"fmt"
)

// ContextKey is used to store and access information from the context
type ContextKey string
type ErrTokenExpired string

func (e ErrTokenExpired) Error() string {
	return string(e)
}

var (
	// ContextReturnURL is the key to use for the ouath return url
	ContextReturnURL = ContextKey("return_url")
	// ContextClientID is the key to use for the spotify client id
	ContextClientID = ContextKey("client_id")
	// ContextClientSecret is the key to use for the spotify client secret
	ContextClientSecret = ContextKey("client_secret")
	// ContextAccessToken is the key to use for the spotify access token
	ContextAccessToken = ContextKey("access_token")
	// ContextRefreshToken TODO
	ContextRefreshToken = ContextKey("refresh_token")
	// ContextTimeRange is the key to use for the spotify search
	ContextTimeRange = ContextKey("time_range")
	// ContextResults is the key to use to retrieve the results
	ContextResults = ContextKey("results")
	// ContextLogger is the key to use to retrieve the logger
	ContextLogger = "logger"
	// EventNeedsRefreshToken holds the key to log when a user needs a to
	// refresh their session
	EventNeedsRefreshToken = "token_needs_refresh"
	// EventNon200Response holds the key to log when an external request
	// comes back with a non-200 response
	EventNon200Response = "non_200_response"
)

// HandleOauth is the handler to use for oauth returns from spotify
func HandleOauth(ctx context.Context, code string) (context.Context, error) {
	tokens, err := requestTokens(ctx, code)
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, ContextAccessToken, tokens[0])
	ctx = context.WithValue(ctx, ContextRefreshToken, tokens[1])

	return ctx, nil
}

// GetTopTracks will perform a search for the user's top tracks with the
// provided limit.
func GetTopTracks(ctx context.Context, limit int32) (context.Context, error) {
	tracks, err := getTopTracks(ctx, limit)
	if err != nil {
		return nil, err
	}

	c := context.WithValue(ctx, ContextResults, tracks)
	return c, nil
}

// GetTopArtists will perform a search for the user's top artists
func GetTopArtists(ctx context.Context) (context.Context, error) {
	artists, err := getTopArtists(ctx)
	if err != nil {
		return nil, err
	}

	c := context.WithValue(ctx, ContextResults, *artists)
	return c, nil
}

// GetGenresForArtists will perform a search for the provided artists IDs
// and then researches the genres assosciated with each one. A mapping
// of genres to occurrences is returned.
func GetGenresForArtists(ctx context.Context, ids []string) (context.Context, error) {
	ret := map[string]int32{}
	artists, err := getArtists(ctx, ids)
	if err != nil {
		return nil, err
	}

	for _, i := range *artists {
		for _, ii := range i.Genres {
			if _, ok := ret[ii]; ok {
				ret[ii]++
			} else {
				ret[ii] = 1
			}
		}
	}

	c := context.WithValue(ctx, ContextResults, getPairs(ret))
	return c, nil
}

// GetGenresForTracks will perform a search for the provided track IDs
// and then researches the genres associated with each one. A mapping
// of genres to occurrences is returned.
func GetGenresForTracks(ctx context.Context, ids []string) (context.Context, error) {
	as := map[string]int32{}
	aids := []string{}
	tracks, err := getTracks(ctx, ids)
	if err != nil {
		return nil, err
	}

	for _, i := range *tracks {
		for _, ii := range i.Artists {
			if _, ok := as[ii.Name]; !ok {
				as[ii.Name] = 1
				aids = append(aids, ii.ID)
			}
		}
	}

	if len(aids) < 1 {
		return nil, errors.New(fmt.Sprint("no artists found for ", len(ids), "tracks"))
	}

	artists, err := getArtists(ctx, aids)
	if err != nil {
		return nil, err
	}

	ret := map[string]int32{}
	for _, i := range *artists {
		for _, ii := range i.Genres {
			if _, ok := ret[ii]; ok {
				ret[ii]++
			} else {
				ret[ii] = 1
			}
		}
	}

	c := context.WithValue(ctx, ContextResults, getPairs(ret))
	return c, nil
}

// GetUserInfo will perform a request to retrieve the user's ID and email.
func GetUserInfo(ctx context.Context) (context.Context, error) {
	info, err := getUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	c := context.WithValue(ctx, ContextResults, info)
	return c, nil
}

func GetAudioFeatures(ctx context.Context, ids []string) (*AudioFeatures, error) {
	af, err := getAudioFeatures(ctx, ids)
	if err != nil {
		return nil, err
	}

	return af, nil
}

func GetUserLibraryTracks(ctx context.Context) (Tracks, error) {
	t, err := getUserLibraryTracks(ctx)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// RefreshToken attempts to get a new access token for the user
func RefreshToken(ctx context.Context) (context.Context, error) {
	tok, err := refreshToken(ctx)
	if err != nil {
		return nil, err
	}

	c := context.WithValue(ctx, ContextResults, tok)
	return c, nil
}
