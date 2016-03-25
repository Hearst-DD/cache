package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NOOPNew(t *testing.T) {
	c := NOOPCache()
	_, ok := c.(*noopCache)
	assert.True(t, ok)
}

func Test_NOOPOps(t *testing.T) {
	c := NOOPCache()
	c.Put("foo", "bar", time.Hour)
	val, ok := c.Get("foo")
	assert.Equal(t, NotFound, ok)
	assert.Nil(t, val)
}
