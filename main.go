package main

import (
	"github.com/aerogo/aero"
	"runtime"
)

var db shrinkflateDb

func main() {
	runtime.GOMAXPROCS(1)

	// create the db instance
	db, ctx, cancel, err := shrinkflateDb{
		host: "localhost",
		port: 27017,
		name: "shrinkflate",
	}.New()
	if err != nil {
		panic(err)
	}
	defer cancel()
	defer func() {
		if err = db.conn.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// initiate the app
	app := aero.New()

	configure(app).Run()
}

func configure(app *aero.Application) *aero.Application {

	controller := shrinkflateController{}

	app.Get("/", controller.Welcome)
	app.Post("/compress", controller.Compress)

	return app
}
