package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"github.com/mike-webster/spotify-views/keys"
	"github.com/mike-webster/spotify-views/logging"
	"github.com/sirupsen/logrus"
)

func loadContextValues(c *gin.Context) {
	logging.GetLogger(c).WithField("event", "attaching context values").Debug()

	vals, err := parseEnvironmentVariables(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	uid, _ := c.Cookie("svid")
	c.Set(string(keys.ContextSpotifyUserID), uid)

	for k, v := range vals {
		key, ok := k.(string)
		if !ok {
			kk, ok := k.(keys.ContextKey)
			if !ok {
				logging.GetLogger(c).WithFields(map[string]interface{}{
					"event": "couldnt_parse_context_field",
					"key":   k,
					"value": v,
				}).Error()
				c.AbortWithError(500, errors.New("couldnt parse context values"))
				return
			}
			key = string(kk)
		}
		c.Set(key, v)
	}
	tok, _ := c.Cookie(cookieKeyToken)
	if len(tok) > 0 {
		c.Set(string(keys.ContextSpotifyAccessToken), tok)
	}
	ref, _ := c.Cookie(cookieKeyRefresh)
	if len(ref) > 0 {
		c.Set(string(keys.ContextSpotifyRefreshToken), ref)
	}
	c.Next()
}

func authenticate(c *gin.Context) {
	logging.GetLogger(c).WithField("event", "authenticating").Debug()
	tok, err := c.Cookie(cookieKeyToken)
	if err == nil {
		logging.GetLogger(c).Info("found token")

		c.Set(string(keys.ContextSpotifyAccessToken), tok)
		c.Next()
	}

	c.Redirect(301, "/?noauth")
	return
}

// consolidate stack on crahes
func recovery(c *gin.Context) {
	logging.GetLogger(c).WithField("event", "attaching_panic_recovery").Debug()

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
	if _redisDB != nil {
		_, err := _redisDB.Ping(c).Result()
		if err == nil {
			c.Set("Redis", _redisDB)
			logging.GetLogger(nil).WithField("event", "found_redis_client").Info()
			c.Next()
			return
		}

		logging.GetLogger(nil).WithField("event", "existing_redis_ping_error").Error(err.Error())
	}

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
		logging.GetLogger(nil).WithField("event", "new_redis_ping_error").Error(err.Error())
	} else {
		c.Set("Redis", rdb)
		_redisDB = rdb
		logging.GetLogger(nil).WithField("event", "new_redis_client").Info()
	}
	c.Next()
}

func parseUserID(c *gin.Context) {
	logging.GetLogger(c).WithField("event", "parsing_user_id").Debug()

	uid, err := c.Cookie("svid")
	if err == nil {
		c.Set(string(keys.ContextSpotifyUserID), uid)
	}
	c.Next()
}

func requestLogger(ctx *gin.Context) {
	logging.GetLogger(ctx).WithField("event", "attaching_request_logger").Debug()

	// don't log requests to these paths when successful
	quiet := []string{
		"/healthcheck",
		"/static",
		"/clouds",
		"/logos",
		"/favicon.ico",
	}
	skip := false
	for _, i := range quiet {
		if strings.HasPrefix(ctx.Request.URL.Path, i) {
			skip = true
		}
	}
	if skip && ctx.Writer.Status() == 200 {
		ctx.Next()
		return
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

	ctx.Set(string(keys.ContextLogger), logger)
	ctx.Next()
}
