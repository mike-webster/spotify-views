package spotify

import "fmt"

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
