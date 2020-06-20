package spotify

import (
	"fmt"
	"strings"
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
	return fmt.Sprintf(`<iframe src="https://open.spotify.com/embed/track/%s" width="300" height="80" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>`, strings.Split(t.URI, ":")[2])
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
