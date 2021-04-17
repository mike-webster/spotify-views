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
	"github.com/mike-webster/spotify-views/keys"
	"github.com/mike-webster/spotify-views/logging"
	sortablemap "github.com/mike-webster/spotify-views/sortablemap"
	spotify "github.com/mike-webster/spotify-views/spotify"
	"github.com/sirupsen/logrus"
)

var (
	queryStringCode                 = "code"
	queryStringError                = "error"
	queryStringTimeRange            = "time_range"
	cookieKeyToken                  = "svauth"
	cookieKeyID                     = "svid"
	cookieKeyRefresh                = "svref"
	keyArtistInfo            string = "artist-cache"
	topTracksLimit           int32  = 25
	topGenresTopTracksLimit  int32  = 50
	wordCloudTopTracksLimit  int32  = 50
	spotifyPlayerHeightShort int32  = 80
	spotifyPlayerHeightTall  int32  = 380
	spotifyPlayerWidth       int32  = 300

	ddlOpts = map[string]string{
		"Recent":         "short_term",
		"In Between":     "medium_term",
		"Going Way Back": "long_term",
	}
)

func handlerOauth(c *gin.Context) {
	logger := logging.GetLogger(c)
	code := c.Query(queryStringCode)
	// TODO: query state verification
	qErr := c.Query(queryStringError)
	if len(qErr) > 0 {
		// the user is a fucker and they denied access
		logger.WithError(errors.New(qErr)).Error("user did not grant access")
		c.Status(500)
		return
	}

	accessTok, refreshTok, err := spotify.HandleOauth(c, code)
	if err != nil {
		logger.WithError(err).Error("error handling spotify oauth")
		c.Status(500)
		return
	}

	c.Set(string(keys.ContextSpotifyAccessToken), fmt.Sprint(accessTok))
	c.Set(string(keys.ContextSpotifyRefreshToken), fmt.Sprint(refreshTok))

	if len(accessTok) < 1 {
		logger.WithError(err).Error("no access token returned from spotify")
		c.Status(500)
		return
	}

	userResultCtx, err := spotify.GetUserInfo(c)
	if err != nil {
		logger.WithError(err).Error("couldnt retrieve userid from spotify")
		c.Status(500)
		return
	}

	ctxResults := keys.GetContextValue(userResultCtx, keys.ContextSpotifyResults)
	if ctxResults == nil {
		logger.Error("no id returned from query")
		c.Status(500)
		return
	}

	info := ctxResults.(map[string]string)

	_, err = data.SaveUser(c, info["id"], info["email"])
	if err != nil {
		logger.WithField("info", info).WithError(err).Error("couldnt save user")
		c.Status(500)
		return
	}

	if len(refreshTok) < 1 {
		logger.Error("no refresh token returned from spotify")
	}

	logger.WithFields(logrus.Fields{
		"event": "user_login",
		"id": info["id"],
		"email": info["email"],
	}).Info("user logged in successfully")

	c.SetCookie(cookieKeyID, fmt.Sprint(info["id"]), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
	c.SetCookie(cookieKeyToken, fmt.Sprint(accessTok), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
	c.SetCookie(cookieKeyRefresh, fmt.Sprint(refreshTok), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
	val, err := c.Cookie("redirect_url")
	if err == nil && len(val) > 0 {
		c.Redirect(http.StatusTemporaryRedirect, val)
	}
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(PathTopTracks, "?", queryStringTimeRange, "=short_term"))
}

func handlerTopTracks(c *gin.Context) {
	logger := logging.GetLogger(c)
	logger.Debug("loading user's top tracks")

	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		mv := ddlOpts[tr]
		c.Set(string(keys.ContextSpotifyTimeRange), mv)
	}

	tracksResultsCtx, err := spotify.GetTopTracks(c, topTracksLimit)
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			// If the error is due to the token being expired, we will have automatically attempted
			// to get a refresh token for the user.  If  that was successful, it will be returned
			// as the error value.  Set the cookie to the new value and redirect the user back to the
			// same path  to start the process again with the new token.
			if len(err.Error()) > 0 {
				c.SetCookie(cookieKeyToken, fmt.Sprint(err.Error()), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
				c.Redirect(http.StatusTemporaryRedirect, PathTopTracks)
				return
			}

			logging.GetLogger(c).Info("couldnt refresh token for user")
		}

		logger.WithError(err).Error("couldnt retrieve top tracks from spotify")
		c.Status(500)
		return
	}

	reqTracks := keys.GetContextValue(tracksResultsCtx, keys.ContextSpotifyResults)
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

	type Result struct {
		Key        string
		Value      string
		Background string
		Width      int
		Height     int
	}

	type tempBag struct {
		Category string
		Type     string
		Opts     []string
		Results  []Result
	}

	r := []Result{}
	for _, i := range tracks {
		r = append(r, Result{
			Key:        i.FindArtist(),
			Value:      i.Name,
			Background: i.FindImage().URL,
			Height:     i.FindImage().Height,
			Width:      i.FindImage().Width,
		})
	}
	data := tempBag{
		Category: "Tracks",
		Type:     "",
		Opts:     []string{"Recent", "In Between", "Going Way Back"},
		Results:  r,
	}

	c.HTML(200, "newtops.tmpl", data)
}

func handlerTopArtists(c *gin.Context) {
	logger := logging.GetLogger(c)

	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		mv := ddlOpts[tr]
		c.Set(string(keys.ContextSpotifyTimeRange), mv)
	}

	artists, err := spotify.GetTopArtists(c)
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			// If the error is due to the token being expired, we will have automatically attempted
			// to get a refresh token for the user.  If  that was successful, it will be returned
			// as the error value.  Set the cookie to the new value and redirect the user back to the
			// same path  to start the process again with the new token.
			if len(err.Error()) > 0 {
				c.SetCookie(cookieKeyToken, fmt.Sprint(err.Error()), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
				c.Redirect(http.StatusTemporaryRedirect, PathTopArtists)
				return
			}

			logging.GetLogger(c).Info("couldnt refresh token for user")
		}

		logger.WithError(err).Error("couldnt retrieve top artists from spotify")
		c.Status(500)
		return
	}

	type Result struct {
		Key        string
		Value      string
		Background string
		Width      int
		Height     int
	}

	type tempBag struct {
		Category string
		Type     string
		Opts     []string
		Results  []Result
	}

	r := []Result{}
	for _, i := range *artists {
		r = append(r, Result{
			Key:        "",
			Value:      i.Name,
			Background: i.FindImage().URL,
			Height:     i.FindImage().Height,
			Width:      i.FindImage().Width,
		})
	}
	data := tempBag{
		Category: "Artists",
		Type:     "",
		Opts:     []string{"Recent", "In Between", "Going Way Back"},
		Results:  r,
	}

	c.HTML(200, "newtops.tmpl", data)
}

func handlerUserLibraryTempo(c *gin.Context) {
	logging.GetLogger(c).WithField("event", "webby_test").Debug()
	t, err := spotify.GetUserLibraryTracks(c)
	if err != nil {
		logging.GetLogger(c).WithError(err).Error()
		c.Status(500)
		return
	}

	af, err := spotify.GetAudioFeatures(c, t.IDs())
	if err != nil {
		c.Status(500)
		return
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
	logger := logging.GetLogger(c)
	var ctx context.Context

	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		mv := ddlOpts[tr]
		c.Set(string(keys.ContextSpotifyTimeRange), mv)
	}

	artists, err := spotify.GetTopArtists(ctx)
	if err != nil {
		logger.WithError(err).Error("couldnt retrieve top artists from spotify")
		c.Status(500)
		return
	}

	reqCtx, err := spotify.GetGenresForArtists(ctx, artists.IDs())
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			// If the error is due to the token being expired, we will have automatically attempted
			// to get a refresh token for the user.  If  that was successful, it will be returned
			// as the error value.  Set the cookie to the new value and redirect the user back to the
			// same path  to start the process again with the new token.
			if len(err.Error()) > 0 {
				c.SetCookie(cookieKeyToken, fmt.Sprint(err.Error()), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
				c.Redirect(http.StatusTemporaryRedirect, PathTopArtistGenres)
				return
			}

			logging.GetLogger(c).Info("couldnt refresh token for user")
		}

		logger.WithError(err).Error("couldnt retrieve genres for artists from spotify")
		c.Status(500)
		return
	}

	reqGenres := keys.GetContextValue(reqCtx, keys.ContextSpotifyResults)
	genres, ok := reqGenres.(spotify.Pairs)
	if !ok {
		logger.WithField("type", reflect.TypeOf(reqGenres)).Error("couldnt parse genres from spotify")
		c.Status(500)
		return
	}

	type Result struct {
		Key        string
		Value      string
		Background string
		Width      int
		Height     int
	}

	type tempBag struct {
		Category string
		Type     string
		Opts     []string
		Results  []Result
	}

	r := []Result{}
	sort.Sort(sort.Reverse(genres))
	for _, i := range genres {
		r = append(r, Result{
			Key:   i.Key,
			Value: fmt.Sprint("( ", i.Value, ")"),
		})
	}
	data := tempBag{
		Category: "Genres",
		Type:     "artists",
		Opts:     []string{"Recent", "In Between", "Going Way Back"},
		Results:  r,
	}

	c.HTML(200, "newtops.tmpl", data)
}

func handlerTopTracksGenres(c *gin.Context) {
	logger := logging.GetLogger(c)

	if tr := c.Query(queryStringTimeRange); len(tr) > 0 {
		mv := ddlOpts[tr]
		c.Set(string(keys.ContextSpotifyTimeRange), mv)
	}

	reqCtx, err := spotify.GetTopTracks(c, topGenresTopTracksLimit)
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			// If the error is due to the token being expired, we will have automatically attempted
			// to get a refresh token for the user.  If  that was successful, it will be returned
			// as the error value.  Set the cookie to the new value and redirect the user back to the
			// same path  to start the process again with the new token.
			if len(err.Error()) > 0 {
				c.SetCookie(cookieKeyToken, fmt.Sprint(err.Error()), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
				c.Redirect(http.StatusTemporaryRedirect, PathTopTracks)
				return
			}
		}

		logger.WithError(err).Error("couldnt retrieve top tracks from spotify")
		c.Status(500)
		return
	}

	reqTracks := keys.GetContextValue(reqCtx, keys.ContextSpotifyResults)
	tracks, ok := reqTracks.(spotify.Tracks)
	if !ok {
		logger.WithField("type", reflect.TypeOf(reqTracks)).Error("couldnt parse tracks from spotify")
		c.Status(500)
		return
	}

	reqCtx, err = spotify.GetGenresForTracks(reqCtx, tracks.IDs())
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			// If the error is due to the token being expired, we will have automatically attempted
			// to get a refresh token for the user.  If  that was successful, it will be returned
			// as the error value.  Set the cookie to the new value and redirect the user back to the
			// same path  to start the process again with the new token.
			if len(err.Error()) > 0 {
				c.SetCookie(cookieKeyToken, fmt.Sprint(err.Error()), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
				c.Redirect(http.StatusTemporaryRedirect, PathTopTracks)
				return
			}

			logging.GetLogger(c).Info("couldnt refresh token for user")
		}

		logger.WithError(err).Error("couldnt retreive top genres for top tracks from spotify")
		c.Status(500)
		return
	}

	reqGenres := keys.GetContextValue(reqCtx, keys.ContextSpotifyResults)
	genres, ok := reqGenres.(spotify.Pairs)
	if !ok {
		logger.WithField("type", reflect.TypeOf(reqGenres)).Error("couldnt parse genres from spotify")
		c.Status(500)
		return
	}

	type Result struct {
		Key        string
		Value      string
		Background string
		Width      int
		Height     int
	}

	type tempBag struct {
		Category string
		Type     string
		Opts     []string
		Results  []Result
	}

	r := []Result{}
	sort.Sort(sort.Reverse(genres))
	for _, i := range genres {
		r = append(r, Result{
			Key:   i.Key,
			Value: fmt.Sprint("( ", i.Value, ")"),
		})
	}
	data := tempBag{
		Category: "Genres",
		Type:     "tracks",
		Opts:     []string{"Recent", "In Between", "Going Way Back"},
		Results:  r,
	}

	c.HTML(200, "newtops.tmpl", data)
}

func handlerWordCloud(c *gin.Context) {
	c.HTML(200, "wordcloud2.tmpl", nil)
}

func handlerWordCloudData(c *gin.Context) {
	logger := logging.GetLogger(c)

	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		c.Set(string(keys.ContextSpotifyTimeRange), tr)
	}

	reqCtx, err := spotify.GetTopTracks(c, wordCloudTopTracksLimit)
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			// If the error is due to the token being expired, we will have automatically attempted
			// to get a refresh token for the user.  If  that was successful, it will be returned
			// as the error value.  Set the cookie to the new value and redirect the user back to the
			// same path  to start the process again with the new token.
			if len(err.Error()) > 0 {
				c.SetCookie(cookieKeyToken, fmt.Sprint(err.Error()), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
				c.Redirect(http.StatusTemporaryRedirect, PathTopTracks)
				return
			}
		}

		logger.WithError(err).Error("couldnt retrieve top tracks from spotify")
		c.Status(500)
		return
	}
	logger.Debug("got top tracks")

	reqTracks := keys.GetContextValue(reqCtx, keys.ContextSpotifyResults)
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

	wordCounts, err := genius.GetLyricCountForSong(c, searches)
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
	returl := keys.GetContextValue(c, keys.ContextSpotifyReturnURL)
	if returl == nil {
		//c.HTML(500, "error.tmpl", nil)
		c.Status(500)
		return
	}

	// TODO Add state
	pathScopes := url.QueryEscape(strings.Join(scopes, " "))
	spotifyURL := fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s&show_dialog=false",
		keys.GetContextValue(c, keys.ContextSpotifyClientID),
		pathScopes,
		returl)

	if len(c.Query("redirectUrl")) > 0 {
		c.SetCookie("redirect_url", c.Query("redirectUrl"), 600, "/", strings.Replace(host, "https://", "", -1), false, true)
	}

	logging.GetLogger(c).WithFields(logrus.Fields{
		"event": "redirect_for_oauth",
		"url": spotifyURL}).Debug("auth redirect")

	c.Redirect(http.StatusTemporaryRedirect, spotifyURL)
}

func handlerHome(c *gin.Context) {
	c.HTML(200, "home.tmpl", nil)
}

func handlerRecommendations(c *gin.Context) {
	type vb struct {
		Names []string
	}
	ret := vb{}
	for _, i := range *getData(c) {
		ret.Names = append(ret.Names, i)
	}
	c.HTML(200, "recommendations.tmpl", ret)
}

func handlerTest(c *gin.Context) {
	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		mv := ddlOpts[tr]
		c.Set(string(keys.ContextSpotifyTimeRange), mv)
	}

	tracksResultsCtx, err := spotify.GetTopTracks(c, topTracksLimit)
	if err != nil {
		logging.GetLogger(c).WithError(err).Error()
		c.Status(500)
		return
	}

	reqTracks := keys.GetContextValue(tracksResultsCtx, keys.ContextSpotifyResults)
	if reqTracks == nil {
		logging.GetLogger(c).Error("no tracks returned from spotify")
		c.Status(500)
		return
	}

	tracks, ok := reqTracks.(spotify.Tracks)
	if !ok {
		logging.GetLogger(c).WithField("type", reflect.TypeOf(reqTracks)).Error("couldnt parse tracks returned from spotify")
		c.Status(500)
		return
	}

	type Result struct {
		Key        string
		Value      string
		Background string
		Width      int
		Height     int
	}

	type tempBag struct {
		Category string
		Type     string
		Opts     []string
		Results  []Result
	}

	r := []Result{}
	for _, i := range tracks {
		r = append(r, Result{
			Key:        i.FindArtist(),
			Value:      i.Name,
			Background: i.FindImage().URL,
			Height:     i.FindImage().Height,
			Width:      i.FindImage().Width,
		})
	}
	data := tempBag{
		Category: "Tracks",
		Type:     "",
		Opts:     []string{"Recent", "In Between", "Going Way Back"},
		Results:  r,
	}

	c.HTML(200, "newtops.tmpl", data)
}

// getLevel1Recs retrieves the user's top artists and their related artists
func getLevel1Recs(ctx context.Context, seeds *spotify.Artists) (*map[string]int, context.Context) {
	artists := map[string]int{}
	artistCache := map[string]string{}

	for _, i := range *seeds {
		iName := strings.ToLower(i.Name)
		if _, ok := artists[iName]; !ok {
			artists[iName] = 1
			continue
		}
		artists[iName]++

		res, err := spotify.GetRelatedArtists(ctx, i.ID)
		if err != nil {
			fmt.Println("fuck3, ", err)
			return nil, ctx
		}

		for _, j := range *res {
			iName := strings.ToLower(j.Name)
			if _, ok := artists[iName]; !ok {
				artists[iName] = 1
				artistCache[iName] = j.ID
				continue
			}
			artists[iName]++
		}
	}

	return &artists, context.WithValue(ctx, keyArtistInfo, &artistCache)
}

func removeUsersKnownArtists(ctx context.Context, lib *map[string]int, recs *map[string]int) *[]string {

	sa := sortablemap.GetSortableMap(*recs)
	sort.Sort(sort.Reverse(sa))

	lm := sortablemap.GetSortableMap(*lib)
	sort.Sort(sort.Reverse(lm))

	whatsLeft := []string{}

	llib := *lib

	for k, v := range llib {
		logging.GetLogger(ctx).WithFields(map[string]interface{}{
			"key":   k,
			"value": v,
		}).Info("checking")
	}

	logging.GetLogger(ctx).WithFields(map[string]interface{}{
		"lib length":  len(llib),
		"recs length": len(*recs),
	}).Info("checking data")
	for k := range *recs {
		if _, ok := llib[strings.ToLower(k)]; !ok {
			logging.GetLogger(ctx).WithField("artist", k).Info("not in lib")
			whatsLeft = append(whatsLeft, strings.ToLower(k))
		} else {
			logging.GetLogger(ctx).WithField("artist", k).Info("excluding artist from lib from recs")
		}
	}

	return &whatsLeft
}

func getLevel2Recs(ctx context.Context, recs *spotify.Artists, orderedRecs *map[string]int) *map[string]int {
	// this method is fucked. I can't retrieve any additional info from spotify
	// with the way I've been passing around this information.  I need the ids for
	// api calls, and all I have is names.
	rrArtists := map[string]int{}
	sorted := sortablemap.GetSortableMap(*orderedRecs)
	sort.Sort(sort.Reverse(sorted))
	chunk := sorted.Take(5)
	ids := []string{}

	cache := *ctx.Value(keyArtistInfo).(*map[string]string)

	for _, i := range chunk {
		found := false
		iName := ""
		for _, j := range *recs {
			iName = strings.ToLower(j.Name)
			if i.Key == iName {
				ids = append(ids, j.ID)
				found = true
				break
			}
		}
		if !found {
			logging.GetLogger(nil).Error(fmt.Sprint("couldnt find artist id; ", iName))
			if v, ok := cache[iName]; ok {
				logging.GetLogger(nil).Info("found artist in cache")
				ids = append(ids, v)
			}
		}
	}

	res, err := spotify.GetRecommendations(ctx, map[string][]string{spotify.KeySeedArtists: ids})
	if err != nil {
		panic(err)
	}

	for _, j := range res.Tracks {
		for _, k := range j.Artists {
			iName := strings.ToLower(k)
			if _, ok := rrArtists[iName]; !ok {
				rrArtists[iName] = 1
				continue
			}

			rrArtists[iName]++
		}

	}

	// fold in the original recs

	for k := range *orderedRecs {
		iName := strings.ToLower(k)
		if _, ok := rrArtists[iName]; !ok {
			rrArtists[iName] = 1
			continue
		}

		rrArtists[iName]++
	}

	return &rrArtists
}

func addToCache(cache map[string]string, artists *spotify.Artists) map[string]string {
	ret := map[string]string{}
	for _, i := range *artists {
		iName := strings.ToLower(i.Name)
		if _, ok := ret[iName]; !ok {
			ret[iName] = i.ID
		}
	}
	return ret
}

func getData(ctx context.Context) *[]string {
	// Leaving off:
	//
	//
	// This is cool... I'm getting recs, and they seem accurate.
	// however; the filtering does not seem to be working.  I'm seeing results of which I'm sure I have songs saved.
	//
	//
	// Ideally, when I get recommendations they wouldn't contain any artists for which I have songs saved.  We're trying to
	// surface new music.

	ctx = context.WithValue(ctx, keys.ContextSpotifyTimeRange, "short_term")
	topArtists1, err := spotify.GetTopArtists(ctx)
	if err != nil {
		fmt.Println("fuck1: ", err)
		return nil
	}

	ctx = context.WithValue(ctx, keys.ContextSpotifyTimeRange, "medium_term")
	topArtists2, err := spotify.GetTopArtists(ctx)
	if err != nil {
		fmt.Println("fuck2: ", err)
		return nil
	}

	ctx = context.WithValue(ctx, keys.ContextSpotifyTimeRange, "long_term")
	topArtists3, err := spotify.GetTopArtists(ctx)
	if err != nil {
		fmt.Println("fuck3: ", err)
		return nil
	}

	sumArtists := append(spotify.Artists{}, *topArtists1...)
	sumArtists = append(sumArtists, *topArtists2...)
	sumArtists = append(sumArtists, *topArtists3...)

	artistCache := map[string]string{}
	artistCache = addToCache(artistCache, &sumArtists)

	ctx = context.WithValue(ctx, keyArtistInfo, &artistCache)

	lvl1, ctx := getLevel1Recs(ctx, &sumArtists)

	// this isn't going to work because I'll only have the top artist information, not their related artists
	lvl2 := getLevel2Recs(ctx, &sumArtists, lvl1)

	lib, err := spotify.GetUserLibraryTracks(ctx)
	if err != nil {
		panic(err)
	}
	libMap := getArtistCountFromLib(&lib)

	return removeUsersKnownArtists(ctx, libMap, lvl2)
}

func getArtistCountFromLib(lib *spotify.Tracks) *map[string]int {
	ret := map[string]int{}
	for _, i := range *lib {
		for _, j := range i.Artists {
			iName := strings.ToLower(j.Name)
			if _, ok := ret[iName]; !ok {
				ret[iName] = 1
				continue
			}

			ret[iName]++
		}
	}

	return &ret
}
