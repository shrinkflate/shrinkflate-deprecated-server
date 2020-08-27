package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type ImageData struct {
	Id       string `bson:"_id"`
	Path     string `bson:"path"`
	Callback string `bson:"callback"`
	Time     int64  `bson:"time"`
}

type shrinkflateDb struct {
	conn     *mongo.Client
	host     string
	port     int
	name     string
	database *mongo.Database
}

func (db shrinkflateDb) New() (*shrinkflateDb, context.Context, context.CancelFunc, error) {
	conn, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", db.host, db.port)))
	if err != nil {
		return &shrinkflateDb{}, nil, nil, err
	}

	db.conn = conn
	db.database = conn.Database(db.name)

	// create the context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// connect the database
	err = db.conn.Connect(ctx)
	if err != nil {
		return &shrinkflateDb{}, nil, nil, err
	}

	// make sure we're connected
	err = db.conn.Ping(ctx, readpref.Primary())
	if err != nil {
		return &shrinkflateDb{}, nil, nil, err
	}

	return &db, ctx, cancel, err
}

func (db shrinkflateDb) StoreImage(path string, callback string) (string, error) {
	collection, ctx, cancel := db.getCollectionWithContext("images")
	defer cancel()

	result, err := collection.InsertOne(ctx, bson.M{
		"path":     path,
		"callback": callback,
		"time":     time.Now().Unix(),
	})
	if err != nil {
		return "", err
	}

	objectId := result.InsertedID.(primitive.ObjectID)
	return objectId.Hex(), nil
}

func (db shrinkflateDb) FindImage(id string) (ImageData, error) {
	collection, ctx, cancel := db.getCollectionWithContext("images")
	defer cancel()

	image := ImageData{}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return image, err
	}

	cursor := collection.FindOne(ctx, bson.M{"_id": objectId})

	err = cursor.Decode(&image)

	return image, err
}

func (db shrinkflateDb) DeleteAllImages() error {
	collection, ctx, cancel := db.getCollectionWithContext("images")
	defer cancel()

	return collection.Drop(ctx)
}

func (db shrinkflateDb) getCollectionWithContext(collection string) (*mongo.Collection, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return db.database.Collection(collection), ctx, cancel
}
