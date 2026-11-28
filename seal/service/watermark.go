package service

import (
	"image"
	"image/color"
)

// WatermarkFunc 在章面上叠加水印文字的回调类型。
type WatermarkFunc func(img image.Image, c color.Color, content string, size, degree int, fontFile string) error
