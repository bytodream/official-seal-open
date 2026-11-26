package utils

import (
	"math"
	"strings"
)

// F32Or 章面生成辅助函数。
func F32Or(p *float32, def float32) float32 {
	if p == nil {
		return def
	}
	return *p
}

// IsBlank 判断字符串是否为空或仅空白。
func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

const millimetresPerInch = 25.4

// ToFloatPixel 将毫米长度按 DPI 转为浮点像素。
func ToFloatPixel(dpi int, lengthMm float32) float32 {
	f := float64(lengthMm) / millimetresPerInch * float64(dpi)
	return float32(math.Round(f*100) / 100)
}

// ToIntPixel 将毫米长度按 DPI 转为整数像素。
func ToIntPixel(dpi int, lengthMm float32) int {
	f := float64(lengthMm) / millimetresPerInch * float64(dpi)
	return int(math.Round(f))
}
