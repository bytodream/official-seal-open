// Package draw 章面绘制聚合入口。
// 章型权限与使用次数由各章型 Draw 负责，所有路径均需通过许可证校验。
package draw

import (
	"fmt"
	"image"

	"github.com/bytodream/official-seal-open/seal/internal/draw/module/communist"
	"github.com/bytodream/official-seal-open/seal/internal/draw/module/doublecir"
	"github.com/bytodream/official-seal-open/seal/internal/draw/module/doublelinecircle"
	"github.com/bytodream/official-seal-open/seal/internal/draw/module/doublelineoval"
	"github.com/bytodream/official-seal-open/seal/internal/draw/module/enterprise"
	"github.com/bytodream/official-seal-open/seal/internal/draw/module/labour"
	"github.com/bytodream/official-seal-open/seal/internal/draw/module/morefootcircle"
	"github.com/bytodream/official-seal-open/seal/internal/draw/module/morefootoval"
	"github.com/bytodream/official-seal-open/seal/internal/draw/module/morefootrectangle"
	"github.com/bytodream/official-seal-open/seal/internal/draw/module/oval"
	"github.com/bytodream/official-seal-open/seal/internal/draw/module/youthleague"
	"github.com/bytodream/official-seal-open/seal/utils"

	"github.com/bytodream/official-seal-open/seal/internal/draw/base"
	"github.com/bytodream/official-seal-open/seal/model"
	"github.com/bytodream/official-seal-open/seal/utils/ioc"
)

func init() {
	communist.Register()
	doublecir.Register()
	doublelinecircle.Register()
	doublelineoval.Register()
	enterprise.Register()
	labour.Register()
	morefootcircle.Register()
	morefootoval.Register()
	morefootrectangle.Register()
	oval.Register()
	youthleague.Register()
}

// SealPaintModel 章面绘制参数（像素级，由 base.SealPaintAttribute 转换而来）。
// 兼容别名：实际定义在公共层 internal/core。
type SealPaintModel = base.SealPaintModel

// PaintOption 绘制模型可选配置（公共层实现）。
// 兼容别名：实际定义在公共层 internal/core。
type PaintOption = base.PaintOption

// WithSampleWatermark 标记预览样章：Draw 收尾时在章面中心叠加「样章」标识（text 为空用默认文字）。
func WithSampleWatermark(text string) PaintOption {
	return base.WithSampleWatermark(text)
}

// BuildSealPaintModel 从 base.SealPaintAttribute 构建（公共层实现）。
func BuildSealPaintModel(attr *base.SealPaintAttribute, transparencyAlpha int, fontFile string, opts ...PaintOption) SealPaintModel {
	return base.BuildSealPaintModel(attr, transparencyAlpha, fontFile, opts...)
}

// BuildSealPaintModelLegacy 兼容旧路径：直接从 model.CustomSealAttribute 构建绘制模型。
// 内部经 BuildSealAttributeBundle 按章型分发构建后转换为绘制模型。
func BuildSealPaintModelLegacy(a *model.CustomSealAttribute, transparencyAlpha int, fontFile string, opts ...PaintOption) (SealPaintModel, error) {
	bundle, err := BuildSealAttributeBundle(a)
	if err != nil {
		return SealPaintModel{}, err
	}
	return base.BuildSealPaintModel(bundle.Paint, transparencyAlpha, fontFile, opts...), nil
}

// ResolveFont 返回字体文件完整路径（公共层实现）。
func ResolveFont(fontPath string) string {
	return base.ResolveFont(fontPath)
}

// Translucent 将章面色带写入指定 Alpha，白色背景像素丢弃（公共层实现）。
func Translucent(src image.Image, transparency int) *image.RGBA {
	return base.Translucent(src, transparency)
}

// TransparencyFromSealString 将 0–100 透明度字符串换算为 Alpha（公共层实现）。
func TransparencyFromSealString(s *string) int {
	return base.TransparencyFromSealString(s)
}

// DrawSeal 按章型分发到各章型模块；未单独实现的章型回退企业圆章绘制。
// 分发/权限校验失败时返回包装 ErrSealCategoryNotLicensed 的错误。
func DrawSeal(m SealPaintModel) (image.Image, error) {
	cmd, err := drawService(m.Category)
	if err != nil {
		return nil, err
	}
	return cmd.Draw(m)
}

// registeredCategories 已单独实现绘制服务的章型；其余章型（enterprise、rectangle、square 等）
// 统一回退企业圆章绘制。
var registeredCategories = []model.SealCategory{
	model.SealCategoryLabour,
	model.SealCategoryCommunist,
	model.SealCategoryYouthLeague,
	model.SealCategoryMoreFootCircle,
	model.SealCategoryDoubleCir,
	model.SealCategoryMoreFootRectangle,
	model.SealCategoryMoreFootOval,
	model.SealCategoryDoublelineCir,
	model.SealCategoryOval,
	model.SealCategoryDoublelineOval,
}

// drawService 按章型名从 ioc 容器获取绘制服务；未单独实现的章型回退企业圆章服务。
// 权限由各章型 Draw 按 m.Category 校验，回退时仍校验调用方真实章型。
func drawService(category model.SealCategory) (base.DrawService, error) {
	name := category.ToString()
	if category.Contains(registeredCategories) {
		comp, err := ioc.DefaultContainer.GetComponent(name)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", utils.ErrSealCategoryNotLicensed, name)
		}
		svc, ok := comp.Value.(base.DrawService)
		if !ok {
			return nil, fmt.Errorf("%w: %s", utils.ErrSealCategoryNotLicensed, name)
		}
		return svc, nil
	}
	// 其余章型回退企业圆章绘制（矩形、方形等未单独实现的章型）。
	comp, err := ioc.DefaultContainer.GetComponent(model.SealCategoryEnterprise.ToString())
	if err != nil {
		return nil, fmt.Errorf("%w: %s", utils.ErrSealCategoryNotLicensed, name)
	}
	svc, ok := comp.Value.(base.DrawService)
	if !ok {
		return nil, fmt.Errorf("%w: %s", utils.ErrSealCategoryNotLicensed, name)
	}
	return svc, nil
}
