package service

import (
	"github.com/bytodream/official-seal-open/seal/internal/draw"
	"github.com/bytodream/official-seal-open/seal/internal/render"
	"github.com/bytodream/official-seal-open/seal/model"
	"github.com/bytodream/official-seal-open/seal/service/dto"
)

// CompanySealCreateService 企业电子公章的生成与预览服务。
type CompanySealCreateService struct {
	ElectronicAbstractSealCreateService
}

// NewCompanySealCreateService 基于公共依赖创建企业电子公章服务。
func NewCompanySealCreateService(base ElectronicAbstractSealCreateService) *CompanySealCreateService {
	return &CompanySealCreateService{ElectronicAbstractSealCreateService: base}
}

// GetAutoSealImage 生成清晰/模糊章面并经 ImageService 上传，返回各 fileKey 与主图大小。
func (s *CompanySealCreateService) GetAutoSealImage(seal *dto.Seal, attribute *dto.CustomSealAttribute) (*dto.SealImage, error) {
	if seal == nil {
		return nil, dto.ErrSealAttributeNotConfigured
	}
	if s.ImageService == nil {
		return nil, dto.ErrImageServiceNotConfigured
	}
	domainSeal, err := seal.ToDomain()
	if err != nil {
		return nil, err
	}
	domainAttr, err := domainAttributeFromService(attribute)
	if err != nil {
		return nil, err
	}
	normalImage, err := s.generateImage(domainSeal, domainAttr)
	if err != nil {
		return nil, domainErr(err)
	}
	attr := dto.ImageAttrDefault
	normalImageBytes, err := s.ImageService.ToImageBytes(dto.FromDomainRaster(normalImage), attr)
	if err != nil {
		return nil, err
	}
	imageKey, err := s.ImageService.UploadBytes(normalImageBytes, attr, true)
	if err != nil {
		return nil, err
	}
	blurImage := render.DefaultThickenImage(normalImage, domainSeal.Radius)
	blurImageKey, err := s.ImageService.UploadRasterImage(dto.FromDomainRaster(blurImage), attr, true)
	if err != nil {
		return nil, err
	}
	out, err := dto.FromDomainSealImage(&model.SealImage{
		ImageKey:     imageKey,
		BlurImageKey: blurImageKey,
		Size:         int64(len(normalImageBytes)),
	})
	if err != nil {
		return nil, err
	}
	return out, nil
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

func (s *CompanySealCreateService) generateImage(seal *model.Seal, attribute *model.CustomSealAttribute) (model.RasterImage, error) {
	if attribute == nil {
		return nil, model.ErrSealAttributeNotConfigured
	}
	return s.commonSealGenerator().GenerateImage(seal, attribute)
}

// generatePreviewImage 生成预览样章：showWatermark 为 true 时由 Draw 在章面中心叠加「样章」标识，
// watermarkText 为空则使用默认文字。
func (s *CompanySealCreateService) generatePreviewImage(seal *model.Seal, attribute *model.CustomSealAttribute, showWatermark bool, watermarkText string) (model.RasterImage, error) {
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

func (s *ElectronicAbstractSealCreateService) commonSealGenerator() *render.CommonSealGenerator {
	return &render.CommonSealGenerator{FontPath: s.FontPath}
}
