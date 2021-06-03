package router

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"image/png"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mike-webster/spotify-views/env"
	"github.com/mike-webster/spotify-views/keys"
	"github.com/mike-webster/spotify-views/logging"
	"github.com/mike-webster/spotify-views/spotify"
	"github.com/psykhi/wordclouds"
)

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
		wordclouds.FontFile("web/fonts/Ubuntu-L.ttf"),
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
	directory := "web/clouds/"
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
	PathCombinedGenres   = "/genres"
	PathLogin            = "/login"
	PathHome             = "/"
	PathWordCloud        = "/wordcloud"
	PathWordCloudData    = "/wordcloud/data"
	PathUserLibraryTempo = "/library/tempo"
	PathRecommendations  = "/tracks/recommendations"
	PathTest             = "/test"
)

func Run(ctx context.Context) {
	ctx = context.WithValue(ctx, keys.ContextMasterKey, os.Getenv("MASTER_KEY"))

	if os.Getenv("GO_ENV") == "test" {
		//testMethod(ctx)
		return
	}

	runServer(ctx)
}

func runServer(ctx context.Context) {
	r := gin.New()
	r.Use(recovery)
	r.Use(setContextLogger)
	r.Use(setTokens)
	r.Use(setSpotifyUserID)
	r.Use(setEnv)
	r.Use(setDependencies)
	r.Use(CORSMiddleware)

	// TODO: remove this once we're done migrating to the react app
	r.StaticFile("/sitemap", "./web/sitemap.xml")
	r.Static("/web/css", "./web")
	r.Static("/web/js", "./web")
	r.Static("/logos/", "./web/logos")
	r.Static("/images/", "./web/images")
	r.StaticFile("/favicon.ico", "./web/logos/favicon.ico")
	r.LoadHTMLGlob("web/templates/*")

	r.GET(PathHome, handlerHome)
	r.GET(PathSpotifyOauth, handlerOauth) // leave this one
	r.GET(PathLogin, handlerLogin)        // remove

	//r.GET(PathTopTracks, authenticate, handlerTopTracks)
	r.GET(PathWordCloud, authenticate, handlerWordCloud)
	r.GET(PathUserLibraryTempo, authenticate, handlerUserLibraryTempo)
	r.GET(PathRecommendations, authenticate, handlerRecommendations) // remove
	r.GET(PathTest, authenticate, handlerTest)
	// TODOEND

	api := r.Group("/api/v1")
	{
		api.GET(PathLogin, handlerLogin)
		api.GET(PathRecommendations, handlerRecommendations)
		api.GET(PathTopTracks, authenticate, handlerTopTracks)
		api.GET(PathTopArtists, authenticate, handlerTopArtists)
		// api.GET(PathTopArtistGenres, authenticate, handlerTopArtistsGenres)
		// api.GET(PathTopTracksGenres, authenticate, handlerTopTracksGenres)
		api.GET(PathCombinedGenres, authenticate, handlerCombinedGenres)
		api.GET(PathWordCloudData, authenticate, handlerWordCloudData)
	}

	env, err := env.ParseEnv()
	if err != nil {
		panic(err)
	}

	r.Run(fmt.Sprint(":", env.Port))
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
	tok := spotify.Token{
		Access:  fmt.Sprint(keys.GetContextValue(ctx, keys.ContextSpotifyAccessToken)),
		Refresh: fmt.Sprint(keys.GetContextValue(ctx, keys.ContextSpotifyRefreshToken)),
	}
	success, err := tok.RefreshMe(requestCtx)
	if err != nil {
		logging.GetLogger(ctx).WithError(err).Error("refresh token attempt failed")
		return "", err
	}

	if !success {
		return "", errors.New("token refresh unsuccessful")
	}

	return tok.Access, nil
}
