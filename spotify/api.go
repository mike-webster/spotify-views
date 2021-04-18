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

func GetAudioFeatures(ctx context.Context, ids []string) (*AudioFeatures, error) {
	if len(ids) > 100 {
		// we need pagination
		ret := AudioFeatures{}
		for i := 0; i < len(ids); i += 100 {
			b := i
			e := i + 100
			if e > len(ids) {
				e = len(ids)
			}
			cids := ids[b:e]
			af, err := getAudioFeatures(ctx, cids)
			if err != nil {
				return nil, err
			}
			ret = append(ret, *af...)
		}
		return &ret, nil
	}

	af, err := getAudioFeatures(ctx, ids)
	if err != nil {
		return nil, err
	}

	return af, nil
}

// GetTopArtists will perform a search for the user's top artists
func GetTopArtists(ctx context.Context) (*Artists, error) {
	artists, err := getTopArtists(ctx)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

// GetRecommendations will perform a request to  retrieve spotify's recommendations for the user
func GetRecommendations(ctx context.Context, seeds map[string][]string) (*Recommendation, error) {
	return getRecommendations(ctx, seeds)
}

// GetRelatedArtists will  perform a request for artists  that spotify believes to be similar
// to the  artist with  the provided ID
func GetRelatedArtists(ctx context.Context, id string) (*[]Artist, error) {
	return getRelatedArtists(ctx, id)
}

// GetUserInfo will perform a request to retrieve the user's ID and email.
func GetUserInfo(ctx context.Context) (context.Context, error) {
	info, err := getUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	c := context.WithValue(ctx, keys.ContextSpotifyResults, info)
	return c, nil
}
