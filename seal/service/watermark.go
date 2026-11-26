package service

import (
	"image/color"

	"github.com/bytodream/official-seal-open/seal/service/dto"
)

// WatermarkFunc 在章面上叠加水印文字的回调类型。
type WatermarkFunc func(img dto.RasterImage, c color.Color, content string, size, degree int, fontFile string) error
