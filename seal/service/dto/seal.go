package dto

import "image"

// SealType 印章类型枚举。
type SealType string

const (
	// SealTypeNotSupported 不支持的印章类型
	SealTypeNotSupported SealType = "NOT_SUPPORTED"
	// SealTypeCompany 企业公章
	SealTypeCompany SealType = "COMPANY"
	// SealTypePersonal 个人签名
	SealTypePersonal SealType = "PERSONAL"
	// SealTypePractice 个人执业章
	SealTypePractice SealType = "PRACTICE"
	// SealTypeLP 法定代表人章
	SealTypeLP SealType = "LP"
)

// BusinessType 签章业务类型。
type BusinessType string

const (
	// BusinessTypePhysics 物理签章
	BusinessTypePhysics BusinessType = "PHYSICS"
	// BusinessTypeElectronic 电子签章
	BusinessTypeElectronic BusinessType = "ELECTRONIC"
	// BusinessTypePhysicsAndElectronic 电子签章和物理签章
	BusinessTypePhysicsAndElectronic BusinessType = "PHYSICSANDELECTRONIC"
)

const (
	// SealSpecCircular42 直径 42mm 圆章规格常量
	SealSpecCircular42 = "CIRCULAR_42"
	// SealSpecCircular50 直径 50mm 多横排圆章规格常量
	SealSpecCircular50 = "CIRCULAR_50"
	// SealSpecOval4530 45×30mm 椭圆多横排章规格常量
	SealSpecOval4530 = "OVAL_45_30"
	// SealSpecOval45303218 45×30mm 双行椭圆章规格常量（内圈 32×18mm）
	SealSpecOval45303218 = "OVAL_45_30_32_18"
	// SealSpecRectangle5020 50×20mm 长方形多横排章规格常量
	SealSpecRectangle5020 = "RECTANGLE_50_20"
)

// Seal 印章实体（服务层视图）。
type Seal struct {
	// Radius 模糊程度参数
	Radius *int `json:"radius,omitempty"`
	// Transparency 透明度配置值（0–100 字符串，与前端约定一致）
	Transparency *string `json:"transparency,omitempty"`
}

// SealImage 章面上传后的各版本 fileKey 与大小。
type SealImage struct {
	// OriginImageKey 原始上传图 fileKey（上传章面时有值）
	OriginImageKey string `json:"orignImageKey,omitempty"`
	// ImageKey 主章面图 fileKey
	ImageKey string `json:"imageKey,omitempty"`
	// BlurImageKey 模糊章面图 fileKey
	BlurImageKey string `json:"blurImageKey,omitempty"`
	// Size 主图字节大小，用于存储计费
	Size int64 `json:"size,omitempty"`
}

// ImageAttr 控制章面上传、缩放与旋转等行为。
type ImageAttr struct {
	// DefaultAttr 是否使用默认属性组合
	DefaultAttr bool `json:"defaultAttr,omitempty"`
	// Transparent 是否按透明图处理
	Transparent bool `json:"transparent,omitempty"`
	// Zoom 是否缩放
	Zoom bool `json:"zoom,omitempty"`
	// Rotate 是否旋转
	Rotate bool `json:"rotate,omitempty"`
	// LimitWidth 限定最大宽度
	LimitWidth int `json:"limitWidth,omitempty"`
	// LimitHeight 限定最大高度
	LimitHeight int `json:"limitHeight,omitempty"`
	// Degree 旋转角度
	Degree int `json:"degree,omitempty"`
	// Transparency 透明度（字符串，上传自定义章时使用）
	Transparency string `json:"transparency,omitempty"`
}

// ImageAttrDefault 默认图像上传属性（仅 DefaultAttr 开启）。
var ImageAttrDefault = ImageAttr{DefaultAttr: true}

// SealModel 创建或更新印章的请求模型。
type SealModel struct {
	// Category 签章类型（物理/电子）
	Category BusinessType `json:"category,omitempty"`
	// ElectronicModel 电子章扩展参数（Category 为 ELECTRONIC 时有效）
	ElectronicModel *ElectronicModel `json:"electronicModel,omitempty"`
}

// ElectronicModel 电子章扩展参数。
type ElectronicModel struct {
	// Width 章面宽度
	Width *int `json:"width,omitempty"`
	// Height 章面高度
	Height *int `json:"height,omitempty"`
	// OpenImageBlur 预览/展示是否使用模糊图
	OpenImageBlur bool `json:"openImageBlur,omitempty"`
	// IsAuto 是否为自动生成章面（相对上传图）
	IsAuto *bool `json:"isAuto,omitempty"`
	// SealSpec 规格枚举名
	SealSpec string `json:"sealSpec,omitempty"`
	// SealImageKey 已有章面 fileKey
	SealImageKey string `json:"sealImageKey,omitempty"`
	// SealBlurImageKey 已有模糊章面 fileKey
	SealBlurImageKey string `json:"sealBlurImageKey,omitempty"`
	// UploadImage Base64 上传章面
	UploadImage string `json:"uploadImage,omitempty"`
	// UploadImageKey 已上传章面 fileKey
	UploadImageKey string `json:"uploadImageKey,omitempty"`
	// AutoAttribute 自定义章面参数（自动生成）
	AutoAttribute *CustomSealAttribute `json:"autoAttribute,omitempty"`
	// IsCreateByAuth 是否认证通过后自动创建
	IsCreateByAuth bool `json:"isCreateByAuth,omitempty"`
	// ShowWatermark 是否叠加「样章」标识；Preview 固定为 true
	ShowWatermark bool `json:"showWatemark,omitempty"`
	// WatermarkText 「样章」标识自定义文字，为空时使用默认文字
	WatermarkText string `json:"watermarkText,omitempty"`
	// Radius 模糊半径参数
	Radius *int `json:"radius,omitempty"`
	// Transparency 透明度
	Transparency *string `json:"transparency,omitempty"`
	// SpecEntityId 规格实体 ID
	SpecEntityId string `json:"specEntityId,omitempty"`
}

// ToSeal 转为 Preview / 生成用的 Seal 视图。
func (m *SealModel) ToSeal() *Seal {
	s := &Seal{}
	if m != nil && m.Category == BusinessTypeElectronic && m.ElectronicModel != nil {
		em := m.ElectronicModel
		s.Radius = em.Radius
		s.Transparency = em.Transparency
	}
	return s
}

// RasterImage 章面位图（image.Image 类型别名）。
type RasterImage = image.Image
