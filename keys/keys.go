package keys

import (
	"context"
)

type ContextKey string

var (
	AppEventErrTokenExpired    = "token_expired"
	AppEventErrDataRetrieval   = "data_retrieval_error"
	AppEventErrRefreshingToken = "error_refreshing_token"
	ErrCantFindValue           = "couldnt find requested value in context"
	ContextHost                = ContextKey("host")
	ContextPort                = ContextKey("port")
	ContextDbHost              = ContextKey("db_host")
	ContextDbUser              = ContextKey("db_user")
	ContextDbPass              = ContextKey("db_pass")
	ContextSecurityKey         = ContextKey("sec_key")
	ContextDatabase            = ContextKey("db_name")
	ContextLyricsToken         = ContextKey("genius_access_token")
	ContextLogger              = ContextKey("logger")
	ContextLoggerFields        = ContextKey("logger_fields")
	// ContextSpotifyReturnURL is the key to use for the ouath return url
	ContextSpotifyReturnURL = ContextKey("return_url")
	// ContextSpotifyClientIDContextSpotifyClientID is the key to use for the spotify client id
	ContextSpotifyClientID = ContextKey("client_id")
	// ContextSpotifyClientSecret is the key to use for the spotify client secret
	ContextSpotifyClientSecret = ContextKey("client_secret")
	// ContextSpotifyAccessToken is the key to use for the spotify access token
	ContextSpotifyAccessToken = ContextKey("access_token")
	// ContextSpotifyRefreshToken TODO
	ContextSpotifyRefreshToken = ContextKey("refresh_token")
	// ContextSpotifyTimeRange is the key to use for the spotify search
	ContextSpotifyTimeRange = ContextKey("time_range")
	// ContextSpotifyResults is the key to use to retrieve the results
	ContextSpotifyResults = ContextKey("results")
	ContextSpotifyUserID  = ContextKey("s_user_id")
	ContextMasterKey      = ContextKey("master_key")
	ContextDependencies   = ContextKey("deps")
	ContextSkipCache      = ContextKey("skip_cache")
)

func GetContextValue(ctx context.Context, key ContextKey) interface{} {
	ret := ctx.Value(key)
	if ret != nil {
		return ret
	}

	ret = ctx.Value(string(key))
	if ret != nil {
		return ret
	}

	return nil
}
