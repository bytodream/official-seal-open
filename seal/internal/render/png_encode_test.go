package render

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

// 调色板编码必须与直接 RGBA 编码逐像素一致，否则印章颜色会漂移。
func TestPNGBytesPaletteRoundTrip(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 96, 64))
	red := color.RGBA{R: 216, G: 48, B: 42, A: 255}
	darkRed := color.RGBA{R: 108, G: 24, B: 21, A: 180}
	for y := 0; y < 64; y++ {
		for x := 0; x < 96; x++ {
			switch {
			case y < 32 && x < 48:
				img.Set(x, y, red)
			case y < 32:
				img.Set(x, y, color.RGBA{})
			default:
				img.Set(x, y, darkRed)
			}
		}
	}

	got, err := PNGBytes(img)
	if err != nil {
		t.Fatal(err)
	}

	var want bytes.Buffer
	if err := png.Encode(&want, img); err != nil {
		t.Fatal(err)
	}

	paletted, err := png.Decode(bytes.NewReader(got))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := paletted.(*image.Paletted); !ok {
		t.Fatalf("颜色数少时应走调色板编码，实际为 %T", paletted)
	}
	direct, err := png.Decode(bytes.NewReader(want.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r1, g1, b1, a1 := direct.At(x, y).RGBA()
			r2, g2, b2, a2 := paletted.At(x, y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				t.Fatalf("像素 (%d,%d) 不一致: 直接编码 %v 调色板 %v",
					x, y, []uint32{r1, g1, b1, a1}, []uint32{r2, g2, b2, a2})
			}
		}
	}
	if len(got) >= want.Len() {
		t.Errorf("调色板编码未变小: 调色板 %d 直接 %d", len(got), want.Len())
	}
}

// 颜色数超过 256 时必须回退到通用编码，不能产出错误的调色板图。
func TestPNGBytesFallbackBeyond256Colors(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x * 8), G: uint8(y * 8), B: uint8(x + y), A: 255})
		}
	}

	got, err := PNGBytes(img)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := png.Decode(bytes.NewReader(got))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := decoded.(*image.Paletted); ok {
		t.Fatal("颜色数超过 256 却仍走了调色板编码")
	}

	var want bytes.Buffer
	if err := png.Encode(&want, img); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want.Bytes()) {
		t.Error("回退路径应与通用编码输出一致")
	}
}

// 调色板图经 AsRGBA 展开的结果必须与通用 Set 路径逐像素一致。
func TestAsRGBAPalettedFastPath(t *testing.T) {
	pal := image.NewPaletted(image.Rect(0, 0, 64, 48), color.Palette{
		color.NRGBA{},
		color.NRGBA{R: 216, G: 48, B: 42, A: 255},
		color.NRGBA{R: 108, G: 24, B: 21, A: 180},
	})
	for y := 0; y < 48; y++ {
		for x := 0; x < 64; x++ {
			switch {
			case y < 16:
				pal.SetColorIndex(x, y, 1)
			case y < 32:
				pal.SetColorIndex(x, y, 2)
			}
		}
	}

	got := AsRGBA(pal)
	want := image.NewRGBA(pal.Bounds())
	for y := 0; y < 48; y++ {
		for x := 0; x < 64; x++ {
			want.Set(x, y, pal.At(x, y))
		}
	}
	for y := 0; y < 48; y++ {
		for x := 0; x < 64; x++ {
			r1, g1, b1, a1 := want.At(x, y).RGBA()
			r2, g2, b2, a2 := got.At(x, y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				t.Fatalf("像素 (%d,%d) 不一致: 通用 %v 快速路径 %v",
					x, y, []uint32{r1, g1, b1, a1}, []uint32{r2, g2, b2, a2})
			}
		}
	}
}

// AsNRGBA 与 PNG 存储语义一致，对半透明像素也必须无损。
func TestAsNRGBAPalettedLossless(t *testing.T) {
	pal := image.NewPaletted(image.Rect(0, 0, 64, 48), color.Palette{
		color.NRGBA{},
		color.NRGBA{R: 216, G: 48, B: 42, A: 255},
		color.NRGBA{R: 108, G: 24, B: 21, A: 180},
	})
	for y := 0; y < 48; y++ {
		for x := 0; x < 64; x++ {
			switch {
			case y < 16:
				pal.SetColorIndex(x, y, 1)
			case y < 32:
				pal.SetColorIndex(x, y, 2)
			}
		}
	}

	got := AsNRGBA(pal)
	for i, c := range pal.Palette {
		n, ok := c.(color.NRGBA)
		if !ok {
			t.Fatalf("调色板项 %d 不是 color.NRGBA: %T", i, c)
		}
		// 找到任意一个使用该索引的像素，比对颜色
		for y := 0; y < 48; y++ {
			found := false
			for x := 0; x < 64; x++ {
				if pal.ColorIndexAt(x, y) == uint8(i) {
					// At().RGBA() 返回预乘后的 16 位值，期望值取同一形式
					wr, wg, wb, wa := n.RGBA()
					r, g, b, a := got.At(x, y).RGBA()
					if r != wr || g != wg || b != wb || a != wa {
						t.Errorf("调色板项 %d 展开后颜色不符: 期望 (%d,%d,%d,%d) 实际 (%d,%d,%d,%d)",
							i, wr>>8, wg>>8, wb>>8, wa>>8, r>>8, g>>8, b>>8, a>>8)
					}
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}
}
