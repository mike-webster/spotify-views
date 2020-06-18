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
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("non-200 response")
	}

	type tempResp struct {
		Items Tracks `json:"items"`
	}

	var ret tempResp
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &ret)
	if err != nil {
		return nil, err
	}
	return ret.Items, nil
}
