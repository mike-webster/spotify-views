package spotify

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/mike-webster/spotify-views/keys"
	"github.com/stretchr/testify/assert"
)

func TestGetTopTracks(t *testing.T) {
	t.Run("TestGetTopTracksRequest", func(t *testing.T) {
		ctx := context.Background()
		tf := TFShort
		t.Run("no token", func(t *testing.T) {
			_, err := getTopTracksRequest(ctx, tf)
			assert.Equal(t, reflect.TypeOf(ErrNoToken("")), reflect.TypeOf(err))
		})

		token := "tok"
		ctx = context.WithValue(ctx, keys.ContextSpotifyAccessToken, token)

		t.Run("token gets stored in header", func(t *testing.T) {
			req, err := getTopTracksRequest(ctx, tf)
			assert.Nil(t, err)
			assert.Equal(t, req.Header.Get("Authorization"), fmt.Sprint("Bearer ", token))
		})
	})

	t.Run("TestParseTopTrackResponse", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			bytes := []byte(getTopTracksPayload)

			as, err := parseTopTrackResponse(&bytes)
			assert.Nil(t, err)
			assert.NotNil(t, as)
		})

		t.Run("bad body", func(t *testing.T) {
			bytes := []byte("fdakslfjda;klfjad;kjadl;")
			_, err := parseTopTrackResponse(&bytes)
			assert.NotNil(t, err)
		})
	})
}

func TestTopTracksForArtists(t *testing.T) {
	t.Run("TestGetTopTracksForArtistRequest", func(t *testing.T) {
		ctx := context.Background()
		id := "1234"
		t.Run("no token", func(t *testing.T) {
			_, err := getTopTracksForArtistRequest(ctx, id)
			assert.Equal(t, reflect.TypeOf(ErrNoToken("")), reflect.TypeOf(err))
		})

		token := "tok"
		ctx = context.WithValue(ctx, keys.ContextSpotifyAccessToken, token)

		t.Run("token gets stored in header", func(t *testing.T) {
			req, err := getTopTracksForArtistRequest(ctx, id)
			assert.Nil(t, err)
			assert.Equal(t, req.Header.Get("Authorization"), fmt.Sprint("Bearer ", token))
		})
	})

	t.Run("TestParseTopTracksForArtistResponse", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			bytes := []byte(getTopTracksArtistsPayload)

			as, err := parseTopTracksForArtistResponse(&bytes)
			assert.Nil(t, err)
			assert.NotNil(t, as)
		})

		t.Run("bad body", func(t *testing.T) {
			bytes := []byte("fdakslfjda;klfjad;kjadl;")
			_, err := parseTopTracksForArtistResponse(&bytes)
			assert.NotNil(t, err)
		})
	})
}

func TestGetTrackGenres(t *testing.T) {
	// TODO: need to be able to mock call first
}

func TestEmbeddedPlayerTrack(t *testing.T) {
	tr := Track{ID: "1234"}
	url := fmt.Sprintf(`<iframe src="https://open.spotify.com/embed/track/%s" width="300" height="80" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>`, tr.ID)
	assert.Equal(t, url, tr.EmbeddedPlayer())
}

func TestTracksIDs(t *testing.T) {
	ts := Tracks{Track{ID: "1"}, Track{ID: "2"}}
	assert.Equal(t, []string{"1", "2"}, ts.IDs())
}

func TestFindArtist(t *testing.T) {
	t.Run("NoArtist", func(t *testing.T) {
		tr := Track{}
		assert.Equal(t, "", tr.FindArtist())
	})

	t.Run("HasArtist", func(t *testing.T) {
		a := Artist{Name: "test"}
		tr := Track{Artists: Artists{a}}
		assert.Equal(t, a.Name, tr.FindArtist())
	})
}

func TestFindTrackImage(t *testing.T) {
	t.Run("NoImages", func(t *testing.T) {
		tr := Track{}
		assert.Equal(t, &Image{}, tr.FindImage())
	})

	t.Run("OneImage", func(t *testing.T) {
		i := Image{URL: "test"}
		tr := Track{Album: Album{Images: []Image{i}}}
		assert.Equal(t, &i, tr.FindImage())
	})

	t.Run("ManyImages", func(t *testing.T) {
		i := Image{URL: "test"}
		ii := Image{URL: "test2"}
		iii := Image{URL: "test3"}
		tr := Track{Album: Album{Images: []Image{i, ii, iii}}}
		assert.Equal(t, &ii, tr.FindImage())
	})
}

func TestFindSpotifyURL(t *testing.T) {
	t.Run("NoSpotifyLink", func(t *testing.T) {
		assert.Equal(t, "", (&Track{}).TrySpotifyURL())
	})
	t.Run("SpotifyLink", func(t *testing.T) {
		tr := Track{Links: map[string]string{"spotify": "test"}}
		assert.Equal(t, "test", tr.TrySpotifyURL())
	})
}

var (
	getTopTracksArtistsPayload = `{
		"tracks": [
		  {
			"album": {
			  "album_type": "album",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				  },
				  "href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				  "id": "5rJVTTK0ucAxQhkUc0nXbH",
				  "name": "Tiny Moving Parts",
				  "type": "artist",
				  "uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
				}
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/1JMg02mQ5nmVKBWWoDUIeo"
			  },
			  "href": "https://api.spotify.com/v1/albums/1JMg02mQ5nmVKBWWoDUIeo",
			  "id": "1JMg02mQ5nmVKBWWoDUIeo",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b273c59d701e358b4612be5289e9",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e02c59d701e358b4612be5289e9",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d00004851c59d701e358b4612be5289e9",
				  "width": 64
				}
			  ],
			  "name": "breathe",
			  "release_date": "2019-09-13",
			  "release_date_precision": "day",
			  "total_tracks": 10,
			  "type": "album",
			  "uri": "spotify:album:1JMg02mQ5nmVKBWWoDUIeo"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				},
				"href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				"id": "5rJVTTK0ucAxQhkUc0nXbH",
				"name": "Tiny Moving Parts",
				"type": "artist",
				"uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
			  }
			],
			"disc_number": 1,
			"duration_ms": 191294,
			"explicit": false,
			"external_ids": {
			  "isrc": "USHR21912603"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/0uqTwEKTbkuhrn8yGSO6b5"
			},
			"href": "https://api.spotify.com/v1/tracks/0uqTwEKTbkuhrn8yGSO6b5",
			"id": "0uqTwEKTbkuhrn8yGSO6b5",
			"is_local": false,
			"is_playable": true,
			"name": "Medicine",
			"popularity": 51,
			"preview_url": "https://p.scdn.co/mp3-preview/6c66ead34e15f46119433775516d513b26c7b5f4?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 3,
			"type": "track",
			"uri": "spotify:track:0uqTwEKTbkuhrn8yGSO6b5"
		  },
		  {
			"album": {
			  "album_type": "album",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				  },
				  "href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				  "id": "5rJVTTK0ucAxQhkUc0nXbH",
				  "name": "Tiny Moving Parts",
				  "type": "artist",
				  "uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
				}
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/6glHUYWR8paVp72N7xmBci"
			  },
			  "href": "https://api.spotify.com/v1/albums/6glHUYWR8paVp72N7xmBci",
			  "id": "6glHUYWR8paVp72N7xmBci",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b273012a2f9f54157ff657ecb0e7",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e02012a2f9f54157ff657ecb0e7",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d00004851012a2f9f54157ff657ecb0e7",
				  "width": 64
				}
			  ],
			  "name": "Pleasant Living",
			  "release_date": "2014-09-09",
			  "release_date_precision": "day",
			  "total_tracks": 12,
			  "type": "album",
			  "uri": "spotify:album:6glHUYWR8paVp72N7xmBci"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				},
				"href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				"id": "5rJVTTK0ucAxQhkUc0nXbH",
				"name": "Tiny Moving Parts",
				"type": "artist",
				"uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
			  }
			],
			"disc_number": 1,
			"duration_ms": 160613,
			"explicit": false,
			"external_ids": {
			  "isrc": "US72W1411223"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/20zs5ECJ4on1YqXHiDZkTN"
			},
			"href": "https://api.spotify.com/v1/tracks/20zs5ECJ4on1YqXHiDZkTN",
			"id": "20zs5ECJ4on1YqXHiDZkTN",
			"is_local": false,
			"is_playable": true,
			"name": "Always Focused",
			"popularity": 47,
			"preview_url": "https://p.scdn.co/mp3-preview/d48366bbeddaa7c51fe420f6072078a45ec1fdc5?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 2,
			"type": "track",
			"uri": "spotify:track:20zs5ECJ4on1YqXHiDZkTN"
		  },
		  {
			"album": {
			  "album_type": "album",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				  },
				  "href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				  "id": "5rJVTTK0ucAxQhkUc0nXbH",
				  "name": "Tiny Moving Parts",
				  "type": "artist",
				  "uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
				}
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/25mFiSG2fnKDliGDtyYNpa"
			  },
			  "href": "https://api.spotify.com/v1/albums/25mFiSG2fnKDliGDtyYNpa",
			  "id": "25mFiSG2fnKDliGDtyYNpa",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b273b96bfc99d8dc4ce5084fdd14",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e02b96bfc99d8dc4ce5084fdd14",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d00004851b96bfc99d8dc4ce5084fdd14",
				  "width": 64
				}
			  ],
			  "name": "Swell",
			  "release_date": "2018-01-26",
			  "release_date_precision": "day",
			  "total_tracks": 10,
			  "type": "album",
			  "uri": "spotify:album:25mFiSG2fnKDliGDtyYNpa"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				},
				"href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				"id": "5rJVTTK0ucAxQhkUc0nXbH",
				"name": "Tiny Moving Parts",
				"type": "artist",
				"uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
			  }
			],
			"disc_number": 1,
			"duration_ms": 198794,
			"explicit": false,
			"external_ids": {
			  "isrc": "US72W1821004"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/4fXuiFVbcZV5TDcPszBeip"
			},
			"href": "https://api.spotify.com/v1/tracks/4fXuiFVbcZV5TDcPszBeip",
			"id": "4fXuiFVbcZV5TDcPszBeip",
			"is_local": false,
			"is_playable": true,
			"name": "Caution",
			"popularity": 44,
			"preview_url": "https://p.scdn.co/mp3-preview/06a9e2ddfaf02b9fa4dd553c715e861442d018c7?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 4,
			"type": "track",
			"uri": "spotify:track:4fXuiFVbcZV5TDcPszBeip"
		  },
		  {
			"album": {
			  "album_type": "single",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				  },
				  "href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				  "id": "5rJVTTK0ucAxQhkUc0nXbH",
				  "name": "Tiny Moving Parts",
				  "type": "artist",
				  "uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
				}
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/59VOzUIZFytEY3AuOf44oC"
			  },
			  "href": "https://api.spotify.com/v1/albums/59VOzUIZFytEY3AuOf44oC",
			  "id": "59VOzUIZFytEY3AuOf44oC",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b273083dbff84e7ba5c9dd52b8e0",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e02083dbff84e7ba5c9dd52b8e0",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d00004851083dbff84e7ba5c9dd52b8e0",
				  "width": 64
				}
			  ],
			  "name": "You Lost Me / Guardians",
			  "release_date": "2020-03-11",
			  "release_date_precision": "day",
			  "total_tracks": 2,
			  "type": "album",
			  "uri": "spotify:album:59VOzUIZFytEY3AuOf44oC"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				},
				"href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				"id": "5rJVTTK0ucAxQhkUc0nXbH",
				"name": "Tiny Moving Parts",
				"type": "artist",
				"uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
			  }
			],
			"disc_number": 1,
			"duration_ms": 200778,
			"explicit": false,
			"external_ids": {
			  "isrc": "USHR22012601"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/32IJAFA0QegOuUvz9Qlknl"
			},
			"href": "https://api.spotify.com/v1/tracks/32IJAFA0QegOuUvz9Qlknl",
			"id": "32IJAFA0QegOuUvz9Qlknl",
			"is_local": false,
			"is_playable": true,
			"name": "You Lost Me",
			"popularity": 42,
			"preview_url": "https://p.scdn.co/mp3-preview/c8e5eb738c21e14b95e5ceeb0f250ff5fa9f1e64?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 1,
			"type": "track",
			"uri": "spotify:track:32IJAFA0QegOuUvz9Qlknl"
		  },
		  {
			"album": {
			  "album_type": "album",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				  },
				  "href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				  "id": "5rJVTTK0ucAxQhkUc0nXbH",
				  "name": "Tiny Moving Parts",
				  "type": "artist",
				  "uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
				}
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/1JMg02mQ5nmVKBWWoDUIeo"
			  },
			  "href": "https://api.spotify.com/v1/albums/1JMg02mQ5nmVKBWWoDUIeo",
			  "id": "1JMg02mQ5nmVKBWWoDUIeo",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b273c59d701e358b4612be5289e9",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e02c59d701e358b4612be5289e9",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d00004851c59d701e358b4612be5289e9",
				  "width": 64
				}
			  ],
			  "name": "breathe",
			  "release_date": "2019-09-13",
			  "release_date_precision": "day",
			  "total_tracks": 10,
			  "type": "album",
			  "uri": "spotify:album:1JMg02mQ5nmVKBWWoDUIeo"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				},
				"href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				"id": "5rJVTTK0ucAxQhkUc0nXbH",
				"name": "Tiny Moving Parts",
				"type": "artist",
				"uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
			  }
			],
			"disc_number": 1,
			"duration_ms": 216069,
			"explicit": false,
			"external_ids": {
			  "isrc": "USHR21912605"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/0yBbztMlRUNnWwTaCkx0bl"
			},
			"href": "https://api.spotify.com/v1/tracks/0yBbztMlRUNnWwTaCkx0bl",
			"id": "0yBbztMlRUNnWwTaCkx0bl",
			"is_local": false,
			"is_playable": true,
			"name": "Vertebrae",
			"popularity": 39,
			"preview_url": "https://p.scdn.co/mp3-preview/45d5130ede984e39e2dd65240595a2251816f639?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 5,
			"type": "track",
			"uri": "spotify:track:0yBbztMlRUNnWwTaCkx0bl"
		  },
		  {
			"album": {
			  "album_type": "album",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				  },
				  "href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				  "id": "5rJVTTK0ucAxQhkUc0nXbH",
				  "name": "Tiny Moving Parts",
				  "type": "artist",
				  "uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
				}
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/0hGe2EeB4e2PqTDwEaeqYn"
			  },
			  "href": "https://api.spotify.com/v1/albums/0hGe2EeB4e2PqTDwEaeqYn",
			  "id": "0hGe2EeB4e2PqTDwEaeqYn",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b2737b0c69ed623054a7d9a17e94",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e027b0c69ed623054a7d9a17e94",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d000048517b0c69ed623054a7d9a17e94",
				  "width": 64
				}
			  ],
			  "name": "Celebrate",
			  "release_date": "2016-05-20",
			  "release_date_precision": "day",
			  "total_tracks": 10,
			  "type": "album",
			  "uri": "spotify:album:0hGe2EeB4e2PqTDwEaeqYn"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				},
				"href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				"id": "5rJVTTK0ucAxQhkUc0nXbH",
				"name": "Tiny Moving Parts",
				"type": "artist",
				"uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
			  }
			],
			"disc_number": 1,
			"duration_ms": 181800,
			"explicit": false,
			"external_ids": {
			  "isrc": "US72W1611462"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/4jtUCPu0LU8y8EHG8tnhom"
			},
			"href": "https://api.spotify.com/v1/tracks/4jtUCPu0LU8y8EHG8tnhom",
			"id": "4jtUCPu0LU8y8EHG8tnhom",
			"is_local": false,
			"is_playable": true,
			"name": "Birdhouse",
			"popularity": 38,
			"preview_url": "https://p.scdn.co/mp3-preview/7ef598a95f329a8d2286395ece315514566f449f?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 3,
			"type": "track",
			"uri": "spotify:track:4jtUCPu0LU8y8EHG8tnhom"
		  },
		  {
			"album": {
			  "album_type": "album",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				  },
				  "href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				  "id": "5rJVTTK0ucAxQhkUc0nXbH",
				  "name": "Tiny Moving Parts",
				  "type": "artist",
				  "uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
				}
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/6glHUYWR8paVp72N7xmBci"
			  },
			  "href": "https://api.spotify.com/v1/albums/6glHUYWR8paVp72N7xmBci",
			  "id": "6glHUYWR8paVp72N7xmBci",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b273012a2f9f54157ff657ecb0e7",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e02012a2f9f54157ff657ecb0e7",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d00004851012a2f9f54157ff657ecb0e7",
				  "width": 64
				}
			  ],
			  "name": "Pleasant Living",
			  "release_date": "2014-09-09",
			  "release_date_precision": "day",
			  "total_tracks": 12,
			  "type": "album",
			  "uri": "spotify:album:6glHUYWR8paVp72N7xmBci"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				},
				"href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				"id": "5rJVTTK0ucAxQhkUc0nXbH",
				"name": "Tiny Moving Parts",
				"type": "artist",
				"uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
			  }
			],
			"disc_number": 1,
			"duration_ms": 179733,
			"explicit": false,
			"external_ids": {
			  "isrc": "US72W1411222"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/01uhrXb8QsuIs1HizoqxdJ"
			},
			"href": "https://api.spotify.com/v1/tracks/01uhrXb8QsuIs1HizoqxdJ",
			"id": "01uhrXb8QsuIs1HizoqxdJ",
			"is_local": false,
			"is_playable": true,
			"name": "Sundress",
			"popularity": 38,
			"preview_url": "https://p.scdn.co/mp3-preview/cb6527f63641697067287fa0a6b3a002f588f17e?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 1,
			"type": "track",
			"uri": "spotify:track:01uhrXb8QsuIs1HizoqxdJ"
		  },
		  {
			"album": {
			  "album_type": "album",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				  },
				  "href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				  "id": "5rJVTTK0ucAxQhkUc0nXbH",
				  "name": "Tiny Moving Parts",
				  "type": "artist",
				  "uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
				}
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/1JMg02mQ5nmVKBWWoDUIeo"
			  },
			  "href": "https://api.spotify.com/v1/albums/1JMg02mQ5nmVKBWWoDUIeo",
			  "id": "1JMg02mQ5nmVKBWWoDUIeo",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b273c59d701e358b4612be5289e9",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e02c59d701e358b4612be5289e9",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d00004851c59d701e358b4612be5289e9",
				  "width": 64
				}
			  ],
			  "name": "breathe",
			  "release_date": "2019-09-13",
			  "release_date_precision": "day",
			  "total_tracks": 10,
			  "type": "album",
			  "uri": "spotify:album:1JMg02mQ5nmVKBWWoDUIeo"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				},
				"href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				"id": "5rJVTTK0ucAxQhkUc0nXbH",
				"name": "Tiny Moving Parts",
				"type": "artist",
				"uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
			  }
			],
			"disc_number": 1,
			"duration_ms": 190524,
			"explicit": false,
			"external_ids": {
			  "isrc": "USHR21912601"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/4eJFCgYa3D1gBgHcaeuh6Y"
			},
			"href": "https://api.spotify.com/v1/tracks/4eJFCgYa3D1gBgHcaeuh6Y",
			"id": "4eJFCgYa3D1gBgHcaeuh6Y",
			"is_local": false,
			"is_playable": true,
			"name": "The Midwest Sky",
			"popularity": 37,
			"preview_url": "https://p.scdn.co/mp3-preview/2fb05b31fcc15b860a78f198b8f273c5cf700146?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 1,
			"type": "track",
			"uri": "spotify:track:4eJFCgYa3D1gBgHcaeuh6Y"
		  },
		  {
			"album": {
			  "album_type": "album",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				  },
				  "href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				  "id": "5rJVTTK0ucAxQhkUc0nXbH",
				  "name": "Tiny Moving Parts",
				  "type": "artist",
				  "uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
				}
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/0hGe2EeB4e2PqTDwEaeqYn"
			  },
			  "href": "https://api.spotify.com/v1/albums/0hGe2EeB4e2PqTDwEaeqYn",
			  "id": "0hGe2EeB4e2PqTDwEaeqYn",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b2737b0c69ed623054a7d9a17e94",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e027b0c69ed623054a7d9a17e94",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d000048517b0c69ed623054a7d9a17e94",
				  "width": 64
				}
			  ],
			  "name": "Celebrate",
			  "release_date": "2016-05-20",
			  "release_date_precision": "day",
			  "total_tracks": 10,
			  "type": "album",
			  "uri": "spotify:album:0hGe2EeB4e2PqTDwEaeqYn"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				},
				"href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				"id": "5rJVTTK0ucAxQhkUc0nXbH",
				"name": "Tiny Moving Parts",
				"type": "artist",
				"uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
			  }
			],
			"disc_number": 1,
			"duration_ms": 262040,
			"explicit": false,
			"external_ids": {
			  "isrc": "US72W1611464"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/2DsELkXyf4eehUIeT4hID1"
			},
			"href": "https://api.spotify.com/v1/tracks/2DsELkXyf4eehUIeT4hID1",
			"id": "2DsELkXyf4eehUIeT4hID1",
			"is_local": false,
			"is_playable": true,
			"name": "Common Cold",
			"popularity": 37,
			"preview_url": "https://p.scdn.co/mp3-preview/c9d388dfcd1a0ca61ca8b7f87f565ac7cdc5072d?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 5,
			"type": "track",
			"uri": "spotify:track:2DsELkXyf4eehUIeT4hID1"
		  },
		  {
			"album": {
			  "album_type": "album",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				  },
				  "href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				  "id": "5rJVTTK0ucAxQhkUc0nXbH",
				  "name": "Tiny Moving Parts",
				  "type": "artist",
				  "uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
				}
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/1JMg02mQ5nmVKBWWoDUIeo"
			  },
			  "href": "https://api.spotify.com/v1/albums/1JMg02mQ5nmVKBWWoDUIeo",
			  "id": "1JMg02mQ5nmVKBWWoDUIeo",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b273c59d701e358b4612be5289e9",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e02c59d701e358b4612be5289e9",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d00004851c59d701e358b4612be5289e9",
				  "width": 64
				}
			  ],
			  "name": "breathe",
			  "release_date": "2019-09-13",
			  "release_date_precision": "day",
			  "total_tracks": 10,
			  "type": "album",
			  "uri": "spotify:album:1JMg02mQ5nmVKBWWoDUIeo"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				},
				"href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				"id": "5rJVTTK0ucAxQhkUc0nXbH",
				"name": "Tiny Moving Parts",
				"type": "artist",
				"uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
			  }
			],
			"disc_number": 1,
			"duration_ms": 168661,
			"explicit": false,
			"external_ids": {
			  "isrc": "USHR21912607"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/5eyCGJPbszc4xak3db5UKL"
			},
			"href": "https://api.spotify.com/v1/tracks/5eyCGJPbszc4xak3db5UKL",
			"id": "5eyCGJPbszc4xak3db5UKL",
			"is_local": false,
			"is_playable": true,
			"name": "Bloody Nose",
			"popularity": 36,
			"preview_url": "https://p.scdn.co/mp3-preview/ac65f18e798a7095aa4b971bd62f2c7d1f1a6ea8?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 7,
			"type": "track",
			"uri": "spotify:track:5eyCGJPbszc4xak3db5UKL"
		  }
		]
	  }`
	getTopTracksPayload = `{
		"items": [
		  {
			"album": {
			  "album_type": "ALBUM",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/1lKZzN2d4IqiEYxyECIEHI"
				  },
				  "href": "https://api.spotify.com/v1/artists/1lKZzN2d4IqiEYxyECIEHI",
				  "id": "1lKZzN2d4IqiEYxyECIEHI",
				  "name": "Hot Mulligan",
				  "type": "artist",
				  "uri": "spotify:artist:1lKZzN2d4IqiEYxyECIEHI"
				}
			  ],
			  "available_markets": [
				"AD",
				"AE",
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BH",
				"BO",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CY",
				"CZ",
				"DE",
				"DK",
				"DO",
				"DZ",
				"EC",
				"EE",
				"EG",
				"ES",
				"FI",
				"FR",
				"GB",
				"GR",
				"GT",
				"HK",
				"HN",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JO",
				"JP",
				"KW",
				"LB",
				"LI",
				"LT",
				"LU",
				"LV",
				"MA",
				"MC",
				"MT",
				"MX",
				"MY",
				"NI",
				"NL",
				"NO",
				"NZ",
				"OM",
				"PA",
				"PE",
				"PH",
				"PL",
				"PS",
				"PT",
				"PY",
				"QA",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"SV",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"UY",
				"VN",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/3wl3zdJVNhLyJfqdXaCRyp"
			  },
			  "href": "https://api.spotify.com/v1/albums/3wl3zdJVNhLyJfqdXaCRyp",
			  "id": "3wl3zdJVNhLyJfqdXaCRyp",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b2738cc9c6e183cef734184e15b7",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e028cc9c6e183cef734184e15b7",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d000048518cc9c6e183cef734184e15b7",
				  "width": 64
				}
			  ],
			  "name": "Pilot",
			  "release_date": "2018-03-23",
			  "release_date_precision": "day",
			  "total_tracks": 11,
			  "type": "album",
			  "uri": "spotify:album:3wl3zdJVNhLyJfqdXaCRyp"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/1lKZzN2d4IqiEYxyECIEHI"
				},
				"href": "https://api.spotify.com/v1/artists/1lKZzN2d4IqiEYxyECIEHI",
				"id": "1lKZzN2d4IqiEYxyECIEHI",
				"name": "Hot Mulligan",
				"type": "artist",
				"uri": "spotify:artist:1lKZzN2d4IqiEYxyECIEHI"
			  }
			],
			"available_markets": [
			  "AD",
			  "AE",
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BH",
			  "BO",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CY",
			  "CZ",
			  "DE",
			  "DK",
			  "DO",
			  "DZ",
			  "EC",
			  "EE",
			  "EG",
			  "ES",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "GT",
			  "HK",
			  "HN",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JO",
			  "JP",
			  "KW",
			  "LB",
			  "LI",
			  "LT",
			  "LU",
			  "LV",
			  "MA",
			  "MC",
			  "MT",
			  "MX",
			  "MY",
			  "NI",
			  "NL",
			  "NO",
			  "NZ",
			  "OM",
			  "PA",
			  "PE",
			  "PH",
			  "PL",
			  "PS",
			  "PT",
			  "PY",
			  "QA",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "SV",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "UY",
			  "VN",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 204708,
			"explicit": false,
			"external_ids": {
			  "isrc": "USZZ81810006"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/5DEF7cgddcZ1dUBk96J1Hx"
			},
			"href": "https://api.spotify.com/v1/tracks/5DEF7cgddcZ1dUBk96J1Hx",
			"id": "5DEF7cgddcZ1dUBk96J1Hx",
			"is_local": false,
			"name": "Deluxe Capacitor",
			"popularity": 38,
			"preview_url": "https://p.scdn.co/mp3-preview/f5fe4fe11933662f2109bed96a5c18182e3090c8?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 1,
			"type": "track",
			"uri": "spotify:track:5DEF7cgddcZ1dUBk96J1Hx"
		  },
		  {
			"album": {
			  "album_type": "ALBUM",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				  },
				  "href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				  "id": "5rJVTTK0ucAxQhkUc0nXbH",
				  "name": "Tiny Moving Parts",
				  "type": "artist",
				  "uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
				}
			  ],
			  "available_markets": [
				"AD",
				"AE",
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BH",
				"BO",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CY",
				"CZ",
				"DE",
				"DK",
				"DO",
				"DZ",
				"EC",
				"EE",
				"EG",
				"ES",
				"FI",
				"FR",
				"GB",
				"GR",
				"GT",
				"HK",
				"HN",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JO",
				"JP",
				"KW",
				"LB",
				"LI",
				"LT",
				"LU",
				"LV",
				"MA",
				"MC",
				"MT",
				"MX",
				"MY",
				"NI",
				"NL",
				"NO",
				"NZ",
				"OM",
				"PA",
				"PE",
				"PH",
				"PL",
				"PS",
				"PT",
				"PY",
				"QA",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"SV",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"UY",
				"VN",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/1JMg02mQ5nmVKBWWoDUIeo"
			  },
			  "href": "https://api.spotify.com/v1/albums/1JMg02mQ5nmVKBWWoDUIeo",
			  "id": "1JMg02mQ5nmVKBWWoDUIeo",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b273c59d701e358b4612be5289e9",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e02c59d701e358b4612be5289e9",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d00004851c59d701e358b4612be5289e9",
				  "width": 64
				}
			  ],
			  "name": "breathe",
			  "release_date": "2019-09-13",
			  "release_date_precision": "day",
			  "total_tracks": 10,
			  "type": "album",
			  "uri": "spotify:album:1JMg02mQ5nmVKBWWoDUIeo"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				},
				"href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				"id": "5rJVTTK0ucAxQhkUc0nXbH",
				"name": "Tiny Moving Parts",
				"type": "artist",
				"uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
			  }
			],
			"available_markets": [
			  "AD",
			  "AE",
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BH",
			  "BO",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CY",
			  "CZ",
			  "DE",
			  "DK",
			  "DO",
			  "DZ",
			  "EC",
			  "EE",
			  "EG",
			  "ES",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "GT",
			  "HK",
			  "HN",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JO",
			  "JP",
			  "KW",
			  "LB",
			  "LI",
			  "LT",
			  "LU",
			  "LV",
			  "MA",
			  "MC",
			  "MT",
			  "MX",
			  "MY",
			  "NI",
			  "NL",
			  "NO",
			  "NZ",
			  "OM",
			  "PA",
			  "PE",
			  "PH",
			  "PL",
			  "PS",
			  "PT",
			  "PY",
			  "QA",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "SV",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "UY",
			  "VN",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 168661,
			"explicit": false,
			"external_ids": {
			  "isrc": "USHR21912607"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/5eyCGJPbszc4xak3db5UKL"
			},
			"href": "https://api.spotify.com/v1/tracks/5eyCGJPbszc4xak3db5UKL",
			"id": "5eyCGJPbszc4xak3db5UKL",
			"is_local": false,
			"name": "Bloody Nose",
			"popularity": 37,
			"preview_url": "https://p.scdn.co/mp3-preview/ac65f18e798a7095aa4b971bd62f2c7d1f1a6ea8?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 7,
			"type": "track",
			"uri": "spotify:track:5eyCGJPbszc4xak3db5UKL"
		  },
		  {
			"album": {
			  "album_type": "ALBUM",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/2TM0qnbJH4QPhGMCdPt7fH"
				  },
				  "href": "https://api.spotify.com/v1/artists/2TM0qnbJH4QPhGMCdPt7fH",
				  "id": "2TM0qnbJH4QPhGMCdPt7fH",
				  "name": "Neck Deep",
				  "type": "artist",
				  "uri": "spotify:artist:2TM0qnbJH4QPhGMCdPt7fH"
				}
			  ],
			  "available_markets": [
				"AD",
				"AE",
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BH",
				"BO",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CY",
				"CZ",
				"DE",
				"DK",
				"DO",
				"DZ",
				"EC",
				"EE",
				"EG",
				"ES",
				"FI",
				"FR",
				"GB",
				"GR",
				"GT",
				"HK",
				"HN",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JO",
				"JP",
				"KW",
				"LB",
				"LI",
				"LT",
				"LU",
				"LV",
				"MA",
				"MC",
				"MT",
				"MX",
				"MY",
				"NI",
				"NL",
				"NO",
				"NZ",
				"OM",
				"PA",
				"PE",
				"PH",
				"PL",
				"PS",
				"PT",
				"PY",
				"QA",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"SV",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"UY",
				"VN",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/3umOBqXWR9VnJTQoe9Qkkj"
			  },
			  "href": "https://api.spotify.com/v1/albums/3umOBqXWR9VnJTQoe9Qkkj",
			  "id": "3umOBqXWR9VnJTQoe9Qkkj",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b27382a3560435840c14a72f6b6e",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e0282a3560435840c14a72f6b6e",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d0000485182a3560435840c14a72f6b6e",
				  "width": 64
				}
			  ],
			  "name": "Life's Not Out To Get You",
			  "release_date": "2015-08-14",
			  "release_date_precision": "day",
			  "total_tracks": 12,
			  "type": "album",
			  "uri": "spotify:album:3umOBqXWR9VnJTQoe9Qkkj"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/2TM0qnbJH4QPhGMCdPt7fH"
				},
				"href": "https://api.spotify.com/v1/artists/2TM0qnbJH4QPhGMCdPt7fH",
				"id": "2TM0qnbJH4QPhGMCdPt7fH",
				"name": "Neck Deep",
				"type": "artist",
				"uri": "spotify:artist:2TM0qnbJH4QPhGMCdPt7fH"
			  }
			],
			"available_markets": [
			  "AD",
			  "AE",
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BH",
			  "BO",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CY",
			  "CZ",
			  "DE",
			  "DK",
			  "DO",
			  "DZ",
			  "EC",
			  "EE",
			  "EG",
			  "ES",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "GT",
			  "HK",
			  "HN",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JO",
			  "JP",
			  "KW",
			  "LB",
			  "LI",
			  "LT",
			  "LU",
			  "LV",
			  "MA",
			  "MC",
			  "MT",
			  "MX",
			  "MY",
			  "NI",
			  "NL",
			  "NO",
			  "NZ",
			  "OM",
			  "PA",
			  "PE",
			  "PH",
			  "PL",
			  "PS",
			  "PT",
			  "PY",
			  "QA",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "SV",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "UY",
			  "VN",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 218954,
			"explicit": false,
			"external_ids": {
			  "isrc": "USHR21579009"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/4oVdhvxZrKQTM9ZsUIZa3S"
			},
			"href": "https://api.spotify.com/v1/tracks/4oVdhvxZrKQTM9ZsUIZa3S",
			"id": "4oVdhvxZrKQTM9ZsUIZa3S",
			"is_local": false,
			"name": "December",
			"popularity": 63,
			"preview_url": "https://p.scdn.co/mp3-preview/327c54e45f33eb9dd172866f6e9ac6a44ad04be0?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 9,
			"type": "track",
			"uri": "spotify:track:4oVdhvxZrKQTM9ZsUIZa3S"
		  },
		  {
			"album": {
			  "album_type": "ALBUM",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				  },
				  "href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				  "id": "5rJVTTK0ucAxQhkUc0nXbH",
				  "name": "Tiny Moving Parts",
				  "type": "artist",
				  "uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
				}
			  ],
			  "available_markets": [
				"AD",
				"AE",
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BH",
				"BO",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CY",
				"CZ",
				"DE",
				"DK",
				"DO",
				"DZ",
				"EC",
				"EE",
				"EG",
				"ES",
				"FI",
				"FR",
				"GB",
				"GR",
				"GT",
				"HK",
				"HN",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JO",
				"JP",
				"KW",
				"LB",
				"LI",
				"LT",
				"LU",
				"LV",
				"MA",
				"MC",
				"MT",
				"MX",
				"MY",
				"NI",
				"NL",
				"NO",
				"NZ",
				"OM",
				"PA",
				"PE",
				"PH",
				"PL",
				"PS",
				"PT",
				"PY",
				"QA",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"SV",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"UY",
				"VN",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/1JMg02mQ5nmVKBWWoDUIeo"
			  },
			  "href": "https://api.spotify.com/v1/albums/1JMg02mQ5nmVKBWWoDUIeo",
			  "id": "1JMg02mQ5nmVKBWWoDUIeo",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b273c59d701e358b4612be5289e9",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e02c59d701e358b4612be5289e9",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d00004851c59d701e358b4612be5289e9",
				  "width": 64
				}
			  ],
			  "name": "breathe",
			  "release_date": "2019-09-13",
			  "release_date_precision": "day",
			  "total_tracks": 10,
			  "type": "album",
			  "uri": "spotify:album:1JMg02mQ5nmVKBWWoDUIeo"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				},
				"href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				"id": "5rJVTTK0ucAxQhkUc0nXbH",
				"name": "Tiny Moving Parts",
				"type": "artist",
				"uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
			  }
			],
			"available_markets": [
			  "AD",
			  "AE",
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BH",
			  "BO",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CY",
			  "CZ",
			  "DE",
			  "DK",
			  "DO",
			  "DZ",
			  "EC",
			  "EE",
			  "EG",
			  "ES",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "GT",
			  "HK",
			  "HN",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JO",
			  "JP",
			  "KW",
			  "LB",
			  "LI",
			  "LT",
			  "LU",
			  "LV",
			  "MA",
			  "MC",
			  "MT",
			  "MX",
			  "MY",
			  "NI",
			  "NL",
			  "NO",
			  "NZ",
			  "OM",
			  "PA",
			  "PE",
			  "PH",
			  "PL",
			  "PS",
			  "PT",
			  "PY",
			  "QA",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "SV",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "UY",
			  "VN",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 191294,
			"explicit": false,
			"external_ids": {
			  "isrc": "USHR21912603"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/0uqTwEKTbkuhrn8yGSO6b5"
			},
			"href": "https://api.spotify.com/v1/tracks/0uqTwEKTbkuhrn8yGSO6b5",
			"id": "0uqTwEKTbkuhrn8yGSO6b5",
			"is_local": false,
			"name": "Medicine",
			"popularity": 52,
			"preview_url": "https://p.scdn.co/mp3-preview/6c66ead34e15f46119433775516d513b26c7b5f4?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 3,
			"type": "track",
			"uri": "spotify:track:0uqTwEKTbkuhrn8yGSO6b5"
		  },
		  {
			"album": {
			  "album_type": "ALBUM",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/1zNfnkHqbNqPMm0LY98Tfx"
				  },
				  "href": "https://api.spotify.com/v1/artists/1zNfnkHqbNqPMm0LY98Tfx",
				  "id": "1zNfnkHqbNqPMm0LY98Tfx",
				  "name": "fredo disco",
				  "type": "artist",
				  "uri": "spotify:artist:1zNfnkHqbNqPMm0LY98Tfx"
				}
			  ],
			  "available_markets": [
				"AD",
				"AE",
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BH",
				"BO",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CY",
				"CZ",
				"DE",
				"DK",
				"DO",
				"DZ",
				"EC",
				"EE",
				"EG",
				"ES",
				"FI",
				"FR",
				"GB",
				"GR",
				"GT",
				"HK",
				"HN",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JO",
				"JP",
				"KW",
				"LB",
				"LI",
				"LT",
				"LU",
				"LV",
				"MA",
				"MC",
				"MT",
				"MX",
				"MY",
				"NI",
				"NL",
				"NO",
				"NZ",
				"OM",
				"PA",
				"PE",
				"PH",
				"PL",
				"PS",
				"PT",
				"PY",
				"QA",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"SV",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"UY",
				"VN",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/6mfgImbbrrxp1e0HqVygbU"
			  },
			  "href": "https://api.spotify.com/v1/albums/6mfgImbbrrxp1e0HqVygbU",
			  "id": "6mfgImbbrrxp1e0HqVygbU",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b273bb3926b0b0a5c53a34296508",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e02bb3926b0b0a5c53a34296508",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d00004851bb3926b0b0a5c53a34296508",
				  "width": 64
				}
			  ],
			  "name": "school spirit",
			  "release_date": "2017-09-29",
			  "release_date_precision": "day",
			  "total_tracks": 7,
			  "type": "album",
			  "uri": "spotify:album:6mfgImbbrrxp1e0HqVygbU"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/1zNfnkHqbNqPMm0LY98Tfx"
				},
				"href": "https://api.spotify.com/v1/artists/1zNfnkHqbNqPMm0LY98Tfx",
				"id": "1zNfnkHqbNqPMm0LY98Tfx",
				"name": "fredo disco",
				"type": "artist",
				"uri": "spotify:artist:1zNfnkHqbNqPMm0LY98Tfx"
			  }
			],
			"available_markets": [
			  "AD",
			  "AE",
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BH",
			  "BO",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CY",
			  "CZ",
			  "DE",
			  "DK",
			  "DO",
			  "DZ",
			  "EC",
			  "EE",
			  "EG",
			  "ES",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "GT",
			  "HK",
			  "HN",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JO",
			  "JP",
			  "KW",
			  "LB",
			  "LI",
			  "LT",
			  "LU",
			  "LV",
			  "MA",
			  "MC",
			  "MT",
			  "MX",
			  "MY",
			  "NI",
			  "NL",
			  "NO",
			  "NZ",
			  "OM",
			  "PA",
			  "PE",
			  "PH",
			  "PL",
			  "PS",
			  "PT",
			  "PY",
			  "QA",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "SV",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "UY",
			  "VN",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 122506,
			"explicit": true,
			"external_ids": {
			  "isrc": "QZ5FN1738374"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/3EOdzmkePMPoC99wVzYsJS"
			},
			"href": "https://api.spotify.com/v1/tracks/3EOdzmkePMPoC99wVzYsJS",
			"id": "3EOdzmkePMPoC99wVzYsJS",
			"is_local": false,
			"name": "saturn suv",
			"popularity": 48,
			"preview_url": "https://p.scdn.co/mp3-preview/053812cd879adf8a6ade188f5bc6f4ad75fb9a67?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 2,
			"type": "track",
			"uri": "spotify:track:3EOdzmkePMPoC99wVzYsJS"
		  },
		  {
			"album": {
			  "album_type": "ALBUM",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/1lKZzN2d4IqiEYxyECIEHI"
				  },
				  "href": "https://api.spotify.com/v1/artists/1lKZzN2d4IqiEYxyECIEHI",
				  "id": "1lKZzN2d4IqiEYxyECIEHI",
				  "name": "Hot Mulligan",
				  "type": "artist",
				  "uri": "spotify:artist:1lKZzN2d4IqiEYxyECIEHI"
				}
			  ],
			  "available_markets": [
				"AD",
				"AE",
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BH",
				"BO",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CY",
				"CZ",
				"DE",
				"DK",
				"DO",
				"DZ",
				"EC",
				"EE",
				"EG",
				"ES",
				"FI",
				"FR",
				"GB",
				"GR",
				"GT",
				"HK",
				"HN",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JO",
				"JP",
				"KW",
				"LB",
				"LI",
				"LT",
				"LU",
				"LV",
				"MA",
				"MC",
				"MT",
				"MX",
				"MY",
				"NI",
				"NL",
				"NO",
				"NZ",
				"OM",
				"PA",
				"PE",
				"PH",
				"PL",
				"PS",
				"PT",
				"PY",
				"QA",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"SV",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"UY",
				"VN",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/3wl3zdJVNhLyJfqdXaCRyp"
			  },
			  "href": "https://api.spotify.com/v1/albums/3wl3zdJVNhLyJfqdXaCRyp",
			  "id": "3wl3zdJVNhLyJfqdXaCRyp",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b2738cc9c6e183cef734184e15b7",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e028cc9c6e183cef734184e15b7",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d000048518cc9c6e183cef734184e15b7",
				  "width": 64
				}
			  ],
			  "name": "Pilot",
			  "release_date": "2018-03-23",
			  "release_date_precision": "day",
			  "total_tracks": 11,
			  "type": "album",
			  "uri": "spotify:album:3wl3zdJVNhLyJfqdXaCRyp"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/1lKZzN2d4IqiEYxyECIEHI"
				},
				"href": "https://api.spotify.com/v1/artists/1lKZzN2d4IqiEYxyECIEHI",
				"id": "1lKZzN2d4IqiEYxyECIEHI",
				"name": "Hot Mulligan",
				"type": "artist",
				"uri": "spotify:artist:1lKZzN2d4IqiEYxyECIEHI"
			  }
			],
			"available_markets": [
			  "AD",
			  "AE",
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BH",
			  "BO",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CY",
			  "CZ",
			  "DE",
			  "DK",
			  "DO",
			  "DZ",
			  "EC",
			  "EE",
			  "EG",
			  "ES",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "GT",
			  "HK",
			  "HN",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JO",
			  "JP",
			  "KW",
			  "LB",
			  "LI",
			  "LT",
			  "LU",
			  "LV",
			  "MA",
			  "MC",
			  "MT",
			  "MX",
			  "MY",
			  "NI",
			  "NL",
			  "NO",
			  "NZ",
			  "OM",
			  "PA",
			  "PE",
			  "PH",
			  "PL",
			  "PS",
			  "PT",
			  "PY",
			  "QA",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "SV",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "UY",
			  "VN",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 205299,
			"explicit": false,
			"external_ids": {
			  "isrc": "USZZ81810007"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/2GavHHPGrI59bNL5ouyUWk"
			},
			"href": "https://api.spotify.com/v1/tracks/2GavHHPGrI59bNL5ouyUWk",
			"id": "2GavHHPGrI59bNL5ouyUWk",
			"is_local": false,
			"name": "All You Wanted By Michelle Branch",
			"popularity": 38,
			"preview_url": "https://p.scdn.co/mp3-preview/2d2a8a41ab2618694774cffd57e03a3c47d4d201?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 2,
			"type": "track",
			"uri": "spotify:track:2GavHHPGrI59bNL5ouyUWk"
		  },
		  {
			"album": {
			  "album_type": "ALBUM",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/6FBDaR13swtiWwGhX1WQsP"
				  },
				  "href": "https://api.spotify.com/v1/artists/6FBDaR13swtiWwGhX1WQsP",
				  "id": "6FBDaR13swtiWwGhX1WQsP",
				  "name": "blink-182",
				  "type": "artist",
				  "uri": "spotify:artist:6FBDaR13swtiWwGhX1WQsP"
				}
			  ],
			  "available_markets": [
				"AD",
				"AE",
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BH",
				"BO",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CY",
				"CZ",
				"DE",
				"DK",
				"DO",
				"DZ",
				"EC",
				"EE",
				"EG",
				"ES",
				"FI",
				"FR",
				"GB",
				"GR",
				"GT",
				"HK",
				"HN",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JO",
				"JP",
				"KW",
				"LB",
				"LI",
				"LT",
				"LU",
				"LV",
				"MA",
				"MC",
				"MT",
				"MX",
				"MY",
				"NI",
				"NL",
				"NO",
				"NZ",
				"OM",
				"PA",
				"PE",
				"PH",
				"PL",
				"PS",
				"PT",
				"PY",
				"QA",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"SV",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"UY",
				"VN",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/4gARZz9eV7zbGbtOjhVTPF"
			  },
			  "href": "https://api.spotify.com/v1/albums/4gARZz9eV7zbGbtOjhVTPF",
			  "id": "4gARZz9eV7zbGbtOjhVTPF",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b273198542728b101474c4afe0a1",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e02198542728b101474c4afe0a1",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d00004851198542728b101474c4afe0a1",
				  "width": 64
				}
			  ],
			  "name": "NINE",
			  "release_date": "2019-09-20",
			  "release_date_precision": "day",
			  "total_tracks": 15,
			  "type": "album",
			  "uri": "spotify:album:4gARZz9eV7zbGbtOjhVTPF"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/6FBDaR13swtiWwGhX1WQsP"
				},
				"href": "https://api.spotify.com/v1/artists/6FBDaR13swtiWwGhX1WQsP",
				"id": "6FBDaR13swtiWwGhX1WQsP",
				"name": "blink-182",
				"type": "artist",
				"uri": "spotify:artist:6FBDaR13swtiWwGhX1WQsP"
			  }
			],
			"available_markets": [
			  "AD",
			  "AE",
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BH",
			  "BO",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CY",
			  "CZ",
			  "DE",
			  "DK",
			  "DO",
			  "DZ",
			  "EC",
			  "EE",
			  "EG",
			  "ES",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "GT",
			  "HK",
			  "HN",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JO",
			  "JP",
			  "KW",
			  "LB",
			  "LI",
			  "LT",
			  "LU",
			  "LV",
			  "MA",
			  "MC",
			  "MT",
			  "MX",
			  "MY",
			  "NI",
			  "NL",
			  "NO",
			  "NZ",
			  "OM",
			  "PA",
			  "PE",
			  "PH",
			  "PL",
			  "PS",
			  "PT",
			  "PY",
			  "QA",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "SV",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "UY",
			  "VN",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 180960,
			"explicit": false,
			"external_ids": {
			  "isrc": "USSM11904207"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/39n1y2JqVOPIdfUeLqzgfl"
			},
			"href": "https://api.spotify.com/v1/tracks/39n1y2JqVOPIdfUeLqzgfl",
			"id": "39n1y2JqVOPIdfUeLqzgfl",
			"is_local": false,
			"name": "Darkside",
			"popularity": 63,
			"preview_url": "https://p.scdn.co/mp3-preview/1e5f5e74dbf0595887a3258ff58bf3dcc98e171f?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 4,
			"type": "track",
			"uri": "spotify:track:39n1y2JqVOPIdfUeLqzgfl"
		  },
		  {
			"album": {
			  "album_type": "SINGLE",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/0w7HLMvZOHatWVbAKee1zF"
				  },
				  "href": "https://api.spotify.com/v1/artists/0w7HLMvZOHatWVbAKee1zF",
				  "id": "0w7HLMvZOHatWVbAKee1zF",
				  "name": "Redbone",
				  "type": "artist",
				  "uri": "spotify:artist:0w7HLMvZOHatWVbAKee1zF"
				}
			  ],
			  "available_markets": [
				"AD",
				"AE",
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BH",
				"BO",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CY",
				"CZ",
				"DE",
				"DK",
				"DO",
				"DZ",
				"EC",
				"EE",
				"EG",
				"ES",
				"FI",
				"FR",
				"GB",
				"GR",
				"GT",
				"HK",
				"HN",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JO",
				"JP",
				"KW",
				"LB",
				"LI",
				"LT",
				"LU",
				"LV",
				"MA",
				"MC",
				"MT",
				"MX",
				"MY",
				"NI",
				"NL",
				"NO",
				"NZ",
				"OM",
				"PA",
				"PE",
				"PH",
				"PL",
				"PS",
				"PT",
				"PY",
				"QA",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"SV",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"UY",
				"VN",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/5Gf5m9M6RiK2lkjpbP0xRu"
			  },
			  "href": "https://api.spotify.com/v1/albums/5Gf5m9M6RiK2lkjpbP0xRu",
			  "id": "5Gf5m9M6RiK2lkjpbP0xRu",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b27346814e1b44e54d806753801e",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e0246814e1b44e54d806753801e",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d0000485146814e1b44e54d806753801e",
				  "width": 64
				}
			  ],
			  "name": "Come and Get Your Love",
			  "release_date": "1973-11-01",
			  "release_date_precision": "day",
			  "total_tracks": 1,
			  "type": "album",
			  "uri": "spotify:album:5Gf5m9M6RiK2lkjpbP0xRu"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/0w7HLMvZOHatWVbAKee1zF"
				},
				"href": "https://api.spotify.com/v1/artists/0w7HLMvZOHatWVbAKee1zF",
				"id": "0w7HLMvZOHatWVbAKee1zF",
				"name": "Redbone",
				"type": "artist",
				"uri": "spotify:artist:0w7HLMvZOHatWVbAKee1zF"
			  }
			],
			"available_markets": [
			  "AD",
			  "AE",
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BH",
			  "BO",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CY",
			  "CZ",
			  "DE",
			  "DK",
			  "DO",
			  "DZ",
			  "EC",
			  "EE",
			  "EG",
			  "ES",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "GT",
			  "HK",
			  "HN",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JO",
			  "JP",
			  "KW",
			  "LB",
			  "LI",
			  "LT",
			  "LU",
			  "LV",
			  "MA",
			  "MC",
			  "MT",
			  "MX",
			  "MY",
			  "NI",
			  "NL",
			  "NO",
			  "NZ",
			  "OM",
			  "PA",
			  "PE",
			  "PH",
			  "PL",
			  "PS",
			  "PT",
			  "PY",
			  "QA",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "SV",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "UY",
			  "VN",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 205933,
			"explicit": false,
			"external_ids": {
			  "isrc": "USSM17300653"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/7GVUmCP00eSsqc4tzj1sDD"
			},
			"href": "https://api.spotify.com/v1/tracks/7GVUmCP00eSsqc4tzj1sDD",
			"id": "7GVUmCP00eSsqc4tzj1sDD",
			"is_local": false,
			"name": "Come and Get Your Love - Single Version",
			"popularity": 78,
			"preview_url": "https://p.scdn.co/mp3-preview/0ad92a28a640b75565af67919ebc78e86ccae7f0?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 1,
			"type": "track",
			"uri": "spotify:track:7GVUmCP00eSsqc4tzj1sDD"
		  },
		  {
			"album": {
			  "album_type": "ALBUM",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/0DK7FqcaL3ks9TfFn9y1sD"
				  },
				  "href": "https://api.spotify.com/v1/artists/0DK7FqcaL3ks9TfFn9y1sD",
				  "id": "0DK7FqcaL3ks9TfFn9y1sD",
				  "name": "Box Car Racer",
				  "type": "artist",
				  "uri": "spotify:artist:0DK7FqcaL3ks9TfFn9y1sD"
				}
			  ],
			  "available_markets": [
				"AD",
				"AE",
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BH",
				"BO",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CY",
				"CZ",
				"DE",
				"DK",
				"DO",
				"DZ",
				"EC",
				"EE",
				"EG",
				"ES",
				"FI",
				"FR",
				"GB",
				"GR",
				"GT",
				"HK",
				"HN",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JO",
				"JP",
				"KW",
				"LB",
				"LI",
				"LT",
				"LU",
				"LV",
				"MA",
				"MC",
				"MT",
				"MX",
				"MY",
				"NI",
				"NL",
				"NO",
				"NZ",
				"OM",
				"PA",
				"PE",
				"PH",
				"PL",
				"PT",
				"PY",
				"QA",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"SV",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"UY",
				"VN",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/3gODo8aZ2dTVIaOr9SqeRE"
			  },
			  "href": "https://api.spotify.com/v1/albums/3gODo8aZ2dTVIaOr9SqeRE",
			  "id": "3gODo8aZ2dTVIaOr9SqeRE",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b2732682e92c4fc9b95df0319c72",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e022682e92c4fc9b95df0319c72",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d000048512682e92c4fc9b95df0319c72",
				  "width": 64
				}
			  ],
			  "name": "Box Car Racer",
			  "release_date": "2002-01-01",
			  "release_date_precision": "day",
			  "total_tracks": 13,
			  "type": "album",
			  "uri": "spotify:album:3gODo8aZ2dTVIaOr9SqeRE"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/0DK7FqcaL3ks9TfFn9y1sD"
				},
				"href": "https://api.spotify.com/v1/artists/0DK7FqcaL3ks9TfFn9y1sD",
				"id": "0DK7FqcaL3ks9TfFn9y1sD",
				"name": "Box Car Racer",
				"type": "artist",
				"uri": "spotify:artist:0DK7FqcaL3ks9TfFn9y1sD"
			  }
			],
			"available_markets": [
			  "AD",
			  "AE",
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BH",
			  "BO",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CY",
			  "CZ",
			  "DE",
			  "DK",
			  "DO",
			  "DZ",
			  "EC",
			  "EE",
			  "EG",
			  "ES",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "GT",
			  "HK",
			  "HN",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JO",
			  "JP",
			  "KW",
			  "LB",
			  "LI",
			  "LT",
			  "LU",
			  "LV",
			  "MA",
			  "MC",
			  "MT",
			  "MX",
			  "MY",
			  "NI",
			  "NL",
			  "NO",
			  "NZ",
			  "OM",
			  "PA",
			  "PE",
			  "PH",
			  "PL",
			  "PT",
			  "PY",
			  "QA",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "SV",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "UY",
			  "VN",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 196093,
			"explicit": true,
			"external_ids": {
			  "isrc": "USMC10200460"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/1pU1mucfoUalVw9apwnDhh"
			},
			"href": "https://api.spotify.com/v1/tracks/1pU1mucfoUalVw9apwnDhh",
			"id": "1pU1mucfoUalVw9apwnDhh",
			"is_local": false,
			"name": "Letters To God",
			"popularity": 49,
			"preview_url": "https://p.scdn.co/mp3-preview/3861db4923cc2cad83b59ec8288d502844c53ec7?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 7,
			"type": "track",
			"uri": "spotify:track:1pU1mucfoUalVw9apwnDhh"
		  },
		  {
			"album": {
			  "album_type": "ALBUM",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/6FBDaR13swtiWwGhX1WQsP"
				  },
				  "href": "https://api.spotify.com/v1/artists/6FBDaR13swtiWwGhX1WQsP",
				  "id": "6FBDaR13swtiWwGhX1WQsP",
				  "name": "blink-182",
				  "type": "artist",
				  "uri": "spotify:artist:6FBDaR13swtiWwGhX1WQsP"
				}
			  ],
			  "available_markets": [
				"AD",
				"AE",
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BH",
				"BO",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CY",
				"CZ",
				"DE",
				"DK",
				"DO",
				"DZ",
				"EC",
				"EE",
				"EG",
				"ES",
				"FI",
				"FR",
				"GB",
				"GR",
				"GT",
				"HK",
				"HN",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JO",
				"JP",
				"KW",
				"LB",
				"LI",
				"LT",
				"LU",
				"LV",
				"MA",
				"MC",
				"MT",
				"MX",
				"MY",
				"NI",
				"NL",
				"NO",
				"NZ",
				"OM",
				"PA",
				"PE",
				"PH",
				"PL",
				"PS",
				"PT",
				"PY",
				"QA",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"SV",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"UY",
				"VN",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/3nHpBmW5wJXGeC3ojBkpey"
			  },
			  "href": "https://api.spotify.com/v1/albums/3nHpBmW5wJXGeC3ojBkpey",
			  "id": "3nHpBmW5wJXGeC3ojBkpey",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b27354a8f4f9158546472fbb7280",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e0254a8f4f9158546472fbb7280",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d0000485154a8f4f9158546472fbb7280",
				  "width": 64
				}
			  ],
			  "name": "Take Off Your Pants And Jacket",
			  "release_date": "2001",
			  "release_date_precision": "year",
			  "total_tracks": 13,
			  "type": "album",
			  "uri": "spotify:album:3nHpBmW5wJXGeC3ojBkpey"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/6FBDaR13swtiWwGhX1WQsP"
				},
				"href": "https://api.spotify.com/v1/artists/6FBDaR13swtiWwGhX1WQsP",
				"id": "6FBDaR13swtiWwGhX1WQsP",
				"name": "blink-182",
				"type": "artist",
				"uri": "spotify:artist:6FBDaR13swtiWwGhX1WQsP"
			  }
			],
			"available_markets": [
			  "AD",
			  "AE",
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BH",
			  "BO",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CY",
			  "CZ",
			  "DE",
			  "DK",
			  "DO",
			  "DZ",
			  "EC",
			  "EE",
			  "EG",
			  "ES",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "GT",
			  "HK",
			  "HN",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JO",
			  "JP",
			  "KW",
			  "LB",
			  "LI",
			  "LT",
			  "LU",
			  "LV",
			  "MA",
			  "MC",
			  "MT",
			  "MX",
			  "MY",
			  "NI",
			  "NL",
			  "NO",
			  "NZ",
			  "OM",
			  "PA",
			  "PE",
			  "PH",
			  "PL",
			  "PS",
			  "PT",
			  "PY",
			  "QA",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "SV",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "UY",
			  "VN",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 171533,
			"explicit": false,
			"external_ids": {
			  "isrc": "USMC10110874"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/1fJFuvU2ldmeAm5nFIHcPP"
			},
			"href": "https://api.spotify.com/v1/tracks/1fJFuvU2ldmeAm5nFIHcPP",
			"id": "1fJFuvU2ldmeAm5nFIHcPP",
			"is_local": false,
			"name": "First Date",
			"popularity": 74,
			"preview_url": "https://p.scdn.co/mp3-preview/093414a6b85391d5b7b0929a775d8989d6d4067e?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 3,
			"type": "track",
			"uri": "spotify:track:1fJFuvU2ldmeAm5nFIHcPP"
		  },
		  {
			"album": {
			  "album_type": "ALBUM",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				  },
				  "href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				  "id": "5rJVTTK0ucAxQhkUc0nXbH",
				  "name": "Tiny Moving Parts",
				  "type": "artist",
				  "uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
				}
			  ],
			  "available_markets": [
				"AD",
				"AE",
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BH",
				"BO",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CY",
				"CZ",
				"DE",
				"DK",
				"DO",
				"DZ",
				"EC",
				"EE",
				"EG",
				"ES",
				"FI",
				"FR",
				"GB",
				"GR",
				"GT",
				"HK",
				"HN",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JO",
				"JP",
				"KW",
				"LB",
				"LI",
				"LT",
				"LU",
				"LV",
				"MA",
				"MC",
				"MT",
				"MX",
				"MY",
				"NI",
				"NL",
				"NO",
				"NZ",
				"OM",
				"PA",
				"PE",
				"PH",
				"PL",
				"PS",
				"PT",
				"PY",
				"QA",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"SV",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"UY",
				"VN",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/1JMg02mQ5nmVKBWWoDUIeo"
			  },
			  "href": "https://api.spotify.com/v1/albums/1JMg02mQ5nmVKBWWoDUIeo",
			  "id": "1JMg02mQ5nmVKBWWoDUIeo",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b273c59d701e358b4612be5289e9",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e02c59d701e358b4612be5289e9",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d00004851c59d701e358b4612be5289e9",
				  "width": 64
				}
			  ],
			  "name": "breathe",
			  "release_date": "2019-09-13",
			  "release_date_precision": "day",
			  "total_tracks": 10,
			  "type": "album",
			  "uri": "spotify:album:1JMg02mQ5nmVKBWWoDUIeo"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/5rJVTTK0ucAxQhkUc0nXbH"
				},
				"href": "https://api.spotify.com/v1/artists/5rJVTTK0ucAxQhkUc0nXbH",
				"id": "5rJVTTK0ucAxQhkUc0nXbH",
				"name": "Tiny Moving Parts",
				"type": "artist",
				"uri": "spotify:artist:5rJVTTK0ucAxQhkUc0nXbH"
			  }
			],
			"available_markets": [
			  "AD",
			  "AE",
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BH",
			  "BO",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CY",
			  "CZ",
			  "DE",
			  "DK",
			  "DO",
			  "DZ",
			  "EC",
			  "EE",
			  "EG",
			  "ES",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "GT",
			  "HK",
			  "HN",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JO",
			  "JP",
			  "KW",
			  "LB",
			  "LI",
			  "LT",
			  "LU",
			  "LV",
			  "MA",
			  "MC",
			  "MT",
			  "MX",
			  "MY",
			  "NI",
			  "NL",
			  "NO",
			  "NZ",
			  "OM",
			  "PA",
			  "PE",
			  "PH",
			  "PL",
			  "PS",
			  "PT",
			  "PY",
			  "QA",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "SV",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "UY",
			  "VN",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 199875,
			"explicit": false,
			"external_ids": {
			  "isrc": "USHR21912610"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/6zvODIwiEpwj7njxUWqj6A"
			},
			"href": "https://api.spotify.com/v1/tracks/6zvODIwiEpwj7njxUWqj6A",
			"id": "6zvODIwiEpwj7njxUWqj6A",
			"is_local": false,
			"name": "Hallmark",
			"popularity": 28,
			"preview_url": "https://p.scdn.co/mp3-preview/a16745b71c462dd5dcd69382acebf5d074a56f09?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 10,
			"type": "track",
			"uri": "spotify:track:6zvODIwiEpwj7njxUWqj6A"
		  },
		  {
			"album": {
			  "album_type": "ALBUM",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/6FBDaR13swtiWwGhX1WQsP"
				  },
				  "href": "https://api.spotify.com/v1/artists/6FBDaR13swtiWwGhX1WQsP",
				  "id": "6FBDaR13swtiWwGhX1WQsP",
				  "name": "blink-182",
				  "type": "artist",
				  "uri": "spotify:artist:6FBDaR13swtiWwGhX1WQsP"
				}
			  ],
			  "available_markets": [
				"CA",
				"JP",
				"MX",
				"US"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/2xSmzzqWM6cqaC92Hf7Dyv"
			  },
			  "href": "https://api.spotify.com/v1/albums/2xSmzzqWM6cqaC92Hf7Dyv",
			  "id": "2xSmzzqWM6cqaC92Hf7Dyv",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b2739bce7409f1fd24101e611603",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e029bce7409f1fd24101e611603",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d000048519bce7409f1fd24101e611603",
				  "width": 64
				}
			  ],
			  "name": "Neighborhoods (Deluxe)",
			  "release_date": "2011-01-01",
			  "release_date_precision": "day",
			  "total_tracks": 14,
			  "type": "album",
			  "uri": "spotify:album:2xSmzzqWM6cqaC92Hf7Dyv"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/6FBDaR13swtiWwGhX1WQsP"
				},
				"href": "https://api.spotify.com/v1/artists/6FBDaR13swtiWwGhX1WQsP",
				"id": "6FBDaR13swtiWwGhX1WQsP",
				"name": "blink-182",
				"type": "artist",
				"uri": "spotify:artist:6FBDaR13swtiWwGhX1WQsP"
			  }
			],
			"available_markets": [
			  "CA",
			  "JP",
			  "MX",
			  "US"
			],
			"disc_number": 1,
			"duration_ms": 200360,
			"explicit": false,
			"external_ids": {
			  "isrc": "USUM71113884"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/0y3fPWWin7NL5aED2FjvXP"
			},
			"href": "https://api.spotify.com/v1/tracks/0y3fPWWin7NL5aED2FjvXP",
			"id": "0y3fPWWin7NL5aED2FjvXP",
			"is_local": false,
			"name": "Wishing Well",
			"popularity": 49,
			"preview_url": "https://p.scdn.co/mp3-preview/834a8537c260ac6467b0baa0c89596fa6f288bbb?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 8,
			"type": "track",
			"uri": "spotify:track:0y3fPWWin7NL5aED2FjvXP"
		  },
		  {
			"album": {
			  "album_type": "SINGLE",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/1ajMGdOWCVoMNo62h042Ud"
				  },
				  "href": "https://api.spotify.com/v1/artists/1ajMGdOWCVoMNo62h042Ud",
				  "id": "1ajMGdOWCVoMNo62h042Ud",
				  "name": "Everyone Leaves",
				  "type": "artist",
				  "uri": "spotify:artist:1ajMGdOWCVoMNo62h042Ud"
				},
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/1lKZzN2d4IqiEYxyECIEHI"
				  },
				  "href": "https://api.spotify.com/v1/artists/1lKZzN2d4IqiEYxyECIEHI",
				  "id": "1lKZzN2d4IqiEYxyECIEHI",
				  "name": "Hot Mulligan",
				  "type": "artist",
				  "uri": "spotify:artist:1lKZzN2d4IqiEYxyECIEHI"
				}
			  ],
			  "available_markets": [
				"AD",
				"AE",
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BH",
				"BO",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CY",
				"CZ",
				"DE",
				"DK",
				"DO",
				"DZ",
				"EC",
				"EE",
				"EG",
				"ES",
				"FI",
				"FR",
				"GB",
				"GR",
				"GT",
				"HK",
				"HN",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JO",
				"JP",
				"KW",
				"LB",
				"LI",
				"LT",
				"LU",
				"LV",
				"MA",
				"MC",
				"MT",
				"MX",
				"MY",
				"NI",
				"NL",
				"NO",
				"NZ",
				"OM",
				"PA",
				"PE",
				"PH",
				"PL",
				"PS",
				"PT",
				"PY",
				"QA",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"SV",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"UY",
				"VN",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/6Tg5fae6DW9jHJ30Nhe6PY"
			  },
			  "href": "https://api.spotify.com/v1/albums/6Tg5fae6DW9jHJ30Nhe6PY",
			  "id": "6Tg5fae6DW9jHJ30Nhe6PY",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b273aa445f9f9e2a78cbacfede0a",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e02aa445f9f9e2a78cbacfede0a",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d00004851aa445f9f9e2a78cbacfede0a",
				  "width": 64
				}
			  ],
			  "name": "Split",
			  "release_date": "2016-03-11",
			  "release_date_precision": "day",
			  "total_tracks": 4,
			  "type": "album",
			  "uri": "spotify:album:6Tg5fae6DW9jHJ30Nhe6PY"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/1lKZzN2d4IqiEYxyECIEHI"
				},
				"href": "https://api.spotify.com/v1/artists/1lKZzN2d4IqiEYxyECIEHI",
				"id": "1lKZzN2d4IqiEYxyECIEHI",
				"name": "Hot Mulligan",
				"type": "artist",
				"uri": "spotify:artist:1lKZzN2d4IqiEYxyECIEHI"
			  }
			],
			"available_markets": [
			  "AD",
			  "AE",
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BH",
			  "BO",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CY",
			  "CZ",
			  "DE",
			  "DK",
			  "DO",
			  "DZ",
			  "EC",
			  "EE",
			  "EG",
			  "ES",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "GT",
			  "HK",
			  "HN",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JO",
			  "JP",
			  "KW",
			  "LB",
			  "LI",
			  "LT",
			  "LU",
			  "LV",
			  "MA",
			  "MC",
			  "MT",
			  "MX",
			  "MY",
			  "NI",
			  "NL",
			  "NO",
			  "NZ",
			  "OM",
			  "PA",
			  "PE",
			  "PH",
			  "PL",
			  "PS",
			  "PT",
			  "PY",
			  "QA",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "SV",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "UY",
			  "VN",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 84347,
			"explicit": false,
			"external_ids": {
			  "isrc": "FR26V1609139"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/47x3JRx7R5YNUBH9TD3OhA"
			},
			"href": "https://api.spotify.com/v1/tracks/47x3JRx7R5YNUBH9TD3OhA",
			"id": "47x3JRx7R5YNUBH9TD3OhA",
			"is_local": false,
			"name": "Legen",
			"popularity": 28,
			"preview_url": "https://p.scdn.co/mp3-preview/52e8d85c52dc5e445197b8b20cc41f5348d7684e?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 4,
			"type": "track",
			"uri": "spotify:track:47x3JRx7R5YNUBH9TD3OhA"
		  },
		  {
			"album": {
			  "album_type": "ALBUM",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/6FBDaR13swtiWwGhX1WQsP"
				  },
				  "href": "https://api.spotify.com/v1/artists/6FBDaR13swtiWwGhX1WQsP",
				  "id": "6FBDaR13swtiWwGhX1WQsP",
				  "name": "blink-182",
				  "type": "artist",
				  "uri": "spotify:artist:6FBDaR13swtiWwGhX1WQsP"
				}
			  ],
			  "available_markets": [
				"AD",
				"AE",
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BH",
				"BO",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CY",
				"CZ",
				"DE",
				"DK",
				"DO",
				"DZ",
				"EC",
				"EE",
				"EG",
				"ES",
				"FI",
				"FR",
				"GB",
				"GR",
				"GT",
				"HK",
				"HN",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JO",
				"JP",
				"KW",
				"LB",
				"LI",
				"LT",
				"LU",
				"LV",
				"MA",
				"MC",
				"MT",
				"MX",
				"MY",
				"NI",
				"NL",
				"NO",
				"NZ",
				"OM",
				"PA",
				"PE",
				"PH",
				"PL",
				"PS",
				"PT",
				"PY",
				"QA",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"SV",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"UY",
				"VN",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/4gARZz9eV7zbGbtOjhVTPF"
			  },
			  "href": "https://api.spotify.com/v1/albums/4gARZz9eV7zbGbtOjhVTPF",
			  "id": "4gARZz9eV7zbGbtOjhVTPF",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b273198542728b101474c4afe0a1",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e02198542728b101474c4afe0a1",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d00004851198542728b101474c4afe0a1",
				  "width": 64
				}
			  ],
			  "name": "NINE",
			  "release_date": "2019-09-20",
			  "release_date_precision": "day",
			  "total_tracks": 15,
			  "type": "album",
			  "uri": "spotify:album:4gARZz9eV7zbGbtOjhVTPF"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/6FBDaR13swtiWwGhX1WQsP"
				},
				"href": "https://api.spotify.com/v1/artists/6FBDaR13swtiWwGhX1WQsP",
				"id": "6FBDaR13swtiWwGhX1WQsP",
				"name": "blink-182",
				"type": "artist",
				"uri": "spotify:artist:6FBDaR13swtiWwGhX1WQsP"
			  }
			],
			"available_markets": [
			  "AD",
			  "AE",
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BH",
			  "BO",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CY",
			  "CZ",
			  "DE",
			  "DK",
			  "DO",
			  "DZ",
			  "EC",
			  "EE",
			  "EG",
			  "ES",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "GT",
			  "HK",
			  "HN",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JO",
			  "JP",
			  "KW",
			  "LB",
			  "LI",
			  "LT",
			  "LU",
			  "LV",
			  "MA",
			  "MC",
			  "MT",
			  "MX",
			  "MY",
			  "NI",
			  "NL",
			  "NO",
			  "NZ",
			  "OM",
			  "PA",
			  "PE",
			  "PH",
			  "PL",
			  "PS",
			  "PT",
			  "PY",
			  "QA",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "SV",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "UY",
			  "VN",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 179000,
			"explicit": false,
			"external_ids": {
			  "isrc": "USSM11903821"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/6GHAAp25FcGdaHgAn8rTif"
			},
			"href": "https://api.spotify.com/v1/tracks/6GHAAp25FcGdaHgAn8rTif",
			"id": "6GHAAp25FcGdaHgAn8rTif",
			"is_local": false,
			"name": "Happy Days",
			"popularity": 57,
			"preview_url": "https://p.scdn.co/mp3-preview/5cd5e5eb9578b2c3d21dade8274bc1a8b1641254?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 2,
			"type": "track",
			"uri": "spotify:track:6GHAAp25FcGdaHgAn8rTif"
		  },
		  {
			"album": {
			  "album_type": "ALBUM",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/1lKZzN2d4IqiEYxyECIEHI"
				  },
				  "href": "https://api.spotify.com/v1/artists/1lKZzN2d4IqiEYxyECIEHI",
				  "id": "1lKZzN2d4IqiEYxyECIEHI",
				  "name": "Hot Mulligan",
				  "type": "artist",
				  "uri": "spotify:artist:1lKZzN2d4IqiEYxyECIEHI"
				}
			  ],
			  "available_markets": [
				"AD",
				"AE",
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BH",
				"BO",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CY",
				"CZ",
				"DE",
				"DK",
				"DO",
				"DZ",
				"EC",
				"EE",
				"EG",
				"ES",
				"FI",
				"FR",
				"GB",
				"GR",
				"GT",
				"HK",
				"HN",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JO",
				"JP",
				"KW",
				"LB",
				"LI",
				"LT",
				"LU",
				"LV",
				"MA",
				"MC",
				"MT",
				"MX",
				"MY",
				"NI",
				"NL",
				"NO",
				"NZ",
				"OM",
				"PA",
				"PE",
				"PH",
				"PL",
				"PS",
				"PT",
				"PY",
				"QA",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"SV",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"UY",
				"VN",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/6uPdMxyaE7aoODPLIIhyvv"
			  },
			  "href": "https://api.spotify.com/v1/albums/6uPdMxyaE7aoODPLIIhyvv",
			  "id": "6uPdMxyaE7aoODPLIIhyvv",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b2731b4435e21c995aa1273a92b3",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e021b4435e21c995aa1273a92b3",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d000048511b4435e21c995aa1273a92b3",
				  "width": 64
				}
			  ],
			  "name": "Opportunities",
			  "release_date": "2017-03-31",
			  "release_date_precision": "day",
			  "total_tracks": 7,
			  "type": "album",
			  "uri": "spotify:album:6uPdMxyaE7aoODPLIIhyvv"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/1lKZzN2d4IqiEYxyECIEHI"
				},
				"href": "https://api.spotify.com/v1/artists/1lKZzN2d4IqiEYxyECIEHI",
				"id": "1lKZzN2d4IqiEYxyECIEHI",
				"name": "Hot Mulligan",
				"type": "artist",
				"uri": "spotify:artist:1lKZzN2d4IqiEYxyECIEHI"
			  }
			],
			"available_markets": [
			  "AD",
			  "AE",
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BH",
			  "BO",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CY",
			  "CZ",
			  "DE",
			  "DK",
			  "DO",
			  "DZ",
			  "EC",
			  "EE",
			  "EG",
			  "ES",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "GT",
			  "HK",
			  "HN",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JO",
			  "JP",
			  "KW",
			  "LB",
			  "LI",
			  "LT",
			  "LU",
			  "LV",
			  "MA",
			  "MC",
			  "MT",
			  "MX",
			  "MY",
			  "NI",
			  "NL",
			  "NO",
			  "NZ",
			  "OM",
			  "PA",
			  "PE",
			  "PH",
			  "PL",
			  "PS",
			  "PT",
			  "PY",
			  "QA",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "SV",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "UY",
			  "VN",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 178836,
			"explicit": false,
			"external_ids": {
			  "isrc": "USZZ81710033"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/5z8ExhfkwcCeimoVfaRQRw"
			},
			"href": "https://api.spotify.com/v1/tracks/5z8ExhfkwcCeimoVfaRQRw",
			"id": "5z8ExhfkwcCeimoVfaRQRw",
			"is_local": false,
			"name": "If You Had Spun out in Your Oldsmobile, This Probably Wouldn't Have Happened",
			"popularity": 35,
			"preview_url": "https://p.scdn.co/mp3-preview/7f0079976f954661a6a69c0ec92fa378b0297cd9?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 1,
			"type": "track",
			"uri": "spotify:track:5z8ExhfkwcCeimoVfaRQRw"
		  },
		  {
			"album": {
			  "album_type": "COMPILATION",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/7xGGqA85UIWX1GoTVM4itC"
				  },
				  "href": "https://api.spotify.com/v1/artists/7xGGqA85UIWX1GoTVM4itC",
				  "id": "7xGGqA85UIWX1GoTVM4itC",
				  "name": "The Staple Singers",
				  "type": "artist",
				  "uri": "spotify:artist:7xGGqA85UIWX1GoTVM4itC"
				}
			  ],
			  "available_markets": [
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CZ",
				"DE",
				"DK",
				"EE",
				"FI",
				"FR",
				"GB",
				"GR",
				"HK",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JP",
				"LT",
				"LU",
				"LV",
				"MX",
				"MY",
				"NL",
				"NO",
				"NZ",
				"PH",
				"PL",
				"PT",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/7tUOJxXojOWdWU2T2ZSge7"
			  },
			  "href": "https://api.spotify.com/v1/albums/7tUOJxXojOWdWU2T2ZSge7",
			  "id": "7tUOJxXojOWdWU2T2ZSge7",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b2737fd00f79db3d1ff67255a585",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e027fd00f79db3d1ff67255a585",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d000048517fd00f79db3d1ff67255a585",
				  "width": 64
				}
			  ],
			  "name": "The Very Best Of The Staple Singers",
			  "release_date": "2007-01-01",
			  "release_date_precision": "day",
			  "total_tracks": 20,
			  "type": "album",
			  "uri": "spotify:album:7tUOJxXojOWdWU2T2ZSge7"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/7xGGqA85UIWX1GoTVM4itC"
				},
				"href": "https://api.spotify.com/v1/artists/7xGGqA85UIWX1GoTVM4itC",
				"id": "7xGGqA85UIWX1GoTVM4itC",
				"name": "The Staple Singers",
				"type": "artist",
				"uri": "spotify:artist:7xGGqA85UIWX1GoTVM4itC"
			  }
			],
			"available_markets": [
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CZ",
			  "DE",
			  "DK",
			  "EE",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "HK",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JP",
			  "LT",
			  "LU",
			  "LV",
			  "MX",
			  "MY",
			  "NL",
			  "NO",
			  "NZ",
			  "PH",
			  "PL",
			  "PT",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 196826,
			"explicit": false,
			"external_ids": {
			  "isrc": "USFI80700331"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/5YLnfy7R2kueN0BRPkjiEG"
			},
			"href": "https://api.spotify.com/v1/tracks/5YLnfy7R2kueN0BRPkjiEG",
			"id": "5YLnfy7R2kueN0BRPkjiEG",
			"is_local": false,
			"name": "I'll Take You There",
			"popularity": 68,
			"preview_url": "https://p.scdn.co/mp3-preview/b66f84a27859593e70ebfe55e447523c5bf09a50?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 11,
			"type": "track",
			"uri": "spotify:track:5YLnfy7R2kueN0BRPkjiEG"
		  },
		  {
			"album": {
			  "album_type": "SINGLE",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/2kCO8LXN1usaOPL3iEE28I"
				  },
				  "href": "https://api.spotify.com/v1/artists/2kCO8LXN1usaOPL3iEE28I",
				  "id": "2kCO8LXN1usaOPL3iEE28I",
				  "name": "Tai Verdes",
				  "type": "artist",
				  "uri": "spotify:artist:2kCO8LXN1usaOPL3iEE28I"
				}
			  ],
			  "available_markets": [
				"US"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/7zPKeqaMfkP34adJfqpnQm"
			  },
			  "href": "https://api.spotify.com/v1/albums/7zPKeqaMfkP34adJfqpnQm",
			  "id": "7zPKeqaMfkP34adJfqpnQm",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b273589d963cc768cf794673b1fe",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e02589d963cc768cf794673b1fe",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d00004851589d963cc768cf794673b1fe",
				  "width": 64
				}
			  ],
			  "name": "DRUGS",
			  "release_date": "2020-10-30",
			  "release_date_precision": "day",
			  "total_tracks": 2,
			  "type": "album",
			  "uri": "spotify:album:7zPKeqaMfkP34adJfqpnQm"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/2kCO8LXN1usaOPL3iEE28I"
				},
				"href": "https://api.spotify.com/v1/artists/2kCO8LXN1usaOPL3iEE28I",
				"id": "2kCO8LXN1usaOPL3iEE28I",
				"name": "Tai Verdes",
				"type": "artist",
				"uri": "spotify:artist:2kCO8LXN1usaOPL3iEE28I"
			  }
			],
			"available_markets": [
			  "US"
			],
			"disc_number": 1,
			"duration_ms": 159546,
			"explicit": true,
			"external_ids": {
			  "isrc": "USQX92004562"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/0ooz46SAYMJsdOPnwW4nKK"
			},
			"href": "https://api.spotify.com/v1/tracks/0ooz46SAYMJsdOPnwW4nKK",
			"id": "0ooz46SAYMJsdOPnwW4nKK",
			"is_local": false,
			"name": "DRUGS",
			"popularity": 70,
			"preview_url": "https://p.scdn.co/mp3-preview/fd1b2e7128422bdd03a433f416585c8e846f9f74?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 1,
			"type": "track",
			"uri": "spotify:track:0ooz46SAYMJsdOPnwW4nKK"
		  },
		  {
			"album": {
			  "album_type": "ALBUM",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/0nq64XZMWV1s7XHXIkdH7K"
				  },
				  "href": "https://api.spotify.com/v1/artists/0nq64XZMWV1s7XHXIkdH7K",
				  "id": "0nq64XZMWV1s7XHXIkdH7K",
				  "name": "The Wonder Years",
				  "type": "artist",
				  "uri": "spotify:artist:0nq64XZMWV1s7XHXIkdH7K"
				}
			  ],
			  "available_markets": [
				"AD",
				"AE",
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BH",
				"BO",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CY",
				"CZ",
				"DE",
				"DK",
				"DO",
				"DZ",
				"EC",
				"EE",
				"EG",
				"ES",
				"FI",
				"FR",
				"GB",
				"GR",
				"GT",
				"HK",
				"HN",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JO",
				"JP",
				"KW",
				"LB",
				"LI",
				"LT",
				"LU",
				"LV",
				"MA",
				"MC",
				"MT",
				"MX",
				"MY",
				"NI",
				"NL",
				"NO",
				"NZ",
				"OM",
				"PA",
				"PE",
				"PH",
				"PL",
				"PS",
				"PT",
				"PY",
				"QA",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"SV",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"UY",
				"VN",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/7akNjSf8E2pyvNUmo6equT"
			  },
			  "href": "https://api.spotify.com/v1/albums/7akNjSf8E2pyvNUmo6equT",
			  "id": "7akNjSf8E2pyvNUmo6equT",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b273b454dbf64d6d56f5af0884f7",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e02b454dbf64d6d56f5af0884f7",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d00004851b454dbf64d6d56f5af0884f7",
				  "width": 64
				}
			  ],
			  "name": "No Closer To Heaven",
			  "release_date": "2015-09-04",
			  "release_date_precision": "day",
			  "total_tracks": 13,
			  "type": "album",
			  "uri": "spotify:album:7akNjSf8E2pyvNUmo6equT"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/0nq64XZMWV1s7XHXIkdH7K"
				},
				"href": "https://api.spotify.com/v1/artists/0nq64XZMWV1s7XHXIkdH7K",
				"id": "0nq64XZMWV1s7XHXIkdH7K",
				"name": "The Wonder Years",
				"type": "artist",
				"uri": "spotify:artist:0nq64XZMWV1s7XHXIkdH7K"
			  }
			],
			"available_markets": [
			  "AD",
			  "AE",
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BH",
			  "BO",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CY",
			  "CZ",
			  "DE",
			  "DK",
			  "DO",
			  "DZ",
			  "EC",
			  "EE",
			  "EG",
			  "ES",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "GT",
			  "HK",
			  "HN",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JO",
			  "JP",
			  "KW",
			  "LB",
			  "LI",
			  "LT",
			  "LU",
			  "LV",
			  "MA",
			  "MC",
			  "MT",
			  "MX",
			  "MY",
			  "NI",
			  "NL",
			  "NO",
			  "NZ",
			  "OM",
			  "PA",
			  "PE",
			  "PH",
			  "PL",
			  "PS",
			  "PT",
			  "PY",
			  "QA",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "SV",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "UY",
			  "VN",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 171363,
			"explicit": false,
			"external_ids": {
			  "isrc": "USHR21566007"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/0qKEQiF1jQShdHivzOkpEZ"
			},
			"href": "https://api.spotify.com/v1/tracks/0qKEQiF1jQShdHivzOkpEZ",
			"id": "0qKEQiF1jQShdHivzOkpEZ",
			"is_local": false,
			"name": "A Song for Ernest Hemingway",
			"popularity": 33,
			"preview_url": "https://p.scdn.co/mp3-preview/dac55d8321a63b3cb663b525ad4dcb093f208917?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 7,
			"type": "track",
			"uri": "spotify:track:0qKEQiF1jQShdHivzOkpEZ"
		  },
		  {
			"album": {
			  "album_type": "SINGLE",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/1lKZzN2d4IqiEYxyECIEHI"
				  },
				  "href": "https://api.spotify.com/v1/artists/1lKZzN2d4IqiEYxyECIEHI",
				  "id": "1lKZzN2d4IqiEYxyECIEHI",
				  "name": "Hot Mulligan",
				  "type": "artist",
				  "uri": "spotify:artist:1lKZzN2d4IqiEYxyECIEHI"
				}
			  ],
			  "available_markets": [
				"AD",
				"AE",
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BH",
				"BO",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CY",
				"CZ",
				"DE",
				"DK",
				"DO",
				"DZ",
				"EC",
				"EE",
				"EG",
				"ES",
				"FI",
				"FR",
				"GB",
				"GR",
				"GT",
				"HK",
				"HN",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JO",
				"JP",
				"KW",
				"LB",
				"LI",
				"LT",
				"LU",
				"LV",
				"MA",
				"MC",
				"MT",
				"MX",
				"MY",
				"NI",
				"NL",
				"NO",
				"NZ",
				"OM",
				"PA",
				"PE",
				"PH",
				"PL",
				"PS",
				"PT",
				"PY",
				"QA",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"SV",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"UY",
				"VN",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/0KxmgMVBqwcFvEeaTTX8UU"
			  },
			  "href": "https://api.spotify.com/v1/albums/0KxmgMVBqwcFvEeaTTX8UU",
			  "id": "0KxmgMVBqwcFvEeaTTX8UU",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b27334d228272503a51a4d93a817",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e0234d228272503a51a4d93a817",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d0000485134d228272503a51a4d93a817",
				  "width": 64
				}
			  ],
			  "name": "Fenton",
			  "release_date": "2015-05-12",
			  "release_date_precision": "day",
			  "total_tracks": 4,
			  "type": "album",
			  "uri": "spotify:album:0KxmgMVBqwcFvEeaTTX8UU"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/1lKZzN2d4IqiEYxyECIEHI"
				},
				"href": "https://api.spotify.com/v1/artists/1lKZzN2d4IqiEYxyECIEHI",
				"id": "1lKZzN2d4IqiEYxyECIEHI",
				"name": "Hot Mulligan",
				"type": "artist",
				"uri": "spotify:artist:1lKZzN2d4IqiEYxyECIEHI"
			  }
			],
			"available_markets": [
			  "AD",
			  "AE",
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BH",
			  "BO",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CY",
			  "CZ",
			  "DE",
			  "DK",
			  "DO",
			  "DZ",
			  "EC",
			  "EE",
			  "EG",
			  "ES",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "GT",
			  "HK",
			  "HN",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JO",
			  "JP",
			  "KW",
			  "LB",
			  "LI",
			  "LT",
			  "LU",
			  "LV",
			  "MA",
			  "MC",
			  "MT",
			  "MX",
			  "MY",
			  "NI",
			  "NL",
			  "NO",
			  "NZ",
			  "OM",
			  "PA",
			  "PE",
			  "PH",
			  "PL",
			  "PS",
			  "PT",
			  "PY",
			  "QA",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "SV",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "UY",
			  "VN",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 276856,
			"explicit": false,
			"external_ids": {
			  "isrc": "QZ4JJ1807769"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/4fvnr4ENE59eUqYHAKFpS0"
			},
			"href": "https://api.spotify.com/v1/tracks/4fvnr4ENE59eUqYHAKFpS0",
			"id": "4fvnr4ENE59eUqYHAKFpS0",
			"is_local": false,
			"name": "Buy a Fire Extinguisher Before You Need a Fire Extinguisher (Acoustic)",
			"popularity": 29,
			"preview_url": "https://p.scdn.co/mp3-preview/fd2aee12ffcb4cbb2884ab21e327b1832bdd801d?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 4,
			"type": "track",
			"uri": "spotify:track:4fvnr4ENE59eUqYHAKFpS0"
		  },
		  {
			"album": {
			  "album_type": "ALBUM",
			  "artists": [
				{
				  "external_urls": {
					"spotify": "https://open.spotify.com/artist/1lKZzN2d4IqiEYxyECIEHI"
				  },
				  "href": "https://api.spotify.com/v1/artists/1lKZzN2d4IqiEYxyECIEHI",
				  "id": "1lKZzN2d4IqiEYxyECIEHI",
				  "name": "Hot Mulligan",
				  "type": "artist",
				  "uri": "spotify:artist:1lKZzN2d4IqiEYxyECIEHI"
				}
			  ],
			  "available_markets": [
				"AD",
				"AE",
				"AR",
				"AT",
				"AU",
				"BE",
				"BG",
				"BH",
				"BO",
				"BR",
				"CA",
				"CH",
				"CL",
				"CO",
				"CR",
				"CY",
				"CZ",
				"DE",
				"DK",
				"DO",
				"DZ",
				"EC",
				"EE",
				"EG",
				"ES",
				"FI",
				"FR",
				"GB",
				"GR",
				"GT",
				"HK",
				"HN",
				"HU",
				"ID",
				"IE",
				"IL",
				"IN",
				"IS",
				"IT",
				"JO",
				"JP",
				"KW",
				"LB",
				"LI",
				"LT",
				"LU",
				"LV",
				"MA",
				"MC",
				"MT",
				"MX",
				"MY",
				"NI",
				"NL",
				"NO",
				"NZ",
				"OM",
				"PA",
				"PE",
				"PH",
				"PL",
				"PS",
				"PT",
				"PY",
				"QA",
				"RO",
				"SA",
				"SE",
				"SG",
				"SK",
				"SV",
				"TH",
				"TN",
				"TR",
				"TW",
				"US",
				"UY",
				"VN",
				"ZA"
			  ],
			  "external_urls": {
				"spotify": "https://open.spotify.com/album/3wl3zdJVNhLyJfqdXaCRyp"
			  },
			  "href": "https://api.spotify.com/v1/albums/3wl3zdJVNhLyJfqdXaCRyp",
			  "id": "3wl3zdJVNhLyJfqdXaCRyp",
			  "images": [
				{
				  "height": 640,
				  "url": "https://i.scdn.co/image/ab67616d0000b2738cc9c6e183cef734184e15b7",
				  "width": 640
				},
				{
				  "height": 300,
				  "url": "https://i.scdn.co/image/ab67616d00001e028cc9c6e183cef734184e15b7",
				  "width": 300
				},
				{
				  "height": 64,
				  "url": "https://i.scdn.co/image/ab67616d000048518cc9c6e183cef734184e15b7",
				  "width": 64
				}
			  ],
			  "name": "Pilot",
			  "release_date": "2018-03-23",
			  "release_date_precision": "day",
			  "total_tracks": 11,
			  "type": "album",
			  "uri": "spotify:album:3wl3zdJVNhLyJfqdXaCRyp"
			},
			"artists": [
			  {
				"external_urls": {
				  "spotify": "https://open.spotify.com/artist/1lKZzN2d4IqiEYxyECIEHI"
				},
				"href": "https://api.spotify.com/v1/artists/1lKZzN2d4IqiEYxyECIEHI",
				"id": "1lKZzN2d4IqiEYxyECIEHI",
				"name": "Hot Mulligan",
				"type": "artist",
				"uri": "spotify:artist:1lKZzN2d4IqiEYxyECIEHI"
			  }
			],
			"available_markets": [
			  "AD",
			  "AE",
			  "AR",
			  "AT",
			  "AU",
			  "BE",
			  "BG",
			  "BH",
			  "BO",
			  "BR",
			  "CA",
			  "CH",
			  "CL",
			  "CO",
			  "CR",
			  "CY",
			  "CZ",
			  "DE",
			  "DK",
			  "DO",
			  "DZ",
			  "EC",
			  "EE",
			  "EG",
			  "ES",
			  "FI",
			  "FR",
			  "GB",
			  "GR",
			  "GT",
			  "HK",
			  "HN",
			  "HU",
			  "ID",
			  "IE",
			  "IL",
			  "IN",
			  "IS",
			  "IT",
			  "JO",
			  "JP",
			  "KW",
			  "LB",
			  "LI",
			  "LT",
			  "LU",
			  "LV",
			  "MA",
			  "MC",
			  "MT",
			  "MX",
			  "MY",
			  "NI",
			  "NL",
			  "NO",
			  "NZ",
			  "OM",
			  "PA",
			  "PE",
			  "PH",
			  "PL",
			  "PS",
			  "PT",
			  "PY",
			  "QA",
			  "RO",
			  "SA",
			  "SE",
			  "SG",
			  "SK",
			  "SV",
			  "TH",
			  "TN",
			  "TR",
			  "TW",
			  "US",
			  "UY",
			  "VN",
			  "ZA"
			],
			"disc_number": 1,
			"duration_ms": 163508,
			"explicit": false,
			"external_ids": {
			  "isrc": "USZZ81810015"
			},
			"external_urls": {
			  "spotify": "https://open.spotify.com/track/22bdXU26ewL0ji9K9acT5J"
			},
			"href": "https://api.spotify.com/v1/tracks/22bdXU26ewL0ji9K9acT5J",
			"id": "22bdXU26ewL0ji9K9acT5J",
			"is_local": false,
			"name": "How Do You Know It’s Not Armadillo Shells?",
			"popularity": 53,
			"preview_url": "https://p.scdn.co/mp3-preview/37324ffc85d86768e1ba3d093d03e026482de7a5?cid=db8c415a273c4ffc9286e12ff24b8b9f",
			"track_number": 10,
			"type": "track",
			"uri": "spotify:track:22bdXU26ewL0ji9K9acT5J"
		  }
		],
		"total": 50,
		"limit": 20,
		"offset": 0,
		"previous": null,
		"href": "https://api.spotify.com/v1/me/top/tracks",
		"next": "https://api.spotify.com/v1/me/top/tracks?limit=20&offset=20"
	  }`
)
