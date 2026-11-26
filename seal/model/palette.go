package model

import "image/color"

// SealRed 印泥红，各规格属性构建与章面绘制共用的默认颜色。
var SealRed = color.RGBA{R: 255, G: 0, B: 0, A: 255}

// DefaultSealTransparency 默认章面透明度 Alpha。
const DefaultSealTransparency = 180
