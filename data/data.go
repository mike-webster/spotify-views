package data

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/mike-webster/spotify-views/keys"
	"github.com/mike-webster/spotify-views/logging"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/jmoiron/sqlx"
)

var (
	_db *Database

	nonceSize = 12
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

func GetRefreshTokenForUser(ctx context.Context, id string) (string, error) {
	ok, err := loadDB(ctx)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", errors.New("weird - couldnt connect to databse")
	}

	query := `SELECT refresh FROM tokens WHERE spotify_id = %v`
	type token struct {
		Refresh string `db:"refresh"`
	}
	tok := token{}
	err = _db.Get(&tok, fmt.Sprintf(query, id))
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

	logging.GetLogger(ctx).WithField("event", "saving_refresh_token").Info(enc)

	// write query
	query := fmt.Sprintf(`INSERT INTO tokens (spotify_id, refresh) VALUES ('%v','%v') ON DUPLICATE KEY UPDATE`,
		id, enc)
	logging.GetLogger(ctx).WithFields(map[string]interface{}{
		"event": "sql_query",
		"query": query,
	}).Info()

	res, err := _db.Exec(query)
	if err != nil {
		return false, err
	}

	rows, _ := res.RowsAffected()
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
	host := keys.GetContextValue(ctx, keys.ContextDbHost)
	user := keys.GetContextValue(ctx, keys.ContextDbUser)
	pass := keys.GetContextValue(ctx, keys.ContextDbPass)
	dbname := keys.GetContextValue(ctx, keys.ContextDatabase)

	if host == nil || user == nil || pass == nil || dbname == nil {
		return "", errors.New("missing connection string info")
	}

	return fmt.Sprintf(`%s:%s@tcp(%s)/%s`, user, pass, host, dbname), nil
}

func encrypt(ctx context.Context, val string) (string, error) {
	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.
	key := fmt.Sprint(keys.GetContextValue(ctx, keys.ContextSecurityKey))
	logging.GetLogger(ctx).WithFields(map[string]interface{}{
		"event": "checking_key",
		"key":   key,
	}).Info()
	secKey := []byte(key)

	fmt.Println("sekkey: ", secKey)
	plaintext := []byte(val)

	block, err := aes.NewCipher(secKey)
	if err != nil {
		return "", err
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// write the nonce first
	summer := bytes.NewBuffer(nonce)

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	fmt.Println("enc: ", string(ciphertext))

	// write the encrypted token
	summer.Write(ciphertext)

	toStore := summer.Bytes()

	logging.GetLogger(ctx).WithField("event", "webby_test").Info(string(toStore))

	// encode the bytes to a string for storing
	ret := hex.EncodeToString(toStore)

	fmt.Println("hex: ", ret)
	return ret, nil
}

func decrypt(ctx context.Context, val string) (string, error) {
	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.
	secKey := keys.GetContextValue(ctx, keys.ContextSecurityKey)
	fmt.Println("seckey: ", secKey)

	key := []byte(fmt.Sprint(secKey))
	fmt.Println("key: ", key)

	ciphertext, err := hex.DecodeString(val)
	fmt.Println("decoded: ", string(ciphertext))
	if err != nil {
		return "", err
	}

	nonce := ciphertext[:nonceSize]
	encrypted := ciphertext[nonceSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}

	fmt.Printf("plaintext: %s\n", string(plaintext))

	return string(plaintext), nil
}
