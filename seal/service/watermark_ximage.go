//go:build watermark

package service

import (
	"image"
	"image/color"
	"math"
	"os"

	"github.com/bytodream/official-seal-open/seal/internal/render"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// DefaultWatermark 为 WatermarkFunc 的默认实现（build tag: watermark）。
func DefaultWatermark(img image.Image, c color.Color, content string, size, degree int, fontFile string) error {
	rgba := render.AsRGBA(img)
	face, err := loadFontFace(fontFile, float64(size))
	if err != nil {
		face = basicfont.Face7x13
	}

	d := &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(c),
		Face: face,
	}

	advance := d.MeasureString(content).Ceil()
	metrics := face.Metrics()
	textH := (metrics.Ascent + metrics.Descent).Ceil()
	x := (rgba.Bounds().Dx() - advance) / 2
	y := (rgba.Bounds().Dy() + textH/2) / 2

	if degree == 0 {
		d.Dot = fixed.P(x, y)
		d.DrawString(content)
		return nil
	}

	cx := float64(rgba.Bounds().Dx()) / 2
	cy := float64(rgba.Bounds().Dy()) / 2
	rad := float64(degree) * math.Pi / 180
	px := float64(x)
	py := float64(y)
	rx := (px-cx)*math.Cos(rad) - (py-cy)*math.Sin(rad) + cx
	ry := (px-cx)*math.Sin(rad) + (py-cy)*math.Cos(rad) + cy
	d.Dot = fixed.P(int(rx), int(ry))
	d.DrawString(content)
	return nil
}

// loadFontFace 读取字体文件并按 size 创建 OpenType font.Face。
func loadFontFace(fontFile string, size float64) (font.Face, error) {
	data, err := os.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}
	tt, err := opentype.Parse(data)
	if err != nil {
		return nil, err
	}
	return opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}
