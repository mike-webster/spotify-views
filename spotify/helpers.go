package spotify

func getPairs(m map[string]int32) Pairs {
	ret := Pairs{}
	for k, v := range m {
		ret = append(ret, Pair{Key: k, Value: v})
	}

	return ret
}

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
