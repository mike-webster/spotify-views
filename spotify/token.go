package spotify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/mike-webster/spotify-views/keys"
)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	Type         string `json:"token_type"`
	Scope        string `json:"scope"`
	Exp          int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type Token struct {
	Access  string
	Refresh string
}

// ----
// API
// ----

func ExchangeOauthCode(ctx context.Context, code string) (*Token, error) {
	req, err := getCodeSwapRequest(ctx, code)
	if err != nil {
		return nil, err
	}

	respBody, err := makeRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return parseTokensFromCodeSwapResponse(respBody)
}

// ----
// Members
// ----

func (t *Token) RefreshMe(ctx context.Context) (bool, error) {
	req, err := getRefreshRequest(ctx, t.Refresh)
	if err != nil {
		return false, err
	}

	resp, err := makeRequest(ctx, req)
	if err != nil {
		return false, err
	}

	tok, err := parseTokenFromRefreshResponse(resp)
	if err != nil {
		return false, err
	}

	t.Access = tok
	return true, nil
}

// ----
// Helpers
// ----

func parseTokenFromRefreshResponse(body *[]byte) (string, error) {
	type tempResp struct {
		AccessToken string `json:"access_token"`
	}

	var b tempResp
	err := json.Unmarshal(*body, &b)
	if err != nil {
		return "", err
	}
	return b.AccessToken, nil
}

func getRefreshRequest(ctx context.Context, refTok string) (*http.Request, error) {
	if len(refTok) < 1 {
		return nil, errors.New("no refresh token provided")
	}

	tokURL := "https://accounts.spotify.com/api/token"
	vals, err := getRefreshParams(ctx)
	if err != nil {
		return nil, err
	}

	body := url.Values{}
	body.Set("grant_type", "refresh_token")
	body.Set("refresh_token", refTok)

	req, err := http.NewRequest("POST", tokURL, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}

	key := base64.URLEncoding.EncodeToString([]byte(fmt.Sprint((*vals)["client_id"], ":", (*vals)["client_secret"])))
	req.Header.Add("Authorization", fmt.Sprint("Basic ", key))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func getRefreshParams(ctx context.Context) (*map[string]string, error) {
	ret := map[string]string{}
	clientID := keys.GetContextValue(ctx, keys.ContextSpotifyClientID)
	if clientID == nil {
		return &map[string]string{}, errors.New("no client id provided")
	}
	ret["client_id"] = fmt.Sprint(clientID)

	clientSecret := keys.GetContextValue(ctx, keys.ContextSpotifyClientSecret)
	if clientSecret == nil {
		return &map[string]string{}, errors.New("no client secret provided")
	}
	ret["client_secret"] = fmt.Sprint(clientSecret)

	return &ret, nil
}

func getCodeSwapParams(ctx context.Context) (*map[string]string, error) {
	ret := map[string]string{}
	returnURL := keys.GetContextValue(ctx, keys.ContextSpotifyReturnURL)
	if returnURL == nil {
		return &map[string]string{}, errors.New("no return url provided")
	}

	strReturnURL, ok := returnURL.(string)
	if !ok {
		return &map[string]string{}, errors.New("return url couldn't be parsed")
	}
	ret["url"] = strReturnURL

	clientID := keys.GetContextValue(ctx, keys.ContextSpotifyClientID)
	if clientID == nil {
		return &map[string]string{}, errors.New("no client id provided")
	}

	strClientID, ok := clientID.(string)
	if !ok {
		return &map[string]string{}, errors.New("client id couldn't be parsed")
	}
	ret["client_id"] = strClientID

	clientSecret := keys.GetContextValue(ctx, keys.ContextSpotifyClientSecret)
	if clientSecret == nil {
		return &map[string]string{}, errors.New("no client secret provided")
	}

	strClientSecret, ok := clientSecret.(string)
	if !ok {
		return &map[string]string{}, errors.New("client secret couldn't be parsed")
	}
	ret["client_secret"] = strClientSecret

	return &ret, nil
}

func getCodeSwapRequest(ctx context.Context, code string) (*http.Request, error) {
	vals, err := getCodeSwapParams(ctx)
	if err != nil {
		return nil, err
	}

	tokenURL := "https://accounts.spotify.com/api/token"
	body := url.Values{}
	body.Set("grant_type", "authorization_code")
	body.Set("code", code)
	body.Set("redirect_uri", (*vals)["url"])
	body.Set("client_id", (*vals)["client_id"])
	body.Set("client_secret", (*vals)["client_secret"])

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func parseTokensFromCodeSwapResponse(resp *[]byte) (*Token, error) {
	var r tokenResponse
	err := json.Unmarshal(*resp, &r)
	if err != nil {
		return nil, err
	}

	return &Token{Access: r.AccessToken, Refresh: r.RefreshToken}, nil
}
