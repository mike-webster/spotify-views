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
