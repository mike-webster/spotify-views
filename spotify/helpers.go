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
)

func requestTokens(ctx context.Context, code string) ([]string, error) {
	returnURL := ctx.Value(ContextReturnURL)
	if returnURL == nil {
		return []string{}, errors.New("no return url provided")
	}

	strReturnURL, ok := returnURL.(string)
	if !ok {
		return []string{}, errors.New("return url couldn't be parsed")
	}

	clientID := ctx.Value(ContextClientID)
	if clientID == nil {
		return []string{}, errors.New("no client id provided")
	}

	strClientID, ok := clientID.(string)
	if !ok {
		return []string{}, errors.New("client id couldn't be parsed")
	}

	clientSecret := ctx.Value(ContextClientSecret)
	if clientSecret == nil {
		return []string{}, errors.New("no client secret provided")
	}

	strClientSecret, ok := clientSecret.(string)
	if !ok {
		return []string{}, errors.New("client secret couldn't be parsed")
	}

	tokenURL := "https://accounts.spotify.com/api/token"
	body := url.Values{}
	body.Set("grant_type", "authorization_code")
	body.Set("code", code)
	body.Set("redirect_uri", strReturnURL)
	body.Set("client_id", strClientID)
	body.Set("client_secret", strClientSecret)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(body.Encode()))
	if err != nil {
		return []string{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	respBody, err := makeRequest(ctx, req, false)
	if err != nil {
		return []string{}, err
	}

	var r tokenResponse
	err = json.Unmarshal(*respBody, &r)
	if err != nil {
		return []string{}, err
	}

	return []string{
		r.AccessToken, r.RefreshToken,
	}, nil
}

func refreshToken(ctx context.Context) (string, error) {
	refTok := ctx.Value(ContextRefreshToken)
	if refTok == nil {
		return "", errors.New("no refresh token provided")
	}
	clientID := ctx.Value(ContextClientID)
	if clientID == nil {
		return "", errors.New("no client id provided")
	}
	clientSecret := ctx.Value(ContextClientSecret)
	if clientSecret == nil {
		return "", errors.New("no client secret is provided")
	}

	body := url.Values{}
	body.Set("grant_type", "refresh_token")
	body.Set("refresh_token", fmt.Sprint(refTok))
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret))))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := makeRequest(ctx, req, false)
	if err != nil {
		return "", err
	}

	type tempResp struct {
		AccessToken string `json:"access_token"`
	}

	var b tempResp
	err = json.Unmarshal(*resp, &b)
	if err != nil {
		return "", err
	}
	return b.AccessToken, nil
}

func getPairs(m map[string]int32) Pairs {
	ret := Pairs{}
	for k, v := range m {
		ret = append(ret, Pair{Key: k, Value: v})
	}

	return ret
}
