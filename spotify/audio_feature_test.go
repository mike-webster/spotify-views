package spotify

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/mike-webster/spotify-views/keys"
	"github.com/stretchr/testify/assert"
)

func TestGetAudioFeatures(t *testing.T) {
	ids := []string{"5DEF7cgddcZ1dUBk96J1Hx", "5eyCGJPbszc4xak3db5UKL"}
	t.Run("TestParseRequestForGetArtist", func(t *testing.T) {
		ctx := context.Background()
		t.Run("no token", func(t *testing.T) {
			_, err := parseRequestForAudioFeatures(ctx, ids)
			assert.Equal(t, reflect.TypeOf(ErrNoToken("")), reflect.TypeOf(err))
		})

		token := "tok"
		ctx = context.WithValue(ctx, keys.ContextSpotifyAccessToken, token)

		t.Run("token gets stored in header", func(t *testing.T) {
			req, err := parseRequestForAudioFeatures(ctx, ids)
			assert.Nil(t, err)
			assert.Equal(t, req.Header.Get("Authorization"), fmt.Sprint("Bearer ", token))
		})
	})

	t.Run("TestParseResponseForGetArtist", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			bytes := []byte(getArtistPayload)

			as, err := parseResponseForAudioFeatures(&bytes)
			assert.Nil(t, err)
			assert.NotNil(t, as)
		})

		t.Run("bad body", func(t *testing.T) {
			bytes := []byte("fdakslfjda;klfjad;kjadl;")
			_, err := parseResponseForAudioFeatures(&bytes)
			assert.NotNil(t, err)
		})
	})

	t.Run("MainMethod", func(t *testing.T) {
		// TODO
	})
}

func TestChunkRangeAudioFeatures(t *testing.T) {
	t.Run("WhenFewerIDsThanTheLimit", func(t *testing.T) {
		ids := []string{"1", "2", "3"}
		beginning, ending := chunkRangeAudioFeatures(0, &ids)
		assert.Equal(t, 0, beginning)
		assert.Equal(t, len(ids), ending)
	})

	t.Run("WhenStartIsLargerThanLimit", func(t *testing.T) {
		// this shouldn't happen so we're just returning the first page
		ids := []string{}
		for i := 0; i < 105; i++ {
			ids = append(ids, fmt.Sprint(i))
		}

		beginning, ending := chunkRangeAudioFeatures(len(ids)+1, &ids)
		assert.Equal(t, 0, beginning)
		assert.Equal(t, audioFeaturesPageLimit, ending)
	})

	t.Run("WhenThePageRequestIsWithinBounds", func(t *testing.T) {
		ids := []string{}
		for i := 0; i < 505; i++ {
			ids = append(ids, fmt.Sprint(i))
		}

		beginning, ending := chunkRangeAudioFeatures(audioFeaturesPageLimit+1, &ids)
		assert.Equal(t, audioFeaturesPageLimit+1, beginning)
		assert.Equal(t, beginning+audioFeaturesPageLimit, ending)
	})
}

var (
	getAudioFeaturesPayload = `{
		"audio_features": [
		  {
			"danceability": 0.432,
			"energy": 0.942,
			"key": 2,
			"loudness": -3.037,
			"mode": 1,
			"speechiness": 0.0709,
			"acousticness": 0.000505,
			"instrumentalness": 0,
			"liveness": 0.333,
			"valence": 0.663,
			"tempo": 156.902,
			"type": "audio_features",
			"id": "5DEF7cgddcZ1dUBk96J1Hx",
			"uri": "spotify:track:5DEF7cgddcZ1dUBk96J1Hx",
			"track_href": "https://api.spotify.com/v1/tracks/5DEF7cgddcZ1dUBk96J1Hx",
			"analysis_url": "https://api.spotify.com/v1/audio-analysis/5DEF7cgddcZ1dUBk96J1Hx",
			"duration_ms": 204708,
			"time_signature": 4
		  },
		  {
			"danceability": 0.553,
			"energy": 0.926,
			"key": 0,
			"loudness": -4.266,
			"mode": 1,
			"speechiness": 0.0428,
			"acousticness": 0.00159,
			"instrumentalness": 0,
			"liveness": 0.196,
			"valence": 0.603,
			"tempo": 162.981,
			"type": "audio_features",
			"id": "5eyCGJPbszc4xak3db5UKL",
			"uri": "spotify:track:5eyCGJPbszc4xak3db5UKL",
			"track_href": "https://api.spotify.com/v1/tracks/5eyCGJPbszc4xak3db5UKL",
			"analysis_url": "https://api.spotify.com/v1/audio-analysis/5eyCGJPbszc4xak3db5UKL",
			"duration_ms": 168661,
			"time_signature": 3
		  }
		]
	  }`
)
