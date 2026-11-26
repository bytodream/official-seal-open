package base

import (
	"github.com/bytodream/official-seal-open/seal/model"

	_ "embed"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// 字体文件名常量，字体目录由服务层配置。

const (
	// FontPingFangRegular 内嵌默认字体文件名（PingFangRegular.otf，随库打包）
	FontPingFangRegular = "PingFangRegular.otf"
)

// defaultFontData 内嵌默认字体（seal/internal/draw/base/font/PingFangRegular.otf，随库打包）。
//
//go:embed font/PingFangRegular.otf
var defaultFontData []byte

var (
	defaultFontOnce sync.Once
	defaultFontFile string
)

// materializeDefaultFont 将内嵌默认字体落盘（进程内仅一次），返回可读取路径；失败返回空串。
func materializeDefaultFont() string {
	defaultFontOnce.Do(func() {
		dir, err := os.MkdirTemp("", "official-seal-font")
		if err != nil {
			return
		}
		p := filepath.Join(dir, FontPingFangRegular)
		if err := os.WriteFile(p, defaultFontData, 0o600); err != nil {
			return
		}
		defaultFontFile = p
	})
	return defaultFontFile
}

// ResolveFont 返回字体文件完整路径。
// 传参优先；为空时读 OFFICIAL_SEAL_FONTS_PATH 环境变量；
// 仍未配置时回退到内嵌默认字体 seal/internal/draw/base/font/PingFangRegular.otf（落盘到临时目录后返回路径）。
func ResolveFont(fontPath string) string {
	if fontPath == "" {
		fontPath = os.Getenv("OFFICIAL_SEAL_FONTS_PATH")
		if fontPath == "" {
			return materializeDefaultFont()
		}
	}
	return fontPath
}

// LoadFontFace 按像素字号与 DPI 加载字体（默认 HintingFull）。
func LoadFontFace(fontFile string, sizePx float64, dpi int) (font.Face, error) {
	return LoadFontFaceHinting(fontFile, sizePx, dpi, font.HintingFull)
}

// LoadFontFaceHinting 按像素字号与 DPI 加载字体。
// 网格拟合、按未 hint 的轮廓光栅化，因此文字环必须用 HintingNone 才能对齐字形轮廓。
func LoadFontFaceHinting(fontFile string, sizePx float64, dpi int, hinting font.Hinting) (font.Face, error) {
	if dpi <= 0 {
		dpi = model.DefaultCustomSealAttributeDPI
	}
	data, err := os.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}
	tt, err := opentype.Parse(data)
	if err != nil {
		return nil, err
	}
	pt := sizePx * 72.0 / float64(dpi)
	return opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    pt,
		DPI:     float64(dpi),
		Hinting: hinting,
	})
}

// StringBoundsAscent 返回字体级逻辑上边界（1x 像素）＝ hhea.ascender + lineGap/2，
// 它对应字符串包围盒的上边界，而非字形墨迹上边界。
// Go 的 Metrics().Ascent 只是 ascender；lineGap = Height - Ascent - Descent。
// 实测 simsun.ttf：0.8594em + 0.1406/2 ≈ 0.9297em（142px 字号差 ~10px）。
func StringBoundsAscent(face font.Face, ss int) float64 {
	m := face.Metrics()
	asc := float64(m.Ascent) / 64
	desc := float64(m.Descent) / 64
	gap := float64(m.Height)/64 - asc - desc
	return (asc + gap/2) / float64(ss)
}
