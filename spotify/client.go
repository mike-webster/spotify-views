package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getTopTracks(ctx context.Context) (Tracks, error) {
	token := ctx.Value("access_token")
	if token == nil {
		return nil, errors.New("no access token provided")
	}
	tr := ctx.Value("time_range")
	strRange := ""
	if tr != nil {
		strRange = tr.(string)
	}
	url := "https://api.spotify.com/v1/me/top/tracks?limit=25"
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

// decided not to make this one method that accepted the resource
// but went against it to simplify use and avoid passing back
// an interface.

func getTopArtists(ctx context.Context) (*Artists, error) {
	token := ctx.Value("access_token")
	if token == nil {
		return nil, errors.New("no access token provided")
	}
	tr := ctx.Value("time_range")
	strRange := ""
	if tr != nil {
		strRange = tr.(string)
	}
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
	}

	type tempResp struct {
		Items Tracks `json:"tracks"`
	}

	var ret tempResp

	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &ret.Items, nil
}

func makeRequest(ctx context.Context, req *http.Request) (*[]byte, error) {
	client := &http.Client{}
	log.Println("requesting: ", req.URL)
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
