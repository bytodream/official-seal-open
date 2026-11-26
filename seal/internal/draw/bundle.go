package draw

import (
	"github.com/bytodream/official-seal-open/seal/internal/draw/base"
	"github.com/bytodream/official-seal-open/seal/internal/draw/module/enterprise"
	"github.com/bytodream/official-seal-open/seal/model"
	"github.com/bytodream/official-seal-open/seal/utils"
)

// BuildSealAttributeBundle 按章型规格构建完整像素级章面属性（集中分发 + 业务覆盖）。
// 流程与旧 spec.ToSealAttributeFull 一致：按章型分发到各模块 GetCircle 构建底稿，
// 成功后恢复多横排上排距离，再应用业务覆盖字段（内容、开关、防伪、logo 等）。
func BuildSealAttributeBundle(a *model.CustomSealAttribute) (*base.SealAttributeBundle, error) {
	if a == nil {
		return nil, model.ErrSealAttributeNotConfigured
	}
	if a.SealCategory != nil && *a.SealCategory == model.SealCategoryMoreFootCircle {
		// 多横排圆章构建前按字号校正上排距离，无论成功失败都还原，避免污染调用方入参。
		counteractMoreFoot(a)
		defer recoverHeadDistance(a)
	}
	bundle, err := buildSealAttributeBundle(a)
	if err != nil {
		return nil, err
	}
	if bundle.Paint != nil {
		applySealAttributeOverrides(a, bundle.Paint)
	}
	return bundle, nil
}

// buildSealAttributeBundle 按章型分发到各模块的 GetCircle 构建底稿。
// 已单独实现的章型经 ioc 取对应模块；其余章型（含 OVAL、enterprise、rectangle、square 等）
// 与旧逻辑一致回退企业构建器（其内部再按规格分发椭圆/双线椭圆）。
func buildSealAttributeBundle(a *model.CustomSealAttribute) (*base.SealAttributeBundle, error) {
	cat := model.SealCategoryEnterprise
	if a.SealCategory != nil {
		cat = *a.SealCategory
	}
	switch cat {
	case model.SealCategoryYouthLeague,
		model.SealCategoryCommunist,
		model.SealCategoryDoublelineOval,
		model.SealCategoryLabour,
		model.SealCategoryDoublelineCir,
		model.SealCategoryMoreFootCircle,
		model.SealCategoryMoreFootOval,
		model.SealCategoryMoreFootRectangle,
		model.SealCategoryDoubleCir:
		svc, err := drawService(cat)
		if err != nil {
			return nil, err
		}
		return svc.GetCircle(a)
	default:
		circle, err := enterprise.GetEnterprise(a.SealSpecEnum(), a.Foot, a.Head, a.EnterpriseCode, a.GetDpi(), a.InnerContent, a.ToDiySealAttribute(), a.IsSealLogo())
		if err != nil {
			return nil, err
		}
		return &base.SealAttributeBundle{Paint: circle}, nil
	}
}

// applySealAttributeOverrides 将业务覆盖字段写回章面属性。
func applySealAttributeOverrides(a *model.CustomSealAttribute, attr *base.SealPaintAttribute) {
	switch {
	case attr.Category == model.SealCategoryDoublelineOval && attr.InnerLine != nil:
		if !utils.IsBlank(a.InnerContent) {
			attr.Content = a.InnerContent
		}
	case attr.Category == model.SealCategoryDoublelineOval && attr.InnerOval != nil:
		if !utils.IsBlank(a.InnerContent) {
			attr.InnerOval.Content = a.InnerContent
		}
		if !utils.IsBlank(a.Content) {
			attr.Content = a.Content
		}
	case attr.Category == model.SealCategoryDoublelineOval && !utils.IsBlank(a.InnerContent):
		// 无 inner 占位时 fallback
		attr.Content = a.InnerContent
	default:
		attr.Content = a.Content
	}
	if attr.InnerCircle != nil && !utils.IsBlank(a.InnerContent) {
		attr.InnerCircle.Content = a.InnerContent
	}
	if attr.InnerOval != nil && !utils.IsBlank(a.InnerContent) {
		attr.InnerOval.Content = a.InnerContent
	}
	attr.FakeGap = a.IsAntifakeGap()
	switch {
	case model.IsMoreFootCategory(attr.Category):
		attr.SealLogo = false
	case attr.Category == model.SealCategoryDoublelineOval && (attr.InnerOval != nil || attr.InnerLine != nil):
		attr.SealLogo = false
	default:
		attr.SealLogo = a.IsSealLogo()
	}
	if a.EdgeWidth != nil {
		attr.EdgeWidthPx = utils.ToFloatPixel(a.GetDpi(), *a.EdgeWidth)
	}
	if a.IsSealLogo() && a.IsAntifake() &&
		attr.Category != model.SealCategoryYouthLeague &&
		attr.Category != model.SealCategoryCommunist &&
		!model.IsMoreFootCategory(attr.Category) {
		// 掩膜文字走默认规则（普通章 "OS"）；共青团/党委章不在此设置防伪参数。
		attr.Antifake = a.BuildAntifakeInfo("")
	}
	if attr.Logo != nil && a.LogoSize != nil {
		attr.Logo.SizePx = utils.ToIntPixel(a.GetDpi(), *a.LogoSize)
	}
	if a.LogoUpmove != nil && attr.Logo != nil {
		attr.Logo.UpMovePx = utils.ToFloatPixel(a.GetDpi(), *a.LogoUpmove)
	}
	if a.InnerEdgeWidth != nil {
		attr.InnerEdgeWidthPx = utils.ToFloatPixel(a.GetDpi(), *a.InnerEdgeWidth)
	}
	if a.InnerMargin != nil {
		attr.InnerMarginPx = utils.ToFloatPixel(a.GetDpi(), *a.InnerMargin)
	}
}

// counteractMoreFoot 校正多横排 headDistance 与字号联动。
func counteractMoreFoot(a *model.CustomSealAttribute) {
	if a.SealCategory != nil && *a.SealCategory == model.SealCategoryMoreFootCircle &&
		a.HeadDistance != nil && a.FontSize != nil {
		v := float32(float64(*a.HeadDistance) + float64(*a.FontSize)*0.75)
		a.HeadDistance = &v
	}
}

// recoverHeadDistance 恢复上排距离默认值。
func recoverHeadDistance(a *model.CustomSealAttribute) {
	if a.SealCategory != nil && *a.SealCategory == model.SealCategoryMoreFootCircle &&
		a.HeadDistance != nil && a.FontSize != nil {
		v := float32(float64(*a.HeadDistance) - float64(*a.FontSize)*0.75)
		a.HeadDistance = &v
	}
}
