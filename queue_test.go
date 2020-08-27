package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShrinkflateQueue_New(t *testing.T) {
	queue, err := getQueue()

	require.NoError(t, err)
	require.NotNil(t, queue.rdb)
}

func TestShrinkflateQueue_Remember(t *testing.T) {
	queue, _ := getQueue()
	_, err := queue.Remember("hello", "world", 0)
	require.NoError(t, err)
}

func TestShrinkflateQueue_Get(t *testing.T) {
	queue, _ := getQueue()

	val, err := queue.Get("hello")
	require.NoError(t, err)
	require.Equal(t, "world", val)
}

func TestShrinkflateQueue_Forget(t *testing.T) {
	queue, _ := getQueue()

	_, err := queue.Forget("hello")
	require.NoError(t, err)

	_, err = queue.Get("hello")
	require.Error(t, err)
}

func getQueue() (shrinkflateQueue, error) {
	return shrinkflateQueue{
		host:     "localhost",
		port:     6379,
		password: "k",
	}.New()
}
