package main

type Compressor interface {
	Compress(image *Image)
}
