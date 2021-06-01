package spotify

import (
	"context"
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

func makeRequest(ctx context.Context, req *http.Request) (*[]byte, error) {
	s := time.Now()
	logger := logging.GetLogger(ctx)
	id := keys.GetContextValue(ctx, keys.ContextSpotifyUserID)
	url := req.URL.String()
	cacheKey := fmt.Sprint(id, "-", url)

	deps := GetDependencies(ctx)
	if deps == nil {
		return nil, errors.New("couldnt find deps")
	}

	// when this key is provided with a value of true we want to skip
	skip := keys.GetContextValue(ctx, keys.ContextSkipCache)
	if skip != true && deps.Cache != nil {
		if deps.Cache != nil {
			val, err := deps.Cache.Get(ctx, cacheKey)
			if err != nil {
				logger.WithError(err).Error("error checking cache")
			} else {
				b := []byte(val)
				logger.WithField("cached_value", val).Debug("using cache")
				return &b, nil
			}
		}
	}

	resp, err := deps.Client.Do(req)
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

	if resp.StatusCode != 200 {
		if resp.StatusCode == 401 {
			logger.WithField("event", EventNeedsRefreshToken).Info()
			return nil, ErrTokenExpired("")
		} else if resp.StatusCode == 429 {
			wait := 10
			hdr := resp.Header["Retry-After"]
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
		return nil, ErrBadRequest(fmt.Sprint("response code: ", resp.StatusCode))
	}

	if skip != true && deps.Cache != nil {
		err := deps.Cache.Set(ctx, cacheKey, string(b))
		if err != nil {
			logger.WithError(err).Error("error setting cache record")
		}
	}

	return &b, nil
}
