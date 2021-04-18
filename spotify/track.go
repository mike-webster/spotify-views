package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mike-webster/spotify-views/keys"
	"github.com/mike-webster/spotify-views/logging"
	"github.com/mike-webster/spotify-views/sortablemap"
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

// ----
// API
// ---

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

func GetTopTracksForArtist(ctx context.Context, id string) (*Tracks, error) {
	req, err := getTopTracksForArtistRequest(ctx, id)
	if err != nil {
		return nil, err
	}

	body, err := makeRequest(ctx, req, true)
	if err != nil {
		return nil, err
	}

	return parseTopTracksForArtistResponse(body)
}

func GetSavedTracks(ctx context.Context) (*Tracks, error) {
	url := "https://api.spotify.com/v1/me/tracks?limit=50&offset=0"
	more := true
	ret := Tracks{}
	for more {
		t, newUrl, tot, err := getChunkOfUserLibraryTracks(ctx, url)
		if err != nil {
			logging.GetLogger(ctx).Warn(err.Error())
			more = false
		}
		url = newUrl

		ret = append(ret, t...)
		if tot == len(ret) {
			more = false
		}
	}
	return &ret, nil
}

// ----
// Helpers
// ----

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

func getTopTracksForArtistRequest(ctx context.Context, id string) (*http.Request, error) {
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
	return req, nil
}

func parseTopTracksForArtistResponse(body *[]byte) (*Tracks, error) {
	ret := Tracks{}
	err := json.Unmarshal(*body, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

// ----
// Members
// ----

func (t *Tracks) GetGenres(ctx context.Context) (*sortablemap.Map, error) {
	as := map[string]int32{}
	aids := []string{}
	ids := t.IDs()

	// collect the artist ids from the tracks
	for _, i := range *t {
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

	// go through artists to collect genres
	artists, err := getArtists(ctx, aids)
	if err != nil {
		return nil, err
	}

	ret := map[string]int{}
	for _, i := range *artists {
		for _, ii := range i.Genres {
			if _, ok := ret[ii]; ok {
				ret[ii]++
			} else {
				ret[ii] = 1
			}
		}
	}
	retmap := sortablemap.GetSortableMap(ret)
	return &retmap, nil
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
