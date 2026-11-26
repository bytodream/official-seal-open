package render

import (
	"fmt"
	"os"

	"github.com/bytodream/official-seal-open/seal/internal/draw"
	"github.com/bytodream/official-seal-open/seal/model"
)

// CommonSealGenerator 纯 Go 企业圆章生成。
type CommonSealGenerator struct {
	// FontPath 字体路径（配置项 fonts.path）
	FontPath string
}

// GenerateImage 由 seal.CompanySealCreateService 调用；章型权限与使用次数由 Draw 层负责。
// opts 透传给绘制模型构建（如 draw.WithSampleWatermark 用于预览样章）。
func (g *CommonSealGenerator) GenerateImage(seal *model.Seal, attribute *model.CustomSealAttribute, opts ...draw.PaintOption) (model.RasterImage, error) {
	if attribute == nil {
		return nil, model.ErrSealAttributeNotConfigured
	}
	fontFile := draw.ResolveFont(g.FontPath)
	if fontFile != "" {
		if _, err := os.Stat(fontFile); err != nil {
			return nil, fmt.Errorf("未找到字体 %s，请设置 OFFICIAL_SEAL_FONTS_PATH", fontFile)
		}
	} else {
		return nil, fmt.Errorf("CommonSealGenerator: FontPath 未配置")
	}

	trans := model.DefaultSealTransparency
	if seal != nil && seal.Transparency != nil {
		trans = draw.TransparencyFromSealString(seal.Transparency)
	}

	cm, err := draw.BuildSealPaintModelLegacy(attribute, trans, fontFile, opts...)
	if err != nil {
		return nil, err
	}
	return draw.DrawSeal(cm)
}
