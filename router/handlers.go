package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
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
		// the user denied access
		logger.WithError(errors.New(qErr)).Error("user did not grant access")
		c.Status(500)
		return
	}

	tok, err := spotify.ExchangeOauthCode(c, code)
	if err != nil {
		logger.WithError(err).Error("error handling spotify oauth")
		c.Status(500)
		return
	}

	c.Set(string(keys.ContextSpotifyAccessToken), tok.Access)
	c.Set(string(keys.ContextSpotifyRefreshToken), tok.Refresh)

	if len(tok.Access) < 1 {
		logger.WithError(err).Error("no access token returned from spotify")
		c.Status(500)
		return
	}

	u, err := spotify.GetUser(c)
	if err != nil {
		logger.WithError(err).Error("couldnt retrieve userid from spotify")
		c.Status(500)
		return
	}

	err = u.Save(c)
	if err != nil {
		logger.WithField("info", *u).WithError(err).Error("couldnt save user")
		c.Status(500)
		return
	}

	if len(tok.Refresh) < 1 {
		logger.Error("no refresh token returned from spotify")
	}

	logger.WithFields(logrus.Fields{
		"event": "user_login",
		"id":    u.ID,
		"email": u.Email,
	}).Info("user logged in successfully")

	c.SetCookie(cookieKeyID, fmt.Sprint(u.ID), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
	c.SetCookie(cookieKeyToken, fmt.Sprint(tok.Access), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
	c.SetCookie(cookieKeyRefresh, fmt.Sprint(tok.Refresh), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
	// val, err := c.Cookie("redirect_url")
	// if err == nil && len(val) > 0 {
	// 	c.Redirect(http.StatusTemporaryRedirect, val)
	// }

	// EXP_REACT is a feature flag that will be used to toggle new behavior
	// without breaking the existing stuff.
	if os.Getenv("EXP_REACT") == "1" {
		host := c.Request.Referer()
		red, _ := c.Cookie("redirect_url")
		if len(red) > 0 {
			host = fmt.Sprint(host, red)
		}
		c.Redirect(http.StatusTemporaryRedirect, host)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(PathTopTracks, "?", queryStringTimeRange, "=short_term"))
}

func handlerTopTracks(c *gin.Context) {
	logger := logging.GetLogger(c)

	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		mv := ddlOpts[tr]
		c.Set(string(keys.ContextSpotifyTimeRange), mv)
	}

	trax, err := spotify.GetTopTracks(c, spotify.GetTimeFrame(ddlOpts[tr]))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			// TODO: try to refresh token and repeat request
			c.Redirect(http.StatusTemporaryRedirect, PathHome)
			return
		}

		logger.WithError(err).Error("couldnt retrieve top tracks from spotify")
		c.Status(500)
		return
	}

	//data := getTopTracksViewBag(trax)

	if os.Getenv("EXP_REACT") == "1" {
		c.JSON(200, trax)
		return
	}

	c.HTML(200, "newtops.tmpl", trax)
}

func getTopTracksViewBag(trax *spotify.Tracks) interface{} {
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
	for _, i := range *trax {
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
	return data
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
	t, err := spotify.GetSavedTracks(c)
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
	for i := 0; i < len(*t); i++ {
		tr := (*t)[i]
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
	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		mv := ddlOpts[tr]
		c.Set(string(keys.ContextSpotifyTimeRange), mv)
	}

	artists, err := spotify.GetTopArtists(c)
	if err != nil {
		logging.GetLogger(c).WithError(err).Error("couldnt retrieve top artists from spotify")
		c.Status(500)
		return
	}

	genres := artists.GetGenres(c)

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
	for _, i := range *genres {
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

	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		mv := ddlOpts[tr]
		c.Set(string(keys.ContextSpotifyTimeRange), mv)
	}

	trax, err := spotify.GetTopTracks(c, spotify.GetTimeFrame(ddlOpts[tr]))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			// TODO: try to refresh token and repeat request
			c.Redirect(http.StatusTemporaryRedirect, PathHome)
			return
		}

		logger.WithError(err).Error("couldnt retrieve top tracks from spotify")
		c.Status(500)
		return
	}

	genres, err := trax.GetGenres(c)
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			// TODO: try to refresh token and repeat request
			c.Redirect(http.StatusTemporaryRedirect, PathHome)
			return
		}

		logger.WithError(err).Error("couldnt retrieve top tracks from spotify")
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
	for _, i := range *genres {
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

	trax, err := spotify.GetTopTracks(c, spotify.GetTimeFrame(ddlOpts[tr]))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			// TODO: try to refresh token and repeat request
			c.Redirect(http.StatusTemporaryRedirect, PathHome)
			return
		}

		logger.WithError(err).Error("couldnt retrieve top tracks from spotify")
		c.Status(500)
		return
	}

	searches := []genius.LyricSearch{}
	for _, i := range *trax {
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
		"url":   spotifyURL}).Debug("auth redirect")

	c.Redirect(http.StatusTemporaryRedirect, spotifyURL)
}

func handlerHome(c *gin.Context) {
	c.HTML(200, "home.tmpl", nil)
}

func handlerRecommendations(c *gin.Context) {
	recs, err := getRecommendations(c)
	if err != nil {
		logging.GetLogger(c).WithField("event", "failed_recs").Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if os.Getenv("EXP_REACT") == "1" {
		c.JSON(200, *recs)
		return
	}

	c.HTML(200, "recommendations.tmpl", *recs)
}

func handlerTest(c *gin.Context) {
	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		mv := ddlOpts[tr]
		c.Set(string(keys.ContextSpotifyTimeRange), mv)
	}

	trax, err := spotify.GetTopTracks(c, spotify.GetTimeFrame(ddlOpts[tr]))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(spotify.ErrTokenExpired("")) {
			// TODO: try to refresh token and repeat request
			c.Redirect(http.StatusTemporaryRedirect, PathHome)
			return
		}

		logging.GetLogger(c).WithError(err).Error("couldnt retrieve top tracks from spotify")
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
	for _, i := range *trax {
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
		// if this artist isn't in the return object already, initialize it at
		if _, ok := artists[iName]; !ok {
			artists[iName] = 1
			// I don't know why we're continuing here and I'm guessing this is causing a bug
			continue
		}
		artists[iName]++

		res, err := i.GetRelatedArtists(ctx)
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
	// recs: the user's top spotify artists
	// orderedRecs: all of the user's recommendations and how many times they were present in the process

	// this method is fucked. I can't retrieve any additional info from spotify
	// with the way I've been passing around this information.  I need the ids for
	// api calls, and all I have is names.

	// get the 5 most recommended names
	rrArtists := map[string]int{}
	sorted := sortablemap.GetSortableMap(*orderedRecs)
	sort.Sort(sort.Reverse(sorted))
	chunk := sorted.Take(5)
	ids := []string{}

	// getting the name/id artist map out of context
	cache := *ctx.Value(keyArtistInfo).(*map[string]string)

	// iterate through the top 5 recommended
	for _, i := range chunk {
		found := false
		iName := ""

		// iterate through the user's top artists
		for _, j := range *recs {
			iName = strings.ToLower(j.Name)

			// if the current top artists is in the top 5 most recommended
			if i.Key == iName {
				// add this artist id to the seeds we're going the get recommendations for
				ids = append(ids, j.ID)
				found = true
				break
			}
		}

		// if we didn't find the current recommendation in the user's top 10 artists
		if !found {
			logging.GetLogger(nil).Error(fmt.Sprint("couldnt find artist id; ", iName))
			if v, ok := cache[iName]; ok {
				// check the cache for the information for the artist in the recommendation
				logging.GetLogger(nil).Info("found artist in cache")
				// add this artist id to the seeds we're going the get recommendations for
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
			iName := strings.ToLower(k.Name)
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

type Recommendation struct {
	Seed        string
	SeedID      string
	SeedResults *spotify.Artists
}

type Recommendations []Recommendation

func (r *Recommendations) GetSeeds() *map[string]int {
	ids := map[string]int{}
	for _, i := range *r {
		if _, ok := ids[i.SeedID]; !ok {
			ids[i.SeedID] = 1
		} else {
			ids[i.SeedID]++
		}

		// go through each reccomendation for each top artist
		for _, j := range *i.SeedResults {
			if _, ok := ids[j.ID]; !ok {
				ids[j.ID] = 1
			} else {
				ids[j.ID]++
			}
		}
	}

	return &ids
}

func getRecommendations(ctx context.Context) (*spotify.Recommendation, error) {
	// Start by getting thier top artist from each time frame
	ctx = context.WithValue(ctx, keys.ContextSpotifyTimeRange, "short_term")
	topArtists1, err := spotify.GetTopArtists(ctx)
	if err != nil {
		logging.GetLogger(ctx).WithField("event", "rec_art_err_1").WithError(err).Error("couldnt get top artists")
		return nil, err
	}

	ctx = context.WithValue(ctx, keys.ContextSpotifyTimeRange, "medium_term")
	topArtists2, err := spotify.GetTopArtists(ctx)
	if err != nil {
		logging.GetLogger(ctx).WithField("event", "rec_art_err_2").WithError(err).Error("couldnt get top artists")
		return nil, err
	}

	ctx = context.WithValue(ctx, keys.ContextSpotifyTimeRange, "long_term")
	topArtists3, err := spotify.GetTopArtists(ctx)
	if err != nil {
		logging.GetLogger(ctx).WithField("event", "rec_art_err_3").WithError(err).Error("couldnt get top artists")
		return nil, err
	}

	// consolidate into one slice
	sumArtists := append(spotify.Artists{}, *topArtists1...)
	sumArtists = append(sumArtists, *topArtists2...)
	sumArtists = append(sumArtists, *topArtists3...)

	logging.GetLogger(ctx).WithField("top_artist_count", len(sumArtists)).Info()

	// iterate through list and find unique artists
	// - if it's a new artist
	// - - store it in the arr
	// - - get related artists from spotify
	// - - iterate through results
	// - - - if it's a new artist
	// - - - - store it in the arr

	recs := Recommendations{}
	// we're iterating through each of the user's top artists to
	// get their related artists
	for _, i := range sumArtists {
		res, err := i.GetRelatedArtists(ctx)
		if err != nil {
			fmt.Println("fuck3, ", err)
			return nil, err
		}

		seriously := spotify.Artists{}

		for _, j := range *res {
			seriously = append(seriously, j)
		}

		recs = append(recs, Recommendation{
			Seed:        i.Name,
			SeedID:      i.ID,
			SeedResults: &seriously,
		})
	}

	logging.GetLogger(ctx).WithFields(logrus.Fields{
		"recs_1": len(*recs[0].SeedResults),
		"recs_2": len(*recs[1].SeedResults),
		"recs_3": len(*recs[2].SeedResults),
	}).Info("recommendation counts")

	seeds := recs.GetSeeds()
	sortedSeeds := sortablemap.GetSortableMap(*seeds)
	sort.Sort(sort.Reverse(sortedSeeds))
	topSeeds := sortedSeeds.Take(5)
	seedIDs := []string{}

	for _, i := range topSeeds {
		seedIDs = append(seedIDs, i.Key)
	}

	res, err := spotify.GetRecommendations(ctx, map[string][]string{spotify.KeySeedArtists: seedIDs})
	if err != nil {
		return nil, err
	}
	fmt.Println("got recommendations")

	return res, nil
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
