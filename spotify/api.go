package spotify

import (
	"context"
	"errors"
	"fmt"
)

// ContextKey is used to store and access information from the context
type ContextKey string

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
func GetTopTracks(ctx context.Context, limit int32) (*Tracks, error) {
	tracks, err := getTopTracks(ctx, limit)
	if err != nil {
		return nil, err
	}

	return &tracks, nil
}

// GetTopArtists will perform a search for the user's top artists
func GetTopArtists(ctx context.Context) (*Artists, error) {
	artists, err := getTopArtists(ctx)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

// GetGenresForArtists will perform a search for the provided artists IDs
// and then researches the genres assosciated with each one. A mapping
// of genres to occurrences is returned.
func GetGenresForArtists(ctx context.Context, ids []string) (*Pairs, error) {
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

	p := getPairs(ret)
	return &p, nil
}

// GetGenresForTracks will perform a search for the provided track IDs
// and then researches the genres associated with each one. A mapping
// of genres to occurrences is returned.
func GetGenresForTracks(ctx context.Context, ids []string) (*Pairs, error) {
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

	p := getPairs(ret)
	return &p, nil
}
