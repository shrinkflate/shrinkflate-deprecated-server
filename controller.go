package main

import (
	"fmt"
	"github.com/aerogo/aero"
	"io"
	"os"
	"path/filepath"
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
	filename := fmt.Sprintf("images/%s", header.Filename)
	dst, err := os.Create(filename)
	_, err = io.Copy(dst, file)
	if err != nil {
		return ctx.String(err.Error())
	}

	id, err := DB.StoreImage(filename, request.Form.Get("callback"))
	if err != nil {
		return ctx.String(err.Error())
	}

	QueueJob(id)

	return ctx.String(id)
}

func (controller shrinkflateController) Welcome(ctx aero.Context) error {
	return ctx.String("Welcome")
}

func (controller shrinkflateController) Download(ctx aero.Context) error {
	id := ctx.Get("id")

	imageData, err := DB.FindImage(id)
	if err != nil {
		ctx.Response().SetHeader("Status", "404 Not Found")
		return ctx.String("We could not find the image you're looking for")
	}

	return ctx.File(fmt.Sprintf("compressed/%s%s", imageData.Id, filepath.Ext(imageData.Path)))
}
