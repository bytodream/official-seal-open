package service

import "github.com/bytodream/official-seal-open/seal/service/dto"

// ImageService 章面图像转码与上传抽象。
type ImageService interface {
	ToImageBytes(img dto.RasterImage, attr dto.ImageAttr) ([]byte, error)
	UploadBytes(data []byte, attr dto.ImageAttr, companySeal bool) (string, error)
	UploadRasterImage(img dto.RasterImage, attr dto.ImageAttr, companySeal bool) (string, error)
}
