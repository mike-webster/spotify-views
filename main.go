package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var scopes = []string{
	"user-modify-playback-state",
	"user-read-playback-state",
	"streaming",
	"app-remote-control",
	"user-top-read",
	"user-read-playback-position",
	"user-read-recently-played",
}
var clientID = ""
var clientSecret = ""
var host = ""
var returnURL = ""

func main() {
	err := parseEnvironmentVariables()
	returnURL = fmt.Sprint("https://", host, "/spotify/oauth")
	if err != nil {
		panic(err)
	}

	runServer()
}

func parseEnvironmentVariables() error {
	clientID = os.Getenv("CLIENT_ID")
	if len(clientID) < 1 {
		return errors.New("no client id provided")
	}
	clientSecret = os.Getenv("CLIENT_SECRET")
	if len(clientSecret) < 1 {
		return errors.New("no client secret provided")
	}
	host = os.Getenv("HOST")
	if len(host) < 1 {
		return errors.New("no host provided")
	}
	return nil
}

func runServer() {
	r := gin.Default()

	r.GET("/spotify/oauth", func(c *gin.Context) {
		code := c.Query("code")
		//state := c.Query("state")
		qErr := c.Query("error")
		if len(qErr) > 0 {
			// the user is a fucker and they denied access
		}
		tokens, err := requestTokens(code)
		if err != nil {
			log.Println("could not retrieve tokens for user; error: ", err)
			c.JSON(500, gin.H{"msg": err})
		}
		log.Println(fmt.Sprint("success - tokens: \n\tAccess: ", tokens[0], "\n\tRefres: ", tokens[1]))
		c.JSON(200, gin.H{"msg": tokens})
	})

	r.GET("/login", func(c *gin.Context) {
		// TODO Add state
		pathScopes := url.QueryEscape(strings.Join(scopes, " "))
		redirectURL := fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=code&client_id=%s&scopes=%s&redirect_uri=%s&show_dialog=false",
			clientID,
			pathScopes,
			returnURL)
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	})

	r.Run()
}

func requestTokens(code string) ([]string, error) {
	type spotifyResponse struct {
		AccessToken  string `json:"access_token"`
		Type         string `json:"token_type"`
		Scope        string `json:"scope"`
		Exp          int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}

	client := &http.Client{}
	url := "https://accounts.spotify.com/api/token"
	contentType := "application/x-www-form-urlencoded"
	body := map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  returnURL,
		"client_id":     clientID,
		"client_secret": clientSecret,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return []string{}, err
	}
	log.Println("json body: ", string(jsonBody))

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return []string{}, err
	}
	req.Header.Add("Content-Type", contentType)

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
	return []string{
		r.AccessToken, r.RefreshToken,
	}, nil
}
