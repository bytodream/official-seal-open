package model

import (
	"encoding/json"
	"strconv"
)

// CreateSealAPIRequest 对应前端/开放接口创建印章 JSON（SealModel 序列化形态）。
// 仅保留章面生成相关字段；OperateType/Owner/Department 等业务字段本库不消费，
// JSON 中出现时由 Go 解析器自动忽略，调用方如需请自行定义承载结构。
type CreateSealAPIRequest struct {
	// Category 签章业务类型（PHYSICS/ELECTRONIC）
	Category string `json:"category"`
	// ElectronicModel 电子章扩展参数
	ElectronicModel *ElectronicModelAPI `json:"electronicModel"`
}

// ElectronicModelAPI 电子章扩展参数（接口 JSON 形态）。
type ElectronicModelAPI struct {
	// SealSpec 规格枚举名（如 CIRCULAR_42）
	SealSpec string `json:"sealSpec"`
	// Width 章面宽度（px，兼容字符串/数字）
	Width FlexString `json:"width"`
	// Height 章面高度（px，兼容字符串/数字）
	Height FlexString `json:"height"`
	// OpenImageBlur 展示是否使用模糊章面
	OpenImageBlur bool `json:"openImageBlur"`
	// Radius 模糊半径参数
	Radius int `json:"radius"`
	// Transparency 透明度（0–100，兼容字符串/数字）
	Transparency json.RawMessage `json:"transparency"`
	// Auto 是否自动生成章面（相对上传图）
	Auto bool `json:"auto"`
	// SpecEntityID 规格实体 ID
	SpecEntityID string `json:"specEntityId"`
	// AutoAttribute 自定义章面参数
	AutoAttribute *CustomSealAttributeAPI `json:"autoAttribute"`
	// ShowWatermark 是否叠加「样章」标识（Preview 固定开启）
	ShowWatermark bool `json:"showWatemark"`
	// WatermarkText 「样章」标识自定义文字，为空使用默认文字
	WatermarkText string `json:"watermarkText"`
}

// CustomSealAttributeAPI 接口/业务层传入的自定义章面参数（JSON 可序列化，尺寸单位 mm）。
type CustomSealAttributeAPI struct {
	// SealCategory 章面品类名（如 enterprise/doubleline_cir，走名称解析）
	SealCategory string `json:"sealCategory"`
	// SealCategoryType 章面品类类型名（与 SealCategory 配合用于已保存规格）
	SealCategoryType string `json:"sealCategoryType"`
	// EdgeWidth 外圈线宽（mm）
	EdgeWidth FlexFloat32 `json:"edgeWidth"`
	// InnerEdgeWidth 内圈线宽（mm）
	InnerEdgeWidth FlexFloat32 `json:"innerEdgeWidth"`
	// SealLogo 是否显示中心图案
	SealLogo *bool `json:"sealLogo"`
	// LogoSize 中心图案尺寸（mm）
	LogoSize FlexFloat32 `json:"logoSize"`
	// Content 环绕主文字
	Content string `json:"content"`
	// InnerContent 内圈文字
	InnerContent string `json:"innerContent"`
	// Head 头排文字（多横排章）
	Head string `json:"head"`
	// Foot 脚排文字
	Foot string `json:"foot"`
	// EnterpriseCode 企业代码（弧排数字）
	EnterpriseCode string `json:"enterpriseCode"`
	// Antifake 是否启用防伪纹
	Antifake *bool `json:"antifake"`
	// AntifakeGap 是否启用防伪留白间隙
	AntifakeGap *bool `json:"antifakeGap"`
	// SealWidth 章面宽度（mm）
	SealWidth FlexFloat32 `json:"sealWidth"`
	// SealHeight 章面高度（mm）
	SealHeight FlexFloat32 `json:"sealHeight"`
	// LineWidth 装饰线宽（mm）
	LineWidth FlexFloat32 `json:"lineWidth"`
	// FontSize 主文字字号（mm）
	FontSize FlexFloat32 `json:"fontSize"`
	// WordLinePadding 主文字字距（mm）
	WordLinePadding FlexFloat32 `json:"wordLinePadding"`
	// EnterpriseCodeFontSize 企业代码字号（mm）
	EnterpriseCodeFontSize FlexFloat32 `json:"enterprisecodeFontSize"`
	// HeadWidth 头排宽度（mm）
	HeadWidth FlexFloat32 `json:"headWidth"`
	// HeadDistance 头排与主文字距离（mm）
	HeadDistance FlexFloat32 `json:"headDistance"`
	// HeadFontSize 头排字号（mm）
	HeadFontSize FlexFloat32 `json:"headFontSize"`
	// FootWidth 脚排宽度（mm）
	FootWidth FlexFloat32 `json:"footWidth"`
	// FootDistance 脚排与主文字距离（mm）
	FootDistance FlexFloat32 `json:"footDistance"`
	// FootFontSize 脚排字号（mm）
	FootFontSize FlexFloat32 `json:"footFontSize"`
	// LogoUpmove 中心图案上移距离（mm）
	LogoUpmove FlexFloat32 `json:"logoUpmove"`
	// LineInterval 装饰线间距（mm）
	LineInterval FlexFloat32 `json:"lineInterval"`
	// LineLength 装饰线长度（mm）
	LineLength FlexFloat32 `json:"lineLength"`
	// InnerOvalWidth 内椭圆宽度（mm）
	InnerOvalWidth FlexFloat32 `json:"innerOvalWidth"`
	// InnerMargin 内椭圆文字边距（mm）
	InnerMargin FlexFloat32 `json:"innerMargin"`
	// InnerOvalHeight 内椭圆高度（mm）
	InnerOvalHeight FlexFloat32 `json:"innerOvalHeight"`
	// InnerOvalFontSize 内椭圆文字字号（mm）
	InnerOvalFontSize FlexFloat32 `json:"innerOvalFontSize"`
	// InnerOvalWordLinePadding 内椭圆文字字距（mm）
	InnerOvalWordLinePadding FlexFloat32 `json:"innerOvalWordLinePadding"`
	// InnerOvalEdgeWidth 内椭圆线宽（mm）
	InnerOvalEdgeWidth FlexFloat32 `json:"innerOvalEdgeWidth"`
	// Foots 多横排脚注列表
	Foots []SealFoot `json:"foots"`
	// Arc 主文字弧度
	Arc FlexInt `json:"arc"`
}

// ParseCreateSealAPIRequest 解析创建印章 JSON。
func ParseCreateSealAPIRequest(data []byte) (*CreateSealAPIRequest, error) {
	var req CreateSealAPIRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// ToSealModel 转为 CompanySealCreateService.Preview 使用的 SealModel。
func (req *CreateSealAPIRequest) ToSealModel() (*SealModel, error) {
	if req == nil {
		return nil, ErrSealAttributeNotConfigured
	}
	model := &SealModel{
		Category: BusinessType(req.Category),
	}
	if req.ElectronicModel == nil {
		return model, nil
	}
	em := req.ElectronicModel
	auto := em.Auto
	transparency, err := ParseFlexTransparencyField(em.Transparency)
	if err != nil {
		return nil, err
	}
	radius := em.Radius
	domainEM := &ElectronicModel{
		SealSpec:      em.SealSpec,
		OpenImageBlur: em.OpenImageBlur,
		Radius:        &radius,
		IsAuto:        &auto,
		SpecEntityId:  em.SpecEntityID,
		ShowWatermark: em.ShowWatermark,
		WatermarkText: em.WatermarkText,
	}
	if transparency != "" {
		domainEM.Transparency = &transparency
	}
	if w := em.Width.String(); w != "" {
		if v, err := strconv.Atoi(w); err == nil {
			domainEM.Width = &v
		}
	}
	if h := em.Height.String(); h != "" {
		if v, err := strconv.Atoi(h); err == nil {
			domainEM.Height = &v
		}
	}
	if em.AutoAttribute != nil {
		attr := em.AutoAttribute.ToCustomSealAttribute(em.SealSpec)
		domainEM.AutoAttribute = attr
	}
	model.ElectronicModel = domainEM
	return model, nil
}

// ToCustomSealAttribute 将 API autoAttribute 转为领域对象。
func (a *CustomSealAttributeAPI) ToCustomSealAttribute(sealSpecFromParent string) *CustomSealAttribute {
	if a == nil {
		return nil
	}
	attr := NewCustomSealAttribute()
	attr.SealCategoryType = a.SealCategoryType
	attr.Content = a.Content
	attr.InnerContent = a.InnerContent
	attr.Head = a.Head
	attr.Foot = a.Foot
	attr.EnterpriseCode = a.EnterpriseCode
	attr.Foots = a.Foots
	if a.SealCategory != "" {
		cat := ConvertSealCategory(a.SealCategory)
		attr.SealCategory = &cat
	}
	if sealSpecFromParent != "" {
		attr.SealSpec = sealSpecFromParent
	}
	attr.EdgeWidth = a.EdgeWidth.Value
	attr.InnerEdgeWidth = a.InnerEdgeWidth.Value
	attr.LogoSize = a.LogoSize.Value
	attr.SealWidth = a.SealWidth.Value
	attr.SealHeight = a.SealHeight.Value
	attr.LineWidth = a.LineWidth.Value
	attr.FontSize = a.FontSize.Value
	attr.WordLinePadding = a.WordLinePadding.Value
	attr.EnterpriseCodeFontSize = a.EnterpriseCodeFontSize.Value
	attr.HeadWidth = a.HeadWidth.Value
	attr.HeadDistance = a.HeadDistance.Value
	attr.HeadFontSize = a.HeadFontSize.Value
	attr.FootWidth = a.FootWidth.Value
	attr.FootDistance = a.FootDistance.Value
	attr.FootFontSize = a.FootFontSize.Value
	attr.LogoUpmove = a.LogoUpmove.Value
	attr.LineInterval = a.LineInterval.Value
	attr.LineLength = a.LineLength.Value
	attr.InnerOvalWidth = a.InnerOvalWidth.Value
	attr.InnerMargin = a.InnerMargin.Value
	attr.InnerOvalHeight = a.InnerOvalHeight.Value
	attr.InnerOvalFontSize = a.InnerOvalFontSize.Value
	attr.InnerOvalWordLinePadding = a.InnerOvalWordLinePadding.Value
	attr.InnerOvalEdgeWidth = a.InnerOvalEdgeWidth.Value
	attr.Arc = a.Arc.Value
	if a.SealLogo != nil {
		attr.SealLogo = a.SealLogo
	}
	if a.Antifake != nil {
		attr.Antifake = a.Antifake
	}
	if a.AntifakeGap != nil {
		attr.AntifakeGap = a.AntifakeGap
	}
	return attr
}
