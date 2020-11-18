package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Artist represents a spotify Artist
type Artist struct {
	Genres     []string `json:"genres"`
	Name       string   `json:"name"`
	Popularity int32    `json:"popularity"`
	URI        string   `json:"uri"`
	ID         string   `json:"ID"`
}

// Artists is a collection of spotify Artist
type Artists []Artist

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

func getArtist(ctx context.Context, id string) (*Artist, error) {
	token := ctx.Value(ContextAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
	}

	url := fmt.Sprint("https://api.spotify.com/v1/artists/", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	body, err := makeRequest(ctx, req, true)
	if err != nil {
		return nil, err
	}

	var ret Artist

	err = json.Unmarshal(*body, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func getArtists(ctx context.Context, ids []string) (*Artists, error) {
	token := ctx.Value(ContextAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
	}

	url := fmt.Sprint("https://api.spotify.com/v1/artists?ids=", strings.Join(ids, ","))

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
		Items Artists `json:"artists"`
	}

	var ret tempResp

	err = json.Unmarshal(*body, &ret)
	if err != nil {
		return nil, err
	}

	return &ret.Items, nil
}

func getRelatedArtists(ctx context.Context, id string) (*[]Artist, error) {
	token := ctx.Value(ContextAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
	}

	url := "https://api.spotify.com/v1/artists/%v/related-artists"

	req, err := http.NewRequest("GET", fmt.Sprintf(url, id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	body, err := makeRequest(ctx, req, false)
	if err != nil {
		return nil, err
	}

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
	err = json.Unmarshal(*body, &rsp)
	if err != nil {
		return nil, err
	}

	ret := []Artist{}
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

func getTopArtists(ctx context.Context) (*Artists, error) {
	token := ctx.Value(ContextAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
	}

	tr := ctx.Value(ContextTimeRange)
	strRange := ""
	if tr != nil {
		strRange = tr.(string)
	}

	// TODO: make this limit a param
	url := "https://api.spotify.com/v1/me/top/artists?limit=25"
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
		Items Artists `json:"items"`
	}

	var ret tempResp
	err = json.Unmarshal(*body, &ret)
	if err != nil {
		return nil, err
	}

	return &ret.Items, nil
}
