package main

import (
	"bytes"
	"fmt"
	"github.com/aerogo/aero"
	"github.com/h2non/bimg"
	"io"
	"strings"
)

type shrinkflateController struct {
}

func (controller shrinkflateController) Compress(ctx aero.Context) error {
	request := ctx.Request().Internal()

	// parse the request
	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		return ctx.String(err.Error())
	}

	// get the file handle
	file, header, err := request.FormFile("image")
	if err != nil {
		return ctx.String(err.Error())
	}
	defer func() {
		_ = file.Close()
	}()

	// copy the file into buffer
	name := strings.Split(header.Filename, ".")
	fmt.Println("Filename", name)

	var buffer bytes.Buffer
	s, err := io.Copy(&buffer, file)
	if err != nil {
		return ctx.String(err.Error())
	}

	// create new image from buffer
	img := bimg.NewImage(buffer.Bytes())
	size, err := img.Size()
	if err != nil {
		return ctx.String(err.Error())
	}

	// resize image
	newImg, err := img.Resize(size.Width, size.Height)
	if err != nil {
		return ctx.String(err.Error())
	}

	// write file
	fileName := fmt.Sprintf("%s%d%s", "images/updated_", s, ".jpg")
	err = bimg.Write(fileName, newImg)
	if err != nil {
		return ctx.String(err.Error())
	}

	return ctx.String(fileName)
}

func (controller shrinkflateController) Welcome(ctx aero.Context) error {
	QueueJob("New ID")
	return ctx.String("Welcome")
}
