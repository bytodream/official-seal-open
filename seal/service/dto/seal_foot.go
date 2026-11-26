package dto

// SealFoot 多横排章面的头/脚文字参数。
type SealFoot struct {
	// Text 横排文字内容
	Text string `json:"text,omitempty"`
	// Width 文字区域宽度
	Width *float32 `json:"width,omitempty"`
	// FontSize 字号
	FontSize *int `json:"fontSize,omitempty"`
	// Distance 与参考边界的距离偏移
	Distance *float32 `json:"distance,omitempty"`
	// DiySpec 可配置的 DIY 规格项列表
	DiySpec []DiySpec `json:"diySpec,omitempty"`
}

// DiySpec 单项 DIY 头/脚排版参数。
type DiySpec struct {
	// Distance 距离偏移
	Distance *float32 `json:"distance,omitempty"`
	// FontSize 字号
	FontSize *int `json:"fontSize,omitempty"`
	// Width 文字区域宽度
	Width *float32 `json:"width,omitempty"`
	// Text 显示文字
	Text string `json:"text,omitempty"`
	// Code 配置项编码
	Code string `json:"code,omitempty"`
	// DefaultVal 默认值
	DefaultVal int `json:"defaultVal,omitempty"`
	// Step 调整步长
	Step float64 `json:"step,omitempty"`
	// IsFontMultilineText 是否为多行文字
	IsFontMultilineText *bool `json:"isFontMultilineText,omitempty"`
}
