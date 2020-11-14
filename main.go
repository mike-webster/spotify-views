package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	data "github.com/mike-webster/spotify-views/data"
	"github.com/mike-webster/spotify-views/spotify"
)

var (
	scopes = []string{
		// "user-modify-playback-state",
		// "user-read-playback-state",
		// "streaming",
		// "app-remote-control",
		"user-top-read",
		"user-read-email",
		"user-library-read",
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
	redisHost    = ""
	redisPort    = ""
	redisPass    = ""
	_redisDB     *redis.Client
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

	if os.Getenv("GO_ENV") == "test" {
		testMethod(ctx)
		return
	}

	err = data.Ping(ctx)
	if err != nil {
		panic(fmt.Sprint("couldnt connect to database; ", err.Error()))
	}

	runServer()
}

func testMethod(ctx context.Context) {
	ctx = context.WithValue(ctx, spotify.ContextAccessToken, os.Getenv("SPOT_TOK"))
	genres, err := spotify.GetGenres(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(genres)
}
