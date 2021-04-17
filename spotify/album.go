package spotify

type Album struct {
	Name   string  `json:"name"`
	Images []Image `json:"images"`
}

func (a *Album) Loc() string {
	if len(a.Images) > 0 {
		return a.Images[0].URL
	}

	return ""
}