// service/dto 与 seal/model 之间显式字段拷贝；新增字段须两侧同步维护。
package dto

import "github.com/bytodream/official-seal-open/seal/model"

func clonePtr[T any](p *T) *T {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}

func sealCategoryToDomain(c *SealCategory) *model.SealCategory {
	if c == nil {
		return nil
	}
	cat := model.SealCategory(*c)
	return &cat
}

func sealCategoryFromDomain(c *model.SealCategory) *SealCategory {
	if c == nil {
		return nil
	}
	cat := SealCategory(*c)
	return &cat
}

func diySpecToDomain(d DiySpec) model.DiySpec {
	return model.DiySpec{
		Distance:            clonePtr(d.Distance),
		FontSize:            clonePtr(d.FontSize),
		Width:               clonePtr(d.Width),
		Text:                d.Text,
		Code:                d.Code,
		DefaultVal:          d.DefaultVal,
		Step:                d.Step,
		IsFontMultilineText: clonePtr(d.IsFontMultilineText),
	}
}

func diySpecFromDomain(d model.DiySpec) DiySpec {
	return DiySpec{
		Distance:            clonePtr(d.Distance),
		FontSize:            clonePtr(d.FontSize),
		Width:               clonePtr(d.Width),
		Text:                d.Text,
		Code:                d.Code,
		DefaultVal:          d.DefaultVal,
		Step:                d.Step,
		IsFontMultilineText: clonePtr(d.IsFontMultilineText),
	}
}

func sealFootToDomain(f SealFoot) model.SealFoot {
	out := model.SealFoot{
		Text:     f.Text,
		Width:    clonePtr(f.Width),
		FontSize: clonePtr(f.FontSize),
		Distance: clonePtr(f.Distance),
	}
	if len(f.DiySpec) > 0 {
		out.DiySpec = make([]model.DiySpec, len(f.DiySpec))
		for i := range f.DiySpec {
			out.DiySpec[i] = diySpecToDomain(f.DiySpec[i])
		}
	}
	return out
}

func sealFootFromDomain(f model.SealFoot) SealFoot {
	out := SealFoot{
		Text:     f.Text,
		Width:    clonePtr(f.Width),
		FontSize: clonePtr(f.FontSize),
		Distance: clonePtr(f.Distance),
	}
	if len(f.DiySpec) > 0 {
		out.DiySpec = make([]DiySpec, len(f.DiySpec))
		for i := range f.DiySpec {
			out.DiySpec[i] = diySpecFromDomain(f.DiySpec[i])
		}
	}
	return out
}

func customSealAttributeToDomain(a *CustomSealAttribute) *model.CustomSealAttribute {
	if a == nil {
		return nil
	}
	out := &model.CustomSealAttribute{
		SealCategory:                  sealCategoryToDomain(a.SealCategory),
		SavedSpecName:                 a.SavedSpecName,
		Content:                       a.Content,
		InnerContent:                  a.InnerContent,
		Foot:                          a.Foot,
		Head:                          a.Head,
		EnterpriseCode:                a.EnterpriseCode,
		EdgeWidth:                     clonePtr(a.EdgeWidth),
		InnerEdgeWidth:                clonePtr(a.InnerEdgeWidth),
		InnerMargin:                   clonePtr(a.InnerMargin),
		LogoSize:                      clonePtr(a.LogoSize),
		SealLogo:                      clonePtr(a.SealLogo),
		Antifake:                      clonePtr(a.Antifake),
		AntifakeCode:                  a.AntifakeCode,
		AntifakeGap:                   clonePtr(a.AntifakeGap),
		SealSpec:                      a.SealSpec,
		Dpi:                           clonePtr(a.Dpi),
		SealCategoryType:              a.SealCategoryType,
		LogoUpmove:                    clonePtr(a.LogoUpmove),
		Arc:                           clonePtr(a.Arc),
		SealHeight:                    clonePtr(a.SealHeight),
		SealWidth:                     clonePtr(a.SealWidth),
		FontSize:                      clonePtr(a.FontSize),
		WordLinePadding:               clonePtr(a.WordLinePadding),
		EnterpriseCodeFontSize:        clonePtr(a.EnterpriseCodeFontSize),
		EnterpriseCodeArc:             clonePtr(a.EnterpriseCodeArc),
		EnterpriseCodeWordLinePadding: clonePtr(a.EnterpriseCodeWordLinePadding),
		HeadParam:                     clonePtrSealFootToDomain(a.HeadParam),
		FootFontSize:                  clonePtr(a.FootFontSize),
		FootWidth:                     clonePtr(a.FootWidth),
		FootDistance:                  clonePtr(a.FootDistance),
		HeadFontSize:                  clonePtr(a.HeadFontSize),
		HeadWidth:                     clonePtr(a.HeadWidth),
		LineLength:                    clonePtr(a.LineLength),
		LineWidth:                     clonePtr(a.LineWidth),
		LineInterval:                  clonePtr(a.LineInterval),
		HeadDistance:                  clonePtr(a.HeadDistance),
		InnerOvalWidth:                clonePtr(a.InnerOvalWidth),
		InnerOvalHeight:               clonePtr(a.InnerOvalHeight),
		InnerOvalEdgeWidth:            clonePtr(a.InnerOvalEdgeWidth),
		InnerOvalFontSize:             clonePtr(a.InnerOvalFontSize),
		InnerOvalWordLinePadding:      clonePtr(a.InnerOvalWordLinePadding),
	}
	if len(a.Foots) > 0 {
		out.Foots = make([]model.SealFoot, len(a.Foots))
		for i := range a.Foots {
			out.Foots[i] = sealFootToDomain(a.Foots[i])
		}
	}
	return out
}

func clonePtrSealFootToDomain(f *SealFoot) *model.SealFoot {
	if f == nil {
		return nil
	}
	sf := sealFootToDomain(*f)
	return &sf
}

func customSealAttributeFromDomain(a *model.CustomSealAttribute) *CustomSealAttribute {
	if a == nil {
		return nil
	}
	out := &CustomSealAttribute{
		SealCategory:                  sealCategoryFromDomain(a.SealCategory),
		SavedSpecName:                 a.SavedSpecName,
		Content:                       a.Content,
		InnerContent:                  a.InnerContent,
		Foot:                          a.Foot,
		Head:                          a.Head,
		EnterpriseCode:                a.EnterpriseCode,
		EdgeWidth:                     clonePtr(a.EdgeWidth),
		InnerEdgeWidth:                clonePtr(a.InnerEdgeWidth),
		InnerMargin:                   clonePtr(a.InnerMargin),
		LogoSize:                      clonePtr(a.LogoSize),
		SealLogo:                      clonePtr(a.SealLogo),
		Antifake:                      clonePtr(a.Antifake),
		AntifakeCode:                  a.AntifakeCode,
		AntifakeGap:                   clonePtr(a.AntifakeGap),
		SealSpec:                      a.SealSpec,
		Dpi:                           clonePtr(a.Dpi),
		SealCategoryType:              a.SealCategoryType,
		LogoUpmove:                    clonePtr(a.LogoUpmove),
		Arc:                           clonePtr(a.Arc),
		SealHeight:                    clonePtr(a.SealHeight),
		SealWidth:                     clonePtr(a.SealWidth),
		FontSize:                      clonePtr(a.FontSize),
		WordLinePadding:               clonePtr(a.WordLinePadding),
		EnterpriseCodeFontSize:        clonePtr(a.EnterpriseCodeFontSize),
		EnterpriseCodeArc:             clonePtr(a.EnterpriseCodeArc),
		EnterpriseCodeWordLinePadding: clonePtr(a.EnterpriseCodeWordLinePadding),
		HeadParam:                     clonePtrSealFootFromDomain(a.HeadParam),
		FootFontSize:                  clonePtr(a.FootFontSize),
		FootWidth:                     clonePtr(a.FootWidth),
		FootDistance:                  clonePtr(a.FootDistance),
		HeadFontSize:                  clonePtr(a.HeadFontSize),
		HeadWidth:                     clonePtr(a.HeadWidth),
		LineLength:                    clonePtr(a.LineLength),
		LineWidth:                     clonePtr(a.LineWidth),
		LineInterval:                  clonePtr(a.LineInterval),
		HeadDistance:                  clonePtr(a.HeadDistance),
		InnerOvalWidth:                clonePtr(a.InnerOvalWidth),
		InnerOvalHeight:               clonePtr(a.InnerOvalHeight),
		InnerOvalEdgeWidth:            clonePtr(a.InnerOvalEdgeWidth),
		InnerOvalFontSize:             clonePtr(a.InnerOvalFontSize),
		InnerOvalWordLinePadding:      clonePtr(a.InnerOvalWordLinePadding),
	}
	if len(a.Foots) > 0 {
		out.Foots = make([]SealFoot, len(a.Foots))
		for i := range a.Foots {
			out.Foots[i] = sealFootFromDomain(a.Foots[i])
		}
	}
	return out
}

func clonePtrSealFootFromDomain(f *model.SealFoot) *SealFoot {
	if f == nil {
		return nil
	}
	sf := sealFootFromDomain(*f)
	return &sf
}

func electronicModelToDomain(em *ElectronicModel) *model.ElectronicModel {
	if em == nil {
		return nil
	}
	return &model.ElectronicModel{
		Width:            clonePtr(em.Width),
		Height:           clonePtr(em.Height),
		OpenImageBlur:    em.OpenImageBlur,
		IsAuto:           clonePtr(em.IsAuto),
		SealSpec:         em.SealSpec,
		SealImageKey:     em.SealImageKey,
		SealBlurImageKey: em.SealBlurImageKey,
		UploadImage:      em.UploadImage,
		UploadImageKey:   em.UploadImageKey,
		AutoAttribute:    customSealAttributeToDomain(em.AutoAttribute),
		IsCreateByAuth:   em.IsCreateByAuth,
		ShowWatermark:    em.ShowWatermark,
		WatermarkText:    em.WatermarkText,
		Radius:           clonePtr(em.Radius),
		Transparency:     clonePtr(em.Transparency),
		SpecEntityId:     em.SpecEntityId,
	}
}

func electronicModelFromDomain(em *model.ElectronicModel) *ElectronicModel {
	if em == nil {
		return nil
	}
	return &ElectronicModel{
		Width:            clonePtr(em.Width),
		Height:           clonePtr(em.Height),
		OpenImageBlur:    em.OpenImageBlur,
		IsAuto:           clonePtr(em.IsAuto),
		SealSpec:         em.SealSpec,
		SealImageKey:     em.SealImageKey,
		SealBlurImageKey: em.SealBlurImageKey,
		UploadImage:      em.UploadImage,
		UploadImageKey:   em.UploadImageKey,
		AutoAttribute:    customSealAttributeFromDomain(em.AutoAttribute),
		IsCreateByAuth:   em.IsCreateByAuth,
		ShowWatermark:    em.ShowWatermark,
		WatermarkText:    em.WatermarkText,
		Radius:           clonePtr(em.Radius),
		Transparency:     clonePtr(em.Transparency),
		SpecEntityId:     em.SpecEntityId,
	}
}

// ToDomain 转为内部 seal/model。
func (s *Seal) ToDomain() (*model.Seal, error) {
	if s == nil {
		return nil, nil
	}
	return &model.Seal{
		Radius:       clonePtr(s.Radius),
		Transparency: clonePtr(s.Transparency),
	}, nil
}

// ToDomain 转为内部 seal/model。
func (a *CustomSealAttribute) ToDomain() (*model.CustomSealAttribute, error) {
	return customSealAttributeToDomain(a), nil
}

// ToDomain 转为内部 seal/model。
func (m *SealModel) ToDomain() (*model.SealModel, error) {
	if m == nil {
		return nil, nil
	}
	return &model.SealModel{
		Category:        model.BusinessType(m.Category),
		ElectronicModel: electronicModelToDomain(m.ElectronicModel),
	}, nil
}

// ToDomain 转为内部 seal/model。
func (a *ImageAttr) ToDomain() (model.ImageAttr, error) {
	if a == nil {
		return model.ImageAttr{}, nil
	}
	return model.ImageAttr{
		DefaultAttr:  a.DefaultAttr,
		Transparent:  a.Transparent,
		Zoom:         a.Zoom,
		Rotate:       a.Rotate,
		LimitWidth:   a.LimitWidth,
		LimitHeight:  a.LimitHeight,
		Degree:       a.Degree,
		Transparency: a.Transparency,
	}, nil
}

// FromDomainSealImage 从内部 model 拷贝服务视图。
func FromDomainSealImage(img *model.SealImage) (*SealImage, error) {
	if img == nil {
		return nil, nil
	}
	return &SealImage{
		OriginImageKey: img.OriginImageKey,
		ImageKey:       img.ImageKey,
		BlurImageKey:   img.BlurImageKey,
		Size:           img.Size,
	}, nil
}

// FromDomainSealModel 从内部 SealModel 拷贝服务视图。
func FromDomainSealModel(m *model.SealModel) (*SealModel, error) {
	if m == nil {
		return nil, nil
	}
	return &SealModel{
		Category:        BusinessType(m.Category),
		ElectronicModel: electronicModelFromDomain(m.ElectronicModel),
	}, nil
}

// ToDomainRaster 章面位图（image.Image，无需拷贝）。
func ToDomainRaster(img RasterImage) model.RasterImage {
	return img
}

// FromDomainRaster 从内部位图转回服务视图类型。
func FromDomainRaster(img model.RasterImage) RasterImage {
	return img
}
