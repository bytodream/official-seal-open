package base

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"sync"

	xdraw "golang.org/x/image/draw"

	"github.com/bytodream/official-seal-open/seal/model"
	"github.com/bytodream/official-seal-open/seal/utils"
)

//go:embed assets/logo.png
var logoPNG []byte

var (
	logoOnce sync.Once
	logoImg  image.Image
	logoErr  error
)

// LoadLogoBitmap 加载 base 层内嵌的共用 Logo 位图资源，各章型统一调用。
func LoadLogoBitmap() (image.Image, error) {
	logoOnce.Do(func() {
		logoImg, logoErr = png.Decode(bytes.NewReader(logoPNG))
	})
	return logoImg, logoErr
}

// RecordUsageUnlessSample 非样章绘制（m.ShowWatermark 为 false）计入一次使用次数；
// 预览样章（m.ShowWatermark 为 true）不计数。
// 各章型在 Draw 成功收尾后调用，额度错误需向上返回。
func RecordUsageUnlessSample(m SealPaintModel) error {
	if m.ShowWatermark {
		return nil
	}
	return utils.RecordUsage()
}

// DrawService 章实现
type DrawService interface {
	// Draw 实现章绘制；实现内部按 m.Category 校验章型权限并在成功收尾后计次
	Draw(m SealPaintModel) (image.Image, error)
	// GetCircle 获取圆章属性
	GetCircle(a *model.CustomSealAttribute) (*SealAttributeBundle, error)
}

func GetDrawServiceName(name model.SealCategory) string {
	return fmt.Sprintf("Draw_%s", name.ToString())
}

// TransferSealBitmapRGB 将位图非近白像素染成印章色（白底透明）。
func TransferSealBitmapRGB(src image.Image, tint color.Color) *image.RGBA {
	b := src.Bounds()
	out := image.NewRGBA(b)
	r, g, bb, _ := tint.RGBA()
	tr, tg, tb := uint8(r>>8), uint8(g>>8), uint8(bb>>8)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			pr, pg, pb, pa := src.At(x, y).RGBA()
			or, og, ob := uint8(pr>>8), uint8(pg>>8), uint8(pb>>8)
			if or <= 225 || og <= 225 || ob <= 225 {
				out.SetRGBA(x, y, color.RGBA{R: tr, G: tg, B: tb, A: uint8(pa >> 8)})
			}
		}
	}
	return out
}

// DrawBitmapLogo 在章面中心绘制位图 Logo，位图资源由 LoadLogoBitmap 统一加载。
// upMove 为向上偏移像素，size 为 Logo 边长像素。
func DrawBitmapLogo(dst *image.RGBA, cx, cy, size int, upMove float32, tint color.Color) error {
	src, err := LoadLogoBitmap()
	if err != nil {
		return err
	}
	colored := TransferSealBitmapRGB(src, tint)
	scaled := image.NewRGBA(image.Rect(0, 0, size, size))
	xdraw.ApproxBiLinear.Scale(scaled, scaled.Bounds(), colored, colored.Bounds(), xdraw.Over, nil)
	topLeft := image.Pt(cx-size/2, cy-size/2-int(upMove))
	xdraw.Draw(dst, scaled.Bounds().Add(topLeft), scaled, image.Point{}, xdraw.Over)
	return nil
}
