package render

import (
	"github.com/bytodream/official-seal-open/seal/internal/draw"
	"github.com/bytodream/official-seal-open/seal/model"
)

// PNGFromCustomAttribute 由 seal.CompanySealCreateService 调用；章型权限与使用次数由 Draw 层负责。
func PNGFromCustomAttribute(a *model.CustomSealAttribute, fontPath string, transparency int) ([]byte, error) {
	return renderCustomSealPNG(a, fontPath, transparency)
}

// renderCustomSealPNG 走 draw.BuildSealPaintModelLegacy → draw.DrawSeal → PNG（包内专用）。
// opts 透传绘制选项（如 draw.WithSampleWatermark 用于预览样章）；权限与计次由 Draw 层负责，不可跳过。
func renderCustomSealPNG(a *model.CustomSealAttribute, fontPath string, transparency int, opts ...draw.PaintOption) ([]byte, error) {
	if a == nil {
		return nil, model.ErrSealAttributeNotConfigured
	}
	fontFile := draw.ResolveFont(fontPath)
	if fontFile == "" {
		return nil, model.ErrSealAttributeNotConfigured
	}
	cm, err := draw.BuildSealPaintModelLegacy(a, transparency, fontFile, opts...)
	if err != nil {
		return nil, err
	}
	img, err := draw.DrawSeal(cm)
	if err != nil {
		return nil, err
	}
	return PNGBytes(img)
}
