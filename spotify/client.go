package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func getTopTracks(ctx context.Context, limit int32) (Tracks, error) {
	token := ctx.Value(ContextAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
	}

	tr := ctx.Value(ContextTimeRange)
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

	body, err := makeRequest(ctx, req)
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

	body, err := makeRequest(ctx, req)
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

	body, err := makeRequest(ctx, req)
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

func getTracks(ctx context.Context, ids []string) (*Tracks, error) {
	token := ctx.Value(ContextAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
	}

	url := fmt.Sprint("https://api.spotify.com/v1/tracks?ids=", strings.Join(ids, ","))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	body, err := makeRequest(ctx, req)
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

func makeRequest(ctx context.Context, req *http.Request) (*[]byte, error) {
	log.Println("making external request: ", req.URL)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Println("unhappy response ", resp.StatusCode, "\n\t", string(b))
		return nil, errors.New("non-200 response")
	}

	return &b, nil
}
