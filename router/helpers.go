package router

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"image/png"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mike-webster/spotify-views/env"
	"github.com/mike-webster/spotify-views/keys"
	"github.com/mike-webster/spotify-views/logging"
	"github.com/mike-webster/spotify-views/spotify"
	"github.com/psykhi/wordclouds"
	"github.com/sirupsen/logrus"
)

func setCookie(c *gin.Context, key string, val string, secure bool, httpOnly bool) {
	host := "localhost"
	if os.Getenv("GO_ENV") == "production" {
		host = "spotify-views.com"
	}
	c.SetCookie(key, val, 3600, "/", host, secure, httpOnly)
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
	PathSpotifyCodeSwap  = "/token"
	PathSpotifyReturn    = "/oauthreturn"
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
	r.Use(logRequests)
	// if os.Getenv("GO_ENV") != "production" {
	// 	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
	// 		// your custom format
	// 		return fmt.Sprintf("%s - [%s] \n\t\"%s %s %s %d %s \"%s\" %s\"\n",
	// 			param.ClientIP,
	// 			param.TimeStamp.Format(time.RFC1123),
	// 			param.Method,
	// 			param.Path,
	// 			param.Request.Proto,
	// 			param.StatusCode,
	// 			param.Latency,
	// 			param.Request.UserAgent(),
	// 			param.ErrorMessage,
	// 		)
	// 	}))
	// }

	r.GET(PathHome, func(c *gin.Context) {
		// healthcheck
		c.Status(200)
	})

	r.GET(PathSpotifyOauth, handlerOauth) // step 2 - code swap
	r.GET(PathLogin, handlerLogin)        // step 1 - user permission

	if os.Getenv("GO_ENV") != "production" {
		r.GET(PathTest, authenticate, handlerTest)
	}

	// TODOEND

	spot := r.Group("/spotify")
	{
		spot.POST(PathSpotifyCodeSwap, handlerSpotifyCodeSwap)
		spot.GET(PathSpotifyReturn, func(c *gin.Context) {
			lgr := logging.GetLogger(ctx)
			code := c.Query(queryStringCode)

			// TODO: query state verification
			qErr := c.Query(queryStringError)
			if len(qErr) > 0 {
				// the user denied access
				lgr.WithError(errors.New(qErr)).Error("user did not grant access")
				c.Status(500)

				return
			}

			tok, err := spotify.ExchangeOauthCode(c, code)
			if err != nil {
				lgr.WithError(err).Error("error handling spotify oauth")
				c.Status(500)
				return
			}

			c.Set(string(keys.ContextSpotifyAccessToken), tok.Access)
			c.Set(string(keys.ContextSpotifyRefreshToken), tok.Refresh)

			if len(tok.Access) < 1 {
				lgr.WithError(err).Error("no access token returned from spotify")
				c.Status(500)
				return
			}

			u, err := spotify.GetUser(c)
			if err != nil {
				lgr.WithError(err).Error("couldnt retrieve userid from spotify")
				c.Status(500)
				return
			}

			err = u.Save(c)
			if err != nil {
				lgr.WithField("info", *u).WithError(err).Error("couldnt save user")
				c.Status(500)
				return
			}

			if len(tok.Refresh) < 1 {
				lgr.Error("no refresh token returned from spotify")
			}

			lgr.WithFields(logrus.Fields{
				"event": "user_login",
				"id":    u.ID,
				"email": u.Email,
			}).Info("user logged in successfully")

			setCookie(c, cookieKeyID, fmt.Sprint(u.ID), false, false)
			setCookie(c, cookieKeyToken, fmt.Sprint(tok.Access), false, false)
			setCookie(c, cookieKeyRefresh, fmt.Sprint(tok.Refresh), false, false)

			host := c.Request.Referer()
			red, _ := c.Cookie("redirect_url")
			if len(red) > 0 {
				host = fmt.Sprint(host, red)
			}
			c.Redirect(http.StatusTemporaryRedirect, host)
			return
		})
	}

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

// what the fudge is this? Why is it taking a gin context like a controller handler but
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
