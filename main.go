package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	data "github.com/mike-webster/spotify-views/data"
	"github.com/mike-webster/spotify-views/keys"
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
	ctx := context.Background()
	vals, err := parseEnvironmentVariables(ctx)
	if err != nil {
		panic(err)
	}
	for k, v := range vals {
		ctx = context.WithValue(ctx, k, v)
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

	// ctx, err := spotify.GetTopArtists(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(ctx.Value(spotify.ContextResults))

	// artists, err := spotify.GetRelatedArtists(ctx, "6FBDaR13swtiWwGhX1WQsP")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(artists)

	//tt, err := spotify.GetTopTracksForArtist(ctx, "6FBDaR13swtiWwGhX1WQsP")
	// if err != nil {
	// 	panic(err)
	// }

	//fmt.Println(getData(ctx))
	fmt.Println(strings.Join(*getData(ctx), ","))
}
