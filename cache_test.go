package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShrinkflateCache_New(t *testing.T) {
	cache, err := cache()

	require.NoError(t, err)
	require.NotNil(t, cache.rdb)
}

func TestShrinkflateCache_Remember(t *testing.T) {
	cache, _ := cache()
	_, err := cache.Remember("hello", "world", 0)
	require.NoError(t, err)
}

func TestShrinkflateCache_Get(t *testing.T) {
	cache, _ := cache()

	val, err := cache.Get("hello")
	require.NoError(t, err)
	require.Equal(t, "world", val)
}

func TestShrinkflateCache_Forget(t *testing.T) {
	cache, _ := cache()

	_, err := cache.Forget("hello")
	require.NoError(t, err)

	_, err = cache.Get("hello")
	require.Error(t, err)
}

func cache() (shrinkflateCache, error) {
	return shrinkflateCache{
		host:     "localhost",
		port:     6379,
		password: "k",
	}.New()
}
