package spotify

import "time"

type item struct {
	DateSaved time.Time `json:"added_at"`
	Track     Track     `json:"track"`
}

type items []item

func (i items) Tracks() Tracks {
	ret := Tracks{}
	for _, j := range i {
		ret = append(ret, j.Track)
	}
	return ret
}
