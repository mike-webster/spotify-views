package spotify

import (
	"fmt"
	"strings"
)

type Track struct {
	Links map[string]string `json:"external_urls"`
	Name  string            `json:"name"`
	URI   string            `json:"uri"`
}

type Tracks []Track

func (t *Track) EmbeddedPlayer() string {
	return fmt.Sprintf(`<iframe src="https://open.spotify.com/embed/track/%s" width="300" height="80" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>`, strings.Split(t.URI, ":")[2])
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
