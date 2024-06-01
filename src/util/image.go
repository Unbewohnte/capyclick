/*
  	capyclick - Capybara clicker game
    Copyright (C) 2024  Kasianov Nikolai Alekseevich (Unbewohnte)

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

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
