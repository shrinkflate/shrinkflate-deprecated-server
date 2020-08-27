package main

import (
	"github.com/aerogo/aero"
	"runtime"
)

var DB *shrinkflateDb
var Cache *shrinkflateCache

func main() {
	runtime.GOMAXPROCS(1)

	// create the DB instance
	db, ctx, cancel, err := shrinkflateDb{
		host: "localhost",
		port: 27017,
		name: "shrinkflate",
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
		host:     "localhost",
		port:     6379,
		password: "k",
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
