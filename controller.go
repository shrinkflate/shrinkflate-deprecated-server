package main

import (
	"fmt"
	"github.com/aerogo/aero"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type shrinkflateController struct {
}

func (controller shrinkflateController) Compress(ctx aero.Context) error {
	allowCors(ctx)
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

	compressor := request.Form.Get("compressor")

	quality, _ := strconv.ParseInt(request.Form.Get("quality"), 10, 64)
	progressive, _ := strconv.ParseBool(request.Form.Get("progressive"))
	progressiveVal := 0
	if progressive {
		progressiveVal = 1
	}

	QueueJob(id, compressor, int(quality), progressiveVal)

	return ctx.String(id)
}

func (controller shrinkflateController) Welcome(ctx aero.Context) error {
	allowCors(ctx)

	fileContent, err := ioutil.ReadFile("public/build/index.html")
	if err != nil {
		return ctx.String("Could not prepare response")
	}

	return ctx.HTML(string(fileContent))
}

func (controller shrinkflateController) RootFiles(ctx aero.Context) error {
	file := ctx.Get("file")

	file = fmt.Sprintf("%s%s", "public/build/", file)

	if !fileExists(file) {
		ctx.Response().SetHeader("status", "404 Not Found")
		return ctx.String("Root file not found")
	}

	return ctx.File(file)
}

func (controller shrinkflateController) JSFiles(ctx aero.Context) error {
	file := ctx.Get("file")

	file = fmt.Sprintf("%s%s", "public/build/static/js/", file)

	if !fileExists(file) {
		ctx.Response().SetHeader("status", "404 Not Found")
		return ctx.String("JS not found")
	}

	return ctx.File(file)
}

func (controller shrinkflateController) CSSFiles(ctx aero.Context) error {
	file := ctx.Get("file")

	file = fmt.Sprintf("%s%s", "public/build/static/css/", file)

	if !fileExists(file) {
		ctx.Response().SetHeader("status", "404 Not Found")
		return ctx.String("CSS not found")
	}

	return ctx.File(file)
}

func (controller shrinkflateController) Download(ctx aero.Context) error {
	allowCors(ctx)
	id := ctx.Get("id")

	imageData, err := DB.FindImage(id)
	if err != nil {
		ctx.Response().SetHeader("Status", "404 Not Found")
		return ctx.String("We could not find the image you're looking for")
	}

	return ctx.File(fmt.Sprintf("compressed/%s%s", imageData.Id, filepath.Ext(imageData.Path)))
}

func allowCors(ctx aero.Context) aero.Context {
	response := ctx.Response()
	response.SetHeader("Access-Control-Allow-Origin", "*")

	return ctx
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
