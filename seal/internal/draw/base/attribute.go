package base

import (
	"errors"
	"image/color"

	"github.com/bytodream/official-seal-open/seal/model"
)

// ErrSealSpecNotSupported 规格尚未移植
var ErrSealSpecNotSupported = errors.New("获取印章属性规格失败: 不支持的 SealSpec")

// SealPaintAttribute 章面绘制属性（各章型规格构建器的统一产物，像素级）。
// 宽高等为像素；边宽、边距等为 float 像素值，由 draw 包消费。
type SealPaintAttribute struct {
	// Dpi 输出分辨率（像素/英寸）
	Dpi int
	// Transparency 章面 Alpha（0–255）
	Transparency int
	// Category 章面品类（决定绘制实现）
	Category model.SealCategory
	// Content 环绕主文字
	Content string
	// Arc 主文字弧度
	Arc int
	// WidthPx 章面宽度（px）
	WidthPx int
	// HeightPx 章面高度（px）
	HeightPx int
	// MarginPx 主文字与外圈边距（px）
	MarginPx float32
	// EdgeWidthPx 外圈线宽（px）
	EdgeWidthPx float32
	// FontSizePx 主文字字号（px）
	FontSizePx int
	// SealLogo 是否绘制中心图案
	SealLogo bool
	// FakeGap 是否绘制防伪留白间隙
	FakeGap bool
	// Logo 中心徽标尺寸与位置
	Logo *SealLogoAttr
	// Foot 脚排文字
	Foot *FootAttr
	// EnterpriseCode 企业代码弧排参数
	EnterpriseCode *EnterpriseCodeAttr
	// Antifake 防伪纹参数（为 nil 时退化绘制五角星）
	Antifake *model.AntifakeInfo
	// InnerEdgeWidthPx 内圈线宽（px）
	InnerEdgeWidthPx float32
	// InnerMarginPx 内圈文字边距（px）
	InnerMarginPx float32
	// SealColor 印泥颜色
	SealColor color.Color
	// YouthLeague 共青团章专用装饰参数
	YouthLeague *YouthLeagueAttr
	// MoreFootHead 多横排章上排头文字
	MoreFootHead *FootAttr
	// MoreFoots 多横排章各段下排文字
	MoreFoots []*FootAttr
	// HeadFoot 双线圆章上排横字（doubleline_cir）
	HeadFoot *FootAttr
	// InnerCircle 双圆章内圈（double_cir）
	InnerCircle *SealPaintAttribute
	// InnerOval 内椭圆文字（OVAL_45_30 及双线椭圆各规格 + innerContent，默认 33×19）
	InnerOval *SealPaintAttribute
	// InnerLine 双行椭圆内环描边（OVAL_45_30_32_18）
	InnerLine *SealPaintAttribute
	// OvalHead 双行椭圆上排横字
	OvalHead *FootAttr
}

// YouthLeagueAttr 共青团章专用装饰参数。
type YouthLeagueAttr struct {
	// LineLengthPx 装饰线长度（px）
	LineLengthPx float32
	// LineWidthPx 装饰线宽（px）
	LineWidthPx float32
	// LineIntervalPx 装饰线间距（px）
	LineIntervalPx float32
	// Head 上排头文字
	Head *FootAttr
}

// SealLogoAttr 中心徽标尺寸与位置（构建前为 mm）。
type SealLogoAttr struct {
	// SizePx 徽标边长（px）
	SizePx int
	// UpMovePx 徽标上移距离（px）
	UpMovePx float32
	// FiveStar 是否以五角星为徽标
	FiveStar bool
}

// FootAttr 单段脚注的文字与排版参数。
type FootAttr struct {
	// Text 脚注文字
	Text string
	// FontSizePx 字号（px）
	FontSizePx int
	// DistancePx 与主文字/上一排的距离（px）
	DistancePx float32
	// WidthPx 排版区宽度（px）
	WidthPx float32
}

// EnterpriseCodeAttr 企业编码弧排参数。
type EnterpriseCodeAttr struct {
	// Code 企业编码数字串
	Code string
	// FontSizePx 字号（px）
	FontSizePx int
	// Arc 编码弧度
	Arc int
	// MarginPx 与主文字的弧向边距（px）
	MarginPx float32
}

// SealAttributeBundle 章面生成器入参集合（当前仅含像素级章面属性 Paint）。
type SealAttributeBundle struct {
	Paint *SealPaintAttribute
}
