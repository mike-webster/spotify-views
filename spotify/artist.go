package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mike-webster/spotify-views/keys"
	"github.com/mike-webster/spotify-views/sortablemap"
)

// Artist represents a spotify Artist
type Artist struct {
	Links      map[string]string `json:"external_urls"`
	Genres     []string          `json:"genres"`
	Name       string            `json:"name"`
	Popularity int32             `json:"popularity"`
	URI        string            `json:"uri"`
	ID         string            `json:"ID"`
	Images     []Image           `json:"images"`
}

// Artists is a collection of spotify Artist
type Artists []Artist

// ----
// API
// ----

func GetArtist(ctx context.Context, id string) (*Artist, error) {
	req, err := parseRequestForGetArtist(ctx, id)
	if err != nil {
		return nil, err
	}

	body, err := makeRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return parseResponseForGetArtist(body)
}

func GetArtists(ctx context.Context, ids []string) (*Artists, error) {
	req, err := parseRequestForGetArtists(ctx, ids)

	body, err := makeRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return parseResponseForGetArtists(body)
}

func GetTopArtists(ctx context.Context) (*Artists, error) {
	req, err := parseRequestForGetTopArtists(ctx)
	if err != nil {
		return nil, err
	}

	body, err := makeRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return parseResponseForGetTopArtists(body)
}

// ----
// Members
// ----

// EmbeddedPlayer  will return the html to use for rendering the embedded spotify
// player iframe
func (a *Artist) EmbeddedPlayer() string {
	return fmt.Sprintf(`<h4 width="300" style="text-align:center">%s</h4><iframe src="https://open.spotify.com/embed/artist/%s" width="300" height="380" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>`, a.Name, a.ID)
}

// IDs returns the ID for each of the artists in the collection of Artists
func (a *Artists) IDs() []string {
	ret := []string{}
	for _, i := range *a {
		ret = append(ret, i.ID)
	}
	return ret
}

func (a *Artists) GetGenres(ctx context.Context) *sortablemap.Map {
	ret := map[string]int{}

	for _, i := range *a {
		for _, ii := range i.Genres {
			if _, ok := ret[ii]; ok {
				ret[ii]++
			} else {
				ret[ii] = 1
			}
		}
	}

	sm := sortablemap.GetSortableMap(ret)
	return &sm
}

func (a *Artist) FindImage() *Image {
	if len(a.Images) < 1 {
		return nil
	}

	// why do I prefer the first one only when there's one?
	if len(a.Images) == 1 {
		return &a.Images[0]
	}

	return &a.Images[1]
}

func (a *Artist) GetRelatedArtists(ctx context.Context) (*Artists, error) {
	req, err := parseRequestForRelatedArtists(ctx, a.ID)
	if err != nil {
		return nil, err
	}

	body, err := makeRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return parseResponseForRelatedArtists(body)
}

// ----
// Helpers
// ----

func parseRequestForGetTopArtists(ctx context.Context) (*http.Request, error) {
	token := keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)
	if token == nil {
		return nil, ErrNoToken("no tok")
	}

	tr := keys.GetContextValue(ctx, keys.ContextSpotifyTimeRange)
	tframe := TFShort
	if val, ok := tr.(string); ok {
		tframe = GetTimeFrame(val)
	}

	// TODO: make this limit a param
	url := "https://api.spotify.com/v1/me/top/artists?limit=25"
	url += fmt.Sprint("&time_range=", tframe.Value())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))
	return req, nil
}

func parseResponseForGetTopArtists(body *[]byte) (*Artists, error) {
	type tempResp struct {
		Items Artists `json:"items"`
	}

	var ret tempResp
	err := json.Unmarshal(*body, &ret)
	if err != nil {
		return nil, err
	}

	return &ret.Items, nil
}

func parseRequestForGetArtists(ctx context.Context, ids []string) (*http.Request, error) {
	token := keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)
	if token == nil {
		return nil, ErrNoToken("no access token provided")
	}

	url := fmt.Sprint("https://api.spotify.com/v1/artists?ids=", strings.Join(ids, ","))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))
	return req, nil
}

func parseResponseForGetArtists(body *[]byte) (*Artists, error) {
	type tempResp struct {
		Items Artists `json:"artists"`
	}

	var ret tempResp

	err := json.Unmarshal(*body, &ret)
	if err != nil {
		return nil, err
	}

	return &ret.Items, nil
}

func parseResponseForGetArtist(body *[]byte) (*Artist, error) {
	var ret Artist

	err := json.Unmarshal(*body, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func parseRequestForGetArtist(ctx context.Context, id string) (*http.Request, error) {
	token := keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)
	if token == nil {
		return nil, ErrNoToken("no access token provided")
	}

	url := fmt.Sprint("https://api.spotify.com/v1/artists/", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	return req, nil
}

func parseRequestForRelatedArtists(ctx context.Context, id string) (*http.Request, error) {
	token := keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)
	if token == nil {
		return nil, ErrNoToken("no access token provided")
	}

	url := "https://api.spotify.com/v1/artists/%v/related-artists"

	req, err := http.NewRequest("GET", fmt.Sprintf(url, id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))
	return req, nil
}

func parseResponseForRelatedArtists(body *[]byte) (*Artists, error) {
	type tRes struct {
		Artists []struct {
			Followers struct {
				Total int64 `json:"total"`
			} `json:"followers"`
			Genres     []string `json:"genres"`
			Link       string   `json:"href"`
			ID         string   `json:"id"`
			Name       string   `json:"name"`
			Popularity int32    `json:"popularity"`
			Type       string   `json:"type"`
		} `json:"artists"`
	}

	rsp := tRes{}
	err := json.Unmarshal(*body, &rsp)
	if err != nil {
		return nil, err
	}

	ret := Artists{}
	for _, i := range rsp.Artists {
		ret = append(ret, Artist{
			Genres:     i.Genres,
			ID:         i.ID,
			Name:       i.Name,
			Popularity: i.Popularity,
		})
	}

	return &ret, nil
}
