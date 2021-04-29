package spotify

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
  )

func TestLog(t *testing.T) {
	t.Run("EmptyImages", func(t *testing.T) {
		a := Album{}
		assert.Equal(t, a.Loc(), "")
	})
	t.Run("MultipleImages", func(t *testing.T) {
		url := "testurl1"
		a := Album{Images: []Image{Image{URL:url}, Image{URL: fmt.Sprint(url, "aa")}}}
		assert.Equal(t, a.Loc(), url)
	})
}