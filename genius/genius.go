package genius

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bbalet/stopwords"
	lyrics "github.com/rhnvrm/lyric-api-go"
)

// ContextKey is used to store and access information from the context
type ContextKey string

// ContextAccessToken is the key to use for the genius access token
var ContextAccessToken = ContextKey("access_token")

// LyricSearch holds the information for which a lyric search is desired
type LyricSearch struct {
	Artist string
	Track  string
}

type tempResp struct {
	Response struct {
		Hits []struct {
			Result struct {
				FullTitle string `json:"full_title"`
				ID        int32  `json:"id"`
				Artist    struct {
					Name string `json:"name"`
					ID   int32  `json:"id"`
				} `json:"primary_artist"`
			} `json:"result"`
		} `json:"hits"`
	} `json:"response"`
}

// GetLyricCountForSong will retrieve the song lyrics for all of the provided searches
// and return a map of each word with a value of how many times it occurred.
func GetLyricCountForSong(ctx context.Context, searches []LyricSearch) (map[string]int, error) {
	token := ctx.Value(ContextAccessToken)
	if token == nil {
		return nil, errors.New("no access token provided")
	}
	maps := []map[string]int{}
	l := lyrics.New(lyrics.WithoutProviders(), lyrics.WithGeniusLyrics(fmt.Sprint(token)))
	for _, i := range searches {
		lyric, err := l.Search(i.Artist, i.Track)
		if err != nil {
			return nil, err
		}
		maps = append(maps, convertToMap(lyric))

	}

	return combineMaps(maps), nil
}

func convertToMap(lyric string) map[string]int {
	ret := map[string]int{}

	log.Println("====\nuntouched:\n", lyric, "\n\n\n\n\n\n\n ")

	treated := strings.TrimSpace(strings.Replace(lyric, "\n", " ", -1))
	pattern := `.*\[{1}.*\].*`
	match, err := regexp.Match(pattern, []byte(treated))
	if err != nil {
		panic(err)
	}
	for match {
		start := strings.Index(treated, "[")
		ending := strings.Index(treated, "]")
		treated = fmt.Sprint(treated[:start], treated[ending+1:])

		match, err = regexp.Match(pattern, []byte(treated))
		if err != nil {
			fmt.Println("error: ", err.Error())
			match = false
		}
	}

	for _, ii := range strings.Split(treated, " ") {
		trimmed := strings.TrimSpace(ii)
		if strings.Replace(trimmed, " ", "", -1) == "" {
			continue
		}

		replacer := strings.NewReplacer(",", "", ".", "", ";", "", ")", "", "Intro", "", "Pre-Chorus", "", "[", "", "?", "", "]", "", "(", "", "Verse", "", "'", "", "Chorus", "")
		lyricsString := replacer.Replace(trimmed)
		cleaned := stopwords.CleanString(lyricsString, "en", true)

		log.Println("parsing: ", ii, " ========> ", cleaned)

		for _, j := range strings.Split(cleaned, " ") {
			if len(strings.TrimSpace(j)) < 2 {
				continue
			}
			if _, ok := ret[j]; !ok {
				ret[j] = 1
			} else {
				ret[j]++
			}
		}

	}

	return ret
}

func combineMaps(maps []map[string]int) map[string]int {
	ret := map[string]int{}
	for _, i := range maps {
		for k, v := range i {
			ret[k] += v
		}
	}

	// filter down to only words with more than one reference
	for k, v := range ret {
		if v == 1 {
			delete(ret, k)
		}
	}

	return ret
}
