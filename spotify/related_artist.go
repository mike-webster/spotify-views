package spotify

type RelatedArtist struct {
	Followers  int64    `json:"followers"`
	Genres     []string `json:"genres"`
	Link       string   `json:"href"`
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Popularity int64    `json:"popularity"`
	Type       string   `json:"type"`
}
