package spotify

import "encoding/json"

type AudioFeature struct {
	ID               string  `json:"id"`
	URI              string  `json:"uri"`
	Danceability     float32 `json:"danceability"`
	Energy           float32 `json:"energy"`
	Key              int     `json:"key"`
	Loudness         float32 `json:"loudness"`
	Mode             int     `json:"mode"`
	Speechiness      float32 `json:"speechiness"`
	Acousticness     float32 `json:"acousticness"`
	Instrumentalness float32 `json:"instrumentalness"`
	Liveness         float32 `json:"liveness"`
	Valence          float32 `json:"valence"`
	Tempo            float32 `json:"tempo"`
	Duration         int64   `json:"duration_ms"`
	TimeSignature    int     `json:"time_signature"`
}

type AudioFeatures []AudioFeature

func (af AudioFeature) String() string {
	str, err := json.Marshal(af)
	if err != nil {
		return ""
	}
	return string(str)
}

func (af *AudioFeatures) String() []string {
	ret := []string{}
	for _, i := range *af {
		ret = append(ret, i.String())
	}
	return ret
}
