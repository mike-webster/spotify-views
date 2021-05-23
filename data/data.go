package data

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/jmoiron/sqlx"
	"github.com/mike-webster/spotify-views/env"
)

func DBInit(ctx context.Context, user, pass string) error {
	env, err := env.ParseSecrets(ctx)
	if err != nil {
		return err
	}
	sql := "./data/create_db.sql"
	f, err := ioutil.ReadFile(sql)
	if err != nil {
		return err
	}

	conStr := fmt.Sprintf(`%s:%s@tcp(%s)/`, user, pass, env.DBHost)
	db, err := sqlx.Connect("mysql", conStr)
	if err != nil {
		return err
	}

	for _, i := range strings.Split(string(f), ";") {
		if len(strings.TrimSpace(i)) < 1 {
			continue
		}
		fmt.Println("executing: ", i)
		_, err = db.Exec(i)
		if err != nil {
			return err
		}
	}

	return nil
}
