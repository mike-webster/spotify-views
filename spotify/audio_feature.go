package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/mike-webster/spotify-views/keys"
	"github.com/mike-webster/spotify-views/logging"
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

func (af AudioFeature) String() string {
	str, err := json.Marshal(af)
	if err != nil {
		return ""
	}
	return string(str)
}

func (af *AudioFeatures) String() []string {
	ret := []string{}
	for _, i := range *af {
		ret = append(ret, i.String())
	}
	return ret
}

func getAudioFeatures(ctx context.Context, ids []string) (*AudioFeatures, error) {
	logger := logging.GetLogger(ctx)
	token := keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
	}

	if len(ids) > 100 {
		logger.WithField("count", len(ids)).Warn("too many tracks passed, reducing to the first 100")
		ids = ids[:100]
	}

	url := fmt.Sprint("https://api.spotify.com/v1/audio-features?ids=", strings.Join(ids, ","))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	body, err := makeRequest(ctx, req, true)
	if err != nil {
		return nil, err
	}

	type tempResp struct {
		Items AudioFeatures `json:"audio_features"`
	}

	ret := tempResp{}
	err = json.Unmarshal(*body, &ret)
	if err != nil {
		logger.WithField("body", string(*body)).Error(err.Error())
		return nil, err
	}

	return &ret.Items, nil
}
