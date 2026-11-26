package base

import (
	"image/color"

	"github.com/bytodream/official-seal-open/seal/model"
)

// SealPaintModel 章面绘制参数（像素级，由 base.SealPaintAttribute 转换而来）。
type SealPaintModel struct {
	// Dpi 输出分辨率（像素/英寸）
	Dpi int
	// Width 画布宽度（px）
	Width int
	// Height 画布高度（px）
	Height int
	// Content 环绕主文字
	Content string
	// Arc 主文字弧度
	Arc int
	// Margin 主文字与外圈边距（px）
	Margin float32
	// EdgeWidth 外圈线宽（px）
	EdgeWidth float32
	// FontSize 主文字字号（px）
	FontSize int
	// SealLogo 是否绘制中心图案
	SealLogo bool
	// LogoSize 中心图案边长（px）
	LogoSize int
	// LogoUpmove 中心图案上移距离（px）
	LogoUpmove float32
	// FakeGap 是否绘制防伪留白间隙
	FakeGap bool
	// Antifake 是否绘制防伪纹
	Antifake bool
	// Transparency 章面 Alpha（0–255）
	Transparency int
	// FontFile 字体文件路径
	FontFile string
	// AntifakeInfo 防伪纹参数（为 nil 时中心退化绘制五角星）
	AntifakeInfo *model.AntifakeInfo
	// SealColor 印泥颜色
	SealColor color.Color
	// Category 章面品类（决定 DrawSeal 分发的绘制实现）
	Category model.SealCategory
	// InnerEdgeWidth 内圈线宽（px）
	InnerEdgeWidth float32
	// InnerMargin 内圈文字边距（px）
	InnerMargin float32
	// Foot 脚排文字
	Foot *FootAttr
	// EnterpriseCode 企业代码弧排参数
	EnterpriseCode *EnterpriseCodeAttr
	// YouthLeague 共青团章专用装饰参数
	YouthLeague *YouthLeagueAttr
	// LogoFiveStar 中心图案是否为五角星
	LogoFiveStar bool
	// MoreFootHead 多横排章上排头文字
	MoreFootHead *FootAttr
	// MoreFoots 多横排章各段下排文字
	MoreFoots []*FootAttr
	// HeadFoot 双行圆章上排横字（doubleline_cir）
	HeadFoot *FootAttr
	// InnerCircle 双圆章内圈（double_cir）
	InnerCircle *SealPaintAttribute
	// InnerOval 椭圆章内椭圆（OVAL_45_30 + innerContent）
	InnerOval *SealPaintAttribute
	// InnerLine 双行椭圆内环描边（OVAL_45_30_32_18）
	InnerLine *SealPaintAttribute
	// OvalHead 双行椭圆上排横字
	OvalHead *FootAttr
	// ShowWatermark 是否在章面上叠加「样章」标识（预览）
	ShowWatermark bool
	// WatermarkText 「样章」标识文字，为空时使用默认文字
	WatermarkText string
}

// PaintOption 绘制模型可选配置。
type PaintOption func(*SealPaintModel)

// WithSampleWatermark 将章面标记为预览样章：Draw 收尾时叠加「样章」标识；text 为空时使用默认文字。
func WithSampleWatermark(text string) PaintOption {
	return func(m *SealPaintModel) {
		m.ShowWatermark = true
		m.WatermarkText = text
	}
}

// BuildSealPaintModel 从 base.SealPaintAttribute 构建。
func BuildSealPaintModel(attr *SealPaintAttribute, transparencyAlpha int, fontFile string, opts ...PaintOption) SealPaintModel {
	m := SealPaintModel{
		Dpi:            attr.Dpi,
		Width:          attr.WidthPx,
		Height:         attr.HeightPx,
		Content:        attr.Content,
		Arc:            attr.Arc,
		Margin:         attr.MarginPx,
		EdgeWidth:      attr.EdgeWidthPx,
		FontSize:       attr.FontSizePx,
		SealLogo:       attr.SealLogo,
		FakeGap:        attr.FakeGap,
		Antifake:       attr.Antifake != nil,
		Transparency:   transparencyAlpha,
		FontFile:       fontFile,
		Foot:           attr.Foot,
		EnterpriseCode: attr.EnterpriseCode,
		YouthLeague:    attr.YouthLeague,
		MoreFootHead:   attr.MoreFootHead,
		MoreFoots:      attr.MoreFoots,
		HeadFoot:       attr.HeadFoot,
		InnerCircle:    attr.InnerCircle,
		InnerOval:      attr.InnerOval,
		InnerLine:      attr.InnerLine,
		OvalHead:       attr.OvalHead,
		// 优先使用章面属性自带的印泥色；未设置时由绘制层回退 model.SealRed。
		SealColor:      attr.SealColor,
		Category:       attr.Category,
		InnerEdgeWidth: attr.InnerEdgeWidthPx,
		InnerMargin:    attr.InnerMarginPx,
	}
	// Transparency 由调用方传入的 transparencyAlpha 决定，覆盖 attr.Transparency 的构建期默认值。
	if attr.Logo != nil {
		m.LogoSize = attr.Logo.SizePx
		m.LogoUpmove = attr.Logo.UpMovePx
		m.LogoFiveStar = attr.Logo.FiveStar
	} else {
		m.LogoFiveStar = true
	}
	if attr.Antifake != nil {
		info := *attr.Antifake
		info.FontFile = fontFile
		m.AntifakeInfo = &info
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&m)
		}
	}
	return m
}
