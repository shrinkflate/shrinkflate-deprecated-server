package main

import "C"
import (
	"fmt"
	"github.com/h2non/bimg"
	"log"
	"os"
	"path/filepath"
)

type LibVipsCompressor struct {
}

func (compressor LibVipsCompressor) Compress(image* Image) {

	//bimg.VipsCacheSetMax(0)
	//bimg.VipsCacheSetMaxMem(0)

	err := image.Load()
	if err != nil {
		log.Println("Failed to load the image from database", image)
		return
	}

	buffer, err := bimg.Read(image.imageData.Path)
	if err != nil {
		log.Println("Failed to read the image file", image)
		return
	}

	// create new image from buffer
	img := bimg.NewImage(buffer)
	size, err := img.Size()
	if err != nil {
		log.Println("Failed to create new image from buffer", image)
		bimg.VipsCacheDropAll()
		bimg.VipsDebugInfo()
		return
	}

	// resize image
	newImg, err := img.Resize(size.Width, size.Height)
	if err != nil {
		log.Println("Failed to resize image", image)
		return
	}

	// write file
	fileName := fmt.Sprintf("compressed/%s%s", image.id, filepath.Ext(image.imageData.Path))
	err = bimg.Write(fileName, newImg)
	if err != nil {
		log.Println("Failed to write file")
		return
	}

	err = os.Remove(image.imageData.Path)

	if err != nil {
		log.Println("Could not delete file")
		return
	}
}
