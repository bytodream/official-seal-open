package model

// DefaultCustomSealAttributeDPI 自定义章面默认 DPI（1000）。
const DefaultCustomSealAttributeDPI = 1000

// CustomSealAttribute 接口/业务层传入的自定义章面参数（JSON 可序列化）。
type CustomSealAttribute struct {
	// SealCategory 章面类型（如 enterprise、more_foot_circle）
	SealCategory *SealCategory `json:"sealCategory,omitempty"`
	// SavedSpecName 已保存的规格名称（与 SealCategoryType 同时存在时走 useSavedSpec）
	SavedSpecName string `json:"savedSpecName,omitempty"`
	// Content 环绕主文字（通常为公司名）
	Content string `json:"content,omitempty"`
	// InnerContent 内圈文字（如双圆章）
	InnerContent string `json:"innerContent,omitempty"`
	// Foot 下排单段文字（非多 Foots 列表时）
	Foot string `json:"foot,omitempty"`
	// Head 上排文字
	Head string `json:"head,omitempty"`
	// EnterpriseCode 企业信息编码（防伪码周围文字等）
	EnterpriseCode string `json:"enterpriseCode,omitempty"`
	// EdgeWidth 外圈线宽
	EdgeWidth *float32 `json:"edgeWidth,omitempty"`
	// InnerEdgeWidth 内圈线宽（圆章）
	InnerEdgeWidth *float32 `json:"innerEdgeWidth,omitempty"`
	// InnerMargin 内边距
	InnerMargin *float32 `json:"innerMargin,omitempty"`
	// LogoSize 中心徽标尺寸
	LogoSize *float32 `json:"logoSize,omitempty"`
	// SealLogo 是否显示中心图案
	SealLogo *bool `json:"sealLogo,omitempty"`
	// Antifake 是否启用防伪条纹
	Antifake *bool `json:"antifake,omitempty"`
	// AntifakeCode 掩膜文字：普通章 "OS"，共青团章 "Q"
	AntifakeCode string `json:"antifakeCode,omitempty"`
	// AntifakeGap 防伪条纹是否留间隙
	AntifakeGap *bool `json:"antifakeGap,omitempty"`
	// SealSpec 规格枚举名字符串
	SealSpec string `json:"sealSpec,omitempty"`
	// Dpi 输出分辨率
	Dpi *int `json:"dpi,omitempty"`
	// SealCategoryType 章面业务分类类型（如 foreignCNEN，与 SavedSpecName 配合）
	SealCategoryType string `json:"sealCategoryType,omitempty"`
	// LogoUpmove 徽标上移量
	LogoUpmove *float32 `json:"logoUpmove,omitempty"`
	// Arc 主文字弧度相关参数
	Arc *int `json:"arc,omitempty"`
	// SealHeight DIY 规格高度（mm）
	SealHeight *float32 `json:"sealHeight,omitempty"`
	// SealWidth DIY 规格宽度（mm）
	SealWidth *float32 `json:"sealWidth,omitempty"`
	// FontSize 主文字字号
	FontSize *float32 `json:"fontSize,omitempty"`
	// WordLinePadding 主文字行距/字距
	WordLinePadding *float32 `json:"wordLinePadding,omitempty"`
	// EnterpriseCodeFontSize 编码区字号
	EnterpriseCodeFontSize *float32 `json:"enterprisecodeFontSize,omitempty"`
	// EnterpriseCodeArc 编码区弧度
	EnterpriseCodeArc *int `json:"enterprisecodeArc,omitempty"`
	// EnterpriseCodeWordLinePadding 编码区字距
	EnterpriseCodeWordLinePadding *float32 `json:"enterprisecodeWordLinePadding,omitempty"`
	// Foots 多段下排文字（多横排章）
	Foots []SealFoot `json:"foots,omitempty"`
	// HeadParam 上排结构化参数（多横排章）
	HeadParam *SealFoot `json:"headParam,omitempty"`
	// FootFontSize 下排文字字号
	FootFontSize *float32 `json:"footFontSize,omitempty"`
	// FootWidth 下排文字宽度
	FootWidth *float32 `json:"footWidth,omitempty"`
	// FootDistance 下排文字距离
	FootDistance *float32 `json:"footDistance,omitempty"`
	// HeadFontSize 上排文字字号
	HeadFontSize *float32 `json:"headFontSize,omitempty"`
	// HeadWidth 上排文字宽度
	HeadWidth *float32 `json:"headWidth,omitempty"`
	// LineLength 辅助线长度
	LineLength *float32 `json:"lineLength,omitempty"`
	// LineWidth 辅助线宽度
	LineWidth *float32 `json:"lineWidth,omitempty"`
	// LineInterval 线间距
	LineInterval *float32 `json:"lineInterval,omitempty"`
	// HeadDistance 上排文字距离（more_foot_circle 会与 FontSize 联动校正）
	HeadDistance *float32 `json:"headDistance,omitempty"`
	// InnerOvalWidth 内椭圆宽
	InnerOvalWidth *float32 `json:"innerOvalWidth,omitempty"`
	// InnerOvalHeight 内椭圆高
	InnerOvalHeight *float32 `json:"innerOvalHeight,omitempty"`
	// InnerOvalEdgeWidth 内椭圆线宽
	InnerOvalEdgeWidth *float32 `json:"innerOvalEdgeWidth,omitempty"`
	// InnerOvalFontSize 内椭圆字号
	InnerOvalFontSize *float32 `json:"innerOvalFontSize,omitempty"`
	// InnerOvalWordLinePadding 内椭圆字距
	InnerOvalWordLinePadding *float32 `json:"innerOvalWordLinePadding,omitempty"`

	// ToSealAttributeFunc 可选钩子：由集成层注入完整 SealAttribute 映射，不可 JSON 序列化。
	ToSealAttributeFunc func(*CustomSealAttribute) (SealAttribute, error) `json:"-"`
}

// NewCustomSealAttribute 创建带默认开关与 DPI 的实例（中心图案、防伪、防伪间隙默认开启）。
func NewCustomSealAttribute() *CustomSealAttribute {
	// 三个 Boolean 必须各自独立取址：encoding/json 反序列化对非空指针是"写回原地址"
	// 共用同一个变量时，antifakeGap:false 会把 sealLogo/antifake 一起改成 false。
	sealLogo, antifake, antifakeGap := true, true, true
	dpi := DefaultCustomSealAttributeDPI
	return &CustomSealAttribute{
		SealLogo:    &sealLogo,
		Antifake:    &antifake,
		AntifakeGap: &antifakeGap,
		Dpi:         &dpi,
	}
}

// SealSpecEnum 返回章面规格名，空值时默认 CIRCULAR_42。
func (a *CustomSealAttribute) SealSpecEnum() string {
	if a == nil || IsBlank(a.SealSpec) {
		return SealSpecCircular42
	}
	return a.SealSpec
}

// UseSavedSpec 判断是否使用已保存规格（需同时有 savedSpecName 与 sealCategoryType）。
func (a *CustomSealAttribute) UseSavedSpec() bool {
	if a == nil {
		return false
	}
	return !IsBlank(a.SavedSpecName) && !IsBlank(a.SealCategoryType)
}

// IsSealLogo 是否显示中心图案，指针为空时默认 true。
func (a *CustomSealAttribute) IsSealLogo() bool {
	if a == nil || a.SealLogo == nil {
		return true
	}
	return *a.SealLogo
}

// IsAntifake 是否启用防伪纹理，指针为空时默认 true。
func (a *CustomSealAttribute) IsAntifake() bool {
	if a == nil || a.Antifake == nil {
		return true
	}
	return *a.Antifake
}

// IsAntifakeGap 是否启用防伪留白间隙，指针为空时默认 true。
func (a *CustomSealAttribute) IsAntifakeGap() bool {
	if a == nil || a.AntifakeGap == nil {
		return true
	}
	return *a.AntifakeGap
}

// GetDpi 返回输出 DPI，未配置时 1000。
func (a *CustomSealAttribute) GetDpi() int {
	if a == nil || a.Dpi == nil {
		return DefaultCustomSealAttributeDPI
	}
	return *a.Dpi
}
