package main

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"image/png"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mike-webster/spotify-views/data"
	"github.com/mike-webster/spotify-views/keys"
	"github.com/mike-webster/spotify-views/logging"
	"github.com/mike-webster/spotify-views/spotify"
	"github.com/psykhi/wordclouds"
)

// TODO: put together a yaml config to parse these for me
// TODO: these context keys should all be the same type, not unique per package. abstract
//       the keys package
func parseEnvironmentVariables(ctx context.Context) (map[interface{}]interface{}, error) {
	ret := map[interface{}]interface{}{}
	clientID = os.Getenv("CLIENT_ID")
	if len(clientID) < 1 {
		return nil, errors.New("no client id provided")
	}
	ret[keys.ContextSpotifyClientID] = clientID

	clientSecret = os.Getenv("CLIENT_SECRET")
	if len(clientSecret) < 1 {
		return nil, errors.New("no client secret provided")
	}
	ret[keys.ContextSpotifyClientSecret] = clientSecret

	host = os.Getenv("HOST")
	if len(host) < 1 {
		return nil, errors.New("no host provided")
	}
	ret[keys.ContextSpotifyReturnURL] = fmt.Sprint("https://www.", host, "/spotify/oauth")

	// TODO: Do we need this in the context? or just set for the main package?
	// consider: the main goal here is to be able to verify everything is working
	// on app start using the context returned from this method.
	lyricsKey = os.Getenv("LYRICS_KEY")
	if len(lyricsKey) < 1 {
		return nil, errors.New("no lyrics key provided")
	}
	ret[keys.ContextLyricsToken] = lyricsKey

	dbHost = os.Getenv("DB_HOST")
	if len(dbHost) < 1 {
		return nil, errors.New("no db host provided")
	}
	ret[keys.ContextDbHost] = dbHost

	dbUser = os.Getenv("DB_USER")
	if len(dbUser) < 1 {
		return nil, errors.New("no db user provided")
	}
	ret[keys.ContextDbUser] = dbUser

	dbPass = os.Getenv("DB_PASS")
	if len(dbPass) < 1 {
		return nil, errors.New("no db pass provided")
	}
	ret[keys.ContextDbPass] = dbPass

	dbName = os.Getenv("DB_NAME")
	if len(dbName) < 1 {
		return nil, errors.New("no db name provided")
	}
	ret[keys.ContextDatabase] = dbName

	secKey = os.Getenv("SEC_KEY")
	if len(secKey) < 1 {
		return nil, errors.New("no sec key provided")
	}
	ret[keys.ContextSecurityKey] = secKey

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
)

func runServer() {
	r := gin.New()
	r.Use(requestLogger())
	r.Use(recovery)
	r.Use(parseUserID)
	r.Use(redisClient)
	r.Use(loadContextValues())
	r.LoadHTMLGlob("templates/*")
	r.GET(PathSpotifyOauth, handlerOauth)
	r.GET(PathTopTracks, handlerTopTracks)
	r.GET(PathTopArtists, handlerTopArtists)
	r.GET(PathTopArtistGenres, handlerTopArtistsGenres)
	r.GET(PathTopTracksGenres, handlerTopTracksGenres)
	r.GET(PathLogin, handlerLogin)
	r.GET(PathHome, handlerHome)
	r.GET(PathWordCloud, handlerWordCloud)
	r.GET(PathWordCloudData, handlerWordCloudData)
	r.GET(PathUserLibraryTempo, handlerUserLibraryTempo)
	r.GET(PathRecommendations, handlerRecommendations)

	r.StaticFile("/sitemap", "./static/sitemap.xml")
	r.Static("/static/css", "./static")
	r.Static("/static/js", "./static")

	r.Run()
}

func refreshToken(c *gin.Context) (bool, error) {
	// need to refresh tokens and try again
	// TODO: we'll probably need a way to stop infinite redirects
	logger := logging.GetLogger(c)

	spotifyID, err := c.Cookie(cookieKeyID)
	if err != nil {
		logger.WithError(err).Error("no userid stored")
		return false, err
	}

	refreshToken, err := data.GetRefreshTokenForUser(c, spotifyID)
	if err != nil {
		logger.WithError(err).Error("no refresh token stored")
		return false, err
	}

	requestCtx := context.WithValue(c, keys.ContextSpotifyRefreshToken, refreshToken)
	refrshResponseCtx, err := spotify.RefreshToken(requestCtx)
	if err != nil {
		logger.WithError(err).Error("refresh token attempt failed")
		return false, err
	}

	refTok := keys.GetContextValue(refrshResponseCtx, keys.ContextSpotifyResults)
	if refTok == nil {
		logger.WithError(err).Error("no token returned from refresh attempt")
		return false, err
	}

	c.SetCookie(cookieKeyToken, fmt.Sprint(refTok), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
	c.Redirect(http.StatusTemporaryRedirect, PathTopTracks)
	return true, nil
}
