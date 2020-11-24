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
		strRange = tr.(string)
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

func getTopTracksForArtist(ctx context.Context, id string) (*[]TopTracks, error) {
	token := keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/artists/%v/top-tracks?country=us", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	body, err := makeRequest(ctx, req, true)
	if err != nil {
		return nil, err
	}

	type tRes struct {
		Tracks []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Popularity int64  `json:"popularity"`
			Artists    []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Link string `json:"href"`
			} `json:"artists"`
		} `json:"tracks"`
	}

	rsp := tRes{}
	err = json.Unmarshal(*body, &rsp)
	if err != nil {
		return nil, err
	}

	ret := []TopTracks{}
	for _, i := range rsp.Tracks {
		artists := []string{}
		for _, j := range i.Artists {
			artists = append(artists, j.Name)
		}
		ret = append(ret, TopTracks{
			ID:         i.ID,
			Name:       i.Name,
			Popularity: i.Popularity,
			Artists:    artists,
		})
	}

	return &ret, nil
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

func getUserLibraryTracks(ctx context.Context) (Tracks, error) {
	logger := logging.GetLogger(nil)
	url := "https://api.spotify.com/v1/me/tracks?limit=50&offset=0"
	more := true
	ret := []Track{}
	for more {
		t, newUrl, tot, err := getChunkOfUserLibraryTracks(ctx, url)
		if err != nil {
			logger.Warn(err.Error())
			more = false
		}
		url = newUrl

		ret = append(ret, t...)
		if tot == len(ret) {
			more = false
		}
	}
	return ret, nil
}
