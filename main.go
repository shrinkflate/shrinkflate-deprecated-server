package main

import (
	"fmt"
	"github.com/aerogo/aero"
	"github.com/h2non/bimg"
	"runtime"
)

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

var n = 0

func configure(app *aero.Application) *aero.Application {
	buffer, err := bimg.Read("images/image_2.jpg")
	if err != nil {
		panic(err)
	}

	app.Get("/", func(ctx aero.Context) error {
		img := bimg.NewImage(buffer)
		size, err := img.Size()
		if err != nil {
			panic(err)
		}

		newImg, err := img.Resize(size.Width, size.Height)
		if err != nil {
			panic(err)
		}

		err = bimg.Write(fmt.Sprintf("%s%d%s", "images/updated_", n, ".jpg"), newImg)
		n += 1
		if err != nil {
			panic(err)
		}

		return ctx.String("Hello World")
	})

	return app
}
