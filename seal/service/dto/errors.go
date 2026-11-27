package dto

import "errors"

// ErrSealAttributeNotConfigured 无法从 CustomSealAttribute 构建章面（参数缺失或无法解析为规格）。
var ErrSealAttributeNotConfigured = errors.New("获取印章属性规格失败: 章面参数未配置或无法构建规格")

// ErrSealCategoryNotSupported 章面品类/章型当前不支持。
var ErrSealCategoryNotSupported = errors.New("获取印章属性规格失败: 不支持的 SealCategory")

// ErrImageServiceNotConfigured 章面图片服务未注入，无法上传生成的位图。
var ErrImageServiceNotConfigured = errors.New("章面图片服务未配置")

// ErrImageUploadNotConfigured DefaultImageService 只提供转码，上传需接入自有存储实现。
var ErrImageUploadNotConfigured = errors.New("章面图片上传未配置：DefaultImageService 仅提供转码，请接入自有存储")
