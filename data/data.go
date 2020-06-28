package data

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/jmoiron/sqlx"
)

var (
	_db                *Database
	ContextHost        = ContextKey("db_host")
	ContextUser        = ContextKey("db_user")
	ContextPass        = ContextKey("db_pass")
	ContextSecurityKey = ContextKey("sec_key")
	ContextDatabase    = ContextKey("db_name")
)

// Database holds a database
type Database struct {
	*sqlx.DB
}
type ContextKey string

// Ping returns true if we can successfully ping the db
func Ping(ctx context.Context) error {
	ok, err := loadDB(ctx)
	if err != nil || !ok {
		return err
	}

	err = _db.Ping()
	if err != nil {
		log.Println("error pinging db: ", err.Error())
		return err
	}

	log.Println("successful db ping")
	return nil
}

func GetRefreshTokenForUser(ctx context.Context, id string) (string, error) {
	ok, err := loadDB(ctx)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", errors.New("weird - couldnt connect to databse")
	}

	query := `SELECT refresh FROM tokens WHERE id = ?`
	type token struct {
		Refresh string `db:"refresh"`
	}
	tok := token{}
	err = _db.Get(&tok, query)
	if err != nil {
		return "", err
	}

	// decode token
	decTok, err := decrypt(ctx, tok.Refresh)
	if err != nil {
		return "", err
	}

	return decTok, nil
}

func SaveRefreshTokenForUser(ctx context.Context, tok string, id string) (bool, error) {
	ok, err := loadDB(ctx)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, errors.New("weird - couldnt connect to databse")
	}

	// encrypt token
	enc, err := encrypt(ctx, tok)
	if err != nil {
		return false, err
	}

	// write query
	query := `INSERT IGNORE INTO tokens (spotify_id, refresh) VALUES (?,?)`
	res := _db.MustExec(query, id, enc)
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func SaveUser(ctx context.Context, id string, email string) (bool, error) {
	ok, err := loadDB(ctx)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, errors.New("weird - couldnt connect to databse")
	}

	query := `INSERT IGNORE INTO users	(spotify_id, email) VALUES (?, ?)`
	res := _db.MustExec(query, id, email)
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
	host := ctx.Value(ContextHost)
	user := ctx.Value(ContextUser)
	pass := ctx.Value(ContextPass)
	dbname := ctx.Value(ContextDatabase)

	if host == nil || user == nil || pass == nil || dbname == nil {
		return "", errors.New("missing connection string info")
	}

	return fmt.Sprintf(`%s:%s@tcp(%s)/%s`, user, pass, host, dbname), nil
}

func encrypt(ctx context.Context, val string) (string, error) {
	secKey := ctx.Value(ContextSecurityKey)
	block, err := aes.NewCipher([]byte(createHash(fmt.Sprint(secKey))))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(val), nil)
	return string(ciphertext), nil
}

func decrypt(ctx context.Context, val string) (string, error) {
	data := []byte(val)
	secKey := ctx.Value(ContextSecurityKey)
	key := []byte(createHash(fmt.Sprint(secKey)))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
