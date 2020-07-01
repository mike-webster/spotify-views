package main

import (
	"context"
	"fmt"

	data "github.com/mike-webster/spotify-views/data"
)

var (
	scopes = []string{
		// "user-modify-playback-state",
		// "user-read-playback-state",
		// "streaming",
		// "app-remote-control",
		"user-top-read",
		"user-read-email",
		// "user-read-playback-position",
		// "user-read-recently-played",
	}
	clientID     = ""
	clientSecret = ""
	host         = ""
	returnURL    = ""
	lyricsKey    = ""
	dbHost       = ""
	dbUser       = ""
	dbPass       = ""
	dbName       = ""
	secKey       = ""
)

// ViewBag is a basic struct to use to pass information to the views
// TODO move this into handlers.go
type ViewBag struct {
	Resource string
	Results  interface{}
}

func main() {
	ctx, err := parseEnvironmentVariables(context.Background())
	if err != nil {
		panic(err)
	}

	err = data.Ping(ctx)
	if err != nil {
		panic(fmt.Sprint("couldnt connect to database; ", err.Error()))
	}

	runServer()
}
