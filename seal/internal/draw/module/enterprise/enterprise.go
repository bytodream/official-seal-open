// enterprise.go 企业圆章（圆形章默认类型：enterprise 及其它未单独实现的圆形章）。
package enterprise

import (
	"image"
	"image/color"
	stddraw "image/draw"

	"github.com/bytodream/official-seal-open/seal/internal/draw/base"
	"github.com/bytodream/official-seal-open/seal/model"
	"github.com/bytodream/official-seal-open/seal/utils"
	"github.com/bytodream/official-seal-open/seal/utils/ioc"
)

var (
	sealCategoryEnterprise = model.SealCategoryEnterprise.ToString()
)

func Register() {
	ioc.Register(NewEnterpriseDrawServiceImpl(), ioc.WithComponentName(sealCategoryEnterprise))
}

// EnterpriseDrawServiceImpl 企业圆章绘制服务接口。
type EnterpriseDrawServiceImpl struct{}

func NewEnterpriseDrawServiceImpl() base.DrawService {
	return &EnterpriseDrawServiceImpl{}
}

// Draw 绘制企业圆章。
func (s *EnterpriseDrawServiceImpl) Draw(m base.SealPaintModel) (image.Image, error) {
	// 检查权限（按调用方章型校验）
	if err := utils.CheckSealCategory(string(m.Category)); err != nil {
		return nil, err
	}

	img := image.NewRGBA(image.Rect(0, 0, m.Width, m.Height))
	stddraw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, stddraw.Src)

	cx, cy := m.Width/2, m.Height/2
	diameter := float64(m.Width) - float64(m.EdgeWidth)
	radius := diameter / 2

	stroke := m.SealColor
	if stroke == nil {
		stroke = model.SealRed
	}
	DrawAnnulus(img, cx, cy, radius, float64(m.EdgeWidth), stroke)

	if m.InnerCircle != nil {
		if err := DrawNestedInnerCircle(img, m, stroke); err != nil {
			return nil, err
		}
	}

	if m.Content != "" {
		if err := DrawArcContent(img, m, radius, stroke); err != nil {
			return nil, err
		}
	}

	if m.Foot != nil && FootValid(m.Foot) {
		if err := DrawFootCenter(img, m, stroke); err != nil {
			return nil, err
		}
	}

	if m.EnterpriseCode != nil && !utils.IsBlank(m.EnterpriseCode.Code) {
		if err := DrawEnterpriseCode(img, m, radius, stroke); err != nil {
			return nil, err
		}
	}

	// 防伪线在所有内容之后、Translucent 之前绘制（白线会切断文字与边框）。
	if m.FakeGap {
		DrawAntifakeGapLines(img, m.Width, m.EdgeWidth)
	}

	if m.SealLogo && m.AntifakeInfo != nil {
		img = Translucent(img, m.Transparency)
		if err := DrawAntifakeLogo(img, m); err != nil {
			return nil, err
		}
		if err := ApplySampleWatermark(img, m); err != nil {
			return nil, err
		}
		// 章面已收尾，成功返回前计入使用次数（前面任何报错都不计数）
		if err := base.RecordUsageUnlessSample(m); err != nil {
			return nil, err
		}
		return img, nil
	}

	if m.SealLogo && m.LogoSize > 0 {
		logoCol := stroke
		if logoCol == nil {
			logoCol = model.SealRed
		}
		if m.LogoFiveStar {
			DrawFiveStar(img, cx, cy-int(m.LogoUpmove), m.LogoSize, logoCol)
		} else if err := base.DrawBitmapLogo(img, cx, cy, m.LogoSize, m.LogoUpmove, logoCol); err != nil {
			return nil, err
		}
	}

	// 章面收尾（样章标识 + 透明化）成功后再计入使用次数（前面任何报错都不计数）
	out, err := Finish(img, m)
	if err != nil {
		return nil, err
	}
	if err := base.RecordUsageUnlessSample(m); err != nil {
		return nil, err
	}
	return out, nil
}

var _ base.DrawService = (*EnterpriseDrawServiceImpl)(nil)
