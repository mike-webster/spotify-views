package data

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mike-webster/spotify-views/keys"
	"github.com/mike-webster/spotify-views/logging"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/jmoiron/sqlx"
)

var (
	_db *Database

	//nonceSize = 12
)

// Database holds a database
type Database struct {
	*sqlx.DB
}

// Ping returns true if we can successfully ping the db
func Ping(ctx context.Context) error {
	ok, err := loadDB(ctx)
	if err != nil || !ok {
		return err
	}

	logger := logging.GetLogger(ctx)
	err = _db.Ping()
	if err != nil {
		logger.WithField("event", "ping_fail").Error(err)

		return err
	}

	logger.WithField("event", "ping_success").Info()
	return nil
}

func SaveUser(ctx context.Context, id string, email string) (bool, error) {
	ok, err := loadDB(ctx)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, errors.New("weird - couldnt connect to databse")
	}

	query := `INSERT INTO users	(spotify_id, email) VALUES (?, ?)`
	res, err := _db.Exec(query, id, email)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
			logging.GetLogger(ctx).WithFields(map[string]interface{}{
				"event": "duplicate_user_insert",
				"email": email,
			}).Error()
			return false, nil
		}
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func loadDB(ctx context.Context) (bool, error) {
	if _db != nil {
		return true, nil
	}

	conn, err := getConnectionString(ctx)
	if err != nil {
		return false, err
	}

	db, err := sqlx.Connect("mysql", conn)
	if err != nil {
		return false, err
	}
	_db = &Database{DB: db}
	return true, nil
}

func getConnectionString(ctx context.Context) (string, error) {
	host := keys.GetContextValue(ctx, keys.ContextDbHost)
	user := keys.GetContextValue(ctx, keys.ContextDbUser)
	pass := keys.GetContextValue(ctx, keys.ContextDbPass)
	dbname := keys.GetContextValue(ctx, keys.ContextDatabase)

	if host == nil || user == nil || pass == nil || dbname == nil {
		return "", errors.New("missing connection string info")
	}

	return fmt.Sprintf(`%s:%s@tcp(%s)/%s`, user, pass, host, dbname), nil
}
