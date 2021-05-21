package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/mike-webster/spotify-views/data"
	"github.com/mike-webster/spotify-views/keys"
	"github.com/mike-webster/spotify-views/router"
)

func main() {
	ctx := context.WithValue(context.Background(), keys.ContextMasterKey, os.Getenv("MASTER_KEY"))
	args := os.Args
	if len(args) > 1 && args[1] == "-db_init" {
		dbInit(ctx, args)
		return
	}
	router.Run(context.Background())
}

func dbInit(ctx context.Context, args []string) {
	if len(args) < 4 {
		panic("incorrect nuber of args, please provide root user and pass")
	}

	user := strings.Replace(args[2], "-u=", "", 1)
	pass := strings.Replace(args[3], "-p=", "", 1)

	err := data.DBInit(ctx, user, pass)
	if err != nil {
		panic(err)
	}

	fmt.Println("db init success!")
}
