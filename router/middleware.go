package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/mike-webster/spotify-views/data"
	"github.com/mike-webster/spotify-views/env"
	"github.com/mike-webster/spotify-views/keys"
	"github.com/mike-webster/spotify-views/logging"
	spotify "github.com/mike-webster/spotify-views/spotify"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

func setTokens(c *gin.Context) {
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
	c.Set(string(keys.ContextMasterKey), os.Getenv("MASTER_KEY"))
	// vals, err := parseEnvironmentVariables(c)

	secrets, err := env.ParseSecrets(c)
	if err != nil {
		entry.WithError(err).Error("error encountered parsing secrets")
		c.AbortWithError(500, err)
		return
	}

	env, err := env.ParseEnv()
	if err != nil {
		entry.WithError(err).Error("error encountered env vars")
		c.AbortWithError(500, err)
		return
	}

	c.Set(string(keys.ContextSpotifyClientID), secrets.ClientID)
	c.Set(string(keys.ContextSpotifyClientSecret), secrets.ClientSecret)
	c.Set(string(keys.ContextDatabase), secrets.DBName)
	c.Set(string(keys.ContextDbHost), secrets.DBHost)
	c.Set(string(keys.ContextDbUser), secrets.DBUser)
	c.Set(string(keys.ContextDbPass), secrets.DBPass)
	c.Set(string(keys.ContextLyricsToken), secrets.LyricsKey)
	c.Set(string(keys.ContextSecurityKey), secrets.SecurityKey)

	// fix for local dev
	if strings.Contains(env.Host, "localhost") {
		c.Set(string(keys.ContextSpotifyReturnURL), fmt.Sprint("http://", env.Host, ":", env.Port, "/spotify/oauth"))
	} else {
		c.Set(string(keys.ContextSpotifyReturnURL), fmt.Sprint("https://www.", env.Host, "/spotify/oauth"))
	}

	c.Next()
}

func authenticate(c *gin.Context) {
	entry := logging.GetLogger(c)
	tok, err := c.Cookie(cookieKeyToken)
	if err != nil {
		entry.WithField("redirect_reason", "authenticate: found no token").WithError(err).Error("redirecting")
		c.Redirect(301, "/?noauth")
		return
	}

	c.Set(string(keys.ContextSpotifyAccessToken), tok)
	c.Next()
}

func setDependencies(c *gin.Context) {
	//var db *sqlx.DB
	host := keys.GetContextValue(c, keys.ContextDbHost)
	user := keys.GetContextValue(c, keys.ContextDbUser)
	pass := keys.GetContextValue(c, keys.ContextDbPass)
	dbname := keys.GetContextValue(c, keys.ContextDatabase)

	if host == nil || user == nil || pass == nil || dbname == nil {
		report := fmt.Sprintf("%v:%v:%v:%v", host == nil, user == nil, pass == nil, dbname == nil)
		logging.GetLogger(c).Warn(fmt.Sprint("missing connection string info: ", report))
	}

	conStr := fmt.Sprintf(`%s:%s@tcp(%s)/%s`, user, pass, host, dbname)
	db, err := data.GetLiveDB(conStr)
	if err != nil {
		logging.GetLogger(c).WithError(err).Error("missing connection string info")
	}

	c.Set(string(keys.ContextDependencies),
		&spotify.Dependencies{
			Client: &http.Client{},
			DB:     db,
		})
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

func parseLoggerValues(c *gin.Context) *logging.LoggerFields {
	reqID, _ := uuid.NewV4()
	uid, _ := c.Cookie("svid")
	return &logging.LoggerFields{
		UserAgent:   c.Request.UserAgent(),
		Referer:     c.Request.Referer(),
		QueryString: c.Request.URL.RawQuery,
		Path:        c.Request.URL.Path,
		Method:      c.Request.Method,
		ClientIP:    c.ClientIP(),
		RequestID:   reqID.String(),
		UserID:      uid,
	}
}
