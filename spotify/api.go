package spotify

import (
	"context"
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
