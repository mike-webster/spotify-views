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
)

var (
	queryStringCode                 = "code"
	queryStringError                = "error"
	queryStringTimeRange            = "time_range"
	cookieKeyToken                  = "svauth"
	cookieKeyID                     = "svid"
	keyArtistInfo            string = "artist-cache"
	topTracksLimit           int32  = 25
	topGenresTopTracksLimit  int32  = 50
	wordCloudTopTracksLimit  int32  = 50
	spotifyPlayerHeightShort int32  = 80
	spotifyPlayerHeightTall  int32  = 380
	spotifyPlayerWidth       int32  = 300
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

	accessTok, refreshTok, err := spotify.HandleOauth(c, code)
	fmt.Println("Tokens: ", accessTok, " - ", refreshTok)
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

	success, err := data.SaveUser(c, info["id"], info["email"])
	if err != nil {
		logger.WithField("info", info).WithError(err).Error("couldnt save user")
		c.Status(500)
		return
	}

	if !success {
		logger.WithField("info", info).Warn("couldnt create user - may have already existed")
	}

	if len(refreshTok) < 1 {
		logger.Error("no refresh token returned from spotify")
	} else {
		success, err := data.SaveRefreshTokenForUser(c, fmt.Sprint(refreshTok), info["id"])
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
	c.SetCookie(cookieKeyToken, fmt.Sprint(accessTok), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
	val, err := c.Cookie("redirect_url")
	if err == nil && len(val) > 0 {
		c.Redirect(http.StatusTemporaryRedirect, val)
	}
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(PathTopTracks, "?", queryStringTimeRange, "=short_term"))
}

func handlerTopTracks(c *gin.Context) {
	logger := logging.GetLogger(c)
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		// TODO: check for refresh token first
		// TODO: maybe move this ^^ check into middleware?
		logger.Debug("no token, redirecting")
		c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathTopTracks)
		return
	}

	requestCtx := context.WithValue(c, keys.ContextSpotifyAccessToken, token)
	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		requestCtx = context.WithValue(requestCtx, keys.ContextSpotifyTimeRange, tr)
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
	logger := logging.GetLogger(c)
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		logger.Debug("no token, redirecting")
		c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathTopArtists)
		return
	}

	requestCtx := context.WithValue(c, keys.ContextSpotifyAccessToken, token)
	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		requestCtx = context.WithValue(requestCtx, keys.ContextSpotifyTimeRange, tr)
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

	reqArtists := keys.GetContextValue(artistResponseCtx, keys.ContextSpotifyResults)
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

func handlerUserLibraryTempo(c *gin.Context) {
	logger := logging.GetLogger(c)
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		logger.Debug("no token, redirecting")
		c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathTopTracks)
		return
	}

	requestCtx := context.WithValue(c, keys.ContextSpotifyAccessToken, token)
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
	logger := logging.GetLogger(c)
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		logger.Debug("no token, redirecting")
		c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathTopArtistGenres)
		return
	}

	ctx := context.WithValue(c, keys.ContextSpotifyAccessToken, token)

	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		ctx = context.WithValue(ctx, keys.ContextSpotifyTimeRange, tr)
	}

	reqCtx, err := spotify.GetTopArtists(ctx)
	if err != nil {
		logger.WithError(err).Error("couldnt retrieve top artists from spotify")
		c.Status(500)
		return
	}

	reqArtists := keys.GetContextValue(reqCtx, keys.ContextSpotifyResults)
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

	reqGenres := keys.GetContextValue(reqCtx, keys.ContextSpotifyResults)
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
	logger := logging.GetLogger(c)
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		logger.Debug("no token, redirecting")
		c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathTopTracksGenres)
		return
	}

	ctx := context.WithValue(c, keys.ContextSpotifyAccessToken, token)

	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		ctx = context.WithValue(ctx, keys.ContextSpotifyTimeRange, tr)
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

	reqTracks := keys.GetContextValue(reqCtx, keys.ContextSpotifyResults)
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

	reqGenres := keys.GetContextValue(reqCtx, keys.ContextSpotifyResults)
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
	logger := logging.GetLogger(c)
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		logger.Debug("no token, redirecting")
		c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathWordCloud)
		return
	}

	ctx := context.WithValue(c, keys.ContextSpotifyAccessToken, token)

	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		ctx = context.WithValue(ctx, keys.ContextSpotifyTimeRange, tr)
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

	ctx = context.WithValue(c, keys.ContextLyricsToken, lyricsKey)

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
	returl := keys.GetContextValue(c, keys.ContextSpotifyReturnURL)
	if returl == nil {
		//c.HTML(500, "error.tmpl", nil)
		c.Status(500)
		return
	}

	// TODO Add state
	pathScopes := url.QueryEscape(strings.Join(scopes, " "))
	spotifyURL := fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s&show_dialog=false",
		clientID,
		pathScopes,
		returl)

	if len(c.Query("redirectUrl")) > 0 {
		c.SetCookie("redirect_url", c.Query("redirectUrl"), 600, "/", strings.Replace(host, "https://", "", -1), false, true)
	}

	logging.GetLogger(nil).Info(fmt.Sprint("redirecting for spotify auth: ", spotifyURL))

	c.Redirect(http.StatusTemporaryRedirect, spotifyURL)
}

func handlerHome(c *gin.Context) {
	c.HTML(200, "home.tmpl", nil)
}

func handlerRecommendations(c *gin.Context) {
	logger := logging.GetLogger(nil)
	token, err := c.Cookie(cookieKeyToken)
	if err != nil {
		logger.Debug("no token, redirecting")
		c.Redirect(http.StatusTemporaryRedirect, PathLogin+"?redirectUrl="+PathTopTracksGenres)
		return
	}

	ctx := context.WithValue(c, keys.ContextSpotifyAccessToken, token)
	ctx = context.WithValue(c, keyArtistInfo, &map[string]string{})
	c.JSON(200, getData(ctx))
}

func handlerTest(c *gin.Context) {
	id, err := c.Cookie("svid")
	if err != nil {
		panic(err)
	}
	refTok, err := data.GetRefreshTokenForUser(c, id)
	if err != nil {
		panic(err)
	}
	ctx := context.WithValue(c, keys.ContextSpotifyRefreshToken, refTok)
	tok, err := spotify.RefreshToken(ctx)
	if err != nil {
		panic(err)
	}
	c.JSON(200, gin.H{"token": tok})
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

		for _, i := range *res {
			iName := strings.ToLower(i.Name)
			if _, ok := artists[iName]; !ok {
				artists[iName] = 1
				artistCache[iName] = i.ID
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
	// this didn't work
	for _, k := range sa {
		found := false
		for _, kk := range lm {
			if strings.ToLower(k.Key) == strings.ToLower(kk.Key) {
				found = true
				break
			}
		}
		if found {
			continue
		}
		whatsLeft = append(whatsLeft, k.Key)
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

	iRsp, err := spotify.GetTopArtists(ctx)
	if err != nil {
		fmt.Println("fuck1: ", err)
		return nil
	}
	ita := keys.GetContextValue(iRsp, keys.ContextSpotifyResults)
	topArtists, ok := ita.(spotify.Artists)
	if !ok {
		fmt.Println("fuck3", reflect.TypeOf(ita))
		return nil
	}

	artistCache := map[string]string{}
	for _, i := range topArtists {
		iName := strings.ToLower(i.Name)
		if _, ok := artistCache[iName]; !ok {
			artistCache[iName] = i.ID
		}
	}
	ctx = context.WithValue(ctx, keyArtistInfo, &artistCache)

	lvl1, ctx := getLevel1Recs(ctx, &topArtists)

	// so... I think ideally, before we return the artists would it be good to filter by their library?
	// I'm thinking - the artist recommendation counter indicates the "strength" of the recommendation
	// because it's just how many times we would "recommend" that artist based on their current top artsts.
	// We can use the inverse, measure the artists in their library with a measure  of tracks saved to
	// indicate strength, in order to determine what we should actually be recommending.  We wouldn't want
	// to be recommending someone such as myself Blink-182 - based on my library you could probably guess
	// I know most of that already.

	// random thought: don't use "top artists" as provided by spotify, calculate top artists based on
	// the user's library.  Use the measure described above - tracks per artist saved - as a measure
	// of popularity

	// another random thought: I wonder if there's a way to find out how many times a user has listened
	// to a track.  If so, we can take this a little further.

	// The more we can refine this and sort of build our own definitions, the closer we'll be able to get
	// to offering our own "featured" recommendations.  We'll want to make sure they're different from
	// what spotify gives you... so think  about how to  measure  the differences.

	// -- end rant

	// this isn't going to work because I'll only have the top artist information, not their related artists
	lvl2 := getLevel2Recs(ctx, &topArtists, lvl1)

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
		iName := strings.ToLower(i.Name)
		if _, ok := ret[iName]; !ok {
			ret[iName] = 1
			continue
		}

		ret[iName]++
	}

	return &ret
}
