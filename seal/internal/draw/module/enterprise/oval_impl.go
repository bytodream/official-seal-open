// 本文件是跨章调用的本地副本：椭圆章属性构造复制自 module/oval 与 module/doublelineoval。
// 目的是让 enterprise 不依赖其它章模块；来源包变更时需同步本文件。
package enterprise

import (
	"github.com/bytodream/official-seal-open/seal/internal/draw/base"
	"github.com/bytodream/official-seal-open/seal/model"
	"github.com/bytodream/official-seal-open/seal/utils"
)

// BuildOval4530Enterprise 组装 45×30 企业椭圆章属性。
// 企业圆章构建器在规格为 OVAL_45_30 时会分发到本函数。
func BuildOval4530Enterprise(footText, enterpriseCode, innerContent string, dpi int, sealLogo bool) (*base.SealPaintAttribute, error) {
	const (
		widthMm   = float32(45)
		heightMm  = float32(30)
		edgeMm    = float32(0.9)
		fontMm    = float32(5.8)
		marginMm  = float32(0.8)
		footFont  = float32(5.5)
		footWidth = float32(23)
		arc       = 220
		starMm    = float32(9)
	)
	footDist := float32(4.5)
	logoUp := float32(0)
	// 脚注排宽默认 23mm。

	var code *base.EnterpriseCodeAttr
	if !utils.IsBlank(enterpriseCode) {
		code = &base.EnterpriseCodeAttr{
			Code:       enterpriseCode,
			FontSizePx: utils.ToIntPixel(dpi, 2.2),
			Arc:        100,
		}
		if !utils.IsBlank(footText) {
			logoUp = 1.5
		}
		footDist = 2.5
	} else if !utils.IsBlank(footText) {
		logoUp = 1
	}

	attr := &base.SealPaintAttribute{
		Dpi:            dpi,
		Transparency:   model.DefaultSealTransparency,
		Category:       model.SealCategoryOval,
		Arc:            arc,
		WidthPx:        utils.ToIntPixel(dpi, widthMm),
		HeightPx:       utils.ToIntPixel(dpi, heightMm),
		MarginPx:       utils.ToFloatPixel(dpi, marginMm),
		EdgeWidthPx:    utils.ToFloatPixel(dpi, edgeMm),
		FontSizePx:     utils.ToIntPixel(dpi, fontMm),
		SealLogo:       sealLogo,
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
	if sealLogo {
		attr.Logo = &base.SealLogoAttr{
			SizePx:   utils.ToIntPixel(dpi, starMm),
			UpMovePx: utils.ToFloatPixel(dpi, logoUp),
			FiveStar: true,
		}
	}

	if !utils.IsBlank(innerContent) {
		// 有内椭圆时外圈字号覆盖为 3.6，Logo 移除。
		attr.FontSizePx = utils.ToIntPixel(dpi, 3.6)
		attr.InnerOval = &base.SealPaintAttribute{
			Dpi:         dpi,
			Content:     innerContent,
			WidthPx:     utils.ToIntPixel(dpi, 33),
			HeightPx:    utils.ToIntPixel(dpi, 19),
			EdgeWidthPx: utils.ToFloatPixel(dpi, 0.5),
			FontSizePx:  utils.ToIntPixel(dpi, 3.6),
			MarginPx:    utils.ToFloatPixel(dpi, 0.5),
			Arc:         180,
		}
		attr.SealLogo = false
		attr.Logo = nil
		// 默认脚注基线 = distance 4.5 + fontSize 5.5 = 中心下 10mm，而内椭圆 33×19
		// 下缘仅 9.5mm，脚注必然压穿内椭圆线；且内外椭圆下环带高仅 4.5mm、两端随
		// 外椭圆收窄，任何字号都无法在带内同时避开两条线。因此有内椭圆时把脚注
		// 移入内椭圆下部空白区（此分支无五角星，内圈弧文 arc=180 全部位于上半部，
		// 下部无遮挡）。4.5/1.5/20 的取值经像素实测：墨迹下缘与内椭圆线两端最紧处仍留约 0.7mm。
		if attr.Foot != nil {
			attr.Foot.FontSizePx = utils.ToIntPixel(dpi, 4.5)
			attr.Foot.DistancePx = utils.ToFloatPixel(dpi, 1.5)
			attr.Foot.WidthPx = utils.ToFloatPixel(dpi, 20)
		}
	}
	return attr, nil
}

type doublelineOvalInner struct {
	width, height, edge, font, margin float32
	arc                               int
}

// defaultDoublelineInnerOval 返回默认内椭圆参数。
func defaultDoublelineInnerOval() doublelineOvalInner {
	return doublelineOvalInner{
		width: 33, height: 19, edge: 0.5, font: 3.6, margin: 0.5, arc: 180,
	}
}

// applyDoublelineHeadCodeFoot 设置双线椭圆的上排、编码与脚注。
func applyDoublelineHeadCodeFoot(
	attr *base.SealPaintAttribute,
	dpi int,
	footText, headText, enterpriseCode string,
	footFontDefault, footDistDefault, footDistWithHead, footWidth float32,
	headFont, headDist, headWidth float32,
	codeFont float32, codeArc int,
	footDistWithCode float32,
	logoUpFoot, logoUpCodeFoot *float32,
) {
	footFont := footFontDefault
	footDist := footDistDefault
	headDistMm := headDist

	if !utils.IsBlank(headText) {
		footFont = headFont
		// footDistWithHead>0 时覆盖默认脚注距离（如 40×30 仅有上排时）。
		if footDistWithHead > 0 {
			footDist = footDistWithHead
		}
	}
	if !utils.IsBlank(enterpriseCode) {
		attr.EnterpriseCode = &base.EnterpriseCodeAttr{
			Code:       enterpriseCode,
			FontSizePx: utils.ToIntPixel(dpi, codeFont),
			Arc:        codeArc,
		}
		footDist = footDistWithCode
		if !utils.IsBlank(headText) {
			headDistMm -= 1.5
		}
	}
	if !utils.IsBlank(headText) {
		attr.OvalHead = base.FootAttrMM(dpi, headText, headFont, headDistMm, headWidth)
	}
	if !utils.IsBlank(footText) {
		attr.Foot = base.FootAttrMM(dpi, footText, footFont, footDist, footWidth)
	}
	if attr.Logo != nil {
		up := float32(0)
		if !utils.IsBlank(enterpriseCode) && !utils.IsBlank(footText) && logoUpCodeFoot != nil {
			up = *logoUpCodeFoot
		} else if !utils.IsBlank(footText) && logoUpFoot != nil {
			up = *logoUpFoot
		}
		attr.Logo.UpMovePx = utils.ToFloatPixel(dpi, up)
	}
}

// attachDoublelineInnerContent 写入双线椭圆内圈文字。
func attachDoublelineInnerContent(
	attr *base.SealPaintAttribute,
	dpi int,
	innerContent string,
	in doublelineOvalInner,
) {
	if utils.IsBlank(innerContent) {
		return
	}
	attr.InnerOval = &base.SealPaintAttribute{
		Dpi:         dpi,
		Content:     innerContent,
		WidthPx:     utils.ToIntPixel(dpi, in.width),
		HeightPx:    utils.ToIntPixel(dpi, in.height),
		EdgeWidthPx: utils.ToFloatPixel(dpi, in.edge),
		FontSizePx:  utils.ToIntPixel(dpi, in.font),
		MarginPx:    utils.ToFloatPixel(dpi, in.margin),
		Arc:         in.arc,
	}
	attr.FontSizePx = utils.ToIntPixel(dpi, 3.6)
	attr.Content = ""
	attr.SealLogo = false
	attr.Logo = nil
}

// baseDoublelineOvalAttr 返回双线椭圆公共底稿属性。
func baseDoublelineOvalAttr(dpi int, width, height, edge, font, margin float32, arc int) *base.SealPaintAttribute {
	return &base.SealPaintAttribute{
		Dpi:          dpi,
		Transparency: model.DefaultSealTransparency,
		Category:     model.SealCategoryDoublelineOval,
		Arc:          arc,
		WidthPx:      utils.ToIntPixel(dpi, width),
		HeightPx:     utils.ToIntPixel(dpi, height),
		MarginPx:     utils.ToFloatPixel(dpi, margin),
		EdgeWidthPx:  utils.ToFloatPixel(dpi, edge),
		FontSizePx:   utils.ToIntPixel(dpi, font),
		FakeGap:      true,
		SealColor:    model.SealRed,
	}
}

// GetDoublelineOval 返回双线椭圆章属性。
func GetDoublelineOval(spec, footText, headText, enterpriseCode, innerContent string, dpi int) (*base.SealPaintAttribute, error) {
	switch spec {
	case model.SealSpecOval5038:
		return buildOval5038Doubleline(footText, headText, enterpriseCode, innerContent, dpi)
	case model.SealSpecOval4935:
		return buildOval4935Doubleline(footText, headText, enterpriseCode, innerContent, dpi)
	case model.SealSpecOval4030:
		return buildOval4030Doubleline(footText, headText, enterpriseCode, innerContent, dpi)
	case model.SealSpecOval3929:
		return buildOval3929Doubleline(footText, headText, enterpriseCode, innerContent, dpi)
	case model.SealSpecOval45303218:
		return buildOval45303218Doubleline(footText, headText, enterpriseCode, innerContent, dpi)
	default:
		return nil, base.ErrSealSpecNotSupported
	}
}

// IsDoublelineOvalSpec 判断规格名是否为双线椭圆系列。
func IsDoublelineOvalSpec(spec string) bool {
	switch spec {
	case model.SealSpecOval5038, model.SealSpecOval4935, model.SealSpecOval4030, model.SealSpecOval3929, model.SealSpecOval45303218:
		return true
	default:
		return false
	}
}

// buildOval4030Doubleline 组装 40×30 双线椭圆属性。
func buildOval4030Doubleline(footText, headText, enterpriseCode, innerContent string, dpi int) (*base.SealPaintAttribute, error) {
	attr := baseDoublelineOvalAttr(dpi, 40, 30, 1, 5, 0.8, 200)
	luFoot := float32(1)
	luBoth := float32(2)
	attr.Logo = &base.SealLogoAttr{SizePx: utils.ToIntPixel(dpi, 11), FiveStar: true}
	attr.SealLogo = true
	applyDoublelineHeadCodeFoot(attr, dpi, footText, headText, enterpriseCode,
		5, 6, 7, 22, 4, 3, 26, 2.2, 100, 2, &luFoot, &luBoth)
	attachDoublelineInnerContent(attr, dpi, innerContent, defaultDoublelineInnerOval())
	return attr, nil
}

// buildOval3929Doubleline 组装 39×29 双线椭圆属性。
func buildOval3929Doubleline(footText, headText, enterpriseCode, innerContent string, dpi int) (*base.SealPaintAttribute, error) {
	attr := baseDoublelineOvalAttr(dpi, 39, 29, 1, 5, 0.8, 200)
	luFoot := float32(1)
	luBoth := float32(2)
	attr.Logo = &base.SealLogoAttr{SizePx: utils.ToIntPixel(dpi, 11), FiveStar: true}
	attr.SealLogo = true
	applyDoublelineHeadCodeFoot(attr, dpi, footText, headText, enterpriseCode,
		5, 6, 0, 22, 3, 3, 26, 2.2, 100, 2, &luFoot, &luBoth)
	attachDoublelineInnerContent(attr, dpi, innerContent, defaultDoublelineInnerOval())
	return attr, nil
}

// buildOval4935Doubleline 组装 49×35 双线椭圆属性。
func buildOval4935Doubleline(footText, headText, enterpriseCode, innerContent string, dpi int) (*base.SealPaintAttribute, error) {
	attr := baseDoublelineOvalAttr(dpi, 49, 35, 1, 6, 0.8, 200)
	luFoot := float32(1)
	luBoth := float32(2)
	attr.Logo = &base.SealLogoAttr{SizePx: utils.ToIntPixel(dpi, 14), FiveStar: true}
	attr.SealLogo = true
	applyDoublelineHeadCodeFoot(attr, dpi, footText, headText, enterpriseCode,
		6, 8, 0, 25, 5, 4, 26, 2.2, 100, 3, &luFoot, &luBoth)
	attachDoublelineInnerContent(attr, dpi, innerContent, defaultDoublelineInnerOval())
	return attr, nil
}

// buildOval5038Doubleline 组装 50×38 双线椭圆属性。
func buildOval5038Doubleline(footText, headText, enterpriseCode, innerContent string, dpi int) (*base.SealPaintAttribute, error) {
	attr := baseDoublelineOvalAttr(dpi, 50, 38, 1, 6, 0.8, 200)
	luFoot := float32(2)
	luBoth := float32(2)
	attr.Logo = &base.SealLogoAttr{SizePx: utils.ToIntPixel(dpi, 11), FiveStar: true}
	attr.SealLogo = true
	applyDoublelineHeadCodeFoot(attr, dpi, footText, headText, enterpriseCode,
		7, 8, 0, 23, 6, 3, 26, 2.2, 100, 7, &luFoot, &luBoth)
	attachDoublelineInnerContent(attr, dpi, innerContent, defaultDoublelineInnerOval())
	return attr, nil
}

// buildOval45303218Doubleline 组装 45×30×32×18 双线椭圆规格。
func buildOval45303218Doubleline(footText, headText, enterpriseCode, innerContent string, dpi int) (*base.SealPaintAttribute, error) {
	attr := baseDoublelineOvalAttr(dpi, 45, 30, 1, 4.5, 0.3, 220)
	attr.Content = innerContent
	attr.SealLogo = false
	attr.InnerLine = &base.SealPaintAttribute{
		Dpi:         dpi,
		WidthPx:     utils.ToIntPixel(dpi, 32),
		HeightPx:    utils.ToIntPixel(dpi, 18),
		EdgeWidthPx: utils.ToFloatPixel(dpi, 1),
	}
	if !utils.IsBlank(headText) {
		attr.OvalHead = base.FootAttrMM(dpi, headText, 5, -6, 20)
	}
	ff := float32(5)
	if !utils.IsBlank(footText) {
		attr.Foot = base.FootAttrMM(dpi, footText, ff, -1, 20)
	}
	if !utils.IsBlank(enterpriseCode) {
		attr.EnterpriseCode = &base.EnterpriseCodeAttr{
			Code:       enterpriseCode,
			FontSizePx: utils.ToIntPixel(dpi, 2.5),
			Arc:        110,
		}
	}
	return attr, nil
}
