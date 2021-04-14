package main

import (
	"errors"
	"io/ioutil"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/mike-webster/spotify-views/keys"
	"github.com/mike-webster/spotify-views/logging"
	"github.com/sirupsen/logrus"
)

func setTokens(c *gin.Context) {
	entry := logging.GetLogger(c)

	tok, err := c.Cookie(cookieKeyToken)
	if err == nil {
		if len(tok) > 0 {
			c.Set(string(keys.ContextSpotifyAccessToken), tok)
		}
	}

	ref, err := c.Cookie(cookieKeyRefresh)
	if err == nil {
		if len(ref) > 0 {
			c.Set(string(keys.ContextSpotifyRefreshToken), ref)
		}
	}

	c.Next()
}

func setEnv(c *gin.Context) {
	entry := logging.GetLogger(c)
	entry.Debug()
	vals, err := parseEnvironmentVariables(c)
	if err != nil {
		entry.WithError(err).Error("error encountered parsing env vars")
		c.AbortWithError(500, err)
		return
	}

	for k, v := range vals {
		key, ok := k.(string)
		if !ok {
			kk, ok := k.(keys.ContextKey)
			if !ok {
				entry.WithFields(logrus.Fields{
					"event": "couldnt_parse_context_field",
					"key":   k,
					"value": v,
				}).Error("problem adding an environment variable to the context - aborting")
				c.AbortWithError(500, errors.New("couldnt parse context values"))
				return
			}
			key = string(kk)
		}
		c.Set(key, v)
	}

	c.Next()
}

func authenticate(c *gin.Context) {
	entry := logging.GetLogger(c)
	entry.Debug()
	tok, err := c.Cookie(cookieKeyToken)
	if err != nil {
		entry.WithField("redirect_reason", "authenticate: found no token").WithError(err).Error("redirecting")
		c.Redirect(301, "/?noauth")
		return
	}

	c.Set(string(keys.ContextSpotifyAccessToken), tok)
	c.Next()
}

// consolidate stack on crahes
func recovery(c *gin.Context) {
	defer func(c *gin.Context) {
		if r := recover(); r != nil {
			b, _ := ioutil.ReadAll(c.Request.Body)
			logging.GetLogger(c).WithFields(logrus.Fields{
				"event":    "ErrPanicked",
				"error":    r,
				"stack":    string(debug.Stack()),
				"path":     c.Request.RequestURI,
				"formbody": string(b),
			}).Error("panic recovered")

			c.HTML(500, "error.tmpl", nil)
		}
	}(c)
	c.Next() // execute all the handlers
}

func setContextLogger(c *gin.Context) {
	lf := parseLoggerValues(c)
	c.Set(string(keys.ContextLoggerFields), lf)
	logging.SetRequestLogger(c, logging.GetLogger(c))
	c.Next()
}

func parseLoggerValues(c * gin.Context) *logging.LoggerFields {
	reqID, _ := uuid.NewV4()
	uid, _ := c.Cookie("svid")
	return &logging.LoggerFields{
		UserAgent: c.Request.UserAgent(),
		Referer: c.Request.Referer(),
		QueryString: c.Request.URL.RawQuery,
		Path: c.Request.URL.Path,
		Method: c.Request.Method,
		ClientIP: c.ClientIP(),
		RequestID: reqID.String(),
		UserID: uid,
	}
}