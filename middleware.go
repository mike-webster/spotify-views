package main

import (
	"github.com/gin-gonic/gin"
	data "github.com/mike-webster/spotify-views/data"
	genius "github.com/mike-webster/spotify-views/genius"
	spotify "github.com/mike-webster/spotify-views/spotify"
)

func LoadContextValues() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(string(spotify.ContextClientID), clientID)
		c.Set(string(spotify.ContextClientSecret), clientSecret)
		c.Set(string(genius.ContextAccessToken), lyricsKey)
		c.Set(string(data.ContextHost), dbHost)
		c.Set(string(data.ContextUser), dbUser)
		c.Set(string(data.ContextPass), dbPass)
		c.Set(string(data.ContextDatabase), dbName)
		c.Set(string(data.ContextSecurityKey), secKey)
		c.Next()
	}
}
