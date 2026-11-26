package render

import (
	"image"
	"image/color"
)

// ImageDiffResult 图像逐像素对比结果。
type ImageDiffResult struct {
	// Width 图像宽度（px）
	Width int
	// Height 图像高度（px）
	Height int
	// MismatchedPixels 超出容差的像素数
	MismatchedPixels int
	// TotalPixels 参与比较的像素总数（尺寸不一致时为 -1）
	TotalPixels int
	// MaxChannelDelta 最大单通道差值
	MaxChannelDelta uint8
}

// CompareImages 逐像素比较；tolerance 为每通道允许偏差（0=严格相等）。
func CompareImages(a, b image.Image, tolerance uint8) ImageDiffResult {
	ar := AsRGBA(a)
	br := AsRGBA(b)
	ba, bb := ar.Bounds(), br.Bounds()
	out := ImageDiffResult{
		Width:  ba.Dx(),
		Height: ba.Dy(),
	}
	if ba != bb {
		out.TotalPixels = -1
		return out
	}
	out.TotalPixels = ba.Dx() * ba.Dy()
	var maxDelta uint8
	for y := ba.Min.Y; y < ba.Max.Y; y++ {
		for x := ba.Min.X; x < ba.Max.X; x++ {
			c1 := color.RGBAModel.Convert(ar.At(x, y)).(color.RGBA)
			c2 := color.RGBAModel.Convert(br.At(x, y)).(color.RGBA)
			d := max4(
				absU8(c1.R, c2.R),
				absU8(c1.G, c2.G),
				absU8(c1.B, c2.B),
				absU8(c1.A, c2.A),
			)
			if d > maxDelta {
				maxDelta = d
			}
			if d > tolerance {
				out.MismatchedPixels++
			}
		}
	}
	out.MaxChannelDelta = maxDelta
	return out
}

// absU8 返回 uint8 差值绝对值。
func absU8(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}

// max4 返回四个 uint8 的最大值。
func max4(a, b, c, d uint8) uint8 {
	m := a
	if b > m {
		m = b
	}
	if c > m {
		m = c
	}
	if d > m {
		m = d
	}
	return m
}
