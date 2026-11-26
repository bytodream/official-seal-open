package render

import (
	"github.com/bytodream/official-seal-open/seal/internal/draw/base"
	"github.com/bytodream/official-seal-open/seal/model"

	"image"
	imagestd "image/draw"
)

// maxSampleTries 单个像素拒绝采样的最大尝试次数，用于保证极端半径下也能终止。
const maxSampleTries = 64

// DefaultThickenImage 使用默认半径对章面做模糊加粗。
func DefaultThickenImage(img model.RasterImage, radius *int) image.Image {
	r := 60
	if radius != nil {
		r = *radius
	}
	num := float64(r*30) * 0.01
	return ThickenImage(img, num)
}

// ThickenImage 按给定半径对章面做模糊加粗。
func ThickenImage(src model.RasterImage, radius float64) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	imagestd.Draw(dst, bounds, src, bounds.Min, imagestd.Src)

	rng := base.NewRand()
	quantum := radius
	w := bounds.Dx()
	h := bounds.Dy()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var xDist, yDist int
			// 拒绝采样最多尝试 maxSampleTries 次，避免半径远大于画布时长时间空转。
			for tries := 0; ; tries++ {
				xDist = int((2.0*radius+1.0)*rng.Float64() - quantum)
				yDist = int((2.0*radius+1.0)*rng.Float64() - quantum)
				sx, sy := x+xDist, y+yDist
				if sx >= 0 && sy >= 0 && sx < w && sy < h {
					break
				}
				if tries >= maxSampleTries {
					// 退化到最近的有效边界像素，保证一定终止。
					xDist, yDist = 0, 0
					break
				}
			}
			dst.Set(x, y, dst.At(x+xDist, y+yDist))
		}
	}
	return dst
}
