package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mike-webster/spotify-views/keys"
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
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

func parseGetUserRequest(ctx context.Context) (*http.Request, error) {
	token := keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
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
