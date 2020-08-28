package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type shrinkflateController struct {
}

func (controller shrinkflateController) Compress(w http.ResponseWriter, r *http.Request) {
	allowCors(w)

	// get the file handle
	file, header, err := r.FormFile("image")
	if err != nil {
		w.Header().Add("status", "412")
		sendResponse([]byte(err.Error()), w)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	// copy the file into buffer
	filename := fmt.Sprintf("images/%s", header.Filename)
	dst, err := os.Create(filename)
	defer func() {
		_ = dst.Close()
	}()
	_, err = io.Copy(dst, file)
	if err != nil {
		w.Header().Add("status", "500")
		sendResponse([]byte(err.Error()), w)
		return
	}

	id, err := DB.StoreImage(filename, r.Form.Get("callback"))
	if err != nil {
		w.Header().Add("status", "412")
		sendResponse([]byte(err.Error()), w)
		return
	}

	compressor := r.Form.Get("compressor")

	quality, _ := strconv.ParseInt(r.Form.Get("quality"), 10, 64)
	progressive, _ := strconv.ParseBool(r.Form.Get("progressive"))
	progressiveVal := 0
	if progressive {
		progressiveVal = 1
	}

	QueueJob(id, compressor, int(quality), progressiveVal)

	sendResponse([]byte(id), w)
}

func (controller shrinkflateController) Welcome(w http.ResponseWriter, _ *http.Request) {
	allowCors(w)

	fileContent, err := ioutil.ReadFile("public/build/index.html")
	if err != nil {
		log.Println("Could not read index file", err)
	}

	sendResponse(fileContent, w)
}

func (controller shrinkflateController) Download(w http.ResponseWriter, r *http.Request) {
	allowCors(w)
	vars := mux.Vars(r)

	imageData, err := DB.FindImage(vars["id"])
	if err != nil {
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}

	ext := filepath.Ext(imageData.Path)

	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Disposition", "attachment; filename="+imageData.Id+ext)

	fileContent, err := ioutil.ReadFile("compressed/" + imageData.Id + ext)
	if err != nil {
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileContent)))
	w.Header().Set("Content-Type", mime.TypeByExtension(ext))

	sendResponse(fileContent, w)
}

func allowCors(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
}

func sendResponse(text []byte, w http.ResponseWriter) {
	_, err := w.Write(text)
	if err != nil {
		log.Println("Could not send response", err)
	}
}
