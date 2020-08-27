package main

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShrinkflateDb_New(t *testing.T) {
	// create the db instance
	db, ctx, cancel, err := shrinkflateDb{
		host: "localhost",
		port: 27017,
		name: "shrinkflate_test",
	}.New()

	require.NoError(t, err)
	defer cancel()
	defer func() {
		err := db.conn.Disconnect(ctx)
		require.NoError(t, err)
	}()
}

func TestShrinkflateDb_StoreImage(t *testing.T) {
	// create the db instance
	db, ctx, cancel, err := shrinkflateDb{
		host: "localhost",
		port: 27017,
		name: "shrinkflate_test",
	}.New()
	require.NoError(t, err)
	defer cancel()
	defer func() {
		err := db.conn.Disconnect(ctx)
		require.NoError(t, err)
	}()

	err = db.DeleteAllImages()
	require.NoError(t, err)

	var id string

	{
		id, err = db.StoreImage("/some/path", "https://helloworld.com/")
		require.NoError(t, err)
		require.NotEqual(t, "", id)
	}

	{
		image, err := db.FindImage(id)
		require.NoError(t, err)
		require.NotNil(t, image._id)
	}

	err = db.database.Drop(context.Background())
	require.NoError(t, err)
}
