package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"github.com/mike-webster/spotify-views/keys"
	"github.com/mike-webster/spotify-views/logging"
	"github.com/sirupsen/logrus"
)

func setTokens(c *gin.Context) {
	entry := logging.GetLogger(c).WithField("event", "loading_token")

	tok, err := c.Cookie(cookieKeyToken)
	if err != nil {
		entry.WithFields(logrus.Fields{
			"event": "err_retrieving_token",
			"error": err}).Warn("token not added to context")
	} else {
		if len(tok) > 0 {
			entry.Debug("found token")
			c.Set(string(keys.ContextSpotifyAccessToken), tok)
		}
	}

	ref, err := c.Cookie(cookieKeyRefresh)
	if err != nil {
		entry.WithFields(logrus.Fields{
			"event": "err_retrieving_refresh_token",
			"error": err}).Warn("refresh not added to context")
	} else {
		if len(ref) > 0 {
			entry.Debug("found refresh token")
			c.Set(string(keys.ContextSpotifyRefreshToken), ref)
		}
	}

	c.Next()
}

func setUserID(c *gin.Context) {
	entry := logging.GetLogger(c).WithField("event", "loading_user_id")

	uid, err := c.Cookie("svid")
	if err != nil {
		if c.Request.URL.Path != "/" {
			entry.WithError(err).Error("couldnt retrieve user id")
		}
	} else {
		entry.Debug("found user_id")
		c.Set(string(keys.ContextSpotifyUserID), uid)
	}

	c.Next()
}

func setEnv(c *gin.Context) {
	entry := logging.GetLogger(c).WithField("event", "loading_env_vars")

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
	entry := logging.GetLogger(c).WithField("event", "authenticating")
	tok, err := c.Cookie(cookieKeyToken)
	if err != nil {
		entry.WithField("redirect_reason", "authenticate: found no token").WithError(err).Error("redirecting")
		c.Redirect(301, "/?noauth")
		return
	}

	entry.Info("found token")

	c.Set(string(keys.ContextSpotifyAccessToken), tok)
	c.Next()
}

// consolidate stack on crahes
func recovery(c *gin.Context) {
	logging.GetLogger(c).WithField("event", "attaching_panic_recovery").Debug()

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

func setContextLogger(c *gin.Context) {
	logging.SetRequestLogger(c, logging.GetLogger(c))
	c.Next()
}

func setRequestID(c *gin.Context) {
	reqID, _ := uuid.NewV4()

	entry := logging.GetLogger(c).WithField("request_id", reqID)
	logging.SetRequestLogger(c, entry)
	c.Next()
}

func setClientIP(c *gin.Context) {
	entry := logging.GetLogger(c).WithField("event", "set_client_ip")

	if len(c.ClientIP()) > 0 {
		entry = entry.WithField("client_ip", c.ClientIP())
		logging.SetRequestLogger(c, entry)
	} else {
		entry.Warn("couldnt find IP")
	}

	c.Next()
}

func setMethod(c *gin.Context) {
	entry := logging.GetLogger(c)

	if len(c.Request.Method) > 0 {
		entry = entry.WithField("method", c.Request.Method)
		logging.SetRequestLogger(c, entry)
	}

	c.Next()
}

func setRequestPath(c *gin.Context) {
	entry := logging.GetLogger(c)

	if len(c.Request.URL.Path) > 0 {
		entry = entry.WithField("path", c.Request.URL.Path)
		logging.SetRequestLogger(c, entry)
	}

	c.Next()
}

func setRequestQuery(c *gin.Context) {
	entry := logging.GetLogger(c)

	if len(c.Request.URL.RawQuery) > 0 {
		entry = entry.WithField("query", c.Request.URL.RawQuery)
		logging.SetRequestLogger(c, entry)
	}

	c.Next()
}

func setReferer(c *gin.Context) {
	entry := logging.GetLogger(c)

	if len(c.Request.Referer()) > 0 {
		entry = entry.WithField("referer", c.Request.Referer())
		logging.SetRequestLogger(c, entry)
	}

	c.Next()
}

func setUserAgent(c *gin.Context) {
	entry := logging.GetLogger(c)

	if len(c.Request.UserAgent()) > 0 {
		entry = entry.WithField("user_agent", c.Request.UserAgent())
		logging.SetRequestLogger(c, entry)
	}

	c.Next()
}

func requestLogger(ctx *gin.Context) {
	logging.GetLogger(ctx).WithField("event", "attaching_request_logger").Debug()

	// don't log requests to these paths when successful
	// quiet := []string{
	// 	"/healthcheck",
	// 	"/static",
	// 	"/clouds",
	// 	"/logos",
	// 	"/favicon.ico",
	// }
	// skip := false
	// for _, i := range quiet {
	// 	if strings.HasPrefix(ctx.Request.URL.Path, i) {
	// 		skip = true
	// 	}
	// }
	// if skip && ctx.Writer.Status() == 200 {
	// 	ctx.Next()
	// 	return
	// }

	// log body if one is given]
	// strBody := ""
	// body, err := ioutil.ReadAll(ctx.Request.Body)
	// if err != nil {
	// 	logging.GetLogger(nil).WithField("error", err).Error("cant read request body")
	// } else {
	// 	// write the body back into the request
	// 	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// 	strBody = string(body)
	// 	strBody = strings.Replace(strBody, "\n", "", -1)
	// 	strBody = strings.Replace(strBody, "\t", "", -1)
	// }

	// if len(strBody) > 0 {
	// 	entry = entry.WithField("request_body", strBody)
	// }

	ctx.Next()
}
