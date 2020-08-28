package main

import "log"

type Image struct {
	id         string
	imageData  ImageData
	compressor Compressor
	Quality int
	Progressive int
}

func (image Image) Compress() {
	err := image.Load()
	if err != nil {
		log.Println("Could not load image", image)
		return
	}
	image.compressor.Compress(&image)
}

func (image *Image) Load() error {
	imageData, err := DB.FindImage(image.id)
	if err != nil {
		return err
	}

	image.imageData = imageData
	return nil
}

type ImageOpts struct {
	id          string
	compressor  string
	quality     int
	progressive int
}
