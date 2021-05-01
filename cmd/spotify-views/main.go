package main

import (
	"context"

	"github.com/mike-webster/spotify-views/router"
)

func main() {
	router.Run(context.Background())
}
