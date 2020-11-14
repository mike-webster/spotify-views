package spotify

type Recommendation struct {
	Tracks []struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`
		Link    string   `json:"href"`
		Artists []string `json:"artists"`
	} `json:"tracks"`
	Seeds []struct {
		ID   string `json:"id"`
		Link string `json:"href"`
		Type string `json:"type"`
	} `json:"seeds"`
}
