package model

// DiySealAttribute DIY 章面尺寸与排版参数（主单位为 mm，指针为空时走默认值）。
type DiySealAttribute struct {
	// SealHeight 章面高度
	SealHeight *float32
	// SealWidth 章面宽度
	SealWidth *float32
	// EdgeWidth 外圈线宽
	EdgeWidth *float32
	// LogoSize 中心图案尺寸
	LogoSize *float32
	// LogoUpmove 中心图案上移距离
	LogoUpmove *float32
	// Arc 主文字弧度
	Arc *int
	// FontSize 主文字字号
	FontSize *float32
	// WordLinePadding 主文字字距
	WordLinePadding *float32
	// EnterpriseCodeFontSize 企业代码字号
	EnterpriseCodeFontSize *float32
	// EnterpriseCodeWordLinePadding 企业代码字距
	EnterpriseCodeWordLinePadding *float32
	// EnterpriseCodeArc 企业代码弧度
	EnterpriseCodeArc *int
	// FootFontSize 脚排字号
	FootFontSize *float32
	// FootWidth 脚排宽度
	FootWidth *float32
	// FootDistance 脚排与主文字距离
	FootDistance *float32
	// HeadFontSize 头排字号
	HeadFontSize *float32
	// HeadWidth 头排宽度
	HeadWidth *float32
	// LineLength 装饰线长度
	LineLength *float32
	// LineWidth 装饰线宽
	LineWidth *float32
	// LineInterval 装饰线间距
	LineInterval *float32
	// HeadDistance 头排与主文字距离
	HeadDistance *float32
	// InnerContent 内圈文字
	InnerContent string
	// InnerOvalEdgeWidth 内椭圆线宽
	InnerOvalEdgeWidth *float32
	// InnerOvalFontSize 内椭圆文字字号
	InnerOvalFontSize *float32
	// InnerOvalWidth 内椭圆宽度
	InnerOvalWidth *float32
	// InnerOvalWordLinePadding 内椭圆文字字距
	InnerOvalWordLinePadding *float32
	// InnerOvalHeight 内椭圆高度
	InnerOvalHeight *float32
	// InnerMargin 内椭圆文字边距
	InnerMargin *float32
	// InnerEdgeWidth 内圈线宽
	InnerEdgeWidth *float32
}

func (a *CustomSealAttribute) ToDiySealAttribute() *DiySealAttribute {
	if a == nil {
		return &DiySealAttribute{}
	}
	return &DiySealAttribute{
		SealHeight:                    a.SealHeight,
		SealWidth:                     a.SealWidth,
		EdgeWidth:                     a.EdgeWidth,
		LogoSize:                      a.LogoSize,
		LogoUpmove:                    a.LogoUpmove,
		Arc:                           a.Arc,
		FontSize:                      a.FontSize,
		WordLinePadding:               a.WordLinePadding,
		EnterpriseCodeFontSize:        a.EnterpriseCodeFontSize,
		EnterpriseCodeWordLinePadding: a.EnterpriseCodeWordLinePadding,
		EnterpriseCodeArc:             a.EnterpriseCodeArc,
		FootFontSize:                  a.FootFontSize,
		FootWidth:                     a.FootWidth,
		FootDistance:                  a.FootDistance,
		HeadFontSize:                  a.HeadFontSize,
		HeadWidth:                     a.HeadWidth,
		LineLength:                    a.LineLength,
		LineWidth:                     a.LineWidth,
		LineInterval:                  a.LineInterval,
		HeadDistance:                  a.HeadDistance,
		InnerContent:                  a.InnerContent,
		InnerOvalEdgeWidth:            a.InnerOvalEdgeWidth,
		InnerOvalFontSize:             a.InnerOvalFontSize,
		InnerOvalWidth:                a.InnerOvalWidth,
		InnerOvalWordLinePadding:      a.InnerOvalWordLinePadding,
		InnerOvalHeight:               a.InnerOvalHeight,
		InnerMargin:                   a.InnerMargin,
		InnerEdgeWidth:                a.InnerEdgeWidth,
	}
}
