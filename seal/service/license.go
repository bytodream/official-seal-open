package service

import (
	"errors"

	"github.com/bytodream/official-seal-open/seal/internal/draw"
	"github.com/bytodream/official-seal-open/seal/internal/draw/base"
	"github.com/bytodream/official-seal-open/seal/model"
	"github.com/bytodream/official-seal-open/seal/service/dto"
	"github.com/bytodream/official-seal-open/seal/utils"
)

var (
	// ErrLicenseRequired 未设置有效授权密钥。
	ErrLicenseRequired = utils.ErrLicenseRequired
	// ErrLicenseExpired 授权已过期。
	ErrLicenseExpired = utils.ErrLicenseExpired
	// ErrSealCategoryNotLicensed 该章型未包含在授权范围内。
	ErrSealCategoryNotLicensed = utils.ErrSealCategoryNotLicensed
	// ErrUsageLimitExceeded 授权使用次数已用尽。
	ErrUsageLimitExceeded = utils.ErrUsageLimitExceeded
)

// License 授权信息（类型别名为 utils.License）。
type License = utils.License

// SetLicenseKey 解析并校验授权密钥，成功后作为当前生效授权。
func SetLicenseKey(content string) error {
	return utils.SetLicenseKey(content)
}

// GetLicenseKey 返回当前生效授权的副本；未设置时返回 nil。
func GetLicenseKey() *License {
	return utils.GetLicenseKey()
}

// EnsureLicenseActive 校验当前授权存在且未过期（不含章型与次数判断）。
func EnsureLicenseActive() error {
	return utils.EnsureLicenseActive()
}

// CheckUsageQuota 校验当前许可证的使用次数是否已用尽（只检查，不消耗次数）。
func CheckUsageQuota() error {
	return utils.CheckUsageQuota()
}

// RecordUsage 为当前许可证记录一次使用（次数限制许可证已使用次数 +1 并持久化，不限制次数时为空操作）。
// 注意：生成流程已由各章型 Draw 在成功收尾后自动计入，业务代码请勿重复调用，否则会双计。
func RecordUsage() error {
	return utils.RecordUsage()
}

// UsageStats 返回当前许可证的使用统计：已使用次数、次数上限（0 表示不限制）。
func UsageStats() (used int64, limit int64, err error) {
	return utils.UsageStats()
}

// domainAttributeFromService 规范化、转换并做许可证与章型校验（服务层唯一路径）。
// 使用次数由绘制层（各章型 Draw）按是否为预览样章决定是否计入，服务层不再计数。
func domainAttributeFromService(attribute *dto.CustomSealAttribute) (*model.CustomSealAttribute, error) {
	if attribute == nil {
		return nil, dto.ErrSealAttributeNotConfigured
	}
	normalizeSealAttribute(attribute)
	domain, err := attribute.ToDomain()
	if err != nil {
		return nil, err
	}
	if domain == nil {
		return nil, dto.ErrSealAttributeNotConfigured
	}
	normalizeSealAttributeDomain(domain)
	if err := assertSealGenerationAllowedDomain(domain); err != nil {
		return nil, domainErr(err)
	}
	return domain, nil
}

func normalizeSealAttribute(attribute *dto.CustomSealAttribute) {
	if attribute == nil {
		return
	}
	if dto.IsBlank(attribute.SealSpec) {
		attribute.SealSpec = dto.SealSpecCircular42
	}
	if attribute.SealCategory == nil {
		cat := dto.SealCategoryEnterprise
		attribute.SealCategory = &cat
	}
}

func normalizeSealAttributeDomain(attribute *model.CustomSealAttribute) {
	if attribute == nil {
		return
	}
	if utils.IsBlank(attribute.SealSpec) {
		attribute.SealSpec = model.SealSpecCircular42
	}
	if attribute.SealCategory == nil {
		cat := model.SealCategoryEnterprise
		attribute.SealCategory = &cat
	}
}

// assertSealGenerationAllowedDomain 服务层快速失败：许可证有效性 + 章面参数可构建。
// 章型权限与使用次数统一由各章型 Draw 负责，避免服务层与绘制层两处规则漂移。
func assertSealGenerationAllowedDomain(attribute *model.CustomSealAttribute) error {
	// 按「请求章型」校验权限：构建结果可能把 Category 覆写为企业/椭圆章，
	// 若只按构建结果校验，会出现请求 A 章型却按 B 章型放行（或误拒）。
	// 此处只校验有效性与章型权限，不检查使用次数（实际计数由各章型 Draw 在收尾时完成）。
	if err := utils.CheckCustomSealAttributePermission(attribute); err != nil {
		return err
	}
	bundle, err := draw.BuildSealAttributeBundle(attribute)
	if err != nil {
		return err
	}
	if bundle == nil || bundle.Paint == nil {
		return dto.ErrSealAttributeNotConfigured
	}
	return nil
}

// domainErr 将内部错误映射为服务层错误（供 errors.Is 判定）。
func domainErr(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, model.ErrSealAttributeNotConfigured):
		return dto.ErrSealAttributeNotConfigured
	case errors.Is(err, model.ErrSealCategoryNotSupported),
		errors.Is(err, base.ErrSealSpecNotSupported):
		return dto.ErrSealCategoryNotSupported
	}
	return err
}
