package spotify

import (
	"context"
	"errors"
	"fmt"
	"log"
)

func HandleOauth(ctx context.Context, code string) (context.Context, error) {
	tokens, err := requestTokens(ctx, code)
	if err != nil {
		log.Println("could not retrieve tokens for user; error: ", err)
		return ctx, err
	}
	log.Println(fmt.Sprint("success - tokens: \n\tAccess: ", tokens[0], "\n\tRefres: ", tokens[1]))
	ctx = context.WithValue(ctx, "access_token", tokens[0])
	ctx = context.WithValue(ctx, "refresh_token", tokens[1])
	return ctx, nil
}

func GetTopTracks(ctx context.Context, limit int32) (*Tracks, error) {
	tracks, err := getTopTracks(ctx, limit)
	if err != nil {
		return nil, err
	}
	return &tracks, nil
}

func GetTopArtists(ctx context.Context) (*Artists, error) {
	artists, err := getTopArtists(ctx)
	if err != nil {
		return nil, err
	}
	return artists, nil
}

func GetGenresForArtists(ctx context.Context, ids []string) (*Pairs, error) {
	log.Println("getting ", len(ids), " artists for genres")
	artists, err := getArtists(ctx, ids)
	if err != nil {
		return nil, err
	}
	log.Println("checking genres for ", len(*artists), " artists")

	ret := map[string]int32{}
	for _, i := range *artists {
		for _, ii := range i.Genres {
			if _, ok := ret[ii]; ok {
				ret[ii]++
			} else {
				ret[ii] = 1
			}
		}
	}

	p := getPairs(ret)
	log.Println("\n\npairs\n", p)
	return &p, nil
}

func GetGenresForTracks(ctx context.Context, ids []string) (*Pairs, error) {
	log.Println("getting ", len(ids), " tracks for genres; ", ids)
	tracks, err := getTracks(ctx, ids)
	if err != nil {
		return nil, err
	}

	as := map[string]int32{}
	aids := []string{}

	log.Println("searching ", len(*tracks), "tracks for distinct artists")

	for _, i := range *tracks {
		log.Println("track: ", i.Name)
		for _, ii := range i.Artists {
			log.Println("artist: ", ii.Name)
			if _, ok := as[ii.Name]; !ok {
				as[ii.Name] = 1
				aids = append(aids, ii.ID)
				fmt.Println("adding ", ii.Name)
			}
		}
	}

	if len(aids) < 1 {
		return nil, errors.New(fmt.Sprint("no artists found for ", len(ids), "tracks"))
	}

	log.Println(len(aids), " distinct artists found from top tracks")

	artists, err := getArtists(ctx, aids)
	if err != nil {
		return nil, err
	}

	log.Println("checking genres for ", len(*artists), " tracks")

	ret := map[string]int32{}
	for _, i := range *artists {
		for _, ii := range i.Genres {
			if _, ok := ret[ii]; ok {
				ret[ii]++
			} else {
				ret[ii] = 1
			}
		}
	}

	p := getPairs(ret)
	log.Println("\n\npairs\n", p)
	return &p, nil
}
