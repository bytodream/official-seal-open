package model

import "errors"

// ErrSealAttributeNotConfigured 无法从 CustomSealAttribute 构建章面属性（入参为空或规格缺失）。
var ErrSealAttributeNotConfigured = errors.New("获取印章属性规格失败: 章面参数未配置或无法构建规格")

// ErrSealCategoryNotSupported 不支持的章面品类/章型。
var ErrSealCategoryNotSupported = errors.New("获取印章属性规格失败: 不支持的 SealCategory")
