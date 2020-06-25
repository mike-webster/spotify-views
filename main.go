package main

import (
	"fmt"
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

// ViewBag is a basic struct to use to pass information to the views
// TODO move this into handlers.go
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
