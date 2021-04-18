package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mike-webster/spotify-views/keys"
)

const (
	audioFeaturesPageLimit = 100
)

type AudioFeature struct {
	ID               string  `json:"id"`
	URI              string  `json:"uri"`
	Danceability     float32 `json:"danceability"`
	Energy           float32 `json:"energy"`
	Key              int     `json:"key"`
	Loudness         float32 `json:"loudness"`
	Mode             int     `json:"mode"`
	Speechiness      float32 `json:"speechiness"`
	Acousticness     float32 `json:"acousticness"`
	Instrumentalness float32 `json:"instrumentalness"`
	Liveness         float32 `json:"liveness"`
	Valence          float32 `json:"valence"`
	Tempo            float32 `json:"tempo"`
	Duration         int64   `json:"duration_ms"`
	TimeSignature    int     `json:"time_signature"`
}

type AudioFeatures []AudioFeature

// ----
// API
// ----

func GetAudioFeatures(ctx context.Context, ids []string) (*AudioFeatures, error) {
	ret := AudioFeatures{}

	for i := 0; i < len(ids); i += audioFeaturesPageLimit {
		begin, ending := chunkRangeAudioFeatures(i, &ids)
		af, err := getAudioFeatures(ctx, ids[begin:ending])
		if err != nil {
			return nil, err
		}

		ret = append(ret, *af...)
	}

	return &ret, nil
}

// ----
// Members
// ----

// ----
// Helpers
// ----

// chunkRangeAudioFeatures helps determine the pagination to use
// while iterating through ids to retrieve additional "audio feature"
// information.
// params:
// - start: the beginning of the current iteration
// - ids: a reference to the ids being iterated
func chunkRangeAudioFeatures(start int, ids *[]string) (int, int) {
	if len(*ids) <= audioFeaturesPageLimit {
		return 0, len(*ids)
	}

	if start > len(*ids) {
		return 0, audioFeaturesPageLimit
	}

	return start, start + audioFeaturesPageLimit
}

func getAudioFeatures(ctx context.Context, ids []string) (*AudioFeatures, error) {
	req, err := parseRequestForAudioFeatures(ctx, ids)
	if err != nil {
		return nil, err
	}

	body, err := makeRequest(ctx, req, true)
	if err != nil {
		return nil, err
	}

	return parseResponseForAudioFeatures(body)
}

func parseRequestForAudioFeatures(ctx context.Context, ids []string) (*http.Request, error) {
	token := keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)
	if token == nil {
		return nil, ErrNoToken("no access token provided")
	}

	url := fmt.Sprint("https://api.spotify.com/v1/audio-features?ids=", strings.Join(ids, ","))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))
	return req, nil
}

func parseResponseForAudioFeatures(body *[]byte) (*AudioFeatures, error) {
	type tempResp struct {
		Items AudioFeatures `json:"audio_features"`
	}

	ret := tempResp{}
	err := json.Unmarshal(*body, &ret)
	if err != nil {
		return nil, err
	}

	return &ret.Items, nil
}
