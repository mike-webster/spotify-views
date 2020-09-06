package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	data "github.com/mike-webster/spotify-views/data"
	"github.com/mike-webster/spotify-views/genius"
	"github.com/mike-webster/spotify-views/logging"
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
	logger := logging.GetLogger(nil)
	code := c.Query(queryStringCode)
	// TODO: query state verification
	qErr := c.Query(queryStringError)
	if len(qErr) > 0 {
		// the user is a fucker and they denied access
		logger.WithError(errors.New(qErr)).Error("user did not grant access")
		c.Status(500)
		return
	}

	requestCtx := context.Background()
	requestCtx = context.WithValue(requestCtx, spotify.ContextClientID, c.MustGet(string(spotify.ContextClientID)))
	requestCtx = context.WithValue(requestCtx, spotify.ContextClientSecret, c.MustGet(string(spotify.ContextClientSecret)))
	requestCtx = context.WithValue(requestCtx, spotify.ContextReturnURL, c.MustGet(string(spotify.ContextReturnURL)))

	oauthResultCtx, err := spotify.HandleOauth(requestCtx, code)
	if err != nil {
		logger.WithError(err).Error("error handling spotify oauth")
		c.Status(500)
		return
	}

	token := oauthResultCtx.Value(spotify.ContextAccessToken)
	if token == nil {
		logger.WithError(err).Error("no access token returned from spotify")
		c.Status(500)
		return
	}

	requestCtx = context.WithValue(requestCtx, spotify.ContextAccessToken, token)
	userResultCtx, err := spotify.GetUserInfo(requestCtx)
	if err != nil {
		logger.WithError(err).Error("couldnt retrieve userid from spotify")
		c.Status(500)
		return
	}

	ctxInfo := userResultCtx.Value(spotify.ContextResults)
	if ctxInfo == nil {
		logger.Error("no id returned from query")
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
		logger.WithField("info", info).WithError(err).Error("couldnt save user")
		c.Status(500)
		return
	}

	if !success {
		logger.WithField("info", info).Warn("couldnt create user - may have already existed")
	}

	refresh := oauthResultCtx.Value(spotify.ContextRefreshToken)
	if refresh == nil {
		logger.Error("no refresh token returned from spotify")
	} else {
		success, err := data.SaveRefreshTokenForUser(requestCtx, fmt.Sprint(refresh), info["id"])
		if err != nil {
			logger.WithField("info", info).WithError(err).Error("couldnt save refresh token for user")
			c.Status(500)
			return
		}
		if !success {
			logger.WithField("info", info).Warn("refresh token not inserted - may have already existed")
		}
	}

	c.SetCookie(cookieKeyID, fmt.Sprint(info["id"]), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
	c.SetCookie(cookieKeyToken, fmt.Sprint(token), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
	val, err := c.Cookie("redirect_url")
	if err == nil && len(val) > 0 {
		c.Redirect(http.StatusTemporaryRedirect, val)
	}
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(PathTopTracks, "?", queryStringTimeRange, "=short_term"))
}

func handlerTopTracks(c *gin.Context) {
	logger := logging.GetLogger(nil)
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		logger.Debug("no token, redirecting")
		c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathTopTracks)
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
				c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathTopTracks)
				return
			}
			if success {
				c.Redirect(http.StatusTemporaryRedirect, PathTopTracks)
				return
			}
		}

		logger.WithError(err).Error("couldnt retrieve top tracks from spotify")
		c.Status(500)
		return
	}

	reqTracks := tracksResultsCtx.Value(spotify.ContextResults)
	if reqTracks == nil {
		logger.Error("no tracks returned from spotify")
		c.Status(500)
		return
	}

	tracks, ok := reqTracks.(spotify.Tracks)
	if !ok {
		logger.WithField("type", reflect.TypeOf(reqTracks)).Error("couldnt parse tracks returned from spotify")
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
	logger := logging.GetLogger(nil)
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		logger.Debug("no token, redirecting")
		c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathTopArtists)
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
				c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathTopArtists)
				return
			}
			if success {
				c.Redirect(http.StatusTemporaryRedirect, PathTopArtists)
				return
			}
		}

		logger.WithError(err).Error("couldnt retrieve top artists from spotify")
		c.Status(500)
		return
	}

	reqArtists := artistResponseCtx.Value(spotify.ContextResults)
	if reqArtists == nil {
		logger.Error("received no artists")
		c.Status(500)
		return
	}

	artists, ok := reqArtists.(spotify.Artists)
	if !ok {
		logger.WithField("type", reflect.TypeOf(reqArtists)).Error("couldnt parse artists returned from spotify")
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

func handlerUserLibrary(c *gin.Context) {
	logger := logging.GetLogger(nil)
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		logger.Debug("no token, redirecting")
		c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathTopTracks)
		return
	}

	requestCtx := context.WithValue(c, spotify.ContextAccessToken, token)
	t, err := spotify.GetUserLibraryTracks(requestCtx)
	if err != nil {
		panic(err)
	}

	aids := []string{}
	for _, i := range *t{
		aids = append(aids, i.Genre)
	}


	af, err := spotify.GetAudioFeatures(requestCtx, t.IDs())
	if err != nil {
		panic(err)
	}
	laf := *af

	// leaving off:
	// I think I wired this wrong, the API hadnler probably needs to repackage
	// the results before returning them, right?
	// What I still need to do is figure out how to gracefully tie all this together.
	// I have to look up artists because tracks don't have genre
	// I have to look up the audio features for tempo

	type item struct {
		Artist string
		Title  string
		Tempo  float32
		Genre string
	}

	type viewBag struct {
		Items []item
	}

	vb := viewBag{}
	for i := 0; i < len(t); i ++ {
		tr := t[i]
		for j := 0; j < len(*af); j++ {
			if tr.ID == ta.ID {
				vb.Items = append(vb.Items, item{
					Artist: tr.Artists[0].Name,
					Title: tr.Name,
					Tempo: ta.Tempo,
					Genre: ,
				})
			}
		}
	}

	// m := map[string]float32{}
	// for i := 0; i < len(t); i++ {
	// 	tr := t[i]
	// 	for j := 0; j < len(*af); j++ {
	// 		ta := laf[j]
	// 		if tr.ID == ta.ID {
	// 			m[fmt.Sprint(tr.Artists[0].Name, ":", tr.Name)] = ta.Tempo
	// 			break
	// 		}
	// 	}
	// }
	// sm := sortablemap.GetSortableFloatMap(m)
	// dir := c.Query("sort")
	// if dir != "asc" {
	// 	sort.Sort(sort.Reverse(sm))
	// } else {
	// 	sort.Sort(sm)
	// }

	// vb := viewBag{}
	// for _, i := range sm {
	// 	artist := strings.Split(i.Key, ":")[0]
	// 	title := strings.Split(i.Key, ":")[1]
	// 	vb.Items = append(vb.Items, item{Artist: artist, Title: title, Tempo: i.Value})
	// }

	c.HTML(200, "library.tmpl", vb)
}

func handlerUserLibraryTempo(c *gin.Context) {
	logger := logging.GetLogger(nil)
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		logger.Debug("no token, redirecting")
		c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathTopTracks)
		return
	}

	requestCtx := context.WithValue(c, spotify.ContextAccessToken, token)
	t, err := spotify.GetUserLibraryTracks(requestCtx)
	if err != nil {
		panic(err)
	}
	af, err := spotify.GetAudioFeatures(requestCtx, t.IDs())
	if err != nil {
		panic(err)
	}
	laf := *af

	type item struct {
		Artist string
		Title  string
		Tempo  float32
	}

	type viewBag struct {
		Items []item
	}

	m := map[string]float32{}
	for i := 0; i < len(t); i++ {
		tr := t[i]
		for j := 0; j < len(*af); j++ {
			ta := laf[j]
			if tr.ID == ta.ID {
				m[fmt.Sprint(tr.Artists[0].Name, ":", tr.Name)] = ta.Tempo
				break
			}
		}
	}
	sm := sortablemap.GetSortableFloatMap(m)
	dir := c.Query("sort")
	if dir != "asc" {
		sort.Sort(sort.Reverse(sm))
	} else {
		sort.Sort(sm)
	}

	vb := viewBag{}
	for _, i := range sm {
		artist := strings.Split(i.Key, ":")[0]
		title := strings.Split(i.Key, ":")[1]
		vb.Items = append(vb.Items, item{Artist: artist, Title: title, Tempo: i.Value})
	}

	c.HTML(200, "library.tmpl", vb)
}

func handlerTopArtistsGenres(c *gin.Context) {
	logger := logging.GetLogger(nil)
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		logger.Debug("no token, redirecting")
		c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathTopArtistGenres)
		return
	}

	ctx := context.WithValue(c, spotify.ContextAccessToken, token)

	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		ctx = context.WithValue(ctx, spotify.ContextTimeRange, tr)
	}

	reqCtx, err := spotify.GetTopArtists(ctx)
	if err != nil {
		logger.WithError(err).Error("couldnt retrieve top artists from spotify")
		c.Status(500)
		return
	}

	reqArtists := reqCtx.Value(spotify.ContextResults)
	artists, ok := reqArtists.(spotify.Artists)
	if !ok {
		logger.WithField("type", reflect.TypeOf(reqArtists)).Error("couldnt parse artists from spotify")
		c.Status(500)
		return
	}

	reqCtx, err = spotify.GetGenresForArtists(ctx, artists.IDs())
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			success, err := refreshToken(c)
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathTopArtistGenres)
				return
			}
			if success {
				c.Redirect(http.StatusTemporaryRedirect, PathTopArtistGenres)
				return
			}
		}

		logger.WithError(err).Error("couldnt retrieve genres for artists from spotify")
		c.Status(500)
		return
	}

	reqGenres := reqCtx.Value(spotify.ContextResults)
	genres, ok := reqGenres.(spotify.Pairs)
	if !ok {
		logger.WithField("type", reflect.TypeOf(reqGenres)).Error("couldnt parse genres from spotify")
		c.Status(500)
		return
	}

	sort.Sort(sort.Reverse(genres))
	vb := ViewBag{Resource: "artist", Results: genres}
	c.HTML(200, "topgenres.tmpl", vb)
}

func handlerTopTracksGenres(c *gin.Context) {
	logger := logging.GetLogger(nil)
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		logger.Debug("no token, redirecting")
		c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathTopTracksGenres)
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
				c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathTopTracksGenres)
				return
			}
			if success {
				c.Redirect(http.StatusTemporaryRedirect, PathTopTracksGenres)
				return
			}
		}

		logger.WithError(err).Error("couldnt retrieve top tracks from spotify")
		c.Status(500)
		return
	}

	reqTracks := reqCtx.Value(spotify.ContextResults)
	tracks, ok := reqTracks.(spotify.Tracks)
	if !ok {
		logger.WithField("type", reflect.TypeOf(reqTracks)).Error("couldnt parse tracks from spotify")
		c.Status(500)
		return
	}

	reqCtx, err = spotify.GetGenresForTracks(ctx, tracks.IDs())
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			success, err := refreshToken(c)
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathTopTracksGenres)
				return
			}
			if success {
				c.Redirect(http.StatusTemporaryRedirect, PathTopTracksGenres)
				return
			}
		}

		logger.WithError(err).Error("couldnt retreive top genres for top tracks from spotify")
		c.Status(500)
		return
	}

	reqGenres := reqCtx.Value(spotify.ContextResults)
	genres, ok := reqGenres.(spotify.Pairs)
	if !ok {
		logger.WithField("type", reflect.TypeOf(reqGenres)).Error("couldnt parse genres from spotify")
		c.Status(500)
		return
	}

	sort.Sort(sort.Reverse(genres))
	vb := ViewBag{Resource: "track", Results: genres}
	c.HTML(200, "topgenres.tmpl", vb)
}

func handlerWordCloud(c *gin.Context) {
	c.HTML(200, "wordcloud2.tmpl", nil)
}

func handlerWordCloudData(c *gin.Context) {
	logger := logging.GetLogger(nil)
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		logger.Debug("no token, redirecting")
		c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathWordCloud)
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
				c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathWordCloud)
				return
			}
			if success {
				c.Redirect(http.StatusTemporaryRedirect, PathWordCloud)
				return
			}
		}

		logger.WithError(err).Error("couldnt retrieve top tracks from spotify")
		c.Status(500)
		return
	}
	logger.Debug("got top tracks")

	reqTracks := reqCtx.Value(spotify.ContextResults)
	tracks, ok := reqTracks.(spotify.Tracks)
	if !ok {
		logger.WithField("type", reflect.TypeOf(reqTracks)).Error("couldnt parse tracks from spotify")
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
		logger.WithError(err).Error("couldnt retrieve word counts")
		c.Status(500)
		return
	}
	logger.Debug("got word counts")

	sm := sortablemap.GetSortableMap(wordCounts)
	sort.Sort(sort.Reverse(sm))

	retMap := sm.Take(50)

	logger.Debug("map sorted")

	type viewData struct {
		Filename string          `json:"filename"`
		Maps     sortablemap.Map `json:"maps"`
	}
	vb := viewData{Maps: retMap}

	c.JSON(200, vb)
}

func handlerLogin(c *gin.Context) {
	// TODO Add state
	pathScopes := url.QueryEscape(strings.Join(scopes, " "))
	spotifyURL := fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s&show_dialog=false",
		clientID,
		pathScopes,
		c.MustGet(string(spotify.ContextReturnURL)))

	if len(c.Query("redirectUrl")) > 0 {
		c.SetCookie("redirect_url", c.Query("redirectUrl"), 600, "/", strings.Replace(host, "https://", "", -1), false, true)
	}

	logging.GetLogger(nil).Info(fmt.Sprint("redirecting for spotify auth: ", spotifyURL))

	c.Redirect(http.StatusTemporaryRedirect, spotifyURL)
}

func handlerHome(c *gin.Context) {
	c.HTML(200, "home.tmpl", nil)
}
