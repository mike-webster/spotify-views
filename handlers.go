package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	data "github.com/mike-webster/spotify-views/data"
	genius "github.com/mike-webster/spotify-views/genius"
	sortablemap "github.com/mike-webster/spotify-views/sortablemap"
	spotify "github.com/mike-webster/spotify-views/spotify"
)

var (
	queryStringCode                = "code"
	queryStringError               = "error"
	queryStringTimeRange           = "time_range"
	cookieKeyToken                 = "svauth"
	cookieKeyID                    = "svid"
	topTracksLimit           int32 = 25
	topGenresTopTracksLimit  int32 = 50
	wordCloudTopTracksLimit  int32 = 50
	spotifyPlayerHeightShort int32 = 80
	spotifyPlayerHeightTall  int32 = 380
	spotifyPlayerWidth       int32 = 300
)

func handlerOauth(c *gin.Context) {
	code := c.Query(queryStringCode)
	// TODO: query state verification
	qErr := c.Query(queryStringError)
	if len(qErr) > 0 {
		// the user is a fucker and they denied access
		log.Println("user did not grant access: ", qErr)
		c.Status(500)
		return
	}

	requestCtx := context.Background()
	requestCtx = context.WithValue(requestCtx, spotify.ContextClientID, c.MustGet(string(spotify.ContextClientID)))
	requestCtx = context.WithValue(requestCtx, spotify.ContextClientSecret, c.MustGet(string(spotify.ContextClientSecret)))
	requestCtx = context.WithValue(requestCtx, spotify.ContextReturnURL, c.MustGet(string(spotify.ContextReturnURL)))

	oauthResultCtx, err := spotify.HandleOauth(requestCtx, code)
	if err != nil {
		log.Println("error handling spotify oauth: ", err.Error())
		c.Status(500)
		return
	}

	token := oauthResultCtx.Value(spotify.ContextAccessToken)
	if token == nil {
		log.Println("no token returned from spotify")
		c.Status(500)
		return
	}

	requestCtx = context.WithValue(requestCtx, spotify.ContextAccessToken, token)
	userResultCtx, err := spotify.GetUserInfo(requestCtx)
	if err != nil {
		log.Println("couldnt retrieve userid from spotify; err: ", err.Error())
		c.Status(500)
		return
	}

	ctxInfo := userResultCtx.Value(spotify.ContextResults)
	if ctxInfo == nil {
		log.Println("no id returned from query")
		c.Status(500)
		return
	}

	info := ctxInfo.(map[string]string)

	requestCtx = context.WithValue(requestCtx, data.ContextDatabase, c.MustGet(string(data.ContextDatabase)))
	requestCtx = context.WithValue(requestCtx, data.ContextHost, c.MustGet(string(data.ContextHost)))
	requestCtx = context.WithValue(requestCtx, data.ContextPass, c.MustGet(string(data.ContextPass)))
	requestCtx = context.WithValue(requestCtx, data.ContextSecurityKey, c.MustGet(string(data.ContextSecurityKey)))
	requestCtx = context.WithValue(requestCtx, data.ContextUser, c.MustGet(string(data.ContextUser)))

	success, err := data.SaveUser(requestCtx, info["id"], info["email"])
	if err != nil {
		log.Println("couldnt save user: ", info, "; ", err.Error())
		c.Status(500)
		return
	}

	if !success {
		log.Println("couldnt create user - may have already existed")
	}

	refresh := oauthResultCtx.Value(spotify.ContextRefreshToken)
	if refresh == nil {
		log.Println("no refresh returned from spotify")
	} else {
		success, err := data.SaveRefreshTokenForUser(requestCtx, fmt.Sprint(refresh), info["id"])
		if err != nil {
			log.Println("couldnt save refresh token for user; err: ", err.Error())
			c.Status(500)
			return
		}
		if !success {
			log.Println("token not inserted - may have already existed")
		}
	}

	c.SetCookie(cookieKeyID, fmt.Sprint(info["id"]), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
	c.SetCookie(cookieKeyToken, fmt.Sprint(token), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(PathTopTracks, "?", queryStringTimeRange, "=short_term"))
}

func handlerTopTracks(c *gin.Context) {
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		log.Println("no token - redirecting to login")
		c.Redirect(http.StatusTemporaryRedirect, PathLogin)
		return
	}

	requestCtx := context.WithValue(c, spotify.ContextAccessToken, token)
	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		requestCtx = context.WithValue(requestCtx, spotify.ContextTimeRange, tr)
	}

	tracksResultsCtx, err := spotify.GetTopTracks(requestCtx, topTracksLimit)
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			success, err := refreshToken(c)
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, PathLogin)
				return
			}
			if success {
				c.Redirect(http.StatusTemporaryRedirect, PathTopTracks)
				return
			}
		}

		log.Println("couldnt retrieve top tracks from spotify: ", err.Error())
		c.Status(500)
		return
	}

	reqTracks := tracksResultsCtx.Value(spotify.ContextResults)
	if reqTracks == nil {
		log.Println("no tracks returned")
		c.Status(500)
		return
	}

	tracks, ok := reqTracks.(spotify.Tracks)
	if !ok {
		log.Println("couldnt parse tracks returned from spotify; ", reflect.TypeOf(reqTracks))
		c.Status(500)
		return
	}

	type tempBag struct {
		Width    int32
		Height   int32
		ID       string
		Name     string
		Resource string
	}

	data := []tempBag{}

	for _, i := range tracks {
		data = append(data, tempBag{
			ID:       i.ID,
			Height:   spotifyPlayerHeightShort,
			Width:    spotifyPlayerWidth,
			Resource: "track",
			Name:     i.Name,
		})
	}

	c.HTML(200, "toptracks.tmpl", data)
}

func handlerTopArtists(c *gin.Context) {
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		log.Println("no token - redirecting to login")
		c.Redirect(http.StatusTemporaryRedirect, PathLogin)
		return
	}

	requestCtx := context.WithValue(c, spotify.ContextAccessToken, token)
	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		requestCtx = context.WithValue(requestCtx, spotify.ContextTimeRange, tr)
	}

	artistResponseCtx, err := spotify.GetTopArtists(requestCtx)
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			success, err := refreshToken(c)
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, PathLogin)
				return
			}
			if success {
				c.Redirect(http.StatusTemporaryRedirect, PathTopArtists)
				return
			}
		}

		log.Println("couldnt retrieve top artists from spotify: ", err.Error())
		c.Status(500)
		return
	}

	reqArtists := artistResponseCtx.Value(spotify.ContextResults)
	if reqArtists == nil {
		log.Println("received no artists")
		c.Status(500)
		return
	}

	artists, ok := reqArtists.(spotify.Artists)
	if !ok {
		log.Println("couldnt parse artists")
		c.Status(500)
		return
	}

	type tempBag struct {
		Width    int32
		Height   int32
		ID       string
		Name     string
		Resource string
	}

	data := []tempBag{}

	for _, i := range artists {
		data = append(data, tempBag{
			ID:       i.ID,
			Height:   380,
			Width:    300,
			Resource: "artist",
			Name:     i.Name,
		})
	}

	c.HTML(200, "topartists.tmpl", data)
}

func handlerTopArtistsGenres(c *gin.Context) {
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		log.Println("no token - redirecting to login")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	ctx := context.WithValue(c, spotify.ContextAccessToken, token)

	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		ctx = context.WithValue(ctx, spotify.ContextTimeRange, tr)
	}

	reqCtx, err := spotify.GetTopArtists(ctx)
	if err != nil {
		log.Println("couldnt retrieve top artists from spotify: ", err.Error())
		c.Status(500)
		return
	}

	reqArtists := reqCtx.Value(spotify.ContextResults)
	if reqArtists == nil {
		log.Println("received no artists")
		c.Status(500)
		return
	}

	artists, ok := reqArtists.(spotify.Artists)
	if !ok {
		log.Println("couldnt parse artists")
		c.Status(500)
		return
	}

	reqCtx, err = spotify.GetGenresForArtists(ctx, artists.IDs())
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			success, err := refreshToken(c)
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, PathLogin)
				return
			}
			if success {
				c.Redirect(http.StatusTemporaryRedirect, PathTopArtistGenres)
				return
			}
		}

		log.Println("couldnt retrieve genres for artists from spotify: ", err.Error())
		c.Status(500)
		return
	}

	reqGenres := reqCtx.Value(spotify.ContextResults)
	if reqGenres == nil {
		log.Println("received no genres")
		c.Status(500)
		return
	}

	genres, ok := reqGenres.(spotify.Pairs)
	if !ok {
		log.Println("couldnt parse genres")
		c.Status(500)
		return
	}

	sort.Sort(sort.Reverse(genres))
	vb := ViewBag{Resource: "artist", Results: genres}
	c.HTML(200, "topgenres.tmpl", vb)
}

func handlerTopTracksGenres(c *gin.Context) {
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		log.Println("no token - redirecting to login")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	ctx := context.WithValue(c, spotify.ContextAccessToken, token)

	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		ctx = context.WithValue(ctx, spotify.ContextTimeRange, tr)
	}

	reqCtx, err := spotify.GetTopTracks(ctx, topGenresTopTracksLimit)
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			success, err := refreshToken(c)
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, PathLogin)
				return
			}
			if success {
				c.Redirect(http.StatusTemporaryRedirect, PathTopTracksGenres)
				return
			}
		}

		log.Println("couldnt retrieve top tracks from spotify: ", err.Error())
		c.Status(500)
		return
	}

	reqTracks := reqCtx.Value(spotify.ContextResults)
	if reqTracks == nil {
		log.Println("received no tracks")
		c.Status(500)
		return
	}

	tracks, ok := reqTracks.(spotify.Tracks)
	if !ok {
		log.Println("couldnt parse tracks")
		c.Status(500)
		return
	}

	reqCtx, err = spotify.GetGenresForTracks(ctx, tracks.IDs())
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			success, err := refreshToken(c)
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, PathLogin)
				return
			}
			if success {
				c.Redirect(http.StatusTemporaryRedirect, PathTopTracksGenres)
				return
			}
		}

		log.Println("couldnt retrieve top genres for top tracks from spotify: ", err.Error())
		c.Status(500)
		return
	}

	reqGenres := reqCtx.Value(spotify.ContextResults)
	if reqGenres == nil {
		log.Println("receive no genres")
		c.Status(500)
		return
	}

	genres, ok := reqGenres.(spotify.Pairs)
	if !ok {
		log.Println("couldnt parse genres")
		c.Status(500)
		return
	}

	sort.Sort(sort.Reverse(genres))
	vb := ViewBag{Resource: "track", Results: genres}
	c.HTML(200, "topgenres.tmpl", vb)
}

func handlerWordCloud(c *gin.Context) {
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		log.Println("no token - redirecting to login")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	ctx := context.WithValue(c, spotify.ContextAccessToken, token)

	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		ctx = context.WithValue(ctx, spotify.ContextTimeRange, tr)
	}

	reqCtx, err := spotify.GetTopTracks(ctx, wordCloudTopTracksLimit)
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			success, err := refreshToken(c)
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, PathLogin)
				return
			}
			if success {
				c.Redirect(http.StatusTemporaryRedirect, PathWordCloud)
				return
			}
		}

		log.Println("couldnt retrieve top tracks from spotify: ", err.Error())
		c.Status(500)
		return
	}
	log.Println("got top tracks")

	reqTracks := reqCtx.Value(spotify.ContextResults)
	if reqTracks == nil {
		log.Println("received no tracks")
		c.Status(500)
		return
	}

	tracks, ok := reqTracks.(spotify.Tracks)
	if !ok {
		log.Println("couldnt parse tracks")
		c.Status(500)
		return
	}

	searches := []genius.LyricSearch{}
	for _, i := range tracks {
		searches = append(searches, genius.LyricSearch{
			Artist: i.Artists[0].Name,
			Track:  i.Name,
		})
	}

	ctx = context.WithValue(c, genius.ContextAccessToken, lyricsKey)

	wordCounts, err := genius.GetLyricCountForSong(ctx, searches)
	if err != nil {
		log.Println("couldn't retrieve word counts: ", err.Error())
		c.Status(500)
		return
	}
	log.Println("got word counts")

	sm := sortablemap.GetSortableMap(wordCounts)
	sort.Sort(sort.Reverse(sm))

	log.Println("map sorted")

	filename := fmt.Sprint(time.Now().Unix(), ".png")
	err = generateWordCloud(ctx, filename, wordCounts)

	log.Println("word cloud generated")
	if err != nil {
		log.Println("couldn't generate error: ", err.Error())
		c.Status(500)
		return
	}

	type viewData struct {
		Filename string
		Maps     sortablemap.Map
	}

	vb := viewData{Filename: filename, Maps: sm}
	c.HTML(200, "wordcloud.tmpl", vb)
}

func handlerLogin(c *gin.Context) {
	// TODO Add state
	pathScopes := url.QueryEscape(strings.Join(scopes, " "))
	spotifyURL := fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s&show_dialog=false",
		clientID,
		pathScopes,
		c.MustGet(string(spotify.ContextReturnURL)))
	c.Redirect(http.StatusTemporaryRedirect, spotifyURL)
}

func handlerHome(c *gin.Context) {
	c.HTML(200, "home.tmpl", nil)
}
