package service

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/bytodream/official-seal-open/seal/internal/render"
	"github.com/bytodream/official-seal-open/seal/service/dto"
	xdraw "golang.org/x/image/draw"
)

// DefaultImageService 提供章面位图的转码实现，处理顺序对齐 Java
// ImageServiceImpl.toImageBytes：透明化 → 缩放 → 旋转 → PNG 编码。
//
// 它不涉及对象存储，未注入上传函数时 HasUploadBytes 返回 false、
// UploadBytes 返回 dto.ErrImageUploadNotConfigured。接入自有存储的推荐方式：
//
//	svc := service.NewCompanySealCreateService()
//	svc.SetUploadBytes(func(data []byte, attr dto.ImageAttr, companySeal bool) (string, error) {
//		return myStorage.Upload(data, "image/png"), nil
//	})
//
// 也可内嵌本类型并覆盖 UploadBytes，但必须同时覆盖 HasUploadBytes 返回 true，
// 否则 GetAutoSealImage 会在渲染前拒绝上传。
//
// 本类型的方法都是指针接收者，赋值时须取地址：svc.SetImageService(&myImageService{})。
type DefaultImageService struct {
	uploadBytes    UploadBytes
	isBlurGenerate bool
}

var (
	_ ImageService             = &DefaultImageService{}
	_ ImageServiceConfigurator = &DefaultImageService{}
)

func newDefaultImageService() *DefaultImageService {
	return &DefaultImageService{}
}

func (s *DefaultImageService) SetUploadBytes(uploadBytes UploadBytes) {
	s.uploadBytes = uploadBytes
}

func (s *DefaultImageService) SetBlurGenerate(isBlurGenerate bool) {
	s.isBlurGenerate = isBlurGenerate
}

func (s *DefaultImageService) IsBlurGenerate() bool {
	return s.isBlurGenerate
}

// ToImageBytes 按 ImageAttr 对章面位图做后处理并编码为 PNG。
//
// 透明度字符串缺省为 "30"、换算按 (100-t)*2.55 截断取整，非法输入直接报错而非静默回退。
func (s *DefaultImageService) ToImageBytes(img image.Image, attr dto.ImageAttr) ([]byte, error) {
	if img == nil {
		return nil, fmt.Errorf("章面位图为空")
	}
	out := img

	if attr.Transparent {
		alpha, err := transparencyAlpha(attr.Transparency)
		if err != nil {
			return nil, err
		}
		out = translucentSeal(out, alpha)
	}

	if attr.Zoom {
		limitHeight := attr.LimitHeight
		if attr.DefaultAttr {
			// Java：高度按 3:1 宽度比例推算，用于「默认属性」下的等比输出。
			if b := out.Bounds(); b.Dx() > 0 {
				limitHeight = b.Dy() * 300 / b.Dx()
			}
		}
		out = zoomImage(out, attr.LimitWidth, limitHeight)
	}

	if attr.Rotate {
		rotated, err := rotateImage(out, attr.Degree)
		if err != nil {
			return nil, err
		}
		out = rotated
	}

	return render.PNGBytes(out)
}

func (s *DefaultImageService) HasUploadBytes() bool {
	return s.uploadBytes != nil
}

// UploadBytes 调用注入的上传函数；未注入时返回 dto.ErrImageUploadNotConfigured。
func (s *DefaultImageService) UploadBytes(data []byte, attr dto.ImageAttr, companySeal bool) (string, error) {
	if s.uploadBytes == nil {
		return "", dto.ErrImageUploadNotConfigured
	}
	return s.uploadBytes(data, attr, companySeal)
}

// transparencyAlpha 把 0–100 的透明度字符串换算为 Alpha。
// 缺省值 "30" 与 Java 一致；换算公式同样按 (100-t)*2.55 截断 —— t=0 得 254 而非 255，
// 与渲染期 base.TransparencyFromSealString 的四舍五入语义不同，不要互换。
func transparencyAlpha(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		s = "30"
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("透明度必须为 0–100 的整数，当前为 %q", s)
	}
	if v < 0 {
		v = 0
	}
	if v > 100 {
		v = 100
	}
	return int(float64(100-v) * 2.55), nil
}

// translucentSeal 复刻 Java ImageUtils.translucentSeal(dest, transparency)：
// 跳过完全透明像素与纯白像素（三通道均 >225），其余像素保留 RGB、把 Alpha 写为给定值，
// 最后按 SrcIn 语义与原图 Alpha 相乘。
func translucentSeal(src image.Image, alpha int) *image.NRGBA {
	if alpha < 0 {
		alpha = 0
	}
	if alpha > 255 {
		alpha = 255
	}
	a := uint32(alpha)

	b := src.Bounds()
	dst := image.NewNRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := color.NRGBAModel.Convert(src.At(x, y)).(color.NRGBA)
			if c.A == 0 {
				continue
			}
			if c.R > 225 && c.G > 225 && c.B > 225 {
				continue
			}
			// AlphaComposite.SrcIn：结果 Alpha = 源 Alpha × 掩码 Alpha。
			c.A = uint8(uint32(c.A) * a / 255)
			dst.SetNRGBA(x, y, c)
		}
	}
	return dst
}

// zoomImage 复刻 Java ImageUtils.zoom：保持宽高比缩放到不超过 limitWidth×limitHeight。
// 单边为 0 时 Java 把它原样交给 Thumbnailator ResampleOp，语义是「该边按另一边的比例推算」，
// 而不是「该边不限制」——后者会拒绝在单边限制下放大。缩放算法本身与原实现
// （Thumbnailator ResampleOp）不同，尺寸一致但像素不保证逐字节相同。
func zoomImage(src image.Image, limitWidth, limitHeight int) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= 0 || h <= 0 {
		return src
	}
	if limitWidth <= 0 && limitHeight > 0 {
		limitWidth = w * limitHeight / h
	}
	if limitHeight <= 0 && limitWidth > 0 {
		limitHeight = h * limitWidth / w
	}
	if limitWidth <= 0 {
		limitWidth = w
	}
	if limitHeight <= 0 {
		limitHeight = h
	}

	scale := math.Min(float64(limitWidth)/float64(w), float64(limitHeight)/float64(h))
	nw := int(math.Round(float64(w) * scale))
	nh := int(math.Round(float64(h) * scale))
	if nw < 1 {
		nw = 1
	}
	if nh < 1 {
		nh = 1
	}
	if nw == w && nh == h {
		return src
	}

	dst := image.NewNRGBA(image.Rect(0, 0, nw, nh))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), src, b, xdraw.Src, nil)
	return dst
}

// rotateImage 复刻 Java ImageUtils.rotate：只接受 90 的整数倍，逐像素搬运、不插值。
func rotateImage(src image.Image, degree int) (image.Image, error) {
	if degree%90 != 0 {
		return nil, fmt.Errorf("旋转角度必须是 90 的倍数，当前为 %d", degree)
	}
	rate := degree / 90 % 4
	if rate < 0 {
		rate += 4
	}
	if rate == 0 {
		return src, nil
	}

	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	var dst *image.NRGBA
	switch rate {
	case 1: // 顺时针 90°
		dst = image.NewNRGBA(image.Rect(0, 0, h, w))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				dst.SetNRGBA(h-1-y, x, srcColor(src, b, x, y))
			}
		}
	case 2: // 180°
		dst = image.NewNRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				dst.SetNRGBA(w-1-x, h-1-y, srcColor(src, b, x, y))
			}
		}
	default: // rate == 3，逆时针 90°
		dst = image.NewNRGBA(image.Rect(0, 0, h, w))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				dst.SetNRGBA(y, w-1-x, srcColor(src, b, x, y))
			}
		}
	}
	return dst, nil
}

// srcColor 按图像左上角为原点取色，避免非零原点导致坐标错位。
func srcColor(src image.Image, b image.Rectangle, x, y int) color.NRGBA {
	return color.NRGBAModel.Convert(src.At(b.Min.X+x, b.Min.Y+y)).(color.NRGBA)
}
