package model

import "strings"

// IsBlank 判断字符串是否为空或仅空白。
// 与 utils.IsBlank 语义一致；model 处于依赖底层，不能反向引用 utils，故保留本地实现。
func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}
