package service

import "github.com/bytodream/official-seal-open/seal/service/dto"

type ToImageBytes func(img dto.RasterImage, attr dto.ImageAttr) ([]byte, error)
type UploadBytes func(bytes []byte, attr dto.ImageAttr, companySeal bool) (string, error)
type UploadRasterImage func(image dto.RasterImage, attr dto.ImageAttr, companySeal bool) (string, error)

// ImageService 章面图像转码与上传抽象。
type ImageService interface {
	ToImageBytes(img dto.RasterImage, attr dto.ImageAttr) ([]byte, error)
	UploadBytes(bytes []byte, attr dto.ImageAttr, companySeal bool) (string, error)
	UploadRasterImage(image dto.RasterImage, attr dto.ImageAttr, companySeal bool) (string, error)
}

// ImageServiceConfigurator 可选能力：把转码或上传替换为自定义函数。
// 未实现该接口的 ImageService 在调用 service 层的 SetXxx 时保持自身行为不变。
type ImageServiceConfigurator interface {
	SetToImageBytes(ToImageBytes)
	SetUploadBytes(UploadBytes)
	SetUploadRasterImage(UploadRasterImage)
}
