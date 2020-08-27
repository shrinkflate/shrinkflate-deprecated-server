package main

import (
	"fmt"
	"github.com/h2non/bimg"
	"log"
	"os"
	"path/filepath"
)

type Image struct {
	id        string
	imageData ImageData
}

func (image Image) Compress() {
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
		return
	}

	// resize image
	newImg, err := img.Resize(size.Width, size.Height)
	if err != nil {
		log.Println("Failed to resize image", image)
		return
	}

	// write file
	fileName := fmt.Sprintf("compressed/%s.%s", image.id, filepath.Ext(image.imageData.Path))
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

func (image *Image) Load() error {
	imageData, err := DB.FindImage(image.id)
	if err != nil {
		return err
	}

	image.imageData = imageData
	return nil
}
