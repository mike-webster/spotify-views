package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func requestTokens(ctx context.Context, code string) ([]string, error) {
	returnURL := ctx.Value("return_url")
	if returnURL == nil {
		return []string{}, errors.New("no return url provided")
	}
	strReturnURL, ok := returnURL.(string)
	if !ok {
		return []string{}, errors.New("return url couldn't be parsed")
	}

	clientID := ctx.Value("client_id")
	if clientID == nil {
		return []string{}, errors.New("no client id provided")
	}
	strClientID, ok := clientID.(string)
	if !ok {
		return []string{}, errors.New("client id couldn't be parsed")
	}

	clientSecret := ctx.Value("client_secret")
	if clientSecret == nil {
		return []string{}, errors.New("no client secret provided")
	}
	strClientSecret, ok := clientSecret.(string)
	if !ok {
		return []string{}, errors.New("client secret couldn't be parsed")
	}

	client := &http.Client{}
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

	resp, err := client.Do(req)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Println(fmt.Sprintf("error -- non 200 response -- Body: %s", b))
		return []string{}, errors.New(fmt.Sprint("non 200 response from spotify: ", resp.Status))
	}

	var r spotifyResponse
	err = json.Unmarshal(b, &r)
	if err != nil {
		return []string{}, err
	}
	log.Println(r.ToString())
	return []string{
		r.AccessToken, r.RefreshToken,
	}, nil
}

func getPairs(m map[string]int32) Pairs {
	ret := Pairs{}
	for k, v := range m {
		ret = append(ret, Pair{Key: k, Value: v})
	}
	return ret
}
