package util

import "image"

// Resizes given image to multiple dimensions
func GenerateIcons(imgSrc image.Image, dimensions [][2]uint) []image.Image {
	var icons []image.Image
	for _, iDimensions := range dimensions {
		icons = append(icons, ResizeImage(imgSrc, iDimensions))
	}

	return icons
}
