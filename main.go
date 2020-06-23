package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
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
var lyricsKey = ""

type ViewBag struct {
	Resource string
	Results  interface{}
}

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
	lyricsKey = os.Getenv("LYRICS_KEY")
	if len(lyricsKey) < 1 {
		return errors.New("no lyrics key provided")
	}
	return nil
}

func runServer() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/spotify/oauth", handlerOauth)
	r.GET("/tracks/top", handlerTopTracks)
	r.GET("/artists/top", handlerTopArtists)
	r.GET("/artists/genres", handlerTopArtistsGenres)
	r.GET("/tracks/genres", handlerTopTracksGenres)
	r.GET("/login", handlerLogin)
	r.GET("/", handlerHome)

	r.Static("/static/css", "./static")
	r.Static("/static/js", "./static")

	r.Run()
}
