package ui

import (
	"bytes"
	"image"
	"image/png"
	"io/ioutil"
)

// simple png image decoder for drawing

func getImg(path string) image.Image {
	var data, err = ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var img image.Image
	img, err = png.Decode(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	return img
}
