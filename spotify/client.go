package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/mike-webster/spotify-views/logging"
	"github.com/sirupsen/logrus"
)

var (
	KeySeedArtists = "seed-artists"
	KeySeedTracks  = "seed-tracks"
	KeySeedGenres  = "seed-genres"
)

func addToCache(ctx context.Context, key string, body *[]byte) error {
	irdb := ctx.Value("Redis")
	if irdb == nil {
		return errors.New("no redis")
	}

	rdb, ok := irdb.(*redis.Client)
	if !ok {
		logging.GetLogger(nil).WithField("event", "redis-cast-error").Error(fmt.Sprint(reflect.TypeOf(irdb)))
	}

	err := rdb.Set(ctx, key, string(*body), 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func calculateRedisKey(ctx context.Context, req *http.Request) (string, error) {
	uid := ctx.Value(string(ContextUserID))
	if len(fmt.Sprint(uid)) < 1 {
		return "", errors.New("no user id in context")
	}
	return fmt.Sprint(uid, "-", req.URL), nil
}

func checkCache(ctx context.Context, key string) (*[]byte, error) {
	irdb := ctx.Value("Redis")
	if irdb == nil {
		return nil, errors.New("no redis")
	}

	rdb, ok := irdb.(*redis.Client)
	if !ok {
		logging.GetLogger(nil).WithField("event", "redis-cast-error").Error(fmt.Sprint(reflect.TypeOf(irdb)))
	}

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			logging.GetLogger(nil).WithField("event", "cache-miss").Debug()
			return nil, nil
		}
		return nil, err
	}

	if len(val) > 0 {
		logging.GetLogger(nil).WithField("event", "cache-hit").Debug()
	}

	bytes := []byte(val)
	return &bytes, nil
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

func getAudioFeatures(ctx context.Context, ids []string) (*AudioFeatures, error) {
	logger := logging.GetLogger(&ctx)
	token := ctx.Value(ContextAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
	}

	if len(ids) > 100 {
		logger.WithField("count", len(ids)).Warn("too many tracks passed, reducing to the first 100")
		ids = ids[:100]
	}

	url := fmt.Sprint("https://api.spotify.com/v1/audio-features?ids=", strings.Join(ids, ","))
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
		Items AudioFeatures `json:"audio_features"`
	}

	ret := tempResp{}
	err = json.Unmarshal(*body, &ret)
	if err != nil {
		logger.WithField("body", string(*body)).Error(err.Error())
		return nil, err
	}

	return &ret.Items, nil
}

func getChunkOfUserLibraryTracks(ctx context.Context, url string) (Tracks, string, int, error) {
	token := ctx.Value(ContextAccessToken)
	if token == nil {
		return nil, "", 0, errors.New("no access token provided")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", 0, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	body, err := makeRequest(ctx, req, true)
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
	token := ctx.Value(ContextAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
	}

	url := "https://api.spotify.com/v1/recommendations/available-genre-seeds"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	body, err := makeRequest(ctx, req, false)
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

func getRecommendations(ctx context.Context, seeds map[string][]string) (*Recommendation, error) {
	// holy shit, this is actually _really_ configurable.  Come back to this
	// and explore the possibilities a little more after v1 is out.

	token := ctx.Value(ContextAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
	}

	qs := "?"
	if artists, ok := seeds[KeySeedArtists]; ok {
		qs = fmt.Sprint(qs, "seed_artists=", strings.Join(artists, ","))
	}

	if tracks, ok := seeds[KeySeedTracks]; ok {
		if len(qs) > 1 {
			qs += "&"
		}

		qs = fmt.Sprint(qs, "seed_tracks=", strings.Join(tracks, ","))
	}

	if genres, ok := seeds[KeySeedGenres]; ok {
		if len(qs) > 1 {
			qs += "&"
		}

		qs = fmt.Sprint(qs, "seed_genres=", strings.Join(genres, ","))
	}

	url := "https://api.spotify.com/v1/recommendations"
	if len(qs) > 1 {
		url = fmt.Sprint(url, qs)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	body, err := makeRequest(ctx, req, false)
	if err != nil {
		return nil, err
	}

	type tApiResponse struct {
		Tracks []struct {
			Artists []struct {
				Link string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"artists"`
			Duration int64  `json:"duration_ms"`
			Link     string `json:"href"`
			ID       string `json:"id"`
			Name     string `json:"name"`
		} `json:"tracks"`
		Seeds []struct {
			InitialPoolSize    int64  `json:"initialPoolSize"`
			AfterFilteringSize int64  `json:"afterFilteringSize"`
			Link               string `json:"href"`
			ID                 string `json:"id"`
			Type               string `json:"type"`
		} `json:"seeds"`
	}

	var rsp tApiResponse
	err = json.Unmarshal(*body, &rsp)
	if err != nil {
		return nil, err
	}

	ret := Recommendation{}
	for _, i := range rsp.Tracks {
		artists := []string{}
		for _, j := range i.Artists {
			artists = append(artists, j.Name)
		}
		ret.Tracks = append(ret.Tracks, struct {
			ID      string   "json:\"id\""
			Name    string   "json:\"name\""
			Link    string   "json:\"href\""
			Artists []string "json:\"artists\""
		}{
			ID:      i.ID,
			Name:    i.Name,
			Link:    i.Link,
			Artists: artists,
		})
	}

	for _, i := range rsp.Seeds {
		ret.Seeds = append(ret.Seeds, struct {
			ID   string "json:\"id\""
			Link string "json:\"href\""
			Type string "json:\"type\""
		}{
			ID:   i.ID,
			Link: i.Link,
			Type: i.Type,
		})
	}

	return &ret, nil
}

func getRelatedArtists(ctx context.Context, id string) (*[]RelatedArtist, error) {
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
			Popularity int64    `json:"popularity"`
			Type       string   `json:"type"`
		} `json:"artists"`
	}

	rsp := tRes{}
	err = json.Unmarshal(*body, &rsp)
	if err != nil {
		return nil, err
	}

	ret := []RelatedArtist{}
	for _, i := range rsp.Artists {
		ret = append(ret, RelatedArtist{
			Followers:  i.Followers.Total,
			Genres:     i.Genres,
			Link:       i.Link,
			ID:         i.ID,
			Name:       i.Name,
			Popularity: i.Popularity,
			Type:       i.Type,
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
	token := ctx.Value(ContextAccessToken)
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

func getUserInfo(ctx context.Context) (map[string]string, error) {
	token := ctx.Value(ContextAccessToken)
	if token == nil {
		return map[string]string{}, errors.New("no access token provided")
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return map[string]string{}, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	body, err := makeRequest(ctx, req, true)
	if err != nil {
		return map[string]string{}, err
	}

	type userResponse struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	ret := userResponse{}
	err = json.Unmarshal(*body, &ret)
	if err != nil {
		return map[string]string{}, err
	}
	return map[string]string{
		"id":    ret.ID,
		"email": ret.Email,
	}, nil
}

func makeRequest(ctx context.Context, req *http.Request, useCache bool) (*[]byte, error) {
	s := time.Now()
	logger := logging.GetLogger(&ctx)
	cacheKey, err := calculateRedisKey(ctx, req)
	if err != nil {
		logger.WithField("event", "redis-key-error").Error(err.Error())
	}

	if useCache {
		val, err := checkCache(ctx, cacheKey)
		if err != nil {
			logger.WithField("event", "redis-error").Error(err.Error())
		}
		if val != nil {
			return val, nil
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	dur := time.Since(s)

	logger.WithFields(logrus.Fields{
		"status":   resp.StatusCode,
		"url":      req.URL,
		"event":    "external_request",
		"duration": dur.String(),
	}).Info()

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 401 {
			logger.WithField("event", EventNeedsRefreshToken).Info()
			return nil, ErrTokenExpired("")
		}

		logger.WithFields(logrus.Fields{
			"event":  EventNon200Response,
			"status": resp.StatusCode,
			"body":   string(b),
		}).Error()
		return nil, errors.New(fmt.Sprint("non-200 response; ", resp.StatusCode))
	}

	if useCache {
		err = addToCache(ctx, cacheKey, &b)
		if err != nil {
			logger.WithField("event", "redis-add-error").Error(err.Error())
		}
	}

	return &b, nil
}
