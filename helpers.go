package main

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"image/png"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/psykhi/wordclouds"
)

func parseEnvironmentVariables() error {
	clientID = os.Getenv("CLIENT_ID")
	if len(clientID) < 1 {
		return errors.New("no client id provided")
	}

	clientSecret = os.Getenv("CLIENT_SECRET")
	if len(clientSecret) < 1 {
		return errors.New("no client secret provided")
	}

	host = os.Getenv("HOST")
	if len(host) < 1 {
		return errors.New("no host provided")
	}

	lyricsKey = os.Getenv("LYRICS_KEY")
	if len(lyricsKey) < 1 {
		return errors.New("no lyrics key provided")
	}

	return nil
}

func generateWordCloud(ctx context.Context, filename string, wordCounts map[string]int) error {
	colors := []color.RGBA{
		{0x17, 0xA5, 0x54, 0xff},
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

func runServer() {
	r := gin.Default()
	r.Use(LoadContextValues())
	r.LoadHTMLGlob("templates/*")
	r.GET("/spotify/oauth", handlerOauth)
	r.GET("/tracks/top", handlerTopTracks)
	r.GET("/artists/top", handlerTopArtists)
	r.GET("/artists/genres", handlerTopArtistsGenres)
	r.GET("/tracks/genres", handlerTopTracksGenres)
	r.GET("/login", handlerLogin)
	r.GET("/", handlerHome)
	r.GET("/wordcloud", handlerWordCloud)

	r.Static("/static/css", "./static")
	r.Static("/static/js", "./static")
	r.Static("/clouds/", "./static/clouds")

	r.Run()
}
