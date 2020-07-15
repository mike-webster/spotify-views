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
	"github.com/mike-webster/spotify-views/genius"
	"github.com/mike-webster/spotify-views/logging"
	"github.com/mike-webster/spotify-views/spotify"
	"github.com/psykhi/wordclouds"
)

func parseEnvironmentVariables(ctx context.Context) (context.Context, error) {
	clientID = os.Getenv("CLIENT_ID")
	if len(clientID) < 1 {
		return nil, errors.New("no client id provided")
	}
	ctx = context.WithValue(ctx, spotify.ContextClientID, clientID)

	clientSecret = os.Getenv("CLIENT_SECRET")
	if len(clientSecret) < 1 {
		return nil, errors.New("no client secret provided")
	}
	ctx = context.WithValue(ctx, spotify.ContextClientSecret, clientSecret)

	host = os.Getenv("HOST")
	if len(host) < 1 {
		return nil, errors.New("no host provided")
	}
	ctx = context.WithValue(ctx, spotify.ContextReturnURL, fmt.Sprint("https://", host, "/spotify/oauth"))
	// TODO: Do we need this in the context? or just set for the main package?
	// consider: the main goal here is to be able to verify everything is working
	// on app start using the context returned from this method.

	lyricsKey = os.Getenv("LYRICS_KEY")
	if len(lyricsKey) < 1 {
		return nil, errors.New("no lyrics key provided")
	}
	ctx = context.WithValue(ctx, genius.ContextAccessToken, lyricsKey)

	dbHost = os.Getenv("DB_HOST")
	if len(dbHost) < 1 {
		return nil, errors.New("no db host provided")
	}
	ctx = context.WithValue(ctx, data.ContextHost, dbHost)

	dbUser = os.Getenv("DB_USER")
	if len(dbUser) < 1 {
		return nil, errors.New("no db user provided")
	}
	ctx = context.WithValue(ctx, data.ContextUser, dbUser)

	dbPass = os.Getenv("DB_PASS")
	if len(dbPass) < 1 {
		return nil, errors.New("no db pass provided")
	}
	ctx = context.WithValue(ctx, data.ContextPass, dbPass)

	dbName = os.Getenv("DB_NAME")
	if len(dbName) < 1 {
		return nil, errors.New("no db name provided")
	}
	ctx = context.WithValue(ctx, data.ContextDatabase, dbName)

	secKey = os.Getenv("SEC_KEY")
	if len(secKey) < 1 {
		return nil, errors.New("no sec key provided")
	}
	ctx = context.WithValue(ctx, data.ContextSecurityKey, secKey)
	ctx = context.WithValue(ctx, logging.ContextLogger, logging.GetLogger(nil))
	return ctx, nil
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
	PathSpotifyOauth    = "/spotify/oauth"
	PathTopTracks       = "/tracks/top"
	PathTopArtists      = "/artists/top"
	PathTopArtistGenres = "/artists/genres"
	PathTopTracksGenres = "/tracks/genres"
	PathLogin           = "/login"
	PathHome            = "/"
	PathWordCloud       = "/wordcloud"
	PathWordCloudData   = "/wordcloud/data"
)

func runServer() {
	r := gin.New()
	r.Use(requestLogger())
	r.Use(recovery)
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

	r.Static("/static/css", "./static")
	r.Static("/static/js", "./static")
	r.Static("/sitemap.xml", "./sitemap.xml")

	r.Run()
}

func refreshToken(c *gin.Context) (bool, error) {
	// need to refresh tokens and try again
	// TODO: we'll probably need a way to stop infinite redirects
	requestCtx := context.WithValue(context.Background(), data.ContextDatabase, c.MustGet(string(data.ContextDatabase)))
	requestCtx = context.WithValue(requestCtx, data.ContextHost, c.MustGet(string(data.ContextHost)))
	requestCtx = context.WithValue(requestCtx, data.ContextPass, c.MustGet(string(data.ContextPass)))
	requestCtx = context.WithValue(requestCtx, data.ContextSecurityKey, c.MustGet(string(data.ContextSecurityKey)))
	requestCtx = context.WithValue(requestCtx, data.ContextUser, c.MustGet(string(data.ContextUser)))
	logger := logging.GetLogger(nil)

	spotifyID, err := c.Cookie(cookieKeyID)
	if err != nil {
		logger.WithError(err).Error("no userid stored")
		return false, err
	}

	refreshToken, err := data.GetRefreshTokenForUser(requestCtx, spotifyID)
	if err != nil {
		logger.WithError(err).Error("no refresh token stored")
		return false, err
	}

	requestCtx = context.WithValue(requestCtx, spotify.ContextRefreshToken, refreshToken)
	refrshResponseCtx, err := spotify.RefreshToken(requestCtx)
	if err != nil {
		logger.WithError(err).Error("refresh token attempt failed")
		return false, err
	}

	refTok := refrshResponseCtx.Value(spotify.ContextResults)
	if refTok == nil {
		logger.WithError(err).Error("no token returned from refresh attempt")
		return false, err
	}

	c.SetCookie(cookieKeyToken, fmt.Sprint(refTok), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
	c.Redirect(http.StatusTemporaryRedirect, PathTopTracks)
	return true, nil
}
