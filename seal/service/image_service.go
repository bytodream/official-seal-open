package service

import (
	"image"

	"github.com/bytodream/official-seal-open/seal/service/dto"
)

// UploadBytes 把编码后的章面图上传到对象存储，返回 fileKey。
type UploadBytes func(data []byte, attr dto.ImageAttr, companySeal bool) (string, error)

// ImageService 章面图像转码与上传抽象。
type ImageService interface {
	// ToImageBytes 按 attr 对章面位图做后处理并编码为 PNG。
	ToImageBytes(img image.Image, attr dto.ImageAttr) ([]byte, error)
	// UploadBytes 上传编码后的章面图并返回 fileKey。
	UploadBytes(data []byte, attr dto.ImageAttr, companySeal bool) (string, error)
}

// ImageServiceConfigurator 可选能力：允许 service 层替换上传实现、控制模糊图生成。
//
// 未实现该接口的 ImageService 仍可用于绘制与转码，但有两个后果：
//   - SetUploadBytes / SetBlurGenerate 被静默忽略（配置不生效）；
//   - HasUploadBytes 恒为 false，GetAutoSealImage 会在渲染前直接返回
//     dto.ErrImageUploadNotConfigured。
type ImageServiceConfigurator interface {
	SetUploadBytes(UploadBytes)
	HasUploadBytes() bool
	SetBlurGenerate(isBlurGenerate bool)
	IsBlurGenerate() bool
}
