package util

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

func NewFont(tt *sfnt.Font, options *opentype.FaceOptions) font.Face {
	newFont, _ := opentype.NewFace(tt, options)
	return newFont
}
