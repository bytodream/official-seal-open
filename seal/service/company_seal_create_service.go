// Package service 电子公章服务实现；对外 API 实体见 service/dto，集成方可 import 顶层 seal 门面。
package service

import (
	"image"
	"reflect"

	"github.com/bytodream/official-seal-open/seal/internal/draw"
	"github.com/bytodream/official-seal-open/seal/internal/render"
	"github.com/bytodream/official-seal-open/seal/model"
	"github.com/bytodream/official-seal-open/seal/service/dto"
)

// CompanySealCreateService 企业电子公章的生成与预览服务。
type CompanySealCreateService struct {
	*electronicAbstractSealCreateService
}

// NewCompanySealCreateService 基于公共依赖创建企业电子公章服务。
func NewCompanySealCreateService() *CompanySealCreateService {
	return &CompanySealCreateService{
		electronicAbstractSealCreateService: newElectronicAbstractSealCreateService(),
	}
}

// SetFontPath 设置章面绘制使用的字体文件路径。
func (e *CompanySealCreateService) SetFontPath(fontPath string) *CompanySealCreateService {
	e.FontPath = fontPath
	return e
}

// SetImageService 注入章面图片服务；传 nil 表示不使用图片服务。
func (e *CompanySealCreateService) SetImageService(imageService ImageService) *CompanySealCreateService {
	e.ImageService = imageService
	return e
}

// configurableImageService 返回实现了 ImageServiceConfigurator 的图片服务；
// 未注入、未实现该接口或持有 nil 指针时返回 nil，调用方的配置会被忽略。
func (e *CompanySealCreateService) configurableImageService() ImageServiceConfigurator {
	c, ok := e.ImageService.(ImageServiceConfigurator)
	if !ok || isNilInterface(c) {
		return nil
	}
	return c
}

// hasImageService 判断是否注入了可用的图片服务（含类型化 nil 的防护）。
func (e *CompanySealCreateService) hasImageService() bool {
	return !isNilInterface(e.ImageService)
}

// isNilInterface 判断接口是否未赋值或持有 nil 指针（typed nil）。
// reflect.ValueOf(nil) 是零值，直接调用 IsNil 会 panic，因此先判 v == nil。
func isNilInterface(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan, reflect.Interface:
		return rv.IsNil()
	default:
		return false
	}
}

// SetUploadBytes 注入上传实现；图片服务未实现 ImageServiceConfigurator 时忽略。
func (e *CompanySealCreateService) SetUploadBytes(uploadBytes UploadBytes) *CompanySealCreateService {
	if c := e.configurableImageService(); c != nil {
		c.SetUploadBytes(uploadBytes)
	}
	return e
}

// HasUploadBytes 报告上传能力是否就绪；GetAutoSealImage 依赖该结果在渲染前短路。
func (e *CompanySealCreateService) HasUploadBytes() bool {
	if c := e.configurableImageService(); c != nil {
		return c.HasUploadBytes()
	}
	return false
}

// SetBlurGenerate 开关模糊章面生成；图片服务未实现 ImageServiceConfigurator 时忽略。
func (e *CompanySealCreateService) SetBlurGenerate(blurGenerate bool) *CompanySealCreateService {
	if c := e.configurableImageService(); c != nil {
		c.SetBlurGenerate(blurGenerate)
	}
	return e
}

// IsBlurGenerate 报告是否生成模糊章面。
func (e *CompanySealCreateService) IsBlurGenerate() bool {
	if c := e.configurableImageService(); c != nil {
		return c.IsBlurGenerate()
	}
	return false
}

// SetWatermarkFunc 设置章面水印回调；当前生成/预览链路未使用该扩展点。
func (e *CompanySealCreateService) SetWatermarkFunc(watermark WatermarkFunc) *CompanySealCreateService {
	e.Watermark = watermark
	return e
}

// GetAutoSealImageBySeal 生成清晰/模糊章面并经 ImageService 上传，返回各 fileKey 与主图大小。
func (s *CompanySealCreateService) GetAutoSealImageBySeal(seal *dto.Seal) (*dto.SealImage, error) {
	return s.GetAutoSealImage(seal, dto.NewCustomSealAttribute())
}

// GetAutoSealImageByCustomSealAttribute 生成清晰/模糊章面并经 ImageService 上传，返回各 fileKey 与主图大小。
func (s *CompanySealCreateService) GetAutoSealImageByCustomSealAttribute(attribute *dto.CustomSealAttribute) (*dto.SealImage, error) {
	radius := 10
	seal := &dto.Seal{Radius: &radius}
	return s.GetAutoSealImage(seal, attribute)
}

// GetAutoSealImageByDefault 生成清晰/模糊章面并经 ImageService 上传，返回各 fileKey 与主图大小。
func (s *CompanySealCreateService) GetAutoSealImageByDefault() (*dto.SealImage, error) {
	radius := 10
	seal := &dto.Seal{Radius: &radius}
	return s.GetAutoSealImage(seal, dto.NewCustomSealAttribute())
}

// GetAutoSealImage 生成清晰/模糊章面并经 ImageService 上传，返回各 fileKey 与主图大小。
func (s *CompanySealCreateService) GetAutoSealImage(seal *dto.Seal, attribute *dto.CustomSealAttribute) (*dto.SealImage, error) {
	if seal == nil || attribute == nil {
		return nil, dto.ErrSealAttributeNotConfigured
	}
	if !s.hasImageService() {
		return nil, dto.ErrImageServiceNotConfigured
	}
	if !s.HasUploadBytes() {
		return nil, dto.ErrImageUploadNotConfigured
	}
	imageBytes, blurImageBytes, err := s.GetSealImage(seal, attribute)
	if err != nil {
		return nil, err
	}
	var (
		imageKey     string
		blurImageKey string
	)
	attr := dto.ImageAttrDefault
	if len(imageBytes) > 0 {
		imageKey, err = s.ImageService.UploadBytes(imageBytes, attr, true)
		if err != nil {
			return nil, err
		}
	}

	if len(blurImageBytes) > 0 {
		blurImageKey, err = s.ImageService.UploadBytes(blurImageBytes, attr, true)
		if err != nil {
			return nil, err
		}
	}

	out, err := dto.FromDomainSealImage(&model.SealImage{
		ImageKey:     imageKey,
		BlurImageKey: blurImageKey,
		Size:         int64(len(imageBytes)),
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetSealImage 生成清晰章面 PNG；IsBlurGenerate 开启时同时返回模糊章面 PNG。两者都不上传。
func (s *CompanySealCreateService) GetSealImage(seal *dto.Seal, attribute *dto.CustomSealAttribute) ([]byte, []byte, error) {
	if seal == nil || attribute == nil {
		return nil, nil, dto.ErrSealAttributeNotConfigured
	}
	if !s.hasImageService() {
		return nil, nil, dto.ErrImageServiceNotConfigured
	}
	domainSeal, err := seal.ToDomain()
	if err != nil {
		return nil, nil, err
	}
	domainAttr, err := domainAttributeFromService(attribute)
	if err != nil {
		return nil, nil, err
	}
	img, err := s.generateImage(domainSeal, domainAttr)
	if err != nil {
		return nil, nil, domainErr(err)
	}
	attr := dto.ImageAttrDefault
	imageBytes, err := s.ImageService.ToImageBytes(img, attr)
	if err != nil {
		return nil, nil, err
	}
	var blurImageBytes []byte
	if s.IsBlurGenerate() {
		blurImage := render.DefaultThickenImage(img, domainSeal.Radius)
		blurImageBytes, err = s.ImageService.ToImageBytes(blurImage, attr)
		if err != nil {
			return nil, nil, err
		}
	}
	return imageBytes, blurImageBytes, nil
}

// Preview 按请求模型生成预览 PNG；固定叠加「样章」标识且不消耗使用额度，
// OpenImageBlur 开启时输出模糊章面。
func (s *CompanySealCreateService) Preview(sm *dto.SealModel) ([]byte, error) {
	if sm == nil {
		return nil, dto.ErrSealAttributeNotConfigured
	}
	em := sm.ElectronicModel
	watermarkText := ""
	if em != nil && em.WatermarkText != "" {
		watermarkText = em.WatermarkText
	}
	seal := sm.ToSeal()
	if em != nil && em.Transparency != nil {
		seal.Transparency = em.Transparency
	}
	domainSeal, err := seal.ToDomain()
	if err != nil {
		return nil, err
	}
	var domainAttr *model.CustomSealAttribute
	if em != nil && em.AutoAttribute != nil {
		// 预览同样做许可证校验；是否计入次数由 Draw 层按样章标识判断。
		domainAttr, err = domainAttributeFromService(em.AutoAttribute)
		if err != nil {
			return nil, err
		}
	}
	// 预览固定叠加「样章」标识，且不记录、不消耗使用额度。
	img, err := s.generatePreviewImage(domainSeal, domainAttr, true, watermarkText)
	if err != nil {
		return nil, domainErr(err)
	}
	if em != nil && em.OpenImageBlur {
		img = render.DefaultThickenImage(img, em.Radius)
	}
	return render.PNGBytes(img)
}

// RenderAttributePNG 按自定义章面参数与给定透明度直接渲染 PNG 字节（不上传）。
func (s *CompanySealCreateService) RenderAttributePNG(attr *dto.CustomSealAttribute, transparency int) ([]byte, error) {
	domainAttr, err := domainAttributeFromService(attr)
	if err != nil {
		return nil, err
	}
	data, err := render.PNGFromCustomAttribute(domainAttr, s.FontPath, transparency)
	if err != nil {
		return nil, domainErr(err)
	}
	return data, nil
}

func (s *CompanySealCreateService) generateImage(seal *model.Seal, attribute *model.CustomSealAttribute) (image.Image, error) {
	if attribute == nil {
		return nil, model.ErrSealAttributeNotConfigured
	}
	return s.commonSealGenerator().GenerateImage(seal, attribute)
}

// generatePreviewImage 生成预览样章：showWatermark 为 true 时由 Draw 在章面中心叠加「样章」标识，
// watermarkText 为空则使用默认文字。
func (s *CompanySealCreateService) generatePreviewImage(seal *model.Seal, attribute *model.CustomSealAttribute, showWatermark bool, watermarkText string) (image.Image, error) {
	if attribute == nil {
		return nil, model.ErrSealAttributeNotConfigured
	}
	if !showWatermark {
		return s.commonSealGenerator().GenerateImage(seal, attribute)
	}
	return s.commonSealGenerator().GenerateImage(seal, attribute, draw.WithSampleWatermark(watermarkText))
}

// InitSealPermission 创建印章前的许可校验（仅检查许可证有效性）。
func (s *CompanySealCreateService) InitSealPermission(_ *dto.Seal) error {
	return EnsureLicenseActive()
}

// IsSupport 判断是否支持该印章类型与业务类型（仅企业电子章）。
func (s *CompanySealCreateService) IsSupport(typ dto.SealType, category dto.BusinessType) bool {
	return typ == dto.SealTypeCompany && category == dto.BusinessTypeElectronic
}

func (s *CompanySealCreateService) commonSealGenerator() *render.CommonSealGenerator {
	return &render.CommonSealGenerator{FontPath: s.FontPath}
}
