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
		if os.Getenv("GO_ENV") == "uat" {
			// what the fuck is this?!
			c.Set(string(keys.ContextSpotifyReturnURL), fmt.Sprint("https://testing-api.", env.Host, "/spotify/oauthreturn"))
		} else {
			c.Set(string(keys.ContextSpotifyReturnURL), fmt.Sprint("https://api.", env.Host, "/spotify/oauth"))
		}

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

	// conStr := fmt.Sprintf(`%s:%s@tcp(%s)/%s`, user, pass, host, dbname)
	// db, err := data.GetLiveDB(conStr)
	// if err != nil {
	// 	logging.GetLogger(c).WithError(err).Error("couldnt connect to database")
	// }

	deps := spotify.Dependencies{
		Client: &http.Client{},
		DB:     nil,
	}

	if os.Getenv("GO_ENV") == "development" {
		ca, err := data.GetLiveCache(c)
		if err != nil {
			logging.GetLogger(c).WithError(err).Error("couldnt get redis cache")
		}
		deps.Cache = ca
	}

	c.Set(string(keys.ContextDependencies),
		&deps)
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
	ip := c.GetHeader("X-Forwarded-For")
	if ip == "" {
		ip = strings.Split(c.Request.RemoteAddr, ":")[0]
	}
	return &logging.LoggerFields{
		UserAgent:   c.Request.UserAgent(),
		Referer:     c.Request.Referer(),
		QueryString: c.Request.URL.RawQuery,
		Path:        c.Request.URL.Path,
		Method:      c.Request.Method,
		ClientIP:    c.ClientIP(),
		NewIP:       ip,
		RequestID:   reqID.String(),
		UserID:      uid,
	}
}

func CORSMiddleware(c *gin.Context) {
	genv := os.Getenv("GO_ENV")
	if genv == "development" || genv == "testing" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	} else {
		ref := c.Request.Referer()
		if strings.HasSuffix(ref, "/") {
			ref = ref[:len(ref)-1]
		}
		acceptable := false
		domains := []string{"testing.spotify-views.com", "www.spotify-views.com", "spotify-views.com"}
		for _, i := range domains {
			if strings.Contains(ref, i) {
				logging.GetLogger(c).WithField("domain", ref).Debug("cors: found acceptable match")
				acceptable = true
				break
			}
		}
		if acceptable {
			c.Writer.Header().Set("Access-Control-Allow-Origin", ref)
		} else {
			logging.GetLogger(c).WithField("host", c.Request.Referer()).Debug("unknown request getting rejected")
		}
	}

	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

func setSpotifyUserID(c *gin.Context) {
	uid, err := c.Cookie(cookieKeyID)
	if err == nil {
		if len(uid) > 0 {
			c.Set(string(keys.ContextSpotifyUserID), uid)
		}
	}

	c.Next()
}

func logRequests(c *gin.Context) {
	logger := logging.GetLogger(c)
	// body, _ := ioutil.ReadAll(c.Request.Body)
	logger.WithField("event", "incoming_request").Info()
	// println(string(body))

	// c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
}
