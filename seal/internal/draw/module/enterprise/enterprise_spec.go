package enterprise

import (
	"github.com/bytodream/official-seal-open/seal/internal/draw/base"
	"github.com/bytodream/official-seal-open/seal/model"
	"github.com/bytodream/official-seal-open/seal/utils"
)

// GetCircle 按自定义属性返回企业章（圆章/椭圆/双线椭圆）章面属性。
func (s *EnterpriseDrawServiceImpl) GetCircle(a *model.CustomSealAttribute) (*base.SealAttributeBundle, error) {
	circle, err := GetEnterprise(a.SealSpecEnum(), a.Foot, a.Head, a.EnterpriseCode, a.GetDpi(), a.InnerContent, a.ToDiySealAttribute(), a.IsSealLogo())
	if err != nil {
		return nil, err
	}
	return &base.SealAttributeBundle{Paint: circle}, nil
}

// enterpriseCirclePreset 企业圆章单规格默认毫米参数。
// 注意脚注距离、文字弧、logo 上移量都是「有无企业代码」两套值。
type enterpriseCirclePreset struct {
	// DiameterMm 章面直径
	DiameterMm float32
	// EdgeMm 外圈线宽
	EdgeMm float32
	// FontMm 主文字字号
	FontMm float32
	// MarginMm 主文字与外圈边距
	MarginMm float32
	// Arc 主文字弧度
	Arc int
	// FootFontMm 脚排字号
	FootFontMm float32
	// FootDistMm 无企业代码时的脚注距离
	FootDistMm float32
	// FootDistCodeMm 有企业代码时的脚注距离
	FootDistCodeMm float32
	// FootWidthMm 脚排宽度
	FootWidthMm float32
	// StarMm 中心五角星尺寸
	StarMm float32
	// CodeFontMm 企业代码字号
	CodeFontMm float32
	// DefaultCodeArc 企业代码默认弧度
	DefaultCodeArc int
	// CodeArcWithFoot 有脚注时的企业代码弧度
	CodeArcWithFoot int
	// FootArc 有脚注时的主文字弧度
	FootArc int
	// UpWithCodeMm 有脚注且有企业代码时中心图案上移距离
	UpWithCodeMm float32
	// UpNoCodeMm 有脚注且无企业代码时中心图案上移距离
	UpNoCodeMm float32
}

var enterprisePresets = map[string]enterpriseCirclePreset{
	"CIRCULAR_19": {19, 0.6, 4, 0.3, 270, 4, 2.5, 2.0, 10, 6.5, 1.0, 85, 75, 240, 0.7, 0.4},
	"CIRCULAR_20": {20, 0.6, 4, 0.3, 270, 4, 2.6, 2.5, 11, 7.5, 1.0, 85, 75, 235, 0.7, 0.4},
	"CIRCULAR_23": {23, 0.6, 5, 0.4, 270, 4, 2.6, 2.5, 11, 7.5, 1.0, 85, 75, 235, 0.7, 0.4},
	"CIRCULAR_27": {27, 0.7, 6, 0.4, 270, 5, 4.0, 2.5, 16, 8.0, 1.0, 85, 75, 235, 0.7, 0.4},
	"CIRCULAR_30": {30, 1.0, 6, 0.4, 270, 5, 4.0, 2.5, 16, 7.0, 2.0, 85, 75, 230, 0.7, 0.4},
	"CIRCULAR_34": {34, 0.9, 5.5, 0.7, 270, 6, 5.0, 5.0, 20, 10.0, 2.0, 85, 75, 230, 1.5, 1.0},
	"CIRCULAR_35": {35, 1.0, 5.5, 0.7, 270, 6, 5.0, 4.0, 20, 12.5, 2.5, 85, 75, 230, 0.7, 0.5},
	"CIRCULAR_36": {36, 1.0, 6, 0.7, 270, 5, 6.0, 5.0, 20, 12.5, 2.5, 85, 75, 230, 0.6, 0.5},
	"CIRCULAR_37": {37, 1.0, 6, 0.7, 270, 5, 6.0, 5.0, 20, 12.5, 2.5, 85, 75, 230, 0.6, 0.5},
	"CIRCULAR_38": {38, 1.0, 6, 0.8, 270, 7, 6.0, 5.0, 24, 12.0, 2.2, 85, 75, 210, 1.5, 1.0},
	"CIRCULAR_39": {39, 1.0, 6, 0.8, 270, 7, 6.0, 4.0, 24, 12.0, 2.2, 85, 75, 230, 1.0, 0.5},
	"CIRCULAR_40": {40, 1.1, 6.5, 0.8, 270, 7, 6.0, 5.0, 24, 13.0, 2.2, 85, 75, 210, 1.5, 1.0},
	"CIRCULAR_41": {41, 1.1, 7, 0.8, 270, 7, 7.0, 5.0, 24, 13.0, 2.2, 85, 75, 230, 1.5, 1.0},
	"CIRCULAR_42": {42, 1.2, 7, 0.8, 270, 6.8, 7.8, 5.8, 22, 14.0, 2.2, 85, 75, 235, 1.5, 1.0},
	"CIRCULAR_45": {45, 1.2, 7.2, 0.8, 270, 6.8, 7.8, 5.8, 22, 14.0, 2.2, 85, 75, 235, 1.5, 1.0},
}

// GetEnterprise 按规格与 DIY 参数返回企业章（圆章/椭圆/双线椭圆）章面属性。
func GetEnterprise(spec, footText, headText, enterpriseCode string, dpi int, innerContent string, diy *model.DiySealAttribute, sealLog bool) (*base.SealPaintAttribute, error) {
	if IsDoublelineOvalSpec(spec) {
		return GetDoublelineOval(spec, footText, headText, enterpriseCode, innerContent, dpi)
	}
	if spec == model.SealSpecOval4530 {
		return BuildOval4530Enterprise(footText, enterpriseCode, innerContent, dpi, sealLog)
	}
	p, ok := enterprisePresets[spec]
	if !ok && spec == model.SealSpecDiy && diy != nil {
		return buildDiyEnterpriseCircle(footText, enterpriseCode, dpi, diy, sealLog)
	}
	if !ok {
		return nil, base.ErrSealSpecNotSupported
	}
	return buildEnterpriseCircle(p, spec, footText, enterpriseCode, dpi, sealLog, diy)
}

// buildEnterpriseCircle 组装章面属性结构。
func buildEnterpriseCircle(p enterpriseCirclePreset, _ string, footText, enterpriseCode string, dpi int, sealLog bool, diy *model.DiySealAttribute) (*base.SealPaintAttribute, error) {
	edge := p.EdgeMm
	font := p.FontMm
	margin := p.MarginMm
	arc := p.Arc
	diameter := p.DiameterMm
	footDist := p.FootDistMm
	footWidth := p.FootWidthMm
	footFont := p.FootFontMm
	star := p.StarMm
	logoUp := float32(0)
	codeArc := p.DefaultCodeArc

	var code *base.EnterpriseCodeAttr
	if !utils.IsBlank(enterpriseCode) {
		codeFont := p.CodeFontMm
		if diy != nil && diy.EnterpriseCodeFontSize != nil {
			codeFont = *diy.EnterpriseCodeFontSize
		}
		code = &base.EnterpriseCodeAttr{
			Code:       enterpriseCode,
			FontSizePx: utils.ToIntPixel(dpi, codeFont),
			Arc:        codeArc,
		}
	}

	if !utils.IsBlank(footText) {
		arc = p.FootArc
		if code != nil {
			footDist = p.FootDistCodeMm
			logoUp = p.UpWithCodeMm
			code.Arc = p.CodeArcWithFoot
		} else {
			logoUp = p.UpNoCodeMm
		}
	}

	attr := &base.SealPaintAttribute{
		Dpi:            dpi,
		Transparency:   model.DefaultSealTransparency,
		Category:       model.SealCategoryEnterprise,
		Arc:            arc,
		WidthPx:        utils.ToIntPixel(dpi, diameter),
		HeightPx:       utils.ToIntPixel(dpi, diameter),
		MarginPx:       utils.ToFloatPixel(dpi, margin),
		EdgeWidthPx:    utils.ToFloatPixel(dpi, edge),
		FontSizePx:     utils.ToIntPixel(dpi, font),
		SealLogo:       sealLog,
		FakeGap:        true,
		SealColor:      model.SealRed,
		EnterpriseCode: code,
	}

	if !utils.IsBlank(footText) {
		attr.Foot = &base.FootAttr{
			Text:       footText,
			FontSizePx: utils.ToIntPixel(dpi, footFont),
			DistancePx: utils.ToFloatPixel(dpi, footDist),
			WidthPx:    utils.ToFloatPixel(dpi, footWidth),
		}
	}
	if sealLog {
		attr.Logo = &base.SealLogoAttr{
			SizePx:   utils.ToIntPixel(dpi, star),
			UpMovePx: utils.ToFloatPixel(dpi, logoUp),
			FiveStar: true,
		}
	}

	applyDiyToCircle(diy, dpi, attr)
	return attr, nil
}

// buildDiyEnterpriseCircle 组装章面属性结构。
func buildDiyEnterpriseCircle(footText, enterpriseCode string, dpi int, diy *model.DiySealAttribute, sealLog bool) (*base.SealPaintAttribute, error) {
	if diy.SealWidth == nil {
		return nil, base.ErrSealSpecNotSupported
	}
	width := *diy.SealWidth
	height := width
	if diy.SealHeight != nil {
		height = *diy.SealHeight
	}
	edge := utils.F32Or(diy.EdgeWidth, 1.2)
	font := utils.F32Or(diy.FontSize, 7)
	margin := utils.F32Or(diy.WordLinePadding, 0.8)
	arc := 270
	if diy.Arc != nil {
		arc = *diy.Arc
	}
	cat := model.SealCategoryEnterprise
	if height != width {
		cat = model.SealCategoryOval
	}
	attr := &base.SealPaintAttribute{
		Dpi:         dpi,
		Category:    cat,
		Arc:         arc,
		WidthPx:     utils.ToIntPixel(dpi, width),
		HeightPx:    utils.ToIntPixel(dpi, height),
		MarginPx:    utils.ToFloatPixel(dpi, margin),
		EdgeWidthPx: utils.ToFloatPixel(dpi, edge),
		FontSizePx:  utils.ToIntPixel(dpi, font),
		SealLogo:    sealLog,
		FakeGap:     true,
		SealColor:   model.SealRed,
	}
	hasCode := !utils.IsBlank(enterpriseCode)
	if !utils.IsBlank(footText) {
		// 脚注字号/距离缺省时参照 CIRCULAR_42 预设（有代码 5.8mm、无代码 7.8mm）；
		// 宽度缺省为 0，表示绘制层按自然字宽居中。
		footDist := float32(7.8)
		if hasCode {
			footDist = 5.8
		}
		attr.Foot = &base.FootAttr{
			Text:       footText,
			FontSizePx: utils.ToIntPixel(dpi, utils.F32Or(diy.FootFontSize, 6.8)),
			DistancePx: utils.ToFloatPixel(dpi, utils.F32Or(diy.FootDistance, footDist)),
			WidthPx:    utils.ToFloatPixel(dpi, utils.F32Or(diy.FootWidth, 0)),
		}
		attr.Arc = 235
	}
	if hasCode {
		codeArc := 90
		if diy.EnterpriseCodeArc != nil {
			codeArc = *diy.EnterpriseCodeArc
		}
		if !utils.IsBlank(footText) && diy.EnterpriseCodeArc == nil {
			codeArc = 75
		}
		attr.EnterpriseCode = &base.EnterpriseCodeAttr{
			Code:       enterpriseCode,
			FontSizePx: utils.ToIntPixel(dpi, utils.F32Or(diy.EnterpriseCodeFontSize, 2.2)),
			Arc:        codeArc,
		}
	}
	if sealLog {
		attr.Logo = &base.SealLogoAttr{
			SizePx:   utils.ToIntPixel(dpi, utils.F32Or(diy.LogoSize, 14)),
			UpMovePx: utils.ToFloatPixel(dpi, utils.F32Or(diy.LogoUpmove, 0)),
			FiveStar: true,
		}
	}
	return attr, nil
}

// applyDiyToCircle 将 DIY 参数覆盖到圆章属性。
func applyDiyToCircle(diy *model.DiySealAttribute, dpi int, attr *base.SealPaintAttribute) {
	if diy == nil || attr == nil {
		return
	}
	if diy.EdgeWidth != nil {
		attr.EdgeWidthPx = utils.ToFloatPixel(dpi, *diy.EdgeWidth)
	}
	if diy.FontSize != nil {
		attr.FontSizePx = utils.ToIntPixel(dpi, *diy.FontSize)
	}
	if diy.WordLinePadding != nil {
		attr.MarginPx = utils.ToFloatPixel(dpi, *diy.WordLinePadding)
	}
	if diy.Arc != nil {
		attr.Arc = *diy.Arc
	}
	if diy.SealWidth != nil {
		attr.WidthPx = utils.ToIntPixel(dpi, *diy.SealWidth)
	}
	if diy.SealHeight != nil {
		attr.HeightPx = utils.ToIntPixel(dpi, *diy.SealHeight)
	}
}
