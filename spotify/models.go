package spotify

import "fmt"

type Track struct {
	Links map[string]string `json:"external_links"`
	Name  string            `json:"name"`
}

type Tracks []Track

func (t *Track) EmbeddedPlayer() string {
	return fmt.Sprintf(`<iframe src="%s" width="300" height="80" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>`, t.Links["spotify"])
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
