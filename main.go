package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	spotify "github.com/mike-webster/spotify-views/spotify"
)

var scopes = []string{
	// "user-modify-playback-state",
	// "user-read-playback-state",
	// "streaming",
	// "app-remote-control",
	"user-top-read",
	// "user-read-playback-position",
	// "user-read-recently-played",
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
	var err error
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/spotify/oauth", func(c *gin.Context) {
		ctx := context.WithValue(c, "return_url", returnURL)
		ctx = context.WithValue(ctx, "client_id", clientID)
		ctx = context.WithValue(ctx, "client_secret", clientSecret)
		code := c.Query("code")
		//state := c.Query("state")
		qErr := c.Query("error")
		if len(qErr) > 0 {
			// the user is a fucker and they denied access
		}
		ctx, err = spotify.HandleOauth(ctx, code)
		if err != nil {
			c.JSON(500, gin.H{"err": err})
			return
		}
		token := ctx.Value("access_token")
		if token == nil {
			c.JSON(500, gin.H{"err": "no access token"})
			return
		}
		c.SetCookie("svauth", fmt.Sprint(token), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
		c.JSON(200, gin.H{"msg": "yay"})
	})

	r.GET("/tracks/top", func(c *gin.Context) {
		token, err := c.Cookie("svauth")
		if err != nil {
			log.Println("no token - redirecting to login")
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		ctx := context.WithValue(c, "access_token", token)
		tracks, err := spotify.GetTopTracks(ctx)
		markups := []string{}
		for _, i := range *tracks {
			markups = append(markups, i.EmbeddedPlayer())
		}
		if err != nil {
			c.JSON(500, gin.H{"err": err})
			return
		}

		c.HTML(200, "toptracks.tmpl", markups)
	})

	r.GET("/login", func(c *gin.Context) {
		// TODO Add state
		pathScopes := url.QueryEscape(strings.Join(scopes, " "))
		redirectURL := fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s&show_dialog=false",
			clientID,
			pathScopes,
			returnURL)
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	})

	r.Run()
}
