package render

import (
	"github.com/bytodream/official-seal-open/seal/model"
)

// ShowcaseItem 文档/示例用章面样例（参数与 all_seals_render_test 一致）。
type ShowcaseItem struct {
	Slug string // 输出文件名（不含扩展名）
	Name string // 中文展示名
	Attr *model.CustomSealAttribute
}

func f32p(v float32) *float32 { return &v }

func catPtr(c model.SealCategory) *model.SealCategory { return &c }

// ShowcaseItems 返回已支持端到端渲染的章型样例列表。
func ShowcaseItems() []ShowcaseItem {
	base := func(cat model.SealCategory, specName string) *model.CustomSealAttribute {
		a := model.NewCustomSealAttribute()
		a.SealCategory = catPtr(cat)
		a.SealSpec = specName
		return a
	}
	return []ShowcaseItem{
		{
			Slug: "enterprise_circular_42",
			Name: "企业圆章 CIRCULAR_42",
			Attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryEnterprise, model.SealSpecCircular42)
				a.Content = "示例科技有限公司"
				a.Foot = "合同专用章"
				a.EnterpriseCode = "110108123456789"
				return a
			}(),
		},
		{
			Slug: "enterprise_circular_34",
			Name: "企业圆章 CIRCULAR_34",
			Attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryEnterprise, "CIRCULAR_34")
				a.Content = "示例科技有限公司"
				return a
			}(),
		},
		{
			Slug: "oval_45_30",
			Name: "企业椭圆章 OVAL_45_30",
			Attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryOval, model.SealSpecOval4530)
				a.Content = "示例科技有限公司"
				a.Foot = "报关专用章"
				a.InnerContent = "检验检疫"
				return a
			}(),
		},
		{
			Slug: "doubleline_oval",
			Name: "双行椭圆章",
			Attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryDoublelineOval, model.SealSpecOval45303218)
				a.Content = "示例科技有限公司"
				a.Head = "合同专用章"
				a.Foot = "财务专用"
				a.EnterpriseCode = "1234567890123"
				a.InnerContent = "环形内圈文字"
				return a
			}(),
		},
		{
			Slug: "doubleline_cir",
			Name: "双行圆章",
			Attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryDoublelineCir, model.SealSpecCircular42)
				a.Content = "示例科技有限公司"
				a.Head = "人事专用章"
				a.Foot = "NO.0001"
				a.EnterpriseCode = "510109123456789"
				return a
			}(),
		},
		{
			Slug: "double_cir",
			Name: "双圆章",
			Attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryDoubleCir, model.SealSpecCircular42)
				a.Content = "示例科技有限公司"
				a.Foot = "财务专用章"
				a.InnerContent = "现金收讫"
				a.EnterpriseCode = "420102123456789"
				return a
			}(),
		},
		{
			Slug: "communist",
			Name: "党委章",
			Attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryCommunist, model.SealSpecCircular42)
				a.Content = "示例科技有限公司委员会"
				a.Foot = "机关委员会"
				a.EnterpriseCode = "110108123456789"
				return a
			}(),
		},
		{
			Slug: "labour",
			Name: "工会章",
			Attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryLabour, model.SealSpecCircular42)
				a.Content = "示例科技有限公司工会委员会"
				a.Foot = "工会委员会"
				a.EnterpriseCode = "110108123456789"
				return a
			}(),
		},
		{
			Slug: "youth_league",
			Name: "共青团章",
			Attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryYouthLeague, model.SealSpecCircular42)
				a.Content = "示例科技有限公司委员会"
				a.Head = "团委"
				a.Foot = "团支部"
				a.EnterpriseCode = "110108123456789"
				return a
			}(),
		},
		{
			Slug: "more_foot_circle",
			Name: "多横排圆章",
			Attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryMoreFootCircle, model.SealSpecCircular50)
				a.Content = "示例科技有限公司"
				a.Head = "合同专用章"
				a.Foots = []model.SealFoot{{Text: "第一分行"}, {Text: "第二分行"}}
				a.EnterpriseCode = "110108123456789"
				return a
			}(),
		},
		{
			Slug: "more_foot_oval",
			Name: "多横排椭圆章",
			Attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryMoreFootOval, model.SealSpecOval4530)
				a.Content = "示例科技有限公司"
				a.Head = "发票专用章"
				a.Foots = []model.SealFoot{{Text: "第一税务分局"}}
				return a
			}(),
		},
		{
			Slug: "more_foot_rectangle",
			Name: "多横排长方形章",
			Attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryMoreFootRectangle, model.SealSpecRectangle5020)
				a.Content = "示例科技有限公司"
				a.Foots = []model.SealFoot{{Text: "人事部门"}, {Text: "财务部门"}}
				return a
			}(),
		},
		{
			Slug: "diy_enterprise",
			Name: "DIY 企业圆章 43mm",
			Attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryEnterprise, model.SealSpecDiy)
				a.Content = "示例科技有限公司"
				a.Foot = "行政专用章"
				a.EnterpriseCode = "370212123456789"
				a.SealWidth = f32p(43)
				a.SealHeight = f32p(43)
				a.FontSize = f32p(7)
				a.EdgeWidth = f32p(1.2)
				return a
			}(),
		},
	}
}

// PNGForShowcase 生成展示样例 PNG，与正式生成同样经过章型权限校验并计入使用次数。
func PNGForShowcase(attr *model.CustomSealAttribute, fontPath string, transparency int) ([]byte, error) {
	return renderCustomSealPNG(attr, fontPath, transparency)
}
