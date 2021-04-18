package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/mike-webster/spotify-views/keys"
)

// Recommendation holds information about a track that spotify is recommending
type Recommendation struct {
	Tracks []Track `json:"tracks"`
	Seeds  []struct {
		ID   string `json:"id"`
		Link string `json:"href"`
		Type string `json:"type"`
	} `json:"seeds"`
}

// ----
// API
// ---

// GetRecommendations will perform a request to  retrieve spotify's recommendations for the user
func GetRecommendations(ctx context.Context, seeds map[string][]string) (*Recommendation, error) {
	req, err := parseRecommendationsRequest(ctx, seeds)
	if err != nil {
		return nil, err
	}

	body, err := makeRequest(ctx, req, false)
	if err != nil {
		return nil, err
	}

	return parseReommendationsResponse(body)
}

// ----
// Members
// ----

// ----
// Helpers
// ----

func getRecommendationsURL(ctx context.Context, seeds map[string][]string) string {
	qs := "?"
	if artists, ok := seeds[KeySeedArtists]; ok {
		qs = fmt.Sprint(qs, "seed_artists=", strings.Join(artists, ","))
	}

	if tracks, ok := seeds[KeySeedTracks]; ok {
		if len(qs) > 1 {
			qs += "&"
		}

		qs = fmt.Sprint(qs, "seed_tracks=", strings.Join(tracks, ","))
	}

	if genres, ok := seeds[KeySeedGenres]; ok {
		if len(qs) > 1 {
			qs += "&"
		}

		qs = fmt.Sprint(qs, "seed_genres=", strings.Join(genres, ","))
	}
	return fmt.Sprint("https://api.spotify.com/v1/recommendations", qs)
}

func parseRecommendationsRequest(ctx context.Context, seeds map[string][]string) (*http.Request, error) {
	token := keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
	}

	url := getRecommendationsURL(ctx, seeds)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	return req, nil
}

func parseReommendationsResponse(body *[]byte) (*Recommendation, error) {
	type tApiResponse struct {
		Tracks []Track `json:"tracks"`
		Seeds  []struct {
			InitialPoolSize    int64  `json:"initialPoolSize"`
			AfterFilteringSize int64  `json:"afterFilteringSize"`
			Link               string `json:"href"`
			ID                 string `json:"id"`
			Type               string `json:"type"`
		} `json:"seeds"`
	}

	var rsp tApiResponse
	err := json.Unmarshal(*body, &rsp)
	if err != nil {
		return nil, err
	}

	ret := Recommendation{Tracks: rsp.Tracks}

	for _, i := range rsp.Seeds {
		ret.Seeds = append(ret.Seeds, struct {
			ID   string "json:\"id\""
			Link string "json:\"href\""
			Type string "json:\"type\""
		}{
			ID:   i.ID,
			Link: i.Link,
			Type: i.Type,
		})
	}

	return &ret, nil
}
