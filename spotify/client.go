package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/mike-webster/spotify-views/keys"
	"github.com/mike-webster/spotify-views/logging"
	"github.com/sirupsen/logrus"
)

var (
	KeySeedArtists = "seed-artists"
	KeySeedTracks  = "seed-tracks"
	KeySeedGenres  = "seed-genres"
)

// TODO: cleanup
type item struct {
	DateSaved time.Time `json:"added_at"`
	Track     Track     `json:"track"`
}

type items []item

func (i items) Tracks() Tracks {
	ret := Tracks{}
	for _, j := range i {
		ret = append(ret, j.Track)
	}
	return ret
}

// ENDTODO

func getChunkOfUserLibraryTracks(ctx context.Context, url string) (Tracks, string, int, error) {
	token := keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)
	if token == nil {
		return nil, "", 0, errors.New("no access token provided")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", 0, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	body, err := makeRequest(ctx, req)
	if err != nil {
		return nil, "", 0, err
	}

	type tempResp struct {
		Link   string `json:"href"`
		Items  items  `json:"items"`
		Limit  int    `json:"limit"`
		Next   string `json:"next"`
		Offset int    `json:"offset"`
		Total  int    `json:"total"`
	}

	var ret tempResp
	err = json.Unmarshal(*body, &ret)
	if err != nil {
		return nil, "", 0, err
	}

	return ret.Items.Tracks(), ret.Next, ret.Total, nil
}

func getGenres(ctx context.Context) ([]string, error) {
	token := keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
	}

	url := "https://api.spotify.com/v1/recommendations/available-genre-seeds"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	body, err := makeRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	type tRes struct {
		Genres []string `json:"genres"`
	}

	rsp := tRes{}
	err = json.Unmarshal(*body, &rsp)
	if err != nil {
		return nil, err
	}

	return rsp.Genres, nil
}

func makeRequest(ctx context.Context, req *http.Request) (*[]byte, error) {
	s := time.Now()
	logger := logging.GetLogger(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	dur := time.Since(s)

	logger.WithFields(logrus.Fields{
		"status":   resp.StatusCode,
		"url":      req.URL,
		"event":    "external_request",
		"duration": dur.String(),
	}).Info("making external request")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	logger.WithField("external_request_response", resp.StatusCode).Debug()

	if resp.StatusCode != 200 {
		if resp.StatusCode == 401 {
			logger.WithField("event", EventNeedsRefreshToken).Info()
			return nil, ErrTokenExpired("")
		} else if resp.StatusCode == 429 {
			wait := 5
			hdr := resp.Header["Retry-Afer"]
			if len(hdr) > 0 {
				wait, err = strconv.Atoi(hdr[0])
				if err != nil {
					panic(err)
				}
			}
			logger.WithField("event", "spotify_rate_limited").Error(fmt.Sprint("waiting ", wait, " seconds"))

			time.Sleep(time.Duration(wait) * time.Second)
			makeRequest(ctx, req)
		}

		logger.WithFields(logrus.Fields{
			"event":  EventNon200Response,
			"status": resp.StatusCode,
			"body":   string(b),
		}).Error()
		return nil, errors.New(fmt.Sprint("non-200 response; ", resp.StatusCode))
	}

	return &b, nil
}
