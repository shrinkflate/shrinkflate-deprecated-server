package main

import (
	"github.com/aerogo/aero"
	"github.com/h2non/bimg"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(2)
	app := aero.New()

	configure(app).Run()
}

func configure(app *aero.Application) *aero.Application {
	buffer, err := bimg.Read("images/image_1.jpg")
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

		err = bimg.Write("images/updated.jpg", newImg)
		if err != nil {
			panic(err)
		}

		return ctx.String("Hello World")
	})

	return app
}
