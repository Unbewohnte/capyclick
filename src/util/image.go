package util

import (
	"bytes"
	"image"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/image/draw"
)

// Opens and decodes image
func OpenImage(path string) (image.Image, error) {
	// Open image
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// Decodes image from image bytes
func DecodeImage(data []byte) (image.Image, error) {
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)

	return img, err
}

// Resizes image to specified dimensions
func ResizeImage(imgSrc image.Image, dimensions [2]uint) image.Image {
	imgDst := image.NewRGBA(image.Rect(0, 0, int(dimensions[0]), int(dimensions[1])))

	draw.NearestNeighbor.Scale(imgDst, imgDst.Rect, imgSrc, imgDst.Rect, draw.Over, nil)

	return imgDst
}

// Resizes image to specified dimensions
func ResizeImageOpen(path string, dimensions [2]uint) (image.Image, error) {
	imgSrc, err := OpenImage(path)
	if err != nil {
		return nil, err
	}

	return ResizeImage(imgSrc, dimensions), nil
}
