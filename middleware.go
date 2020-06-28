package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	data "github.com/mike-webster/spotify-views/data"
	genius "github.com/mike-webster/spotify-views/genius"
	spotify "github.com/mike-webster/spotify-views/spotify"
	"github.com/sirupsen/logrus"
)

func loadContextValues() gin.HandlerFunc {
	return func(c *gin.Context) {
		parseEnvironmentVariables(c)
		c.Set(string(spotify.ContextClientID), clientID)
		c.Set(string(spotify.ContextClientSecret), clientSecret)
		c.Set(string(genius.ContextAccessToken), lyricsKey)
		c.Set(string(data.ContextHost), dbHost)
		c.Set(string(data.ContextUser), dbUser)
		c.Set(string(data.ContextPass), dbPass)
		c.Set(string(data.ContextDatabase), dbName)
		c.Set(string(data.ContextSecurityKey), secKey)
		c.Set(string(spotify.ContextReturnURL), fmt.Sprint("https://", host, "/spotify/oauth"))
		c.Next()
	}
}

// consolidate stack on crahes
func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				b, _ := ioutil.ReadAll(c.Request.Body)

				defaultLogger(c).WithFields(logrus.Fields{
					"event":    "ErrPanicked",
					"error":    r,
					"stack":    string(debug.Stack()),
					"path":     c.Request.RequestURI,
					"formbody": string(b),
				}).Error("panic recovered")

				c.AbortWithStatus(500)
			}
		}()
		c.Next() // execute all the handlers
	}
}

func requestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func(ctx *gin.Context) logrus.FieldLogger {
			// don't log requests to these paths when successful
			quiet := []string{
				"/healthcheck",
				"/static",
				"/clouds",
			}
			skip := false
			for _, i := range quiet {
				if strings.HasPrefix(ctx.Request.URL.Path, i) {
					skip = true
				}
			}
			if skip && ctx.Writer.Status() == 200 {
				return nil
			}

			// log body if one is given]
			strBody := ""
			body, err := ioutil.ReadAll(ctx.Request.Body)
			if err != nil {
				defaultLogger(ctx).WithField("error", err).Error("cant read request body")
			} else {
				// write the body back into the request
				ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

				strBody = string(body)
				strBody = strings.Replace(strBody, "\n", "", -1)
				strBody = strings.Replace(strBody, "\t", "", -1)
			}

			logger := defaultLogger(ctx).WithFields(logrus.Fields{
				"client_ip":    ctx.ClientIP(),
				"event":        "http.in",
				"method":       ctx.Request.Method,
				"path":         ctx.GetString("originalPath"),
				"query":        ctx.Request.URL.RawQuery,
				"referer":      ctx.Request.Referer(),
				"status":       ctx.Writer.Status(),
				"user_agent":   ctx.Request.UserAgent(),
				"request_body": strBody,
			})

			if len(ctx.Errors) > 0 {
				logger.Error(strings.TrimSpace(ctx.Errors.String()))
			} else {
				logger.Info()
			}

			sBody := struct {
				Repository struct {
					Name string `json:"name"`
				} `json:"repository"`
			}{}

			err = json.Unmarshal([]byte(strBody), &sBody)
			if err != nil {
				logger.WithField("error", err).Error("couldn't unmarshall body to get repo name")
			}

			reqID, _ := uuid.NewV4()
			ctx.Set("logger", defaultLogger(ctx).WithFields(logrus.Fields{
				"request_id": reqID,
			}))

			return logger
		}(ctx)
	}
}
