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

import "image"

// Resizes given image to multiple dimensions
func GenerateIcons(imgSrc image.Image, dimensions [][2]uint) []image.Image {
	var icons []image.Image
	for _, iDimensions := range dimensions {
		icons = append(icons, ResizeImage(imgSrc, iDimensions))
	}

	return icons
}
