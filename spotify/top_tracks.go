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

type TopTracks struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Popularity int64    `json:"popularity"`
	Artists    []string `json:"artists"`
}

func getTopTracks(ctx context.Context, limit int32) (Tracks, error) {
	token := keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
	}

	tr := keys.GetContextValue(ctx, keys.ContextSpotifyTimeRange)
	strRange := ""
	if tr != nil {
		castRange, ok := tr.(string)
		if !ok {

			if err, ok := tr.(error); ok {
				// the value wasn't in the context as a string
				logging.GetLogger(ctx).WithError(err).Info("couldnt find time range in context")
			}
		}
		strRange = castRange
	}

	url := fmt.Sprint("https://api.spotify.com/v1/me/top/tracks?limit=", limit)
	if len(strRange) > 0 {
		url += fmt.Sprint("&time_range=", strRange)
	}

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
		Items Tracks `json:"items"`
	}

	var ret tempResp
	err = json.Unmarshal(*body, &ret)
	if err != nil {
		return nil, err
	}

	return ret.Items, nil
}

func getTracks(ctx context.Context, ids []string) (*Tracks, error) {
	token := keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
	}

	url := fmt.Sprint("https://api.spotify.com/v1/tracks?ids=", strings.Join(ids, ","))

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
		Items Tracks `json:"tracks"`
	}

	var ret tempResp
	err = json.Unmarshal(*body, &ret)
	if err != nil {
		return nil, err
	}

	return &ret.Items, nil
}
