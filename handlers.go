package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	genius "github.com/mike-webster/spotify-views/genius"
	sortablemap "github.com/mike-webster/spotify-views/sortablemap"
	spotify "github.com/mike-webster/spotify-views/spotify"
)

var (
	queryStringCode                = "code"
	queryStringError               = "error"
	queryStringTimeRange           = "time_range"
	cookieKeyToken                 = "svauth"
	topTracksLimit           int32 = 25
	topGenresTopTracksLimit  int32 = 50
	wordCloudTopTracksLimit  int32 = 50
	spotifyPlayerHeightShort int32 = 80
	spotifyPlayerHeightTall  int32 = 380
	spotifyPlayerWidth       int32 = 300
)

func handlerOauth(c *gin.Context) {
	ctx := context.WithValue(c, spotify.ContextReturnURL, returnURL)
	ctx = context.WithValue(ctx, spotify.ContextClientID, clientID)
	ctx = context.WithValue(ctx, spotify.ContextClientSecret, clientSecret)
	code := c.Query(queryStringCode)
	// TODO: query state verification
	qErr := c.Query(queryStringError)
	if len(qErr) > 0 {
		// the user is a fucker and they denied access
		log.Println("user did not grant access: ", qErr)
		c.Status(500)
		return
	}

	ctx, err := spotify.HandleOauth(ctx, code)
	if err != nil {
		log.Println("error handling spotify oauth: ", err.Error())
		c.Status(500)
		return
	}

	token := ctx.Value(spotify.ContextAccessToken)
	if token == nil {
		log.Println("no token returned from spotify")
		c.Status(500)
		return
	}

	c.SetCookie(cookieKeyToken, fmt.Sprint(token), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprint("/tracks/top?", queryStringTimeRange, "=short_term"))
}

func handlerTopTracks(c *gin.Context) {
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

	tracks, err := spotify.GetTopTracks(ctx, topTracksLimit)
	if err != nil {
		log.Println("couldnt retrieve top tracks from spotify: ", err.Error())
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

	for _, i := range *tracks {
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
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	ctx := context.WithValue(c, spotify.ContextAccessToken, token)

	tr := c.Query(queryStringTimeRange)
	if len(tr) > 0 {
		ctx = context.WithValue(ctx, spotify.ContextTimeRange, tr)
	}

	artists, err := spotify.GetTopArtists(ctx)
	if err != nil {
		log.Println("couldnt retrieve top artists from spotify: ", err.Error())
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

	for _, i := range *artists {
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

	artists, err := spotify.GetTopArtists(ctx)
	if err != nil {
		log.Println("couldnt retrieve top artists from spotify: ", err.Error())
		c.Status(500)
		return
	}

	genres, err := spotify.GetGenresForArtists(ctx, artists.IDs())
	if err != nil {
		log.Println("couldnt retrieve genres for artists from spotify: ", err.Error())
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

	tracks, err := spotify.GetTopTracks(ctx, topGenresTopTracksLimit)
	if err != nil {
		log.Println("couldnt retrieve top tracks from spotify: ", err.Error())
		c.Status(500)
		return
	}

	genres, err := spotify.GetGenresForTracks(ctx, tracks.IDs())
	if err != nil {
		log.Println("couldnt retrieve top genres for top tracks from spotify: ", err.Error())
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

	tracks, err := spotify.GetTopTracks(ctx, wordCloudTopTracksLimit)
	if err != nil {
		log.Println("couldnt retrieve top tracks from spotify: ", err.Error())
		c.Status(500)
		return
	}

	searches := []genius.LyricSearch{}
	for _, i := range *tracks {
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

	sm := sortablemap.GetSortableMap(wordCounts)
	sort.Sort(sort.Reverse(sm))
	filename := fmt.Sprint(time.Now().Unix(), ".png")
	err = generateWordCloud(ctx, filename, wordCounts)
	if err != nil {
		log.Println("couldn't generate error: ", err.Error())
		c.Status(500)
		return
	}

	// displaying the image
	readBack, err := os.Open(fmt.Sprint("static/clouds/", filename))
	if err != nil {
		log.Println("couldn't read file back: ", err.Error())
		c.Status(500)
		return
	}

	defer readBack.Close()
	id, _, err := image.Decode(readBack)
	if err != nil {
		log.Println("couldnt decode")
		c.Status(500)
		return
	}

	var buff bytes.Buffer
	png.Encode(&buff, id)

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
	redirectURL := fmt.Sprintf("https://accounts.spotify.com/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s&show_dialog=false",
		clientID,
		pathScopes,
		returnURL)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func handlerHome(c *gin.Context) {
	c.HTML(200, "home.tmpl", nil)
}
