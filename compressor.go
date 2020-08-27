package main

import "fmt"

type Image struct {
	id string
}

var n = 0

func (image Image) Compress() {
	fmt.Printf("%d. Compressing: %s\n", n, image.id)
	n += 1
}
