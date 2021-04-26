package spotify

type ErrTokenExpired string
type ErrBadRequest string
type ErrNoToken string

func (e ErrTokenExpired) Error() string {
	return string(e)
}
func (e ErrNoToken) Error() string {
	return string(e)
}
func (e ErrBadRequest) Error() string {
	return string(e)
}

var (

	// EventNeedsRefreshToken holds the key to log when a user needs a to
	// refresh their session
	EventNeedsRefreshToken = "token_needs_refresh"
	// EventNon200Response holds the key to log when an external request
	// comes back with a non-200 response
	EventNon200Response = "non_200_response"
)

type TimeFrame int

const (
	TFShort TimeFrame = iota
	TFMedium
	TFLong
)

func (t TimeFrame) Value() string {
	switch t {
	case TFShort:
		return "short_term"
	case TFMedium:
		return "medium_term"
	case TFLong:
		return "long_term"
	default:
		return ""
	}
}

func GetTimeFrame(str string) TimeFrame {
	switch str {
	case "medium_term":
		return TFMedium
	case "long_term":
		return TFLong
	default:
		return TFShort
	}
}
