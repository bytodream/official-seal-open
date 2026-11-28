// Package seal 企业电子公章对外门面；集成方应只依赖本包 API。
package seal

import (
	"github.com/bytodream/official-seal-open/seal/service"
	"github.com/bytodream/official-seal-open/seal/service/dto"
)

// 服务与授权相关的导出别名。
type (
	// CompanySealCreateService 企业印章创建服务。
	CompanySealCreateService = service.CompanySealCreateService
	// ImageService 图像服务接口。
	ImageService = service.ImageService
	// ImageServiceConfigurator 图像服务的可选配置能力。
	ImageServiceConfigurator = service.ImageServiceConfigurator
	// WatermarkFunc 章面水印回调函数类型。
	WatermarkFunc = service.WatermarkFunc
	// License 授权信息。
	License = service.License
)

// 印章数据模型相关的导出别名。
type (
	// Seal 印章基础信息。
	Seal = dto.Seal
	// SealImage 印章图像。
	SealImage = dto.SealImage
	// SealModel 印章模型。
	SealModel = dto.SealModel
	// ElectronicModel 电子印章参数。
	ElectronicModel = dto.ElectronicModel
	// CustomSealAttribute 自定义印章属性。
	CustomSealAttribute = dto.CustomSealAttribute
	// SealType 印章类型。
	SealType = dto.SealType
	// BusinessType 业务类型。
	BusinessType = dto.BusinessType
	// SealCategory 印章品类/章型。
	SealCategory = dto.SealCategory
	// SealFoot 印章附加信息。
	SealFoot = dto.SealFoot
	// DiySpec 自定义规格描述。
	DiySpec = dto.DiySpec
	// ImageAttr 图像上传属性。
	ImageAttr = dto.ImageAttr
)

// 授权与错误相关的导出函数/变量。
var (
	// SetLicenseKey 设置授权密钥。
	SetLicenseKey = service.SetLicenseKey
	// GetLicenseKey 获取当前授权信息；未设置时返回 nil。
	GetLicenseKey = service.GetLicenseKey
	// EnsureLicenseActive 校验授权是否处于有效状态，无效时返回错误。
	EnsureLicenseActive = service.EnsureLicenseActive
	// ErrLicenseRequired 未设置授权密钥。
	ErrLicenseRequired = service.ErrLicenseRequired
	// ErrLicenseExpired 授权已过期。
	ErrLicenseExpired = service.ErrLicenseExpired
	// ErrSealCategoryNotLicensed 该章型未包含在授权范围内。
	ErrSealCategoryNotLicensed = service.ErrSealCategoryNotLicensed
	// ErrUsageLimitExceeded 授权使用次数已用尽。
	ErrUsageLimitExceeded = service.ErrUsageLimitExceeded
	// ImageAttrDefault 默认图像上传属性。
	ImageAttrDefault = dto.ImageAttrDefault
	// ErrSealAttributeNotConfigured 章面参数未配置或无法构建规格。
	ErrSealAttributeNotConfigured = dto.ErrSealAttributeNotConfigured
	// ErrSealCategoryNotSupported 章面品类/章型当前不支持。
	ErrSealCategoryNotSupported = dto.ErrSealCategoryNotSupported
	// ErrImageServiceNotConfigured 章面图片服务未注入，无法上传生成的位图。
	ErrImageServiceNotConfigured = dto.ErrImageServiceNotConfigured
	// ErrImageUploadNotConfigured 章面图上传实现未配置。
	ErrImageUploadNotConfigured = dto.ErrImageUploadNotConfigured
)

// 印章类型、业务类型与各章型规格/品类的导出常量（语义同 dto 包原常量）。
const (
	// SealTypeCompany 企业公章
	SealTypeCompany = dto.SealTypeCompany
	// BusinessTypeElectronic 电子签章
	BusinessTypeElectronic = dto.BusinessTypeElectronic

	// SealSpecCircular42 直径 42mm 圆章规格常量
	SealSpecCircular42 = dto.SealSpecCircular42
	// SealSpecCircular50 直径 50mm 多横排圆章规格常量
	SealSpecCircular50 = dto.SealSpecCircular50
	// SealSpecOval4530 45×30mm 椭圆多横排章规格常量
	SealSpecOval4530 = dto.SealSpecOval4530
	// SealSpecOval45303218 45×30mm 双行椭圆章规格常量（内圈 32×18mm）
	SealSpecOval45303218 = dto.SealSpecOval45303218
	// SealSpecRectangle5020 50×20mm 长方形多横排章规格常量
	SealSpecRectangle5020 = dto.SealSpecRectangle5020

	// SealCategoryCommunist 党委公章
	SealCategoryCommunist = dto.SealCategoryCommunist
	// SealCategoryEnterprise 普通公章
	SealCategoryEnterprise = dto.SealCategoryEnterprise
	// SealCategoryLabour 工会公章
	SealCategoryLabour = dto.SealCategoryLabour
	// SealCategoryYouthLeague 共青团
	SealCategoryYouthLeague = dto.SealCategoryYouthLeague
	// SealCategoryCPPCC 中国人民政治协商会议章
	SealCategoryCPPCC = dto.SealCategoryCPPCC
	// SealCategoryOval 椭圆章
	SealCategoryOval = dto.SealCategoryOval
	// SealCategoryRectangleTimestamp 长方形时间戳章
	SealCategoryRectangleTimestamp = dto.SealCategoryRectangleTimestamp
	// SealCategoryCSJShusp 长三角一网通办
	SealCategoryCSJShusp = dto.SealCategoryCSJShusp
	// SealCategoryMoreFootCircle 多横排文字圆
	SealCategoryMoreFootCircle = dto.SealCategoryMoreFootCircle
	// SealCategoryMoreFootRectangle 多横排文字长方形章
	SealCategoryMoreFootRectangle = dto.SealCategoryMoreFootRectangle
	// SealCategoryDoubleCir 双圆章
	SealCategoryDoubleCir = dto.SealCategoryDoubleCir
	// SealCategoryPublicSecurity 公检法司
	SealCategoryPublicSecurity = dto.SealCategoryPublicSecurity
	// SealCategoryRectangle 普通长方形章
	SealCategoryRectangle = dto.SealCategoryRectangle
	// SealCategorySquare 正方形章
	SealCategorySquare = dto.SealCategorySquare
	// SealCategoryDoublelineOval 双行文字椭圆
	SealCategoryDoublelineOval = dto.SealCategoryDoublelineOval
	// SealCategoryDoubleOval 双椭圆章
	SealCategoryDoubleOval = dto.SealCategoryDoubleOval
	// SealCategoryArchives 档案章
	SealCategoryArchives = dto.SealCategoryArchives
	// SealCategoryRectangleLand 规土局方章
	SealCategoryRectangleLand = dto.SealCategoryRectangleLand
	// SealCategoryDoublelineCir 双行文字圆章
	SealCategoryDoublelineCir = dto.SealCategoryDoublelineCir
	// SealCategoryMoreFootOval 多横排文字椭圆
	SealCategoryMoreFootOval = dto.SealCategoryMoreFootOval
	// SealCategoryPurposeSeal 用途章
	SealCategoryPurposeSeal = dto.SealCategoryPurposeSeal
)

// NewCompanySealCreateService 基于给定的电子印章创建服务构造企业印章创建服务。
func NewCompanySealCreateService() *CompanySealCreateService {
	return service.NewCompanySealCreateService()
}

// NewCustomSealAttribute 创建带默认开关与 DPI 的 CustomSealAttribute 实例。
func NewCustomSealAttribute() *CustomSealAttribute {
	return dto.NewCustomSealAttribute()
}
