package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mike-webster/spotify-views/keys"
)

const (
	fetchLimitTopTracks = 25
)

// Track represents a spotify track
type Track struct {
	Links      map[string]string `json:"external_urls"`
	Name       string            `json:"name"`
	URI        string            `json:"uri"`
	ID         string            `json:"id"`
	Popularity int64             `json:"popularity"`
	Artists    []Artist          `json:"artists"`
	Album      Album             `json:"album"`
}

// Tracks is a collection of spotify Tracks
type Tracks []Track

func GetTopTracks(ctx context.Context, timeframe TimeFrame) (*Tracks, error) {
	req, err := getTopTracksRequest(ctx, timeframe)
	if err != nil {
		return nil, err
	}

	body, err := makeRequest(ctx, req, true)
	if err != nil {
		return nil, err
	}

	return parseTopTrackResponse(body)
}

func getTopTracksRequest(ctx context.Context, timeframe TimeFrame) (*http.Request, error) {
	token := keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
	}

	url := fmt.Sprint("https://api.spotify.com/v1/me/top/tracks?limit=", fetchLimitTopTracks)
	url += fmt.Sprint("&time_range=", timeframe.Value())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	return req, nil
}

func parseTopTrackResponse(body *[]byte) (*Tracks, error) {
	type tempResp struct {
		Items Tracks `json:"items"`
	}

	var ret tempResp
	err := json.Unmarshal(*body, &ret)
	if err != nil {
		return nil, err
	}

	return &ret.Items, nil
}

// EmbeddedPlayer will return the html to use for rendering the embedded spotify
// player iframe
func (t *Track) EmbeddedPlayer() string {
	return fmt.Sprintf(`<iframe src="https://open.spotify.com/embed/track/%s" width="300" height="80" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>`, t.ID)
}

// IDs returns the ID for each of the tracks in the collection of Tracks
func (t *Tracks) IDs() []string {
	ret := []string{}
	for _, i := range *t {
		ret = append(ret, i.ID)
	}
	return ret
}

func (t *Track) FindArtist() string {
	if len(t.Artists) < 1 {
		return ""
	}

	return t.Artists[0].Name
}

func (t *Track) FindImage() *Image {
	if len(t.Album.Images) < 1 {
		//return nil
		return &Image{URL: ""}
	}

	if len(t.Album.Images) == 1 {
		return &t.Album.Images[0]
	}

	return &t.Album.Images[1]
}

func (t *Track) TrySpotifyURL() string {
	if len(t.Links) > 0 {
		return t.Links["spotify"]
	}

	return ""
}
