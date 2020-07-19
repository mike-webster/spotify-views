package spotify

import (
	"encoding/json"
	"fmt"
)

// TODO: maybe define an interface to handle shared code?

// Track represents a spotify track
type Track struct {
	Links   map[string]string `json:"external_urls"`
	Name    string            `json:"name"`
	URI     string            `json:"uri"`
	ID      string            `json:"id"`
	Artists []struct {
		Link string `json:"href"`
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"artists"`
}

// Tracks is a collection of spotify Tracks
type Tracks []Track

// EmbeddedPlayer will return the html to use for rendering the embedded spotify
// player iframe
func (t *Track) EmbeddedPlayer() string {
	return fmt.Sprintf(`<iframe src="https://open.spotify.com/embed/track/%s" width="300" height="80" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>`, t.ID)
}

// IDs returns the ID for each of the tracks in the collection of Tracks
func (t *Tracks) IDs() []string {
	ret := []string{}
	for _, i := range *t {
		ret = append(ret, i.ID)
	}
	return ret
}

// Artist represents a spotify Artist
type Artist struct {
	Genres     []string `json:"genres"`
	Name       string   `json:"name"`
	Popularity int32    `json:"popularity"`
	URI        string   `json:"uri"`
	ID         string   `json:"ID"`
}

// Artists is a collection of spotify Artist
type Artists []Artist

// EmbeddedPlayer  will return the html to use for rendering the embedded spotify
// player iframe
func (a *Artist) EmbeddedPlayer() string {
	return fmt.Sprintf(`<h4 width="300" style="text-align:center">%s</h4><iframe src="https://open.spotify.com/embed/artist/%s" width="300" height="380" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>`, a.Name, a.ID)
}

// IDs returns the ID for each of the artists in the collection of Artists
func (a *Artists) IDs() []string {
	ret := []string{}
	for _, i := range *a {
		ret = append(ret, i.ID)
	}
	return ret
}

type AudioFeature struct {
	Danceability     float32 `json:"danceability"`
	Energy           float32 `json:"energy"`
	Key              int     `json:"key"`
	Loudness         float32 `json:"loudness"`
	Mode             int     `json:"mode"`
	Speechiness      float32 `json:"speechiness"`
	Acousticness     float32 `json:"acousticness"`
	Instrumentalness float32 `json:"instrumentalness"`
	Liveness         float32 `json:"liveness"`
	Valence          float32 `json:"valence"`
	Tempo            float32 `json:"tempo"`
	Duration         int64   `json:"duration_ms"`
	TimeSignature    int     `json:"time_signature"`
}

type AudioFeatures []AudioFeature

func (af AudioFeature) String() string {
	str, err := json.Marshal(af)
	if err != nil {
		return ""
	}
	return string(str)
}

type spotifyResponse struct {
	AccessToken  string `json:"access_token"`
	Type         string `json:"token_type"`
	Scope        string `json:"scope"`
	Exp          int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (sr *spotifyResponse) ToString() string {
	ret := ""
	ret += "-- Spotify Response\n"
	ret += fmt.Sprint("\tAccessToken: ", sr.AccessToken, "\n")
	ret += fmt.Sprint("\tScope: ", sr.Scope, "\n")
	return ret
}

// TODO replace this with using the sortable map. it's the same thing.

// Pair is an outdated way to sort a map
type Pair struct {
	Key   string
	Value int32
}

// Pairs is a collection of Pair
type Pairs []Pair

func (p Pairs) Len() int           { return len(p) }
func (p Pairs) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p Pairs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Contains returns the index of the element in the collection, if it exists.
func (p Pairs) Contains(key string) int {
	for i, ii := range p {
		if ii.Key == key {
			return i
		}
	}
	return -1
}

// ToMap returns a map representation of the Pairs
func (p Pairs) ToMap() map[string]int32 {
	ret := map[string]int32{}
	for _, i := range p {
		ret[i.Key] = i.Value
	}
	return ret
}
