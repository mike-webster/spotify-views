package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/mike-webster/spotify-views/keys"
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
	uid := keys.GetContextValue(ctx, keys.ContextSpotifyUserID)
	if uid == nil {
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

func getUserInfo(ctx context.Context) (map[string]string, error) {
	token := keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)
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
	logger := logging.GetLogger(ctx)

	cacheKey := ""

	if useCache && false {
		cacheKey, err := calculateRedisKey(ctx, req)
		if err != nil {
			logger.WithField("event", "redis-key-error").Error(err.Error())
		}
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
			tok, err := refreshToken(ctx)
			if err != nil {
				logging.GetLogger(ctx).WithError(err).Error("auto refreshing token failed")
				return nil, ErrTokenExpired("")
			}
			return nil, ErrTokenExpired(tok)
		}

		logger.WithFields(logrus.Fields{
			"event":  EventNon200Response,
			"status": resp.StatusCode,
			"body":   string(b),
		}).Error()
		return nil, errors.New(fmt.Sprint("non-200 response; ", resp.StatusCode))
	}

	if useCache && false {
		err = addToCache(ctx, cacheKey, &b)
		if err != nil {
			logger.WithField("event", "redis-add-error").Error(err.Error())
		}
	}

	return &b, nil
}
