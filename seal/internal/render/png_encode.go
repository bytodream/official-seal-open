package render

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

// 印章图颜色极少（印泥红、透明底、少量文字黑），整图通常不到 200 色。
// 按 32 位 RGBA 直接编码会让每像素白占 4 字节，体积是必要值的数倍
// 颜色数不超过 256 时改走 8 位调色板（PNG color type 3），完全无损且编码更快。
const maxPaletteColors = 256

// PNGBytes 将章面位图编码为 PNG 字节；颜色数不超过 256 时优先走调色板以减小体积。
func PNGBytes(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if paletted, ok := toPaletted(img, maxPaletteColors); ok {
		if err := png.Encode(&buf, paletted); err == nil {
			return buf.Bytes(), nil
		}
		buf.Reset()
	}
	// 颜色过多（如彩色 logo）时回退到通用编码。
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// toPaletted 在颜色数不超过 maxColors 时把图像转成 8 位调色板图。
// 返回 false 表示颜色过多，调用方应回退到通用编码。
func toPaletted(img image.Image, maxColors int) (*image.Paletted, bool) {
	b := img.Bounds()
	if b.Dx() < 1 || b.Dy() < 1 {
		return nil, false
	}

	// 统一按非预乘色取点：PNG 的调色板与 tRNS 都是非预乘语义
	// 而 image.RGBA 内部是预乘存储，直接搬运会二次预乘、颜色变暗。
	var pixAt func(x, y int) color.NRGBA
	switch s := img.(type) {
	case *image.RGBA:
		pixAt = func(x, y int) color.NRGBA {
			row := s.Pix[(y-b.Min.Y)*s.Stride:]
			i := (x - b.Min.X) * 4
			return unpremultiply(row[i], row[i+1], row[i+2], row[i+3])
		}
	case *image.NRGBA:
		pixAt = func(x, y int) color.NRGBA {
			row := s.Pix[(y-b.Min.Y)*s.Stride:]
			i := (x - b.Min.X) * 4
			return color.NRGBA{R: row[i], G: row[i+1], B: row[i+2], A: row[i+3]}
		}
	default:
		pixAt = func(x, y int) color.NRGBA {
			return color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
		}
	}

	key := func(c color.NRGBA) uint32 {
		return uint32(c.R)<<24 | uint32(c.G)<<16 | uint32(c.B)<<8 | uint32(c.A)
	}

	index := make(map[uint32]uint8, maxColors)
	palette := make(color.Palette, 0, maxColors)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			k := key(pixAt(x, y))
			if _, ok := index[k]; ok {
				continue
			}
			if len(palette) >= maxColors {
				return nil, false
			}
			index[k] = uint8(len(palette))
			palette = append(palette, color.NRGBA{
				R: uint8(k >> 24), G: uint8(k >> 16), B: uint8(k >> 8), A: uint8(k),
			})
		}
	}

	dst := image.NewPaletted(b, palette)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		row := dst.Pix[(y-b.Min.Y)*dst.Stride:]
		for x := b.Min.X; x < b.Max.X; x++ {
			row[x-b.Min.X] = index[key(pixAt(x, y))]
		}
	}
	return dst, true
}

// unpremultiply 把预乘 RGBA 还原成非预乘 NRGBA。
// 反预乘算法与标准库 image/png 对 *image.RGBA 的处理一致（常量 m = 0x101*0xffff）。
// 保证调色板路径与直接编码路径逐字节一致，避免舍入差 1 带来的像素差异。
func unpremultiply(r, g, b, a uint8) color.NRGBA {
	switch a {
	case 0:
		return color.NRGBA{}
	case 0xff:
		return color.NRGBA{R: r, G: g, B: b, A: a}
	}
	const m = 0x101 * 0xffff
	aa := uint32(a) * 0x101
	return color.NRGBA{
		R: uint8((uint32(r) * m / aa) >> 8),
		G: uint8((uint32(g) * m / aa) >> 8),
		B: uint8((uint32(b) * m / aa) >> 8),
		A: a,
	}
}

// AsRGBA 将任意图像转为可绘制的 *image.RGBA
func AsRGBA(img image.Image) *image.RGBA {
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}
	if p, ok := img.(*image.Paletted); ok {
		return palettedToRGBA(p)
	}
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(x, y, img.At(x, y))
		}
	}
	return dst
}

// palettedToRGBA 把调色板图展开成预乘 RGBA。调色板颜色本身就是非预乘的
// color.NRGBA，直接查表展开，避免逐像素走 At/Set 的接口与色彩模型开销。
// 注意：image.RGBA 是预乘存储，读取时会反预乘，半透明像素因此有 ±1~2 的
// 固有舍入 —— 需要精确半透明色时请用 AsNRGBA。
func palettedToRGBA(p *image.Paletted) *image.RGBA {
	b := p.Bounds()
	dst := image.NewRGBA(b)
	w, h := b.Dx(), b.Dy()
	tbl := make([][4]uint8, len(p.Palette))
	for i, c := range p.Palette {
		r, g, bl, a := c.RGBA()
		tbl[i] = [4]uint8{uint8(r >> 8), uint8(g >> 8), uint8(bl >> 8), uint8(a >> 8)}
	}
	for y := 0; y < h; y++ {
		row := p.Pix[y*p.Stride:]
		out := dst.Pix[y*dst.Stride:]
		for x := 0; x < w; x++ {
			c := tbl[row[x]]
			i := x * 4
			out[i], out[i+1], out[i+2], out[i+3] = c[0], c[1], c[2], c[3]
		}
	}
	return dst
}

// AsNRGBA 将任意图像转为非预乘的 *image.NRGBA。
// 与 PNG 的存储语义一致，转换完全无损 —— 调色板颜色本身就是 color.NRGBA
// 直接查表拷贝字节即可。
func AsNRGBA(img image.Image) *image.NRGBA {
	if n, ok := img.(*image.NRGBA); ok {
		return n
	}
	if p, ok := img.(*image.Paletted); ok {
		return palettedToNRGBA(p)
	}
	bounds := img.Bounds()
	dst := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(x, y, img.At(x, y))
		}
	}
	return dst
}

// palettedToNRGBA 把调色板图展开成非预乘 NRGBA。
func palettedToNRGBA(p *image.Paletted) *image.NRGBA {
	b := p.Bounds()
	dst := image.NewNRGBA(b)
	w, h := b.Dx(), b.Dy()
	tbl := make([][4]uint8, len(p.Palette))
	for i, c := range p.Palette {
		if n, ok := c.(color.NRGBA); ok {
			tbl[i] = [4]uint8{n.R, n.G, n.B, n.A}
			continue
		}
		r, g, bl, a := c.RGBA()
		tbl[i] = [4]uint8{uint8(r >> 8), uint8(g >> 8), uint8(bl >> 8), uint8(a >> 8)}
	}
	for y := 0; y < h; y++ {
		row := p.Pix[y*p.Stride:]
		out := dst.Pix[y*dst.Stride:]
		for x := 0; x < w; x++ {
			c := tbl[row[x]]
			i := x * 4
			out[i], out[i+1], out[i+2], out[i+3] = c[0], c[1], c[2], c[3]
		}
	}
	return dst
}
