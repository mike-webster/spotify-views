package main

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"image/png"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/mike-webster/spotify-views/keys"
	"github.com/mike-webster/spotify-views/logging"
	"github.com/mike-webster/spotify-views/spotify"
	"github.com/psykhi/wordclouds"
)

// TODO: put together a yaml config to parse these for me
// TODO: these context keys should all be the same type, not unique per package. abstract
//       the keys package
func parseEnvironmentVariables(ctx context.Context) (map[interface{}]interface{}, error) {
	type env struct {
		ClientID     string `envconfig:"CLIENT_ID"`
		ClientSecret string `envconfig:"CLIENT_SECRET"`
		Host         string `envconfig:"HOST"`
		Port         string `envconfig:"PORT"`
		LyricsKey    string `envconfig:"LYRICS_KEY"`
		DbHost       string `envconfig:"DB_HOST"`
		DbUser       string `envconfig:"DB_USER"`
		DbPass       string `envconfig:"DB_PASS"`
		DbName       string `envconfig:"DB_NAME"`
		SecKey       string `envconfig:"SEC_KEY"`
		RedisHost    string `envconfig:"REDIS_HOST"`
		RedisPort    string `envconfig:"REDIS_PORT"`
		RedisPass    string `envconfig:"REDIS_PASS"`
	}
	e := env{}
	envconfig.MustProcess("", &e)

	ret := map[interface{}]interface{}{}
	if len(e.ClientID) < 1 {
		return nil, errors.New("no client id provided")
	}
	ret[keys.ContextSpotifyClientID] = e.ClientID

	if len(e.ClientSecret) < 1 {
		return nil, errors.New("no client secret provided")
	}
	ret[keys.ContextSpotifyClientSecret] = e.ClientSecret

	if len(e.Host) < 1 {
		return nil, errors.New("no host provided")
	}
	ret[keys.ContextHost] = e.Host
	ret[keys.ContextSpotifyReturnURL] = fmt.Sprint("https://www.", e.Host, "/spotify/oauth")

	if len(e.Port) < 1 {
		return nil, errors.New("no port provided")
	}
	ret[keys.ContextPort] = e.Port

	// TODO: Do we need this in the context? or just set for the main package?
	// consider: the main goal here is to be able to verify everything is working
	// on app start using the context returned from this method.

	if len(e.LyricsKey) < 1 {
		return nil, errors.New("no lyrics key provided")
	}
	ret[keys.ContextLyricsToken] = e.LyricsKey

	if len(e.DbHost) < 1 {
		return nil, errors.New("no db host provided")
	}
	ret[keys.ContextDbHost] = e.DbHost

	if len(e.DbUser) < 1 {
		return nil, errors.New("no db user provided")
	}
	ret[keys.ContextDbUser] = e.DbUser

	if len(e.DbPass) < 1 {
		return nil, errors.New("no db pass provided")
	}
	ret[keys.ContextDbPass] = e.DbPass

	if len(e.DbName) < 1 {
		return nil, errors.New("no db name provided")
	}
	ret[keys.ContextDatabase] = e.DbName

	if len(e.SecKey) < 1 {
		return nil, errors.New("no sec key provided")
	}
	ret[keys.ContextSecurityKey] = e.SecKey

	ret["port"] = e.Port

	ret["redis-host"] = os.Getenv("REDIS_HOST")
	ret["redis-port"] = os.Getenv("REDIS_PORT")
	ret["redis-pass"] = os.Getenv("REDIS_PASS")

	return ret, nil
}

func generateWordCloud(ctx context.Context, filename string, wordCounts map[string]int) error {
	colors := []color.RGBA{
		//{0x17, 0xA5, 0x54, 0xff},
		{0x1E, 0xD7, 0x60, 0xFF},
	}

	rgbaColors := []color.Color{}
	for _, i := range colors {
		rgbaColors = append(rgbaColors, i)
	}

	w := wordclouds.NewWordcloud(
		wordCounts,
		wordclouds.FontFile("fonts/Ubuntu-L.ttf"),
		wordclouds.Height(2048),
		wordclouds.Width(2048),
		wordclouds.Debug(),
		wordclouds.FontMaxSize(300),
		wordclouds.FontMinSize(30),
		wordclouds.Colors(rgbaColors),
		wordclouds.RandomPlacement(false),
		wordclouds.MaskBoxes([]*wordclouds.Box{}),
	)

	img := w.Draw()
	directory := "static/clouds/"
	of, err := os.Create(fmt.Sprint(directory, filename))
	if err != nil {
		return err
	}

	err = png.Encode(of, img)
	if err != nil {
		return err
	}
	return of.Close()
}

var (
	PathSpotifyOauth     = "/spotify/oauth"
	PathTopTracks        = "/tracks/top"
	PathTopArtists       = "/artists/top"
	PathTopArtistGenres  = "/artists/genres"
	PathTopTracksGenres  = "/tracks/genres"
	PathLogin            = "/login"
	PathHome             = "/"
	PathWordCloud        = "/wordcloud"
	PathWordCloudData    = "/wordcloud/data"
	PathUserLibraryTempo = "/library/tempo"
	PathRecommendations  = "/tracks/recommendations"
	PathTest             = "/test"
)

func runServer(ctx context.Context) {
	r := gin.New()
	r.Use(recovery)
	//r.Use(requestLogger)
	//r.Use(redisClient)
	r.Use(setTokens)
	r.Use(setUserID)
	r.Use(setEnv)
	r.Use(setContextLogger)
	r.Use(setRequestID)
	r.Use(setClientIP)
	r.Use(setMethod)
	r.Use(setRequestPath)
	r.Use(setRequestQuery)
	r.Use(setReferer)
	r.Use(setUserAgent)

	r.StaticFile("/sitemap", "./static/sitemap.xml")
	r.Static("/static/css", "./static")
	r.Static("/static/js", "./static")
	r.Static("/logos/", "./static/logos")
	r.Static("/images/", "./static/images")
	r.StaticFile("/favicon.ico", "./static/logos/favicon.ico")
	r.LoadHTMLGlob("templates/*")

	r.GET(PathHome, handlerHome)
	r.GET(PathSpotifyOauth, handlerOauth)
	r.GET(PathLogin, handlerLogin)

	r.GET(PathTopTracks, authenticate, handlerTopTracks)
	r.GET(PathTopArtists, authenticate, handlerTopArtists)
	r.GET(PathTopArtistGenres, authenticate, handlerTopArtistsGenres)
	r.GET(PathTopTracksGenres, authenticate, handlerTopTracksGenres)
	r.GET(PathWordCloud, authenticate, handlerWordCloud)
	r.GET(PathWordCloudData, authenticate, handlerWordCloudData)
	r.GET(PathUserLibraryTempo, authenticate, handlerUserLibraryTempo)
	r.GET(PathRecommendations, authenticate, handlerRecommendations)
	r.GET(PathTest, authenticate, handlerTest)

	r.Run()
}

// what the fuck is this? Why is it taking a gin context like a controller handler but
// still getting invoked as if its a helper function?
func refreshToken(ctx context.Context) (string, error) {
	// need to refresh tokens and try again
	// TODO: we'll probably need a way to stop infinite redirects
	refreshToken := keys.GetContextValue(ctx, keys.ContextSpotifyRefreshToken)
	if refreshToken == nil {
		return "", errors.New("No refresh token provided")
	}

	requestCtx := context.WithValue(ctx, keys.ContextSpotifyRefreshToken, refreshToken)
	newTok, err := spotify.RefreshToken(requestCtx)
	if err != nil {
		logging.GetLogger(ctx).WithError(err).Error("refresh token attempt failed")
		return "", err
	}

	return newTok, nil
}
