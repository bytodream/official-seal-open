package render

import (
	"bytes"
	"errors"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytodream/official-seal-open/seal/internal/draw"
	"github.com/bytodream/official-seal-open/seal/internal/draw/base"
	"github.com/bytodream/official-seal-open/seal/model"
	"github.com/bytodream/official-seal-open/seal/utils"
)

// testFontPath 获取测试字体路径，未配置时跳过。
func testFontPath(t *testing.T) string {
	t.Helper()
	return draw.ResolveFont("")
}

// renderCase 单个章型的端到端渲染用例（参数显式构造）。
type renderCase struct {
	// name 用例名
	name string
	// attr 显式构造的自定义章面参数
	attr *model.CustomSealAttribute
	// wantW/wantH 期望画布尺寸（px），为 0 时只校验尺寸合法
	wantW, wantH int
}

// allRenderCases 全部章型的显式参数用例。
func allRenderCases() []renderCase {
	const dpi = model.DefaultCustomSealAttributeDPI
	w := func(mm float32) int { return utils.ToIntPixel(dpi, mm) }

	base := func(cat model.SealCategory, specName string) *model.CustomSealAttribute {
		a := model.NewCustomSealAttribute()
		a.SealCategory = catPtr(cat)
		a.SealSpec = specName
		return a
	}

	return []renderCase{
		{
			name: "企业圆章/CIRCULAR_42",
			attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryEnterprise, model.SealSpecCircular42)
				a.Content = "XX测试科技有限公司"
				a.Foot = "合同专用章"
				a.EnterpriseCode = "110108123456789"
				return a
			}(),
			wantW: w(42), wantH: w(42),
		},
		{
			name: "企业圆章/CIRCULAR_34",
			attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryEnterprise, "CIRCULAR_34")
				a.Content = "XX测试科技有限公司"
				return a
			}(),
			wantW: w(34), wantH: w(34),
		},
		{
			name: "企业椭圆章/OVAL_45_30",
			attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryOval, model.SealSpecOval4530)
				a.Content = "XX测试科技有限公司"
				a.Foot = "报关专用章"
				a.InnerContent = "检验检疫"
				return a
			}(),
			wantW: w(45), wantH: w(30),
		},
		{
			name: "双行椭圆章/OVAL_45_30_32_18",
			attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryDoublelineOval, model.SealSpecOval45303218)
				a.Content = "XX测试科技有限公司"
				a.Head = "合同专用章"
				a.Foot = "财务专用"
				a.EnterpriseCode = "1234567890123"
				a.InnerContent = "环形内圈文字"
				return a
			}(),
			wantW: w(45), wantH: w(30),
		},
		{
			name: "双行圆章/CIRCULAR_42",
			attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryDoublelineCir, model.SealSpecCircular42)
				a.Content = "XX测试科技有限公司"
				a.Head = "人事专用章"
				a.Foot = "NO.0001"
				a.EnterpriseCode = "510109123456789"
				return a
			}(),
			wantW: w(42), wantH: w(42),
		},
		{
			name: "双圆章/CIRCULAR_42",
			attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryDoubleCir, model.SealSpecCircular42)
				a.Content = "XX测试科技有限公司"
				a.Foot = "财务专用章"
				a.InnerContent = "现金收讫"
				a.EnterpriseCode = "420102123456789"
				return a
			}(),
			wantW: w(42), wantH: w(42),
		},
		{
			name: "党委章/CIRCULAR_42",
			attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryCommunist, model.SealSpecCircular42)
				a.Content = "XX测试科技有限公司委员会"
				a.Foot = "机关委员会"
				a.EnterpriseCode = "110108123456789"
				return a
			}(),
			wantW: w(42), wantH: w(42),
		},
		{
			name: "工会章/CIRCULAR_42",
			attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryLabour, model.SealSpecCircular42)
				a.Content = "XX测试科技有限公司工会委员会"
				a.Foot = "工会委员会"
				a.EnterpriseCode = "110108123456789"
				return a
			}(),
			wantW: w(42), wantH: w(42),
		},
		{
			name: "共青团章/CIRCULAR_42",
			attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryYouthLeague, model.SealSpecCircular42)
				a.Content = "XX测试科技有限公司委员会"
				a.Head = "团委"
				a.Foot = "团支部"
				a.EnterpriseCode = "110108123456789"
				return a
			}(),
			wantW: w(42), wantH: w(42),
		},
		{
			name: "多横排圆章/CIRCULAR_50",
			attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryMoreFootCircle, model.SealSpecCircular50)
				a.Content = "XX测试科技有限公司"
				a.Head = "合同专用章"
				a.Foots = []model.SealFoot{{Text: "第一分行"}, {Text: "第二分行"}}
				a.EnterpriseCode = "110108123456789"
				return a
			}(),
			wantW: w(50), wantH: w(50),
		},
		{
			name: "多横排椭圆章/OVAL_45_30",
			attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryMoreFootOval, model.SealSpecOval4530)
				a.Content = "XX测试科技有限公司"
				a.Head = "发票专用章"
				a.Foots = []model.SealFoot{{Text: "第一税务分局"}}
				return a
			}(),
			wantW: w(45), wantH: w(30),
		},
		{
			name: "多横排长方形章/RECTANGLE_50_20",
			attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryMoreFootRectangle, model.SealSpecRectangle5020)
				a.Content = "XX测试科技有限公司"
				a.Foots = []model.SealFoot{{Text: "人事部门"}, {Text: "财务部门"}}
				return a
			}(),
			wantW: w(50), wantH: w(20),
		},
		{
			name: "DIY企业圆章/DIY_SPEC_43mm",
			attr: func() *model.CustomSealAttribute {
				a := base(model.SealCategoryEnterprise, model.SealSpecDiy)
				a.Content = "XX测试科技有限公司"
				a.Foot = "行政专用章"
				a.EnterpriseCode = "370212123456789"
				a.SealWidth = f32p(43)
				a.SealHeight = f32p(43)
				a.FontSize = f32p(7)
				a.EdgeWidth = f32p(1.2)
				return a
			}(),
			wantW: w(43), wantH: w(43),
		},
	}
}

// assertPNGBytes 校验 PNG 魔数并解码返回尺寸。
func assertPNGBytes(t *testing.T, data []byte) (w, h int) {
	t.Helper()
	if len(data) < 8 {
		t.Fatalf("png too small: %d bytes", len(data))
	}
	if !bytes.HasPrefix(data, []byte("\x89PNG\r\n\x1a\n")) {
		t.Fatal("invalid PNG signature")
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	b := img.Bounds()
	return b.Dx(), b.Dy()
}

// skipIfLicenseLimited 测试许可证未覆盖该章型或次数用尽时跳过，不作为绘制逻辑失败。
func skipIfLicenseLimited(t *testing.T, err error) {
	t.Helper()
	if errors.Is(err, utils.ErrLicenseRequired) || errors.Is(err, utils.ErrLicenseExpired) ||
		errors.Is(err, utils.ErrSealCategoryNotLicensed) || errors.Is(err, utils.ErrUsageLimitExceeded) {
		t.Skipf("测试许可证限制，跳过：%v", err)
	}
}

// TestRenderAllSealCategoriesPNG 全部章型端到端渲染为 PNG。
func TestRenderAllSealCategoriesPNG(t *testing.T) {
	// 次数记录隔离到临时文件，避免污染 ~/.official-seal/usage.json
	t.Setenv("OS_USAGE_FILE", filepath.Join(t.TempDir(), "usage.json"))
	fontPath := testFontPath(t)
	applyTestLicense(t)
	for _, tc := range allRenderCases() {
		t.Run(tc.name, func(t *testing.T) {
			data, err := renderCustomSealPNG(tc.attr, fontPath, model.DefaultSealTransparency)
			if err != nil {
				skipIfLicenseLimited(t, err)
				t.Fatalf("renderCustomSealPNG: %v", err)
			}
			savePng(t, tc.name, data)
			w, h := assertPNGBytes(t, data)
			if tc.wantW > 0 && (w != tc.wantW || h != tc.wantH) {
				t.Fatalf("size = %dx%d, want %dx%d", w, h, tc.wantW, tc.wantH)
			}
		})
	}
}

// TestRenderCustomSealPNG_Transparency 自定义透明度生效（输出仍为合法 PNG）。
func TestRenderCustomSealPNG_Transparency(t *testing.T) {
	t.Setenv("OS_USAGE_FILE", filepath.Join(t.TempDir(), "usage.json"))
	fontPath := testFontPath(t)
	a := model.NewCustomSealAttribute()
	cat := model.SealCategoryEnterprise
	a.SealCategory = &cat
	a.SealSpec = model.SealSpecCircular42
	a.Content = "透明度校验专用公司名称"
	data, err := renderCustomSealPNG(a, fontPath, 100)
	if err != nil {
		skipIfLicenseLimited(t, err)
		t.Fatalf("renderCustomSealPNG: %v", err)
	}
	savePng(t, "TestrenderCustomSealPNG_Transparency", data)
	assertPNGBytes(t, data)
}

// TestRenderCustomSealPNG_NilAttribute nil 参数返回 ErrSealAttributeNotConfigured。
func TestRenderCustomSealPNG_NilAttribute(t *testing.T) {
	fontPath := testFontPath(t)
	if _, err := renderCustomSealPNG(nil, fontPath, model.DefaultSealTransparency); err != model.ErrSealAttributeNotConfigured {
		t.Fatalf("err = %v, want ErrSealAttributeNotConfigured", err)
	}
}

// TestRenderCustomSealPNG_UnsupportedSpec 不支持的规格返回 ErrSealSpecNotSupported。
func TestRenderCustomSealPNG_UnsupportedSpec(t *testing.T) {
	fontPath := testFontPath(t)
	a := model.NewCustomSealAttribute()
	cat := model.SealCategoryCommunist
	a.SealCategory = &cat
	a.SealSpec = model.SealSpecOval4530
	if _, err := renderCustomSealPNG(a, fontPath, model.DefaultSealTransparency); err != base.ErrSealSpecNotSupported {
		t.Fatalf("err = %v, want ErrSealSpecNotSupported", err)
	}
}

func savePng(t *testing.T, name string, data []byte) {
	t.Helper()
	// out 目录被 .gitignore 忽略，干净检出时需自行创建，否则 WriteFile 直接失败。
	if err := os.MkdirAll("./out", 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	previewOutputPath := fmt.Sprintf("./out/%s_%d.png", strings.ReplaceAll(name, "/", "_"), time.Now().Unix())
	if err := os.WriteFile(previewOutputPath, data, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	abs, _ := filepath.Abs(previewOutputPath)
	fmt.Printf("已写入: %s (%d bytes)\n", abs, len(data))
}

// builtinTestLicense 内置测试许可证：仅用于本地/CI 渲染测试的授权校验，
// 可通过环境变量 OS_LICENSE 覆盖（生产与集成环境请务必使用环境变量）。
const builtinTestLicense = "AgEAy7qZ7+eVbLwP0H7HR0MUd9IppoSqwvSgdr6IWTD8yzN5gesf5PgNMi4zkX58KTJB9Kk5F0YFtv7Pfr1TmoaDwv3ksIDMU63RNYEt8HAuM4j4r2b0kC52E96SgSmamG1aw3Gw2lDKDT+wdkzwl71ZuOx2b7prCcecTELdv9nJ+tyDJqYj79p1R1O0CzAHL1Ra1jmhzj7RVKwJNg5Ua7+cIzeyFp+pGfkTAuO0cHxoiT+ag5O+bnIgZdzSNMCAlLdrJt0s5qipt99dGzAjRTs4cDOBozMRE5YcFfG09IRjaMM4dpmSVTECHUonZ2RvyYpjNiLM2IsZo0Hzt1V26qUM6+97SwrB2tRnN9PveYzWZo5MKQRfyCWgYVQ+ORx8i48rg4iBbd0s0vjCwstLR83zeaE4/hnZL4uZzkdK3THT2lUSghfEYCdeGvXDkQF7rTsiNmUoJGEkU+j8wdo1c2f1I0zZ8fFFSYD0t6RbHlgagyRoHQF6145sOizJSgTXG68RP3EbabACcJ4FXglD7S/YFniDCg9tbwEchLUZyALALdERVyY7y8gbVYlfowJdeTvBA0Zvjp6SvrdjQK+5dPcncFTRN6vboghAcExbKDKYEWzn30Kt6lE7xMvxmjnnES5/LW/Py30mtGD8ULgFSmF0MWuTTTvTWuhW3awHfpA5Ku4pxmq2llTxjpd8EaAZyLfCTGbg"

// applyTestLicense 为需要许可证的测试装载授权：优先 OS_LICENSE，缺失时回退内置测试密钥。
func applyTestLicense(t *testing.T) {
	t.Helper()
	key := os.Getenv("OS_LICENSE")
	if key == "" {
		key = builtinTestLicense
	}
	if err := utils.SetLicenseKey(key); err != nil {
		t.Skipf("测试许可证不可用，跳过：%v", err)
	}
}
