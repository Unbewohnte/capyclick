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

package resources

import (
	"bytes"
	"embed"
	"image"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

//go:embed resources/*
var ResourcesFS embed.FS

// Reads file with given filename from embedded resources FS and returns its contents
func Get(filename string) []byte {
	data, err := ResourcesFS.ReadFile("resources/" + filename)
	if err != nil {
		return nil
	}

	return data
}

// Returns a decoded image from an image file
func ImageFromFile(filename string) image.Image {
	data := Get(filename)
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

func GetFont(fontFile string) *sfnt.Font {
	tt, err := opentype.Parse(Get(fontFile))
	if err != nil {
		return nil
	}

	return tt
}

func GetAudioPlayer(audioContext *audio.Context, audioFile string) *audio.Player {
	data := bytes.NewReader(Get(audioFile))
	player, err := audioContext.NewPlayer(data)
	if err != nil {
		return nil
	}

	return player
}
