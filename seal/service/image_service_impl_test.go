package service

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/bytodream/official-seal-open/seal/service/dto"
)

// newTestRaster 构造 w×h 白底位图，并在给定坐标画一个红色像素。
func newTestRaster(w, h, markX, markY int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetNRGBA(x, y, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}
	if markX >= 0 && markY >= 0 {
		img.SetNRGBA(markX, markY, color.NRGBA{R: 255, A: 255})
	}
	return img
}

// decodePNG 解码 ToImageBytes 输出并返回图像。
func decodePNG(t *testing.T, data []byte) image.Image {
	t.Helper()
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("输出不是合法 PNG: %v", err)
	}
	return img
}

// nrgbAt 按非预乘语义取点，便于直接断言 RGB 与 Alpha。
func nrgbAt(img image.Image, x, y int) color.NRGBA {
	return color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
}

// TestTransparencyAlpha 覆盖 Java (100-t)*2.55 截断语义与缺省值。
func TestTransparencyAlpha(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want int
	}{
		{"缺省按 30", "", 178},
		{"0 得 254 而非 255", "0", 254},
		{"30", "30", 178},
		{"100 全透明", "100", 0},
		{"负数钳制到 0", "-10", 254},
		{"超过 100 钳制", "150", 0},
		{"带空白", " 50 ", 127},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := transparencyAlpha(tc.in)
			if err != nil {
				t.Fatalf("transparencyAlpha(%q) 报错: %v", tc.in, err)
			}
			if got != tc.want {
				t.Fatalf("transparencyAlpha(%q) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}

	if _, err := transparencyAlpha("abc"); err == nil {
		t.Fatal("非法透明度应报错，实际未报错")
	}
}

// TestToImageBytesTransparent 断言白色背景被丢弃、印泥色写入指定 Alpha。
func TestToImageBytesTransparent(t *testing.T) {
	src := newTestRaster(4, 4, 2, 2)
	got, err := newDefaultImageService().ToImageBytes(src, dto.ImageAttr{
		Transparent:  true,
		Transparency: "30",
	})
	if err != nil {
		t.Fatalf("ToImageBytes 报错: %v", err)
	}

	img := decodePNG(t, got)
	if c := nrgbAt(img, 2, 2); c.R != 255 || c.A != 178 {
		t.Fatalf("印泥像素 = %+v, want R=255 A=178", c)
	}
	if c := nrgbAt(img, 0, 0); c.A != 0 {
		t.Fatalf("白色背景像素 A = %d, want 0", c.A)
	}
}

// TestToImageBytesTransparentDefault 透明度缺省时应按 30 处理。
func TestToImageBytesTransparentDefault(t *testing.T) {
	got, err := newDefaultImageService().ToImageBytes(newTestRaster(2, 2, 0, 0), dto.ImageAttr{Transparent: true})
	if err != nil {
		t.Fatalf("ToImageBytes 报错: %v", err)
	}
	if c := nrgbAt(decodePNG(t, got), 0, 0); c.A != 178 {
		t.Fatalf("缺省透明度 A = %d, want 178", c.A)
	}
}

// TestToImageBytesTransparentKeepsSourceAlpha 半透明边缘按 SrcIn 与掩码 Alpha 相乘，
// 而不是像渲染期那样把 Alpha 一律覆盖为给定值。
func TestToImageBytesTransparentKeepsSourceAlpha(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	img.SetNRGBA(0, 0, color.NRGBA{R: 255, A: 255}) // 不透明印泥
	img.SetNRGBA(1, 0, color.NRGBA{R: 255, A: 128}) // 抗锯齿边缘

	got, err := newDefaultImageService().ToImageBytes(img, dto.ImageAttr{
		Transparent:  true,
		Transparency: "30",
	})
	if err != nil {
		t.Fatalf("ToImageBytes 报错: %v", err)
	}
	out := decodePNG(t, got)
	if c := nrgbAt(out, 0, 0); c.A != 178 {
		t.Fatalf("不透明像素 A = %d, want 178", c.A)
	}
	// 128 × 178 / 255 = 89
	if c := nrgbAt(out, 1, 0); c.A != 89 {
		t.Fatalf("半透明像素 A = %d, want 89（SrcIn 乘算）", c.A)
	}
}

// TestToImageBytesRotate 覆盖 90/180/270 与负角度的像素落点。
func TestToImageBytesRotate(t *testing.T) {
	// 3×2，标记在左上角 (0,0)。
	const w, h = 3, 2
	cases := []struct {
		degree       int
		wantW, wantH int
		wantX, wantY int
	}{
		{90, h, w, h - 1, 0},
		{180, w, h, w - 1, h - 1},
		{270, h, w, 0, w - 1},
		{-90, h, w, 0, w - 1},
		{360, w, h, 0, 0},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%ddeg", tc.degree), func(t *testing.T) {
			got, err := newDefaultImageService().ToImageBytes(newTestRaster(w, h, 0, 0), dto.ImageAttr{
				Rotate: true,
				Degree: tc.degree,
			})
			if err != nil {
				t.Fatalf("ToImageBytes(%d°) 报错: %v", tc.degree, err)
			}
			img := decodePNG(t, got)
			b := img.Bounds()
			if b.Dx() != tc.wantW || b.Dy() != tc.wantH {
				t.Fatalf("%d° 尺寸 = %dx%d, want %dx%d", tc.degree, b.Dx(), b.Dy(), tc.wantW, tc.wantH)
			}
			if c := nrgbAt(img, tc.wantX, tc.wantY); c.R != 255 || c.G != 0 {
				t.Fatalf("%d° 标记点 (%d,%d) = %+v, want 红色", tc.degree, tc.wantX, tc.wantY, c)
			}
		})
	}

	if _, err := (newDefaultImageService()).ToImageBytes(newTestRaster(w, h, 0, 0), dto.ImageAttr{
		Rotate: true,
		Degree: 45,
	}); err == nil {
		t.Fatal("非 90 倍数应报错，实际未报错")
	}
}

// TestToImageBytesZoomDefaultAttr 复刻 Java 默认属性的高度推算：h*300/w。
func TestToImageBytesZoomDefaultAttr(t *testing.T) {
	// 200×100，DefaultAttr 下 limitHeight = 100*300/200 = 150，
	// 与 LimitWidth=50 共同决定 scale = min(0.25, 1.5) = 0.25 → 50×25。
	got, err := newDefaultImageService().ToImageBytes(newTestRaster(200, 100, 5, 5), dto.ImageAttr{
		DefaultAttr: true,
		Zoom:        true,
		LimitWidth:  50,
	})
	if err != nil {
		t.Fatalf("ToImageBytes 报错: %v", err)
	}
	b := decodePNG(t, got).Bounds()
	if b.Dx() != 50 || b.Dy() != 25 {
		t.Fatalf("缩放尺寸 = %dx%d, want 50x25", b.Dx(), b.Dy())
	}
}

// TestToImageBytesZoomLimit 显式宽高时按较小的比例等比缩放。
func TestToImageBytesZoomLimit(t *testing.T) {
	// 100×100，limit 50×40 → scale = 0.4 → 40×40。
	got, err := newDefaultImageService().ToImageBytes(newTestRaster(100, 100, 5, 5), dto.ImageAttr{
		Zoom:        true,
		LimitWidth:  50,
		LimitHeight: 40,
	})
	if err != nil {
		t.Fatalf("ToImageBytes 报错: %v", err)
	}
	b := decodePNG(t, got).Bounds()
	if b.Dx() != 40 || b.Dy() != 40 {
		t.Fatalf("缩放尺寸 = %dx%d, want 40x40", b.Dx(), b.Dy())
	}
}

// TestToImageBytesNoOp 未开启任何后处理时应原样编码，尺寸与像素不变。
func TestToImageBytesNoOp(t *testing.T) {
	src := newTestRaster(8, 6, 3, 3)
	got, err := newDefaultImageService().ToImageBytes(src, dto.ImageAttr{})
	if err != nil {
		t.Fatalf("ToImageBytes 报错: %v", err)
	}
	img := decodePNG(t, got)
	if b := img.Bounds(); b.Dx() != 8 || b.Dy() != 6 {
		t.Fatalf("尺寸 = %dx%d, want 8x6", b.Dx(), b.Dy())
	}
	if c := nrgbAt(img, 3, 3); c.R != 255 || c.A != 255 {
		t.Fatalf("标记像素 = %+v, want R=255 A=255", c)
	}
}

// TestToImageBytesNilImage 空位图应返回明确错误而不是 panic。
func TestToImageBytesNilImage(t *testing.T) {
	if _, err := (newDefaultImageService()).ToImageBytes(nil, dto.ImageAttr{}); err == nil {
		t.Fatal("空位图应报错，实际未报错")
	}
}

// TestDefaultImageServiceUploadNotConfigured 上传能力需由集成方实现。
func TestDefaultImageServiceUploadNotConfigured(t *testing.T) {
	svc := newDefaultImageService()
	if _, err := svc.UploadBytes([]byte{1}, dto.ImageAttr{}, true); err != dto.ErrImageUploadNotConfigured {
		t.Fatalf("UploadBytes err = %v, want ErrImageUploadNotConfigured", err)
	}
}

// TestToImageBytesZoomSingleSideZero 单边限制为 0 时按另一边比例推算，
// 与 Java 把 0 交给 Thumbnailator 的语义一致：可放大。
func TestToImageBytesZoomSingleSideZero(t *testing.T) {
	cases := []struct {
		name           string
		limitW, limitH int
		wantW, wantH   int
	}{
		{"按高度放大", 0, 300, 600, 300},
		{"按高度缩小", 0, 100, 200, 100},
		{"按宽度缩小", 200, 0, 200, 100},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := newDefaultImageService().ToImageBytes(newTestRaster(400, 200, 0, 0), dto.ImageAttr{
				Zoom:        true,
				LimitWidth:  c.limitW,
				LimitHeight: c.limitH,
			})
			if err != nil {
				t.Fatalf("ToImageBytes 报错: %v", err)
			}
			img := decodePNG(t, got)
			if b := img.Bounds(); b.Dx() != c.wantW || b.Dy() != c.wantH {
				t.Fatalf("尺寸 = %dx%d, want %dx%d", b.Dx(), b.Dy(), c.wantW, c.wantH)
			}
		})
	}
}
