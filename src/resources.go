package main

import (
	"bytes"
	"embed"
	"image"

	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

//go:embed resources/*
var ResourcesFS embed.FS

// Reads file with given filename from embedded resources FS and returns its contents
func ResourceGet(filename string) []byte {
	data, err := ResourcesFS.ReadFile("resources/" + filename)
	if err != nil {
		return nil
	}

	return data
}

// Returns a decoded image from an image file
func ImageFromFile(filename string) image.Image {
	data := ResourceGet(filename)
	if data == nil {
		return nil
	}

	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil
	}

	return img
}

func ResourceGetFont(fontFile string) *sfnt.Font {
	tt, err := opentype.Parse(ResourceGet(fontFile))
	if err != nil {
		return nil
	}

	return tt
}
