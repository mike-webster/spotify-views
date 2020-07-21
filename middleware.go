package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	data "github.com/mike-webster/spotify-views/data"
	genius "github.com/mike-webster/spotify-views/genius"
	"github.com/mike-webster/spotify-views/logging"
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
		ctx := context.WithValue(c, "dumb", "im just doing this to switch the type to use with the logging package")
		c.Set("logger", logging.GetLogger(&ctx))
		c.Next()
	}
}

// consolidate stack on crahes
func recovery(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			b, _ := ioutil.ReadAll(c.Request.Body)
			logging.GetLogger(nil).WithFields(logrus.Fields{
				"event":    "ErrPanicked",
				"error":    r,
				"stack":    string(debug.Stack()),
				"path":     c.Request.RequestURI,
				"formbody": string(b),
			}).Error("panic recovered")

			c.HTML(500, "error.tmpl", nil)
		}
	}()
	c.Next() // execute all the handlers
}

func redisClient(c *gin.Context) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASS")
	addr := fmt.Sprintf("%v:%v", host, port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: fmt.Sprint(password),
		DB:       0,
	})
	_, err := rdb.Ping(c).Result()
	if err != nil {
		logging.GetLogger(nil).WithField("event", "redis_crash").Error(err.Error())
	} else {
		c.Set("Redis", rdb)
	}
	c.Next()
}

func parseUserID(c *gin.Context) {
	uid, err := c.Cookie("svid")
	if err == nil {
		c.Set(string(spotify.ContextUserID), uid)
	}
	c.Next()
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
				logging.GetLogger(nil).WithField("error", err).Error("cant read request body")
			} else {
				// write the body back into the request
				ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

				strBody = string(body)
				strBody = strings.Replace(strBody, "\n", "", -1)
				strBody = strings.Replace(strBody, "\t", "", -1)
			}

			reqID, _ := uuid.NewV4()
			logger := logging.GetLogger(nil).WithFields(logrus.Fields{
				"client_ip":    ctx.ClientIP(),
				"event":        "http.in",
				"method":       ctx.Request.Method,
				"path":         ctx.Request.URL.Path,
				"query":        ctx.Request.URL.RawQuery,
				"referer":      ctx.Request.Referer(),
				"status":       ctx.Writer.Status(),
				"user_agent":   ctx.Request.UserAgent(),
				"request_body": strBody,
				"request_id":   reqID,
			})

			if len(ctx.Errors) > 0 {
				logger.Error(strings.TrimSpace(ctx.Errors.String()))
			} else {
				logger.Info()
			}

			ctx.Set(string(logging.ContextLogger), logger)

			return logger
		}(ctx)
	}
}
