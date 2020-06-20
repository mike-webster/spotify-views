package spotify

import (
	"fmt"
)

// TODO: maybe define an interface to handle shared code?

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

type Tracks []Track

func (t *Track) EmbeddedPlayer() string {
	return fmt.Sprintf(`<iframe src="https://open.spotify.com/embed/track/%s" width="300" height="80" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>`, t.ID)
}

func (t *Tracks) IDs() []string {
	ret := []string{}
	for _, i := range *t {
		ret = append(ret, i.ID)
	}
	return ret
}

type Artist struct {
	Genres     []string `json:"genres"`
	Name       string   `json:"name"`
	Popularity int32    `json:"popularity"`
	URI        string   `json:"uri"`
	ID         string   `json:"ID"`
}

type Artists []Artist

func (a *Artist) EmbeddedPlayer() string {
	return fmt.Sprintf(`<h4 width="300" style="text-align:center">%s</h4><iframe src="https://open.spotify.com/embed/artist/%s" width="300" height="380" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>`, a.Name, a.ID)
}

func (a *Artists) IDs() []string {
	ret := []string{}
	for _, i := range *a {
		ret = append(ret, i.ID)
	}
	return ret
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

type Pair struct {
	Key   string
	Value int32
}
type Pairs []Pair

func (p Pairs) Len() int           { return len(p) }
func (p Pairs) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p Pairs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Pairs) Contains(key string) int {
	for i, ii := range p {
		if ii.Key == key {
			return i
		}
	}
	return -1
}
func (p Pairs) ToMap() map[string]int32 {
	ret := map[string]int32{}
	for _, i := range p {
		ret[i.Key] = i.Value
	}
	return ret
}
