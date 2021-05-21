package encrypt

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"github.com/mike-webster/spotify-views/keys"
)

func Encrypt(ctx context.Context, data []byte) (*[]byte, error) {
	hash := keys.GetContextValue(ctx, keys.ContextMasterKey)
	if hash == nil {
		return nil, errors.New("missing master key")
	}
	block, _ := aes.NewCipher([]byte(fmt.Sprint(hash)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return &ciphertext, nil
}

func Decrypt(ctx context.Context, data []byte) (*[]byte, error) {
	hash := keys.GetContextValue(ctx, keys.ContextMasterKey)
	if hash == nil {
		return nil, errors.New("missing master key")
	}
	block, err := aes.NewCipher([]byte(fmt.Sprint(hash)))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return &plaintext, nil
}
