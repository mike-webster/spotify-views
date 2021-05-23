package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/mike-webster/spotify-views/keys"
	"github.com/mike-webster/spotify-views/logging"
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

var (
	ErrNoContext      = "no context provided"
	ErrFieldTooShort  = "field too short: "
	ErrMissingDeps    = "no dependencies in context"
	ErrNoRowsAffected = "no rows affected"
)

func (u *User) Save(ctx context.Context) error {
	if ctx == nil {
		return errors.New(ErrNoContext)
	}

	deps := GetDependencies(ctx)
	if deps == nil {
		return errors.New(ErrMissingDeps)
	}

	if len(u.ID) < 1 {
		return errors.New(fmt.Sprint(ErrFieldTooShort, "ID"))
	}

	if len(u.Email) < 1 {
		return errors.New(fmt.Sprint(ErrFieldTooShort, "Email"))
	}

	query := `INSERT INTO users	(spotify_id, email) VALUES ('p1', 'p2')`
	query = strings.Replace(query, "p1", u.ID, 1)
	query = strings.Replace(query, "p2", u.Email, 1)
	res, err := deps.DB.Exec(ctx, query)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
			logging.GetLogger(ctx).WithFields(map[string]interface{}{
				"event": "duplicate_user_insert",
				"email": u.Email,
			}).Error()
			return nil
		}
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows > 0 {
		return nil
	}

	return errors.New(ErrNoRowsAffected)
}

func GetUser(ctx context.Context) (*User, error) {
	req, err := parseGetUserRequest(ctx)
	if err != nil {
		return nil, err
	}

	body, err := makeRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return parseGetUserResponse(body)
}

func GetSavedTracks(ctx context.Context) (*Tracks, error) {
	url := "https://api.spotify.com/v1/me/tracks?limit=50&offset=0"
	more := true
	ret := Tracks{}
	for more {
		t, newUrl, tot, err := getChunkOfUserLibraryTracks(ctx, url)
		if err != nil {
			logging.GetLogger(ctx).Warn(err.Error())
			more = false
		}
		url = newUrl

		ret = append(ret, *t...)
		if tot == len(ret) {
			more = false
		}
	}
	return &ret, nil
}

func getChunkOfUserLibraryTracks(ctx context.Context, url string) (*Tracks, string, int, error) {
	req, err := parseGetSavedTracksRequest(ctx, url)
	if err != nil {
		return nil, "", 0, err
	}

	body, err := makeRequest(ctx, req)
	if err != nil {
		return nil, "", 0, err
	}

	return parseGetSavedTracksResponse(body)
}

func parseGetSavedTracksRequest(ctx context.Context, url string) (*http.Request, error) {
	token := keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)
	if token == nil {
		return nil, ErrNoToken("no access token provided")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))
	return req, nil
}

func parseGetSavedTracksResponse(body *[]byte) (*Tracks, string, int, error) {
	type tempResp struct {
		Link   string `json:"href"`
		Items  items  `json:"items"`
		Limit  int    `json:"limit"`
		Next   string `json:"next"`
		Offset int    `json:"offset"`
		Total  int    `json:"total"`
	}

	var ret tempResp
	err := json.Unmarshal(*body, &ret)
	if err != nil {
		return nil, "", 0, err
	}
	rettr := ret.Items.Tracks()
	return &rettr, ret.Next, ret.Total, nil
}

func parseGetUserRequest(ctx context.Context) (*http.Request, error) {
	token := keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)
	if token == nil {
		return nil, ErrNoToken("no access token provided")
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))
	return req, nil
}

func parseGetUserResponse(body *[]byte) (*User, error) {
	ret := User{}
	err := json.Unmarshal(*body, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
