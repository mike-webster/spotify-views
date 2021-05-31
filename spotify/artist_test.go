package spotify

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/mike-webster/spotify-views/keys"
	"github.com/stretchr/testify/assert"
)

type TestHttpClient struct {
	Response   *http.Response
	ShouldErr  bool
	ErrMessage string
}

func (c *TestHttpClient) Do(eq *http.Request) (*http.Response, error) {
	if c.ShouldErr {
		return nil, errors.New(c.ErrMessage)
	}

	return c.Response, nil
}

func TestEmbeddedPlayer(t *testing.T) {
	t.Run("CheckURL", func(t *testing.T) {
		a := Artist{Name: "testname", ID: "testid"}
		exp := fmt.Sprintf(`<h4 width="300" style="text-align:center">%s</h4><iframe src="https://open.spotify.com/embed/artist/%s" width="300" height="380" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>`, a.Name, a.ID)
		assert.Equal(t, exp, a.EmbeddedPlayer())
	})
}

func TestArtistIDs(t *testing.T) {
	objs := Artists{
		Artist{ID: "1234"},
		Artist{ID: "2345"},
		Artist{ID: "3456"},
		Artist{ID: "4567"},
		Artist{ID: "5678"},
	}

	exp := []string{}
	for _, i := range objs {
		exp = append(exp, i.ID)
	}

	assert.Equal(t, exp, objs.IDs())
}

func TestGetArtistGenres(t *testing.T) {
	objs := Artists{
		Artist{Genres: []string{"emo", "rock", "rap"}},
		Artist{Genres: []string{"hip hop", "rap", "mumble"}},
		Artist{Genres: []string{"country", "shit", "trash", "noise"}},
		Artist{Genres: []string{"classical"}},
	}

	exp := map[string]int{}
	exp["emo"] = 1
	exp["rock"] = 1
	exp["rap"] = 2
	exp["hip hop"] = 1
	exp["mumble"] = 1
	exp["country"] = 1
	exp["shit"] = 1
	exp["trash"] = 1
	exp["noise"] = 1
	exp["classical"] = 1

	genres := objs.GetGenres(context.Background())

	for k, v := range exp {
		t.Run(fmt.Sprint("CheckingExpectedValues_", k), func(t *testing.T) {
			res, err := genres.Value(k)
			assert.Nil(t, err)
			assert.Equal(t, v, res)
		})
	}
}

func TestArtistFindImage(t *testing.T) {
	t.Run("NilWhenEmpty", func(t *testing.T) {
		a := Artist{}
		assert.Nil(t, a.FindImage())
	})

	t.Run("FirstWhenOnlyOne", func(t *testing.T) {
		a := Artist{Images: []Image{{URL: "test"}}}
		assert.Equal(t, &a.Images[0], a.FindImage())
	})

	t.Run("SecondWhenMoreThanOne", func(t *testing.T) {
		a := Artist{Images: []Image{{URL: "test"}, {URL: "test2"}}}
		assert.Equal(t, &a.Images[1], a.FindImage())
	})
}

func getTestDependencies(ctx context.Context, code int, body string) context.Context {
	deps := Dependencies{
		Client: &TestHttpClient{
			Response: &http.Response{
				StatusCode: code,
				Body:       ioutil.NopCloser(strings.NewReader(body)),
			},
		},
	}

	return context.WithValue(ctx, keys.ContextDependencies, &deps)
}

func TestGetArtist(t *testing.T) {
	id := "6FBDaR13swtiWwGhX1WQsP"
	t.Run("TestParseRequestForGetArtist", func(t *testing.T) {
		ctx := context.Background()
		t.Run("no token", func(t *testing.T) {
			_, err := parseRequestForGetArtist(ctx, id)
			assert.Equal(t, reflect.TypeOf(ErrNoToken("")), reflect.TypeOf(err))
		})

		token := "tok"
		ctx = context.WithValue(ctx, keys.ContextSpotifyAccessToken, token)

		t.Run("token gets stored in header", func(t *testing.T) {
			req, err := parseRequestForGetArtist(ctx, id)
			assert.Nil(t, err)
			assert.Equal(t, req.Header.Get("Authorization"), fmt.Sprint("Bearer ", token))
		})
	})

	t.Run("TestParseResponseForGetArtist", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			bytes := []byte(getArtistPayload)

			as, err := parseResponseForGetArtist(&bytes)
			assert.Nil(t, err)
			assert.NotNil(t, as)
		})

		t.Run("bad body", func(t *testing.T) {
			bytes := []byte("fdakslfjda;klfjad;kjadl;")
			_, err := parseResponseForGetArtist(&bytes)
			assert.NotNil(t, err)
		})
	})

	t.Run("MainMethod", func(t *testing.T) {
		t.Run("HappyPath", func(t *testing.T) {
			ctx := getTestDependencies(context.Background(), 200, "{}")
			ctx = context.WithValue(ctx, keys.ContextSpotifyAccessToken, "test")

			_, err := GetArtist(ctx, "test")
			assert.Equal(t, nil, err)
		})

		t.Run("BadRequest", func(t *testing.T) {
			ctx := getTestDependencies(context.Background(), 400, `{"err":"bad_request"}`)
			ctx = context.WithValue(ctx, keys.ContextSpotifyAccessToken, "test")

			_, err := GetArtist(ctx, "test")
			assert.NotEqual(t, nil, err)
		})
	})
}

func TestGetArtists(t *testing.T) {
	ids := []string{"6FBDaR13swtiWwGhX1WQsP", "1lKZzN2d4IqiEYxyECIEHI"}
	t.Run("TestParseRequestForGetArtists", func(t *testing.T) {
		ctx := context.Background()
		t.Run("no token", func(t *testing.T) {
			_, err := parseRequestForGetArtists(ctx, ids)
			assert.Equal(t, reflect.TypeOf(ErrNoToken("")), reflect.TypeOf(err))
		})

		token := "tok"
		ctx = context.WithValue(ctx, keys.ContextSpotifyAccessToken, token)

		t.Run("token gets stored in header", func(t *testing.T) {
			req, err := parseRequestForGetArtists(ctx, ids)
			assert.Nil(t, err)
			assert.Equal(t, req.Header.Get("Authorization"), fmt.Sprint("Bearer ", token))
		})
	})

	t.Run("TestParseResponseForGetArtists", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			bytes := []byte(getArtistsPayload)

			as, err := parseResponseForGetArtists(&bytes)
			assert.Nil(t, err)
			assert.True(t, len(*as) == 2, len(*as))
		})

		t.Run("bad body", func(t *testing.T) {
			bytes := []byte("fdakslfjda;klfjad;kjadl;")
			_, err := parseResponseForGetArtists(&bytes)
			assert.NotNil(t, err)
		})
	})

	t.Run("MainMethod", func(t *testing.T) {
		t.Run("HappyPath", func(t *testing.T) {
			ctx := getTestDependencies(context.Background(), 200, "{}")
			ctx = context.WithValue(ctx, keys.ContextSpotifyAccessToken, "test")

			_, err := GetArtists(ctx, ids)
			assert.Equal(t, nil, err)
		})

		t.Run("BadRequest", func(t *testing.T) {
			ctx := getTestDependencies(context.Background(), 400, `{"err":"bad_request"}`)
			ctx = context.WithValue(ctx, keys.ContextSpotifyAccessToken, "test")

			_, err := GetArtists(ctx, ids)
			assert.NotEqual(t, nil, err)
		})
	})
}

func TestGetTopArtists(t *testing.T) {
	t.Run("TestParseRequestForGetTopArtists", func(t *testing.T) {
		ctx := context.Background()
		t.Run("no token", func(t *testing.T) {
			_, err := parseRequestForGetTopArtists(ctx)
			assert.Equal(t, reflect.TypeOf(ErrNoToken("")), reflect.TypeOf(err))
		})

		token := "tok"
		ctx = context.WithValue(ctx, keys.ContextSpotifyAccessToken, token)
		t.Run("no timerange provided defaults to short", func(t *testing.T) {
			req, err := parseRequestForGetTopArtists(ctx)
			assert.Nil(t, err)
			assert.True(t, strings.Contains(req.URL.RawQuery, "short_term"))
		})

		ctx = context.WithValue(ctx, keys.ContextSpotifyTimeRange, "medium_term")
		t.Run("timerange provided works", func(t *testing.T) {
			req, err := parseRequestForGetTopArtists(ctx)
			assert.Nil(t, err)
			assert.True(t, strings.Contains(req.URL.RawQuery, "medium_term"))
		})

		t.Run("token gets stored in header", func(t *testing.T) {
			req, err := parseRequestForGetTopArtists(ctx)
			assert.Nil(t, err)
			assert.Equal(t, req.Header.Get("Authorization"), fmt.Sprint("Bearer ", token))
		})
	})

	t.Run("TestParseResponseForGetTopArtists", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			bytes := []byte(getTopArtistsPayload)

			as, err := parseResponseForGetTopArtists(&bytes)
			assert.Nil(t, err)
			assert.True(t, len(*as) == 20, len(*as))
		})

		t.Run("bad body", func(t *testing.T) {
			bytes := []byte("fdakslfjda;klfjad;kjadl;")
			_, err := parseResponseForGetTopArtists(&bytes)
			assert.NotNil(t, err)
		})
	})

	t.Run("MainMethod", func(t *testing.T) {
		t.Run("HappyPath", func(t *testing.T) {
			ctx := getTestDependencies(context.Background(), 200, "{}")
			ctx = context.WithValue(ctx, keys.ContextSpotifyAccessToken, "test")

			_, err := GetTopArtists(ctx)
			assert.Equal(t, nil, err)
		})

		t.Run("BadRequest", func(t *testing.T) {
			ctx := getTestDependencies(context.Background(), 400, `{"err":"bad_request"}`)
			ctx = context.WithValue(ctx, keys.ContextSpotifyAccessToken, "test")

			_, err := GetTopArtists(ctx)
			assert.NotEqual(t, nil, err)
		})
	})
}

func TestGetRelatedArtists(t *testing.T) {
	id := "5fEKZRCUa0JApec5Xy095q"
	t.Run("TestParseRequestForGetRelatedArtists", func(t *testing.T) {
		ctx := context.Background()
		t.Run("no token", func(t *testing.T) {
			_, err := parseRequestForRelatedArtists(ctx, id)
			assert.Equal(t, reflect.TypeOf(ErrNoToken("")), reflect.TypeOf(err))
		})

		token := "tok"
		ctx = context.WithValue(ctx, keys.ContextSpotifyAccessToken, token)

		t.Run("token gets stored in header", func(t *testing.T) {
			req, err := parseRequestForRelatedArtists(ctx, id)
			assert.Nil(t, err)
			assert.Equal(t, req.Header.Get("Authorization"), fmt.Sprint("Bearer ", token))
		})
	})

	t.Run("TestParseResponseForGetArtists", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			bytes := []byte(getRelatedArtistsPayload)

			as, err := parseResponseForRelatedArtists(&bytes)
			assert.Nil(t, err)
			assert.True(t, len(*as) == 20, len(*as))
		})

		t.Run("bad body", func(t *testing.T) {
			bytes := []byte("fdakslfjda;klfjad;kjadl;")
			_, err := parseResponseForRelatedArtists(&bytes)
			assert.NotNil(t, err)
		})
	})

	t.Run("MainMethod", func(t *testing.T) {
		t.Run("HappyPath", func(t *testing.T) {
			ctx := getTestDependencies(context.Background(), 200, "{}")
			ctx = context.WithValue(ctx, keys.ContextSpotifyAccessToken, "test")
			as := &Artist{ID: "test"}

			_, err := as.GetRelatedArtists(ctx)
			assert.Equal(t, nil, err)
		})

		t.Run("BadRequest", func(t *testing.T) {
			ctx := getTestDependencies(context.Background(), 400, `{"err":"bad_request"}`)
			ctx = context.WithValue(ctx, keys.ContextSpotifyAccessToken, "test")
			as := &Artist{ID: "test"}

			_, err := as.GetRelatedArtists(ctx)
			assert.NotEqual(t, nil, err)
		})
	})
}

var (
	getTopArtistsPayload = `{
		"items": [
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/6FBDaR13swtiWwGhX1WQsP"
			},
			"followers": {
			  "href": null,
			  "total": 6422687
			},
			"genres": [
			  "pop punk",
			  "punk",
			  "socal pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/6FBDaR13swtiWwGhX1WQsP",
			"id": "6FBDaR13swtiWwGhX1WQsP",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/ab6761610000e5ebbf402d5a7cbaac5ab2cccd79",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/ab67616100005174bf402d5a7cbaac5ab2cccd79",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/ab6761610000f178bf402d5a7cbaac5ab2cccd79",
				"width": 160
			  }
			],
			"name": "blink-182",
			"popularity": 81,
			"type": "artist",
			"uri": "spotify:artist:6FBDaR13swtiWwGhX1WQsP"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/1lKZzN2d4IqiEYxyECIEHI"
			},
			"followers": {
			  "href": null,
			  "total": 60212
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/1lKZzN2d4IqiEYxyECIEHI",
			"id": "1lKZzN2d4IqiEYxyECIEHI",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/10bca51ac67e0b763c8ebff5b2c33c69ea2be0ff",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/31e79c2077ab3b9a44a150d1e96bee47f8df1534",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/33214e2a41dc93019ddf09def75fff3dfdb57660",
				"width": 160
			  }
			],
			"name": "Hot Mulligan",
			"popularity": 56,
			"type": "artist",
			"uri": "spotify:artist:1lKZzN2d4IqiEYxyECIEHI"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
			},
			"followers": {
			  "href": null,
			  "total": 107428
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "emo",
			  "midwest emo"
			],
			"href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
			"id": "5rJVTTK0ucAxQhkUc0nXbH",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/82ebe2932c0af13a80a6b21a0df713bea1b32baf",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/1f69254014eab5f42cf53ffbb91ded2255dfbf4d",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/ed6be1cf19997f26b1e70830b3986d59c24255b2",
				"width": 160
			  }
			],
			"name": "Tiny Moving Parts",
			"popularity": 50,
			"type": "artist",
			"uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/4UXqAaa6dQYAk18Lv7PEgX"
			},
			"followers": {
			  "href": null,
			  "total": 8741272
			},
			"genres": [
			  "emo",
			  "modern rock",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/4UXqAaa6dQYAk18Lv7PEgX",
			"id": "4UXqAaa6dQYAk18Lv7PEgX",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/078a111caaad88290dfa51d1aebf76f305f63dbf",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/abf638dd7e70c34518f4c2972c8d3c934456713d",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/194afe6aa362b3b1c3ea51744776d96e800be9c9",
				"width": 160
			  }
			],
			"name": "Fall Out Boy",
			"popularity": 83,
			"type": "artist",
			"uri": "spotify:artist:4UXqAaa6dQYAk18Lv7PEgX"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/0nq64XZMWV1s7XHXIkdH7K"
			},
			"followers": {
			  "href": null,
			  "total": 220852
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "emo",
			  "philly indie",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/0nq64XZMWV1s7XHXIkdH7K",
			"id": "0nq64XZMWV1s7XHXIkdH7K",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/98180f99debca14933645e9b01f0eb6167c0a9f5",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/45456be87f29a5cde73d9d6ae979019ef5d3c4f2",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/46d071ca8f178f8adcf41f97448c2a1dffd63d20",
				"width": 160
			  }
			],
			"name": "The Wonder Years",
			"popularity": 56,
			"type": "artist",
			"uri": "spotify:artist:0nq64XZMWV1s7XHXIkdH7K"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/6dEtLwgmSI0hmfwTSjy8cw"
			},
			"followers": {
			  "href": null,
			  "total": 224379
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "chicago pop punk",
			  "pixie",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/6dEtLwgmSI0hmfwTSjy8cw",
			"id": "6dEtLwgmSI0hmfwTSjy8cw",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/6f8c1161b911111f6ab9879094a671dc12f1fcc4",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/090914ce23df7b91b99a569cdb703d73eb64e92b",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/4dff248e0b251cfa639f658c34973f4cef841adb",
				"width": 160
			  }
			],
			"name": "Real Friends",
			"popularity": 54,
			"type": "artist",
			"uri": "spotify:artist:6dEtLwgmSI0hmfwTSjy8cw"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/1zNfnkHqbNqPMm0LY98Tfx"
			},
			"followers": {
			  "href": null,
			  "total": 20341
			},
			"genres": [
			  "alternative emo",
			  "diy emo"
			],
			"href": "https://api.spotify.com/v1/artists/1zNfnkHqbNqPMm0LY98Tfx",
			"id": "1zNfnkHqbNqPMm0LY98Tfx",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/45c074dbaf52e45072dd4094fb025557cd53051f",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/8aace54c02986add836c8818cd48c0f86f6c9c93",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/c56e44f7b1c150ae96ba54568e7c723173ffd23a",
				"width": 160
			  }
			],
			"name": "fredo disco",
			"popularity": 47,
			"type": "artist",
			"uri": "spotify:artist:1zNfnkHqbNqPMm0LY98Tfx"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/4NiJW4q9ichVqL1aUsgGAN"
			},
			"followers": {
			  "href": null,
			  "total": 1864020
			},
			"genres": [
			  "metalcore",
			  "pop punk",
			  "screamo"
			],
			"href": "https://api.spotify.com/v1/artists/4NiJW4q9ichVqL1aUsgGAN",
			"id": "4NiJW4q9ichVqL1aUsgGAN",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/25e98c889e5cadd0b053c5e70dc8a660facd880c",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/f79bef47494cc26ba3dc47c3e499d8a5af7f65b9",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/7fa31df9354d8406dc364f3aa9c7facf4b73266e",
				"width": 160
			  }
			],
			"name": "A Day To Remember",
			"popularity": 76,
			"type": "artist",
			"uri": "spotify:artist:4NiJW4q9ichVqL1aUsgGAN"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/1qqdO7xMptucPDMopsOdkr"
			},
			"followers": {
			  "href": null,
			  "total": 258248
			},
			"genres": [
			  "anthem emo",
			  "pixie",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/1qqdO7xMptucPDMopsOdkr",
			"id": "1qqdO7xMptucPDMopsOdkr",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/9deaf129e22e2da8cbe9791e76794baa95bc8b96",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/13379a7aa44da87cbfce16ec112a26ccb0594995",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/4c6e6e9b97367e59874464306bff0ad9641fc46f",
				"width": 160
			  }
			],
			"name": "State Champs",
			"popularity": 60,
			"type": "artist",
			"uri": "spotify:artist:1qqdO7xMptucPDMopsOdkr"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/3WfJ1OtrWI7RViX9DMyEGy"
			},
			"followers": {
			  "href": null,
			  "total": 1174189
			},
			"genres": [
			  "emo",
			  "neon pop punk",
			  "pop emo",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/3WfJ1OtrWI7RViX9DMyEGy",
			"id": "3WfJ1OtrWI7RViX9DMyEGy",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/ab6761610000e5ebe110010d1bd030263c26d718",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/ab67616100005174e110010d1bd030263c26d718",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/ab6761610000f178e110010d1bd030263c26d718",
				"width": 160
			  }
			],
			"name": "Mayday Parade",
			"popularity": 67,
			"type": "artist",
			"uri": "spotify:artist:3WfJ1OtrWI7RViX9DMyEGy"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/24XtlMhEMNdi822vi0MhY1"
			},
			"followers": {
			  "href": null,
			  "total": 690825
			},
			"genres": [
			  "emo",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/24XtlMhEMNdi822vi0MhY1",
			"id": "24XtlMhEMNdi822vi0MhY1",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/ab6761610000e5eb962c06d6a4e0b5f3ef815083",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/ab67616100005174962c06d6a4e0b5f3ef815083",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/ab6761610000f178962c06d6a4e0b5f3ef815083",
				"width": 160
			  }
			],
			"name": "Taking Back Sunday",
			"popularity": 65,
			"type": "artist",
			"uri": "spotify:artist:24XtlMhEMNdi822vi0MhY1"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/20JZFwl6HVl6yg8a4H3ZqK"
			},
			"followers": {
			  "href": null,
			  "total": 10632209
			},
			"genres": [
			  "baroque pop",
			  "emo",
			  "modern rock"
			],
			"href": "https://api.spotify.com/v1/artists/20JZFwl6HVl6yg8a4H3ZqK",
			"id": "20JZFwl6HVl6yg8a4H3ZqK",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/58518a04cdd1f20a24cf0545838f3a7b775f8080",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/e71f9ba6573c95041ecd71f766788668f1ceb998",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/5b0544b8a3898be8693915ef7bffb216f781d23d",
				"width": 160
			  }
			],
			"name": "Panic! At The Disco",
			"popularity": 83,
			"type": "artist",
			"uri": "spotify:artist:20JZFwl6HVl6yg8a4H3ZqK"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/6olE6TJLqED3rqDCT0FyPh"
			},
			"followers": {
			  "href": null,
			  "total": 12500735
			},
			"genres": [
			  "alternative rock",
			  "grunge",
			  "permanent wave",
			  "rock"
			],
			"href": "https://api.spotify.com/v1/artists/6olE6TJLqED3rqDCT0FyPh",
			"id": "6olE6TJLqED3rqDCT0FyPh",
			"images": [
			  {
				"height": 1057,
				"url": "https://i.scdn.co/image/84282c28d851a700132356381fcfbadc67ff498b",
				"width": 1000
			  },
			  {
				"height": 677,
				"url": "https://i.scdn.co/image/a4e10b79a642e9891383448cbf37d7266a6883d6",
				"width": 640
			  },
			  {
				"height": 211,
				"url": "https://i.scdn.co/image/42ae0f180f16e2f21c1f2212717fc436f5b95451",
				"width": 200
			  },
			  {
				"height": 68,
				"url": "https://i.scdn.co/image/e797ad36d56c3fc8fa06c6fe91263a15bf8391b8",
				"width": 64
			  }
			],
			"name": "Nirvana",
			"popularity": 84,
			"type": "artist",
			"uri": "spotify:artist:6olE6TJLqED3rqDCT0FyPh"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/0DK7FqcaL3ks9TfFn9y1sD"
			},
			"followers": {
			  "href": null,
			  "total": 316653
			},
			"genres": [
			  "pop punk",
			  "punk",
			  "skate punk"
			],
			"href": "https://api.spotify.com/v1/artists/0DK7FqcaL3ks9TfFn9y1sD",
			"id": "0DK7FqcaL3ks9TfFn9y1sD",
			"images": [
			  {
				"height": 1000,
				"url": "https://i.scdn.co/image/420feb7fab6622848736cff1aa11739eae13849f",
				"width": 1000
			  },
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/677591bff2b6391532f1e14bbeb32f043be3f307",
				"width": 640
			  },
			  {
				"height": 200,
				"url": "https://i.scdn.co/image/85d1309dd1ab243f07ba898583d00be0409f394d",
				"width": 200
			  },
			  {
				"height": 64,
				"url": "https://i.scdn.co/image/b1d98d120878dae011615b826fbd848370f57c97",
				"width": 64
			  }
			],
			"name": "Box Car Racer",
			"popularity": 52,
			"type": "artist",
			"uri": "spotify:artist:0DK7FqcaL3ks9TfFn9y1sD"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/0q7m6CrCFIwGBPhYPxr55O"
			},
			"followers": {
			  "href": null,
			  "total": 16428
			},
			"genres": [
			  "neon pop punk",
			  "pop emo",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/0q7m6CrCFIwGBPhYPxr55O",
			"id": "0q7m6CrCFIwGBPhYPxr55O",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/1b377f55ba9a83871f4bcac9ad65133e52e4bcd7",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/9f49a02e19989f5ece6b211483af3e00061d5e56",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/bf2c314f63f847eeab9d0b0ad2c30e45e2c2720d",
				"width": 160
			  }
			],
			"name": "You, Me, And Everyone We Know",
			"popularity": 34,
			"type": "artist",
			"uri": "spotify:artist:0q7m6CrCFIwGBPhYPxr55O"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/6Xq9CIMYWK4RCrMVtfEOM0"
			},
			"followers": {
			  "href": null,
			  "total": 61401
			},
			"genres": [
			  "anthem emo",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/6Xq9CIMYWK4RCrMVtfEOM0",
			"id": "6Xq9CIMYWK4RCrMVtfEOM0",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/43f5bc0879a9becbe54a07fe624e0fbac53a497b",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/99a413eebbd4af3592c7be03cbf4b96da6ed47bf",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/e9db072dec2351dca63ad0e63656f99702ec5fbd",
				"width": 160
			  }
			],
			"name": "Grayscale",
			"popularity": 52,
			"type": "artist",
			"uri": "spotify:artist:6Xq9CIMYWK4RCrMVtfEOM0"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/4iJLPqClelZOBCBifm8Fzv"
			},
			"followers": {
			  "href": null,
			  "total": 1696353
			},
			"genres": [
			  "modern rock",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/4iJLPqClelZOBCBifm8Fzv",
			"id": "4iJLPqClelZOBCBifm8Fzv",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/32c9cc8f05267eb133588c9fee58a00875b42965",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/39ecfaa2cb6574e81d95dc09b1ae9d17b078a5a5",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/8e9ed2401afe2cdfb7daf86f0d6eff564c006cf1",
				"width": 160
			  }
			],
			"name": "Pierce The Veil",
			"popularity": 71,
			"type": "artist",
			"uri": "spotify:artist:4iJLPqClelZOBCBifm8Fzv"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/7FBcuc1gsnv6Y1nwFtNRCb"
			},
			"followers": {
			  "href": null,
			  "total": 5826339
			},
			"genres": [
			  "emo",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/7FBcuc1gsnv6Y1nwFtNRCb",
			"id": "7FBcuc1gsnv6Y1nwFtNRCb",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/bab47daddd2c01b0ee83e54536aa7e2c77ba7c14",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/5182b00b3542c77889971c8618c4a46eded49e2a",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/5cfd00d5a655ea9b7466f2a0e3372505cb399870",
				"width": 160
			  }
			],
			"name": "My Chemical Romance",
			"popularity": 80,
			"type": "artist",
			"uri": "spotify:artist:7FBcuc1gsnv6Y1nwFtNRCb"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/1ZRtiZJWTAW1LEwKe7F3zz"
			},
			"followers": {
			  "href": null,
			  "total": 23266
			},
			"genres": [
			  "neon pop punk",
			  "pop emo",
			  "pop punk",
			  "screamo"
			],
			"href": "https://api.spotify.com/v1/artists/1ZRtiZJWTAW1LEwKe7F3zz",
			"id": "1ZRtiZJWTAW1LEwKe7F3zz",
			"images": [
			  {
				"height": 400,
				"url": "https://i.scdn.co/image/d738930ded0b900fef5956a9844345cedece8047",
				"width": 600
			  },
			  {
				"height": 133,
				"url": "https://i.scdn.co/image/479ed2a2922967012f21b3d322e432a990e97214",
				"width": 200
			  },
			  {
				"height": 43,
				"url": "https://i.scdn.co/image/42fd722056a5a44a4a6a131ae532da637fab6650",
				"width": 64
			  }
			],
			"name": "Just Surrender",
			"popularity": 33,
			"type": "artist",
			"uri": "spotify:artist:1ZRtiZJWTAW1LEwKe7F3zz"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/3jwm6OBdUY5xzFiFIPhMHu"
			},
			"followers": {
			  "href": null,
			  "total": 170662
			},
			"genres": [
			  "emo",
			  "neon pop punk",
			  "pop emo",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/3jwm6OBdUY5xzFiFIPhMHu",
			"id": "3jwm6OBdUY5xzFiFIPhMHu",
			"images": [
			  {
				"height": 600,
				"url": "https://i.scdn.co/image/566bb9fee235601e08a026fb84193610b68e578a",
				"width": 400
			  },
			  {
				"height": 300,
				"url": "https://i.scdn.co/image/2c235aff504845ea0ed9bead6428e5a22c775cea",
				"width": 200
			  },
			  {
				"height": 96,
				"url": "https://i.scdn.co/image/fc9bf900f8a1c3974afba23cb3aee937bef857b7",
				"width": 64
			  }
			],
			"name": "The Academy Is...",
			"popularity": 49,
			"type": "artist",
			"uri": "spotify:artist:3jwm6OBdUY5xzFiFIPhMHu"
		  }
		],
		"total": 50,
		"limit": 20,
		"offset": 0,
		"href": "https://api.spotify.com/v1/me/top/artists",
		"previous": null,
		"next": "https://api.spotify.com/v1/me/top/artists?limit=20&offset=20"
	  }`

	getRelatedArtistsPayload = `{
		"artists": [
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/0gLjJuczGWhqKVMmVpIT52"
			},
			"followers": {
			  "href": null,
			  "total": 219531
			},
			"genres": [
			  "alternative emo",
			  "anthem emo"
			],
			"href": "https://api.spotify.com/v1/artists/0gLjJuczGWhqKVMmVpIT52",
			"id": "0gLjJuczGWhqKVMmVpIT52",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/c6adfe401fd807a4fa83b1b372e376a1c9b27079",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/fab1eb535cd28339c867e48bc17b6b5da30d70df",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/7928511b54d382eb8d30287bcdecae6137e9a90a",
				"width": 160
			  }
			],
			"name": "Turnover",
			"popularity": 60,
			"type": "artist",
			"uri": "spotify:artist:0gLjJuczGWhqKVMmVpIT52"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/0znuUIjvP0LXEslfaq0Nor"
			},
			"followers": {
			  "href": null,
			  "total": 170818
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/0znuUIjvP0LXEslfaq0Nor",
			"id": "0znuUIjvP0LXEslfaq0Nor",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/3a71462b13a0df8d962cf9ad6b9c56068fbf6a8a",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/cce55b6de9647caceeffd8a809ae7d426a21896c",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/ddd1ec75d765898567c3bea7ae6aa7ecdfcc3233",
				"width": 160
			  }
			],
			"name": "Citizen",
			"popularity": 60,
			"type": "artist",
			"uri": "spotify:artist:0znuUIjvP0LXEslfaq0Nor"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/49b68DLRK5eCbtJf7Xx4Cc"
			},
			"followers": {
			  "href": null,
			  "total": 48040
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "emo",
			  "midwest emo",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/49b68DLRK5eCbtJf7Xx4Cc",
			"id": "49b68DLRK5eCbtJf7Xx4Cc",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/ab6761610000e5eb79308deeb1964b064b4332bf",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/ab6761610000517479308deeb1964b064b4332bf",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/ab6761610000f17879308deeb1964b064b4332bf",
				"width": 160
			  }
			],
			"name": "Free Throw",
			"popularity": 49,
			"type": "artist",
			"uri": "spotify:artist:49b68DLRK5eCbtJf7Xx4Cc"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
			},
			"followers": {
			  "href": null,
			  "total": 107428
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "emo",
			  "midwest emo"
			],
			"href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
			"id": "5rJVTTK0ucAxQhkUc0nXbH",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/82ebe2932c0af13a80a6b21a0df713bea1b32baf",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/1f69254014eab5f42cf53ffbb91ded2255dfbf4d",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/ed6be1cf19997f26b1e70830b3986d59c24255b2",
				"width": 160
			  }
			],
			"name": "Tiny Moving Parts",
			"popularity": 50,
			"type": "artist",
			"uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/5rOjuB5uYAoDMHgZM6CFBB"
			},
			"followers": {
			  "href": null,
			  "total": 41789
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/5rOjuB5uYAoDMHgZM6CFBB",
			"id": "5rOjuB5uYAoDMHgZM6CFBB",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/8d3a606700e4cd32812ffcf686470ab2edc02fd9",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/761389b7562e2b0ef951c8e596136beedd46a100",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/75c7a3b9a64d304378b4bed234687091985a8269",
				"width": 160
			  }
			],
			"name": "Mat Kerekes",
			"popularity": 47,
			"type": "artist",
			"uri": "spotify:artist:5rOjuB5uYAoDMHgZM6CFBB"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/7jXQfCP0xRnnind08ie0zT"
			},
			"followers": {
			  "href": null,
			  "total": 132465
			},
			"genres": [
			  "alternative emo",
			  "emo",
			  "new england emo"
			],
			"href": "https://api.spotify.com/v1/artists/7jXQfCP0xRnnind08ie0zT",
			"id": "7jXQfCP0xRnnind08ie0zT",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/63de7da20035abb798dc0c8b3e7210e0df7eac31",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/b7e4300e1861a2c4955aa4364ace5193c1fe71ff",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/f82773a75730fe167b01783a1ddb66950eadd527",
				"width": 160
			  }
			],
			"name": "Sorority Noise",
			"popularity": 54,
			"type": "artist",
			"uri": "spotify:artist:7jXQfCP0xRnnind08ie0zT"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/0nq64XZMWV1s7XHXIkdH7K"
			},
			"followers": {
			  "href": null,
			  "total": 220852
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "emo",
			  "philly indie",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/0nq64XZMWV1s7XHXIkdH7K",
			"id": "0nq64XZMWV1s7XHXIkdH7K",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/98180f99debca14933645e9b01f0eb6167c0a9f5",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/45456be87f29a5cde73d9d6ae979019ef5d3c4f2",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/46d071ca8f178f8adcf41f97448c2a1dffd63d20",
				"width": 160
			  }
			],
			"name": "The Wonder Years",
			"popularity": 56,
			"type": "artist",
			"uri": "spotify:artist:0nq64XZMWV1s7XHXIkdH7K"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/45V5yoo6fI5r3m7kei0onQ"
			},
			"followers": {
			  "href": null,
			  "total": 16362
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/45V5yoo6fI5r3m7kei0onQ",
			"id": "45V5yoo6fI5r3m7kei0onQ",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/d80f777bff65235a82b83dcd57277c9a668ada33",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/8e24e3b74aaf565843d78d8035b87ae375812f5c",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/58d1389fd5f2e86d18d66df18549ee0662785f77",
				"width": 160
			  }
			],
			"name": "Elder Brother",
			"popularity": 35,
			"type": "artist",
			"uri": "spotify:artist:45V5yoo6fI5r3m7kei0onQ"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/7ptm7G8z8VVvwBnDq8fAmD"
			},
			"followers": {
			  "href": null,
			  "total": 69448
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "atlanta punk",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/7ptm7G8z8VVvwBnDq8fAmD",
			"id": "7ptm7G8z8VVvwBnDq8fAmD",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/3900cf2a69f74232c89f54bdc1601df63f7d5313",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/9ff15ff0372f3480918e835cd49626becc2f6f52",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/82ac09475fafdc2a87def9b086f5b425dccb414e",
				"width": 160
			  }
			],
			"name": "Microwave",
			"popularity": 51,
			"type": "artist",
			"uri": "spotify:artist:7ptm7G8z8VVvwBnDq8fAmD"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/3kzNckjE5FzHQhe4pJiLKa"
			},
			"followers": {
			  "href": null,
			  "total": 65902
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "emo",
			  "indie punk",
			  "indie rock",
			  "midwest emo",
			  "new england emo",
			  "worcester ma indie"
			],
			"href": "https://api.spotify.com/v1/artists/3kzNckjE5FzHQhe4pJiLKa",
			"id": "3kzNckjE5FzHQhe4pJiLKa",
			"images": [
			  {
				"height": 1000,
				"url": "https://i.scdn.co/image/2b9ff272ef806f5c555bdc87668aa3350258753f",
				"width": 1000
			  },
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/98da84b703fcfa2a991c4614180ba23707c95a9e",
				"width": 640
			  },
			  {
				"height": 200,
				"url": "https://i.scdn.co/image/3311fec0a3309d6a93fbab0836424130a4a5b9ff",
				"width": 200
			  },
			  {
				"height": 64,
				"url": "https://i.scdn.co/image/458b8b1eb76409b7509520fdcae4739678fe3765",
				"width": 64
			  }
			],
			"name": "The Hotelier",
			"popularity": 44,
			"type": "artist",
			"uri": "spotify:artist:3kzNckjE5FzHQhe4pJiLKa"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/6PsktPFR0UZptKdSqmlS5h"
			},
			"followers": {
			  "href": null,
			  "total": 148871
			},
			"genres": [
			  "alternative emo",
			  "diy emo",
			  "indie punk",
			  "north carolina emo"
			],
			"href": "https://api.spotify.com/v1/artists/6PsktPFR0UZptKdSqmlS5h",
			"id": "6PsktPFR0UZptKdSqmlS5h",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/53d8511d821a392bf883f7db16bbbf8f42f11041",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/53df1667a06c2803d6dbb5ff3c5c915dbf3bad14",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/660356d49dcc8db40417fc1456161a610db641a0",
				"width": 160
			  }
			],
			"name": "Mom Jeans.",
			"popularity": 57,
			"type": "artist",
			"uri": "spotify:artist:6PsktPFR0UZptKdSqmlS5h"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/1lKZzN2d4IqiEYxyECIEHI"
			},
			"followers": {
			  "href": null,
			  "total": 60212
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/1lKZzN2d4IqiEYxyECIEHI",
			"id": "1lKZzN2d4IqiEYxyECIEHI",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/10bca51ac67e0b763c8ebff5b2c33c69ea2be0ff",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/31e79c2077ab3b9a44a150d1e96bee47f8df1534",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/33214e2a41dc93019ddf09def75fff3dfdb57660",
				"width": 160
			  }
			],
			"name": "Hot Mulligan",
			"popularity": 56,
			"type": "artist",
			"uri": "spotify:artist:1lKZzN2d4IqiEYxyECIEHI"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/5vV4gEs3O35SdrdwhvhYwe"
			},
			"followers": {
			  "href": null,
			  "total": 20765
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "chicago pop punk",
			  "easycore",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/5vV4gEs3O35SdrdwhvhYwe",
			"id": "5vV4gEs3O35SdrdwhvhYwe",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/536024d6af2fa67552636a7d855baac8293e52f9",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/9eaf106293920a976c1aaf7dd764304c72398c3d",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/433643a250742a404342551c3903468e5acb9625",
				"width": 160
			  }
			],
			"name": "Homesafe",
			"popularity": 34,
			"type": "artist",
			"uri": "spotify:artist:5vV4gEs3O35SdrdwhvhYwe"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/2dfxY7YDuYCUtWFzWTS6IR"
			},
			"followers": {
			  "href": null,
			  "total": 95611
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "emo",
			  "stl indie"
			],
			"href": "https://api.spotify.com/v1/artists/2dfxY7YDuYCUtWFzWTS6IR",
			"id": "2dfxY7YDuYCUtWFzWTS6IR",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/ab6761610000e5eb0610a2e881c33f18218ee3fd",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/ab676161000051740610a2e881c33f18218ee3fd",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/ab6761610000f1780610a2e881c33f18218ee3fd",
				"width": 160
			  }
			],
			"name": "Foxing",
			"popularity": 49,
			"type": "artist",
			"uri": "spotify:artist:2dfxY7YDuYCUtWFzWTS6IR"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/6KPdmtIl0LA5mRFSqseWhI"
			},
			"followers": {
			  "href": null,
			  "total": 76393
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "aussie emo",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/6KPdmtIl0LA5mRFSqseWhI",
			"id": "6KPdmtIl0LA5mRFSqseWhI",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/9b8ddec1c043396afd936dffedc5e48b3173f70e",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/871271589b5c62195fafcb16983ea99f6b05f72d",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/b54bb30098199f57316b06f5cf37f4cf5dd39772",
				"width": 160
			  }
			],
			"name": "Trophy Eyes",
			"popularity": 50,
			"type": "artist",
			"uri": "spotify:artist:6KPdmtIl0LA5mRFSqseWhI"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/24M8W1AklCxyWTKjrJZDQ8"
			},
			"followers": {
			  "href": null,
			  "total": 54200
			},
			"genres": [
			  "alternative emo",
			  "charlotte nc indie",
			  "indie garage rock",
			  "midwest emo",
			  "north carolina emo"
			],
			"href": "https://api.spotify.com/v1/artists/24M8W1AklCxyWTKjrJZDQ8",
			"id": "24M8W1AklCxyWTKjrJZDQ8",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/29d189be619e9d310d6806f35088a6b49cdea8e2",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/dc3fb3aa2d61783d411adca297949c95e9179f20",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/54ce52e66d8a5d2a104ed763ce6f5542e91c257d",
				"width": 160
			  }
			],
			"name": "It Looks Sad.",
			"popularity": 47,
			"type": "artist",
			"uri": "spotify:artist:24M8W1AklCxyWTKjrJZDQ8"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/2kbAovdYb7krLSGdOrBMRu"
			},
			"followers": {
			  "href": null,
			  "total": 29430
			},
			"genres": [
			  "alternative emo",
			  "modern alternative rock",
			  "oc indie"
			],
			"href": "https://api.spotify.com/v1/artists/2kbAovdYb7krLSGdOrBMRu",
			"id": "2kbAovdYb7krLSGdOrBMRu",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/977af33a411a1d4f4c3f9b8bbdd51cf10cbf6fbe",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/f111b5717e41c15a264707f6bf9dfd5d9cc5cb01",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/77f91a4a799f80fe6bcce3b84ecf1c6076715534",
				"width": 160
			  }
			],
			"name": "Super Whatevr",
			"popularity": 48,
			"type": "artist",
			"uri": "spotify:artist:2kbAovdYb7krLSGdOrBMRu"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/4UXiNDHAiv8DOSLkp0GbSm"
			},
			"followers": {
			  "href": null,
			  "total": 135663
			},
			"genres": [
			  "acoustic rock"
			],
			"href": "https://api.spotify.com/v1/artists/4UXiNDHAiv8DOSLkp0GbSm",
			"id": "4UXiNDHAiv8DOSLkp0GbSm",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/02ea3e23fd0e14463a7a69b0b8349773b2549ac2",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/c0f610f3465e6ca613cc9185065cb1904b86daa6",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/791fae52f4f77ae87f45b17e13a206b6c18956d6",
				"width": 160
			  }
			],
			"name": "Front Porch Step",
			"popularity": 52,
			"type": "artist",
			"uri": "spotify:artist:4UXiNDHAiv8DOSLkp0GbSm"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/5LMPXUMhWXshBPjrqvZOfv"
			},
			"followers": {
			  "href": null,
			  "total": 104751
			},
			"genres": [
			  "alternative emo",
			  "indie garage rock"
			],
			"href": "https://api.spotify.com/v1/artists/5LMPXUMhWXshBPjrqvZOfv",
			"id": "5LMPXUMhWXshBPjrqvZOfv",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/4d144fb4585ce4e4b069749ebc221e9fc2ec801e",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/7885db74984df05a3b4f35c18e15c337094bc525",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/eef18d8bbe2248f785eaa23d26905e99eb4dcb8d",
				"width": 160
			  }
			],
			"name": "Remo Drive",
			"popularity": 53,
			"type": "artist",
			"uri": "spotify:artist:5LMPXUMhWXshBPjrqvZOfv"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/1u3l44NMyNo8xe8ykBZtFp"
			},
			"followers": {
			  "href": null,
			  "total": 17008
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "new england emo",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/1u3l44NMyNo8xe8ykBZtFp",
			"id": "1u3l44NMyNo8xe8ykBZtFp",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/d0c12795c5b73393325e0f00575c11839eb3b556",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/97f0542d4e96e6819672a9cd7ca6ef5e57518c99",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/ffd934e64f8469945910caedc02e47a0ea2a4bf8",
				"width": 160
			  }
			],
			"name": "Somos",
			"popularity": 33,
			"type": "artist",
			"uri": "spotify:artist:1u3l44NMyNo8xe8ykBZtFp"
		  }
		]
	  }`

	getArtistsPayload = `{
		"artists": [
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/6FBDaR13swtiWwGhX1WQsP"
			},
			"followers": {
			  "href": null,
			  "total": 6424396
			},
			"genres": [
			  "pop punk",
			  "punk",
			  "socal pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/6FBDaR13swtiWwGhX1WQsP",
			"id": "6FBDaR13swtiWwGhX1WQsP",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/ab6761610000e5ebbf402d5a7cbaac5ab2cccd79",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/ab67616100005174bf402d5a7cbaac5ab2cccd79",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/ab6761610000f178bf402d5a7cbaac5ab2cccd79",
				"width": 160
			  }
			],
			"name": "blink-182",
			"popularity": 81,
			"type": "artist",
			"uri": "spotify:artist:6FBDaR13swtiWwGhX1WQsP"
		  },
		  {
			"external_urls": {
			  "spotify": "https://open.spotify.com/artist/1lKZzN2d4IqiEYxyECIEHI"
			},
			"followers": {
			  "href": null,
			  "total": 60212
			},
			"genres": [
			  "alternative emo",
			  "anthem emo",
			  "pop punk"
			],
			"href": "https://api.spotify.com/v1/artists/1lKZzN2d4IqiEYxyECIEHI",
			"id": "1lKZzN2d4IqiEYxyECIEHI",
			"images": [
			  {
				"height": 640,
				"url": "https://i.scdn.co/image/10bca51ac67e0b763c8ebff5b2c33c69ea2be0ff",
				"width": 640
			  },
			  {
				"height": 320,
				"url": "https://i.scdn.co/image/31e79c2077ab3b9a44a150d1e96bee47f8df1534",
				"width": 320
			  },
			  {
				"height": 160,
				"url": "https://i.scdn.co/image/33214e2a41dc93019ddf09def75fff3dfdb57660",
				"width": 160
			  }
			],
			"name": "Hot Mulligan",
			"popularity": 56,
			"type": "artist",
			"uri": "spotify:artist:1lKZzN2d4IqiEYxyECIEHI"
		  }
		]
	  }`

	getArtistPayload = `{
		"external_urls": {
		  "spotify": "https://open.spotify.com/artist/6FBDaR13swtiWwGhX1WQsP"
		},
		"followers": {
		  "href": null,
		  "total": 6424396
		},
		"genres": [
		  "pop punk",
		  "punk",
		  "socal pop punk"
		],
		"href": "https://api.spotify.com/v1/artists/6FBDaR13swtiWwGhX1WQsP",
		"id": "6FBDaR13swtiWwGhX1WQsP",
		"images": [
		  {
			"height": 640,
			"url": "https://i.scdn.co/image/ab6761610000e5ebbf402d5a7cbaac5ab2cccd79",
			"width": 640
		  },
		  {
			"height": 320,
			"url": "https://i.scdn.co/image/ab67616100005174bf402d5a7cbaac5ab2cccd79",
			"width": 320
		  },
		  {
			"height": 160,
			"url": "https://i.scdn.co/image/ab6761610000f178bf402d5a7cbaac5ab2cccd79",
			"width": 160
		  }
		],
		"name": "blink-182",
		"popularity": 81,
		"type": "artist",
		"uri": "spotify:artist:6FBDaR13swtiWwGhX1WQsP"
	  }`
)
