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

func handlerOauth(c *gin.Context) {
	ctx := context.WithValue(c, "return_url", returnURL)
	ctx = context.WithValue(ctx, "client_id", clientID)
	ctx = context.WithValue(ctx, "client_secret", clientSecret)
	code := c.Query("code")
	//state := c.Query("state")
	qErr := c.Query("error")
	if len(qErr) > 0 {
		// the user is a fucker and they denied access
	}
	ctx, err := spotify.HandleOauth(ctx, code)
	if err != nil {
		c.JSON(500, gin.H{"err": err})
		return
	}
	token := ctx.Value("access_token")
	if token == nil {
		c.JSON(500, gin.H{"err": "no access token"})
		return
	}
	c.SetCookie("svauth", fmt.Sprint(token), 3600, "/", strings.Replace(host, "https://", "", -1), false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/tracks/top")
}

func handlerTopTracks(c *gin.Context) {
	token, err := c.Cookie("svauth")
	if err != nil {
		log.Println("no token - redirecting to login")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	ctx := context.WithValue(c, "access_token", token)

	tr := c.Query("time_range")
	if len(tr) > 0 {
		ctx = context.WithValue(ctx, "time_range", tr)
	}

	tracks, err := spotify.GetTopTracks(ctx, 25)
	if err != nil {
		c.JSON(500, gin.H{"err": err})
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
			Height:   80,
			Width:    300,
			Resource: "track",
			Name:     i.Name,
		})
	}

	c.HTML(200, "toptracks.tmpl", data)
}

func handlerTopArtists(c *gin.Context) {
	log.Println("getting token")
	token, err := c.Cookie("svauth")
	if err != nil {
		log.Println("no token - redirecting to login")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	ctx := context.WithValue(c, "access_token", token)

	tr := c.Query("time_range")
	if len(tr) > 0 {
		ctx = context.WithValue(ctx, "time_range", tr)
	}
	log.Println("getting artists")
	artists, err := spotify.GetTopArtists(ctx)
	if err != nil {
		c.JSON(500, gin.H{"err": err})
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
	log.Println("getting token")
	token, err := c.Cookie("svauth")
	if err != nil {
		log.Println("no token - redirecting to login")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	ctx := context.WithValue(c, "access_token", token)

	tr := c.Query("time_range")
	if len(tr) > 0 {
		ctx = context.WithValue(ctx, "time_range", tr)
	}
	log.Println("getting artists")
	artists, err := spotify.GetTopArtists(ctx)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"err": err})
		return
	}
	log.Println("getting genres")
	genres, err := spotify.GetGenresForArtists(ctx, artists.IDs())
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	sort.Sort(sort.Reverse(genres))
	log.Println(genres)
	vb := ViewBag{Resource: "artist", Results: genres}
	c.HTML(200, "topgenres.tmpl", vb)
}

func handlerTopTracksGenres(c *gin.Context) {
	log.Println("getting token")
	token, err := c.Cookie("svauth")
	if err != nil {
		log.Println("no token - redirecting to login")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	ctx := context.WithValue(c, "access_token", token)

	tr := c.Query("time_range")
	if len(tr) > 0 {
		ctx = context.WithValue(ctx, "time_range", tr)
	}
	log.Println("getting top tracks")
	tracks, err := spotify.GetTopTracks(ctx, 50)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	log.Println("getting genres")
	genres, err := spotify.GetGenresForTracks(ctx, tracks.IDs())
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	sort.Sort(sort.Reverse(genres))
	log.Println(genres)
	vb := ViewBag{Resource: "track", Results: genres}
	c.HTML(200, "topgenres.tmpl", vb)
}

func handlerWordCloud(c *gin.Context) {
	token, err := c.Cookie("svauth")
	if err != nil {
		log.Println("no token - redirecting to login")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	ctx := context.WithValue(c, "access_token", token)

	tr := c.Query("time_range")
	if len(tr) > 0 {
		ctx = context.WithValue(ctx, "time_range", tr)
	}

	tracks, err := spotify.GetTopTracks(ctx, 25)
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	searches := []genius.LyricSearch{}
	for _, i := range *tracks {
		searches = append(searches, genius.LyricSearch{
			Artist: i.Artists[0].Name,
			Track:  i.Name,
		})
	}

	ctx = context.WithValue(c, "lyrics_token", lyricsKey)

	wordCounts, err := genius.GetLyricCountForSong(ctx, searches)
	if err != nil {
		log.Println("couldn't retrieve word counts: ", err.Error())
		c.Status(500)
		return
	}

	log.Println("\n\ncombined results:\n\n", wordCounts)
	sm := sortablemap.GetSortableMap(wordCounts)
	sort.Sort(sort.Reverse(sm))
	for _, i := range sm {
		log.Println("item: ", i.Key, " -- ", i.Value)
	}

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
		Maps     sortablemap.SortableMap
	}
	vb := viewData{Filename: filename, Maps: sm}
	c.HTML(200, "wordcloud.tmpl", vb)
	// c.Data(200, "image/png", buff.Bytes())
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
