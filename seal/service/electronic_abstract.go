package service

// ElectronicAbstractSealCreateService 电子章创建服务的公共依赖。
// 承载各电子章创建服务共用的字体路径、图像上传与水印；章面生成固定为 internal/render.CommonSealGenerator。
type ElectronicAbstractSealCreateService struct {
	// FontPath 字体路径（配置项 fonts.path）
	FontPath string
	// ImageService 章面图片服务
	ImageService ImageService
	// Watermark 章面水印回调（保留扩展点，当前生成/预览链路未使用；「样章」标识由绘制层内部实现）
	Watermark WatermarkFunc
}
