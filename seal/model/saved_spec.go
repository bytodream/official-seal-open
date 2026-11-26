package model

import "errors"

// ErrSavedSpecNotFound 预设规格不存在
var ErrSavedSpecNotFound = errors.New("预设规格不存在")

// SavedSpecLoader 按组织与规格名加载已保存的自定义章面属性。
type SavedSpecLoader interface {
	LoadSavedSpec(orgID int64, categoryType, savedSpecName string, attr *CustomSealAttribute, electronicModel *SealModel) (*CustomSealAttribute, error)
}

// ComputeSpecAttributes 按已保存规格合并自定义属性。
func ComputeSpecAttributes(orgID int64, attr *CustomSealAttribute, electronicModel *SealModel, loader SavedSpecLoader) (*CustomSealAttribute, error) {
	if attr == nil || orgID == 0 {
		return attr, nil
	}
	if !attr.UseSavedSpec() {
		return attr, nil
	}
	if loader == nil {
		return nil, ErrSavedSpecNotFound
	}
	return loader.LoadSavedSpec(orgID, attr.SealCategoryType, attr.SavedSpecName, attr, electronicModel)
}
