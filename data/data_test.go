package data

import (
	"context"
	"testing"

	"github.com/bmizerany/assert"
)

func TestEncryption(t *testing.T) {
	key := "testsecretkey"
	ctx := context.WithValue(context.Background(), ContextSecurityKey, key)
	val := "encryptthis"

	encrypted, err := encrypt(ctx, val)
	t.Run("EncryptSuccess", func(t *testing.T) {
		assert.Equal(t, err, nil)
		assert.NotEqual(t, val, encrypted)
	})

	dec, err := decrypt(ctx, encrypted)
	t.Run("DecryptSuccess", func(t *testing.T) {
		assert.Equal(t, nil, err)
		assert.Equal(t, dec, val)
	})
}
