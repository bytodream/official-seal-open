package draw_test

import (
	"testing"

	"github.com/bytodream/official-seal-open/seal/internal/draw"
	"github.com/bytodream/official-seal-open/seal/internal/draw/base"
	"github.com/bytodream/official-seal-open/seal/model"
	"github.com/bytodream/official-seal-open/seal/utils"
)

// fp 返回 float32 指针。
func fp(v float32) *float32 { return &v }

// ip 返回 int 指针。
func ip(v int) *int { return &v }

// cp 返回章面品类指针。
func cp(c model.SealCategory) *model.SealCategory { return &c }

// sealSpecCase 单个章型的显式参数用例（不走 JSON/fixture，全部字段在代码中构造）。
type sealSpecCase struct {
	// name 用例名（章型/规格）
	name string
	// attr 显式构造的自定义章面参数
	attr *model.CustomSealAttribute
	// check 章面属性断言
	check func(t *testing.T, p *base.SealPaintAttribute)
}

// buildSealAttr 显式构造基础参数：默认开关（中心图案/防伪/防伪间隙开启）+ DPI 1000。
func buildSealAttr(cat model.SealCategory, specName string) *model.CustomSealAttribute {
	a := model.NewCustomSealAttribute()
	a.SealCategory = cp(cat)
	a.SealSpec = specName
	return a
}

// allSealSpecCases 覆盖全部支持的章型 × 规格。
func allSealSpecCases() []sealSpecCase {
	const dpi = model.DefaultCustomSealAttributeDPI

	return []sealSpecCase{
		{
			name: "企业圆章/CIRCULAR_42_含脚注与企业代码",
			attr: func() *model.CustomSealAttribute {
				a := buildSealAttr(model.SealCategoryEnterprise, model.SealSpecCircular42)
				a.Content = "XX测试科技有限公司"
				a.Foot = "合同专用章"
				a.EnterpriseCode = "110108123456789"
				return a
			}(),
			check: func(t *testing.T, p *base.SealPaintAttribute) {
				if p.Category != model.SealCategoryEnterprise {
					t.Fatalf("category = %v", p.Category)
				}
				if p.WidthPx != utils.ToIntPixel(dpi, 42) || p.HeightPx != utils.ToIntPixel(dpi, 42) {
					t.Fatalf("size = %dx%d", p.WidthPx, p.HeightPx)
				}
				if p.Content != "XX测试科技有限公司" {
					t.Fatalf("content = %q", p.Content)
				}
				if p.Foot == nil || p.Foot.Text != "合同专用章" {
					t.Fatalf("foot = %+v", p.Foot)
				}
				if p.EnterpriseCode == nil || p.EnterpriseCode.Arc != 75 {
					t.Fatalf("code = %+v（有脚注时弧度应为 75）", p.EnterpriseCode)
				}
				if p.Logo == nil || !p.Logo.FiveStar {
					t.Fatalf("logo = %+v", p.Logo)
				}
				if p.Antifake == nil {
					t.Fatal("antifake info missing")
				}
				if !p.FakeGap {
					t.Fatal("fake gap should be on")
				}
			},
		},
		{
			name: "企业圆章/CIRCULAR_34_无脚注无代码",
			attr: func() *model.CustomSealAttribute {
				a := buildSealAttr(model.SealCategoryEnterprise, "CIRCULAR_34")
				a.Content = "XX测试科技有限公司"
				return a
			}(),
			check: func(t *testing.T, p *base.SealPaintAttribute) {
				if p.Category != model.SealCategoryEnterprise {
					t.Fatalf("category = %v", p.Category)
				}
				if p.WidthPx != utils.ToIntPixel(dpi, 34) || p.HeightPx != utils.ToIntPixel(dpi, 34) {
					t.Fatalf("size = %dx%d", p.WidthPx, p.HeightPx)
				}
				if p.Foot != nil {
					t.Fatalf("foot should be nil, got %+v", p.Foot)
				}
				if p.EnterpriseCode != nil {
					t.Fatalf("code should be nil, got %+v", p.EnterpriseCode)
				}
			},
		},
		{
			name: "企业椭圆章/OVAL_45_30_含内圈文字",
			attr: func() *model.CustomSealAttribute {
				a := buildSealAttr(model.SealCategoryOval, model.SealSpecOval4530)
				a.Content = "XX测试科技有限公司"
				a.Foot = "报关专用章"
				a.InnerContent = "检验检疫"
				return a
			}(),
			check: func(t *testing.T, p *base.SealPaintAttribute) {
				if p.Category != model.SealCategoryOval {
					t.Fatalf("category = %v", p.Category)
				}
				if p.WidthPx != utils.ToIntPixel(dpi, 45) || p.HeightPx != utils.ToIntPixel(dpi, 30) {
					t.Fatalf("size = %dx%d", p.WidthPx, p.HeightPx)
				}
				if p.InnerOval == nil || p.InnerOval.Content != "检验检疫" {
					t.Fatalf("inner oval = %+v", p.InnerOval)
				}
				if p.InnerOval.WidthPx != utils.ToIntPixel(dpi, 33) || p.InnerOval.HeightPx != utils.ToIntPixel(dpi, 19) {
					t.Fatalf("inner oval size = %dx%d", p.InnerOval.WidthPx, p.InnerOval.HeightPx)
				}
			},
		},
		{
			name: "双行椭圆章/OVAL_45_30_32_18",
			attr: func() *model.CustomSealAttribute {
				a := buildSealAttr(model.SealCategoryDoublelineOval, model.SealSpecOval45303218)
				a.Content = "XX测试科技有限公司"
				a.Head = "合同专用章"
				a.Foot = "财务专用"
				a.EnterpriseCode = "1234567890123"
				a.InnerContent = "环形内圈文字"
				return a
			}(),
			check: func(t *testing.T, p *base.SealPaintAttribute) {
				if p.Category != model.SealCategoryDoublelineOval {
					t.Fatalf("category = %v", p.Category)
				}
				if p.InnerLine == nil {
					t.Fatal("inner line missing")
				}
				if p.OvalHead == nil || p.OvalHead.Text != "合同专用章" {
					t.Fatalf("oval head = %+v", p.OvalHead)
				}
				if p.Content != "环形内圈文字" {
					t.Fatalf("content = %q", p.Content)
				}
				if p.EnterpriseCode == nil || p.EnterpriseCode.Arc != 110 {
					t.Fatalf("code = %+v", p.EnterpriseCode)
				}
			},
		},
		{
			name: "双行椭圆章/OVAL_40_30_内椭圆",
			attr: func() *model.CustomSealAttribute {
				a := buildSealAttr(model.SealCategoryDoublelineOval, model.SealSpecOval4030)
				a.Content = "XX测试科技有限公司"
				a.Head = "业务专用章"
				a.InnerContent = "内椭圆文字"
				return a
			}(),
			check: func(t *testing.T, p *base.SealPaintAttribute) {
				if p.Category != model.SealCategoryDoublelineOval {
					t.Fatalf("category = %v", p.Category)
				}
				if p.InnerOval == nil || p.InnerOval.Content != "内椭圆文字" {
					t.Fatalf("inner oval = %+v", p.InnerOval)
				}
				if p.WidthPx != utils.ToIntPixel(dpi, 40) {
					t.Fatalf("width = %d", p.WidthPx)
				}
			},
		},
		{
			name: "双行圆章/CIRCULAR_42",
			attr: func() *model.CustomSealAttribute {
				a := buildSealAttr(model.SealCategoryDoublelineCir, model.SealSpecCircular42)
				a.Content = "XX测试科技有限公司"
				a.Head = "人事专用章"
				a.Foot = "NO.0001"
				a.EnterpriseCode = "510109123456789"
				return a
			}(),
			check: func(t *testing.T, p *base.SealPaintAttribute) {
				if p.Category != model.SealCategoryDoublelineCir {
					t.Fatalf("category = %v", p.Category)
				}
				if p.HeadFoot == nil || p.HeadFoot.Text != "人事专用章" {
					t.Fatalf("head foot = %+v", p.HeadFoot)
				}
				if p.WidthPx != utils.ToIntPixel(dpi, 42) || p.HeightPx != utils.ToIntPixel(dpi, 42) {
					t.Fatalf("size = %dx%d", p.WidthPx, p.HeightPx)
				}
			},
		},
		{
			name: "双圆章/CIRCULAR_42",
			attr: func() *model.CustomSealAttribute {
				a := buildSealAttr(model.SealCategoryDoubleCir, model.SealSpecCircular42)
				a.Content = "XX测试科技有限公司"
				a.Foot = "财务专用章"
				a.InnerContent = "现金收讫"
				a.EnterpriseCode = "420102123456789"
				return a
			}(),
			check: func(t *testing.T, p *base.SealPaintAttribute) {
				if p.Category != model.SealCategoryDoubleCir {
					t.Fatalf("category = %v", p.Category)
				}
				if p.InnerCircle == nil || p.InnerCircle.Content != "现金收讫" {
					t.Fatalf("inner circle = %+v", p.InnerCircle)
				}
				if p.WidthPx != utils.ToIntPixel(dpi, 42) || p.HeightPx != utils.ToIntPixel(dpi, 42) {
					t.Fatalf("size = %dx%d", p.WidthPx, p.HeightPx)
				}
			},
		},
		{
			name: "党委章/CIRCULAR_42",
			attr: func() *model.CustomSealAttribute {
				a := buildSealAttr(model.SealCategoryCommunist, model.SealSpecCircular42)
				a.Content = "XX测试科技有限公司委员会"
				a.Foot = "机关委员会"
				a.EnterpriseCode = "110108123456789"
				return a
			}(),
			check: func(t *testing.T, p *base.SealPaintAttribute) {
				if p.Category != model.SealCategoryCommunist {
					t.Fatalf("category = %v", p.Category)
				}
				if p.WidthPx != utils.ToIntPixel(dpi, 42) || p.HeightPx != utils.ToIntPixel(dpi, 42) {
					t.Fatalf("size = %dx%d", p.WidthPx, p.HeightPx)
				}
				if p.Foot == nil || p.Foot.Text != "机关委员会" {
					t.Fatalf("foot = %+v", p.Foot)
				}
				if p.Logo == nil || p.Logo.SizePx <= 0 || p.Logo.FiveStar {
					// 党委章中心为党徽图案（assets/communist_party.png），非五角星
					t.Fatalf("logo = %+v", p.Logo)
				}
			},
		},
		{
			name: "工会章/CIRCULAR_42",
			attr: func() *model.CustomSealAttribute {
				a := buildSealAttr(model.SealCategoryLabour, model.SealSpecCircular42)
				a.Content = "XX测试科技有限公司工会委员会"
				a.Foot = "工会委员会"
				a.EnterpriseCode = "110108123456789"
				return a
			}(),
			check: func(t *testing.T, p *base.SealPaintAttribute) {
				if p.Category != model.SealCategoryLabour {
					t.Fatalf("category = %v", p.Category)
				}
				if p.WidthPx != utils.ToIntPixel(dpi, 42) || p.HeightPx != utils.ToIntPixel(dpi, 42) {
					t.Fatalf("size = %dx%d", p.WidthPx, p.HeightPx)
				}
			},
		},
		{
			name: "共青团章/CIRCULAR_42",
			attr: func() *model.CustomSealAttribute {
				a := buildSealAttr(model.SealCategoryYouthLeague, model.SealSpecCircular42)
				a.Content = "XX测试科技有限公司委员会"
				a.Head = "团委"
				a.Foot = "团支部"
				a.EnterpriseCode = "110108123456789"
				return a
			}(),
			check: func(t *testing.T, p *base.SealPaintAttribute) {
				if p.Category != model.SealCategoryYouthLeague {
					t.Fatalf("category = %v", p.Category)
				}
				if p.YouthLeague == nil || p.YouthLeague.LineLengthPx <= 0 || p.YouthLeague.LineWidthPx <= 0 {
					t.Fatalf("youth league deco = %+v", p.YouthLeague)
				}
				if p.WidthPx != utils.ToIntPixel(dpi, 42) || p.HeightPx != utils.ToIntPixel(dpi, 42) {
					t.Fatalf("size = %dx%d", p.WidthPx, p.HeightPx)
				}
			},
		},
		{
			name: "多横排圆章/CIRCULAR_50",
			attr: func() *model.CustomSealAttribute {
				a := buildSealAttr(model.SealCategoryMoreFootCircle, model.SealSpecCircular50)
				a.Content = "XX测试科技有限公司"
				a.Head = "合同专用章"
				a.Foots = []model.SealFoot{{Text: "第一分行"}, {Text: "第二分行"}}
				a.EnterpriseCode = "110108123456789"
				return a
			}(),
			check: func(t *testing.T, p *base.SealPaintAttribute) {
				if p.Category != model.SealCategoryMoreFootCircle {
					t.Fatalf("category = %v", p.Category)
				}
				if p.MoreFootHead == nil || p.MoreFootHead.Text != "合同专用章" {
					t.Fatalf("more foot head = %+v", p.MoreFootHead)
				}
				if len(p.MoreFoots) != 2 || p.MoreFoots[0].Text != "第一分行" || p.MoreFoots[1].Text != "第二分行" {
					t.Fatalf("more foots = %+v", p.MoreFoots)
				}
				if p.SealLogo {
					t.Fatal("more foot seal should not draw center logo")
				}
				if p.WidthPx != utils.ToIntPixel(dpi, 50) || p.HeightPx != utils.ToIntPixel(dpi, 50) {
					t.Fatalf("size = %dx%d", p.WidthPx, p.HeightPx)
				}
			},
		},
		{
			name: "多横排椭圆章/OVAL_45_30",
			attr: func() *model.CustomSealAttribute {
				a := buildSealAttr(model.SealCategoryMoreFootOval, model.SealSpecOval4530)
				a.Content = "XX测试科技有限公司"
				a.Head = "发票专用章"
				a.Foots = []model.SealFoot{{Text: "第一税务分局"}}
				return a
			}(),
			check: func(t *testing.T, p *base.SealPaintAttribute) {
				if p.Category != model.SealCategoryMoreFootOval {
					t.Fatalf("category = %v", p.Category)
				}
				if len(p.MoreFoots) != 1 {
					t.Fatalf("more foots = %+v", p.MoreFoots)
				}
				if p.WidthPx != utils.ToIntPixel(dpi, 45) || p.HeightPx != utils.ToIntPixel(dpi, 30) {
					t.Fatalf("size = %dx%d", p.WidthPx, p.HeightPx)
				}
			},
		},
		{
			name: "多横排长方形章/RECTANGLE_50_20",
			attr: func() *model.CustomSealAttribute {
				a := buildSealAttr(model.SealCategoryMoreFootRectangle, model.SealSpecRectangle5020)
				a.Content = "XX测试科技有限公司"
				a.Foots = []model.SealFoot{{Text: "人事部门"}, {Text: "财务部门"}}
				return a
			}(),
			check: func(t *testing.T, p *base.SealPaintAttribute) {
				if p.Category != model.SealCategoryMoreFootRectangle {
					t.Fatalf("category = %v", p.Category)
				}
				if len(p.MoreFoots) != 2 {
					t.Fatalf("more foots = %+v", p.MoreFoots)
				}
				if p.WidthPx != utils.ToIntPixel(dpi, 50) || p.HeightPx != utils.ToIntPixel(dpi, 20) {
					t.Fatalf("size = %dx%d", p.WidthPx, p.HeightPx)
				}
			},
		},
		{
			name: "DIY企业圆章/DIY_SPEC_43mm",
			attr: func() *model.CustomSealAttribute {
				a := buildSealAttr(model.SealCategoryEnterprise, model.SealSpecDiy)
				a.Content = "XX测试科技有限公司"
				a.Foot = "行政专用章"
				a.EnterpriseCode = "370212123456789"
				a.SealWidth = fp(43)
				a.SealHeight = fp(43)
				a.FontSize = fp(7)
				a.EdgeWidth = fp(1.2)
				return a
			}(),
			check: func(t *testing.T, p *base.SealPaintAttribute) {
				if p.Category != model.SealCategoryEnterprise {
					t.Fatalf("category = %v", p.Category)
				}
				if p.WidthPx != utils.ToIntPixel(dpi, 43) || p.HeightPx != utils.ToIntPixel(dpi, 43) {
					t.Fatalf("size = %dx%d", p.WidthPx, p.HeightPx)
				}
			},
		},
	}
}

// TestToSealAttributeFull_AllSealCategories 覆盖全部章型的显式参数构建。
func TestToSealAttributeFull_AllSealCategories(t *testing.T) {
	for _, tc := range allSealSpecCases() {
		t.Run(tc.name, func(t *testing.T) {
			b, err := draw.BuildSealAttributeBundle(tc.attr)
			if err != nil {
				t.Fatalf("ToSealAttributeFull: %v", err)
			}
			if b == nil || b.Paint == nil {
				t.Fatal("empty paint attribute")
			}
			p := b.Paint
			if p.Dpi != model.DefaultCustomSealAttributeDPI {
				t.Fatalf("dpi = %d", p.Dpi)
			}
			if p.WidthPx <= 0 || p.HeightPx <= 0 {
				t.Fatalf("invalid size %dx%d", p.WidthPx, p.HeightPx)
			}
			tc.check(t, p)
		})
	}
}

// TestToSealAttributeFull_UnsupportedSpec 不支持的规格返回 base.ErrSealSpecNotSupported。
func TestToSealAttributeFull_UnsupportedSpec(t *testing.T) {
	a := buildSealAttr(model.SealCategoryCommunist, model.SealSpecOval4530)
	if _, err := draw.BuildSealAttributeBundle(a); err != base.ErrSealSpecNotSupported {
		t.Fatalf("err = %v, want base.ErrSealSpecNotSupported", err)
	}
}

// TestToSealAttributeFull_NilAttribute nil 参数返回 ErrSealAttributeNotConfigured。
func TestToSealAttributeFull_NilAttribute(t *testing.T) {
	if _, err := draw.BuildSealAttributeBundle(nil); err != model.ErrSealAttributeNotConfigured {
		t.Fatalf("err = %v, want ErrSealAttributeNotConfigured", err)
	}
}

// TestToSealAttributeFull_CoverageGuard 保证用例覆盖全部可构建章型。
func TestToSealAttributeFull_CoverageGuard(t *testing.T) {
	want := map[model.SealCategory]bool{
		model.SealCategoryEnterprise:        false,
		model.SealCategoryOval:              false,
		model.SealCategoryDoublelineOval:    false,
		model.SealCategoryDoublelineCir:     false,
		model.SealCategoryDoubleCir:         false,
		model.SealCategoryCommunist:         false,
		model.SealCategoryLabour:            false,
		model.SealCategoryYouthLeague:       false,
		model.SealCategoryMoreFootCircle:    false,
		model.SealCategoryMoreFootOval:      false,
		model.SealCategoryMoreFootRectangle: false,
	}
	for _, tc := range allSealSpecCases() {
		if tc.attr.SealCategory != nil {
			want[*tc.attr.SealCategory] = true
		}
	}
	for cat, covered := range want {
		if !covered {
			t.Fatalf("category %q not covered by test cases", cat)
		}
	}
}
