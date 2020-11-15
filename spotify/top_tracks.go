package spotify

type TopTracks struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Popularity int64    `json:"popularity"`
	Artists    []string `json:"artists"`
}
