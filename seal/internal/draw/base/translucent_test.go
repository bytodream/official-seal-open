package base

import (
	"image"
	"image/color"
	"testing"

	"github.com/bytodream/official-seal-open/seal/model"
)

// TestTransparencyFromSealString 覆盖 0–100 边界、越界钳制与非法输入回退。
func TestTransparencyFromSealString(t *testing.T) {
	str := func(s string) *string { return &s }
	cases := []struct {
		name string
		in   *string
		want int
	}{
		{"nil 使用默认", nil, model.DefaultSealTransparency},
		{"0 完全不透明", str("0"), 255},
		{"100 全透明", str("100"), 0},
		{"50 约为一半", str("50"), 128},
		{"负数钳制到 0", str("-20"), 255},
		{"超过 100 钳制到 100", str("150"), 0},
		{"非法字符串回退默认", str("abc"), model.DefaultSealTransparency},
		{"空白字符串回退默认", str("  "), model.DefaultSealTransparency},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := TransparencyFromSealString(tc.in); got != tc.want {
				t.Fatalf("TransparencyFromSealString(%v) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

// TestTranslucentClampsAlpha 越界透明度写入时钳制，避免 uint8 回绕。
func TestTranslucentClampsAlpha(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 1, 1))
	src.Set(0, 0, color.RGBA{R: 255, A: 255})

	for _, alpha := range []int{-127, 0, 255, 300, 1000} {
		out := Translucent(src, alpha)
		got := out.RGBAAt(0, 0).A
		want := uint8(alpha)
		if alpha < 0 {
			want = 0
		}
		if alpha > 255 {
			want = 255
		}
		if got != want {
			t.Fatalf("Translucent(alpha=%d) A = %d, want %d", alpha, got, want)
		}
	}
}
