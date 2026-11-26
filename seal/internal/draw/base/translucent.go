package base

import (
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/bytodream/official-seal-open/seal/model"
)

// Translucent 将章面色带写入指定 Alpha，白色背景像素丢弃。
// 避免把预乘样本当成非预乘写回导致 PNG 编码反预乘时溢出（红色 255 会退化成 15）。
func Translucent(src image.Image, transparency int) *image.RGBA {
	b := src.Bounds()
	out := image.NewRGBA(b)
	// 钳制到 uint8 范围，避免越界值回绕（如 -127 → 129）。
	if transparency < 0 {
		transparency = 0
	}
	if transparency > 255 {
		transparency = 255
	}
	alpha := uint8(transparency)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := color.NRGBAModel.Convert(src.At(x, y)).(color.NRGBA)
			if c.A == 0 {
				continue
			}
			if c.R > 225 && c.G > 225 && c.B > 225 {
				continue
			}
			out.Set(x, y, color.NRGBA{R: c.R, G: c.G, B: c.B, A: alpha})
		}
	}
	return out
}

// TransparencyFromSealString 将 0–100 透明度字符串换算为 Alpha（0 不透明、100 全透明）。
// 非法或越界输入钳制到区间内，避免产生回绕的 Alpha 值。
func TransparencyFromSealString(s *string) int {
	if s == nil {
		return model.DefaultSealTransparency
	}
	t, err := strconv.Atoi(strings.TrimSpace(*s))
	if err != nil {
		return model.DefaultSealTransparency
	}
	if t < 0 {
		t = 0
	}
	if t > 100 {
		t = 100
	}
	return int(math.Round(float64(100-t) * 255 / 100))
}
