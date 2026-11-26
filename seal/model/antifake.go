package model

import "image/color"

// AntifakeInfo 印章中心防伪纹图案参数。
type AntifakeInfo struct {
	Content  string      // 掩膜文字：普通章 "OS"，共青团章 "Q"
	Line     int         // 光栅线数，决定竖/横/斜纹间距
	LineSize int         // 单条光栅线像素宽度
	Dpi      int         // 防伪图生成 DPI（与印章 DPI 缩放配合）
	Color    color.Color // 纹路与文字颜色
	FontFile string      // 宋体等字体路径
}

// BuildAntifakeInfo 从 CustomSealAttribute 构建防伪参数。
// AntifakeCode 为掩膜文字：非空时优先使用；为空时普通章默认 "OS"，
// 共青团章（SealCategoryYouthLeague）默认 "Q"（首次调用会回写 AntifakeCode）。
// 未开启防伪或属性为空时返回 nil。
func (a *CustomSealAttribute) BuildAntifakeInfo(fontFile string) *AntifakeInfo {
	if a == nil || !a.IsAntifake() {
		return nil
	}
	if IsBlank(a.AntifakeCode) {
		a.AntifakeCode = "OS"
		if a.SealCategory != nil && *a.SealCategory == SealCategoryYouthLeague {
			a.AntifakeCode = "Q"
		}
	}
	return &AntifakeInfo{
		Content:  a.AntifakeCode,
		Line:     100,
		LineSize: 3,
		Dpi:      1200,
		Color:    SealRed,
		FontFile: fontFile,
	}
}

// Space 计算相邻光栅线中心距（像素）。
func (info *AntifakeInfo) Space() float64 {
	if info.Line <= 1 {
		return float64(info.Dpi)
	}
	return float64(info.Dpi-info.Line*info.LineSize) / float64(info.Line-1)
}
