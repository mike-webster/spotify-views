package spotify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	t.Run("expired", func(t *testing.T) {
		err := ErrTokenExpired("expired")
		assert.Equal(t, "expired", err.Error())
	})

	t.Run("missing", func(t *testing.T) {
		err := ErrNoToken("missing")
		assert.Equal(t, "missing", err.Error())
	})
}

func TestTimeFrames(t *testing.T) {
	t.Run("Value", func(t *testing.T) {
		t.Run("short", func(t *testing.T) {
			assert.Equal(t, "short_term", TFShort.Value())
		})
		t.Run("medium", func(t *testing.T) {
			assert.Equal(t, "medium_term", TFMedium.Value())
		})
		t.Run("long", func(t *testing.T) {
			assert.Equal(t, "long_term", TFLong.Value())
		})
	})

	t.Run("GetTimeFrame", func(t *testing.T) {
		t.Run("medium", func(t *testing.T) {
			assert.Equal(t, TFMedium, GetTimeFrame("medium_term"))
		})
		t.Run("long", func(t *testing.T) {
			assert.Equal(t, TFLong, GetTimeFrame("long_term"))
		})
		t.Run("short", func(t *testing.T) {
			assert.Equal(t, TFShort, GetTimeFrame("short_term"))
		})
		t.Run("default", func(t *testing.T) {
			assert.Equal(t, TFShort, GetTimeFrame("fdasfdasf"))
		})
	})
}
