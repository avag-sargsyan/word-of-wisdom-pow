package cache_test

import (
	"fmt"
	"testing"

	"github.com/avag-sargsyan/word-of-wisdom-pow/pkg/cache"
	"github.com/stretchr/testify/assert"
)

// Black-box testing for Cache package
func TestCache(t *testing.T) {
	store := cache.NewStore()
	key := 123
	secret := "test"

	store.Put(key, secret)

	s, err := store.Get(key)

	assert.NoError(t, err)
	assert.Equal(t, s, secret)

	store.Delete(key)

	s1, err1 := store.Get(key)

	assert.Equal(t, s1, "")
	assert.Error(t, err1, fmt.Errorf("has no entry for the given key: %d", key))
}
