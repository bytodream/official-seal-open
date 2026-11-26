package base

import (
	"github.com/bytodream/official-seal-open/seal/model"
	"github.com/bytodream/official-seal-open/seal/utils"
)

// FootInput 多横排章的单段横排文字入参（mm 单位，由规格层转像素）。
type FootInput struct {
	// Text 横排文字内容
	Text string
	// FontSize 字号（mm）
	FontSize float32
	// Distance 距离偏移（mm）
	Distance float32
	// Width 文字区域宽度（mm）
	Width float32
}

// BuildMoreFootFootInputs 从自定义属性装配多横排 Foot/Head 入参。
// HeadParam 优先于 Head 字段；各段缺省字号 5mm、宽度 0（绘制层按自然字宽处理）。
func BuildMoreFootFootInputs(a *model.CustomSealAttribute, dpi int) (head *FootInput, foots []FootInput) {
	if a == nil {
		return nil, nil
	}
	const defaultFootFont = float32(5)
	const defaultFootWidth = float32(0)

	if a.HeadParam != nil {
		head = &FootInput{
			Text:     a.HeadParam.Text,
			FontSize: utils.F32Or(a.HeadParam.FontSizeMM(), defaultFootFont),
			Distance: utils.F32Or(a.HeadParam.Distance, 0),
			Width:    utils.F32Or(a.HeadParam.Width, defaultFootWidth),
		}
	} else if !utils.IsBlank(a.Head) {
		head = &FootInput{
			Text:     a.Head,
			FontSize: utils.F32Or(a.HeadFontSize, defaultFootFont),
			Distance: utils.F32Or(a.HeadDistance, 0),
			Width:    utils.F32Or(a.HeadWidth, defaultFootWidth),
		}
	}

	for _, sf := range a.Foots {
		foots = append(foots, FootInput{
			Text:     sf.Text,
			FontSize: utils.F32Or(sf.FontSizeMM(), defaultFootFont),
			Distance: utils.F32Or(sf.Distance, 0),
			Width:    utils.F32Or(sf.Width, defaultFootWidth),
		})
	}
	return head, foots
}

// FootAttrMM 将毫米字段转为像素 FootAttr；宽度缺省（<=0）时取 20mm。
func FootAttrMM(dpi int, text string, fontMm, distanceMm, widthMm float32) *FootAttr {
	if widthMm <= 0 {
		widthMm = 20
	}
	return &FootAttr{
		Text:       text,
		FontSizePx: utils.ToIntPixel(dpi, fontMm),
		DistancePx: utils.ToFloatPixel(dpi, distanceMm),
		WidthPx:    utils.ToFloatPixel(dpi, widthMm),
	}
}
