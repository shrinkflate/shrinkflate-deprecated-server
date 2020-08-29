package main

import (
	"fmt"
	"github.com/discordapp/lilliput"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// get ready to resize image,
// using 8192x8192 maximum resize buffer size
var imageOps = lilliput.NewImageOps(8192)

type LilliputCompressor struct {
}

func (compressor LilliputCompressor) Compress(image *Image) {
	// decoder wants []byte, so read the whole file into a buffer
	inputBuf, err := ioutil.ReadFile(image.imageData.Path)
	if err != nil {
		fmt.Printf("failed to read input file, %s\n", err)
		return
	}

	decoder, err := lilliput.NewDecoder(inputBuf)
	// this error reflects very basic checks,
	// mostly just for the magic bytes of the file to match known image formats
	if err != nil {
		fmt.Printf("error decoding image, %s\n", err)
		return
	}
	defer decoder.Close()

	header, err := decoder.Header()
	// this error is much more comprehensive and reflects
	// format errors
	if err != nil {
		fmt.Printf("error reading image header, %s\n", err)
		return
	}

	// create a buffer to store the output image, 50MB in this case
	outputImg := make([]byte, 50*1024*1024)

	// use user supplied filename to guess output type if provided
	// otherwise don't transcode (use existing type)
	outputType := "." + strings.ToLower(decoder.Description())

	resizeMethod := lilliput.ImageOpsFit

	resizeMethod = lilliput.ImageOpsNoResize

	opts := &lilliput.ImageOptions{
		FileType:             outputType,
		Width:                header.Width(),
		Height:               header.Height(),
		ResizeMethod:         resizeMethod,
		NormalizeOrientation: true,
		EncodeOptions: map[int]int{
			lilliput.JpegQuality:     image.Quality,
			lilliput.JpegProgressive: image.Progressive,
		},
	}

	// resize and transcode image
	outputImg, err = imageOps.Transform(decoder, opts, outputImg)
	if err != nil {
		fmt.Printf("error transforming image, %s\n", err)
		return
	}
	filename := fmt.Sprintf("compressed/%s%s", image.id, filepath.Ext(image.imageData.Path))
	err = ioutil.WriteFile(filename, outputImg, 0777)
	if err != nil {
		fmt.Printf("error writing out resized image, %s\n", err)
		return
	}

	fmt.Printf("image written to %s\n", filename)

	err = os.Remove(image.imageData.Path)
	if err != nil {
		log.Println("Could not remove file", image.imageData.Path)
	}
}
