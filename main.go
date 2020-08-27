package main

import (
	"github.com/aerogo/aero"
	"github.com/joho/godotenv"
	"os"
)

var DB *shrinkflateDb
var Cache *shrinkflateCache

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// create the DB instance
	db, ctx, cancel, err := shrinkflateDb{
		host: os.Getenv("MONGO_HOST"),
		port: os.Getenv("MONGO_PORT"),
		name: os.Getenv("MONGO_DB"),
	}.New()

	DB = db

	if err != nil {
		panic(err)
	}
	defer cancel()
	defer func() {
		if err = DB.conn.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	cache, err := shrinkflateCache{
		host:     os.Getenv("REDIS_HOST"),
		port:     os.Getenv("REDIS_PORT"),
		password: os.Getenv("REDIS_PASS"),
	}.New()

	Cache = cache

	err = PrepareQueueHandler()
	if err != nil {
		panic(err)
	}

	// initiate the app
	app := aero.New()

	configure(app).Run()
}

func configure(app *aero.Application) *aero.Application {

	controller := shrinkflateController{}

	app.Get("/", controller.Welcome)
	app.Post("/compress", controller.Compress)
	app.Get("/download/:id", controller.Download)

	return app
}
