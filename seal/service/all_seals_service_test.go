package service

import (
	"bytes"
	"image/png"
	"os"
	"testing"

	"github.com/bytodream/official-seal-open/seal/internal/draw"
	"github.com/bytodream/official-seal-open/seal/service/dto"
)

// assertServicePNG 校验服务层输出为合法 PNG。
func assertServicePNG(t *testing.T, data []byte) {
	t.Helper()
	if len(data) < 8 || !bytes.HasPrefix(data, []byte("\x89PNG\r\n\x1a\n")) {
		t.Fatalf("invalid PNG output: %d bytes", len(data))
	}
	if _, err := png.Decode(bytes.NewReader(data)); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

// sealCatPtr 返回章面品类指针。
func sealCatPtr(c dto.SealCategory) *dto.SealCategory { return &c }

// serviceTestFontPath 获取测试字体目录，未配置时跳过。
func serviceTestFontPath(t *testing.T) string {
	t.Helper()
	fontPath := os.Getenv("OFFICIAL_SEAL_FONTS_PATH")
	if draw.ResolveFont(fontPath) == "" {
		t.Skip("需要 OFFICIAL_SEAL_FONTS_PATH 指向包含 ***.ttf 的字体路径")
		return ""
	}
	return fontPath
}

// newRenderSealService 构造使用真实 CommonSealGenerator 的公司章服务。
func newRenderSealService(fontPath string) *CompanySealCreateService {
	return NewCompanySealCreateService(ElectronicAbstractSealCreateService{
		FontPath:     fontPath,
		ImageService: stubImageService{},
	})
}

// serviceSealCase 单个章型的服务层用例（SealModel 与章面参数全部显式构造）。
type serviceSealCase struct {
	// name 用例名
	name string
	// model 显式构造的创建印章请求模型
	model func() *dto.SealModel
}

// buildSealModel 显式构造电子章 SealModel（透明度 30、模糊半径 60、自动生成章面）。
func buildSealModel(specName string, attr *dto.CustomSealAttribute) *dto.SealModel {
	autoTrue := true
	radius := 60
	transparency := "30"
	return &dto.SealModel{
		Category: dto.BusinessTypeElectronic,
		ElectronicModel: &dto.ElectronicModel{
			SealSpec:      specName,
			OpenImageBlur: false,
			Radius:        &radius,
			Transparency:  &transparency,
			IsAuto:        &autoTrue,
			AutoAttribute: attr,
		},
	}
}

// buildServiceAttr 显式构造基础章面参数。
func buildServiceAttr(cat dto.SealCategory, specName string) *dto.CustomSealAttribute {
	a := dto.NewCustomSealAttribute()
	a.SealCategory = sealCatPtr(cat)
	a.SealSpec = specName
	return a
}

// allServiceSealCases 全部章型的服务层用例。
func allServiceSealCases() []serviceSealCase {
	return []serviceSealCase{
		{
			name: "企业圆章/CIRCULAR_42",
			model: func() *dto.SealModel {
				a := buildServiceAttr(dto.SealCategoryEnterprise, dto.SealSpecCircular42)
				a.Content = "XX测试科技有限公司"
				a.Foot = "合同专用章"
				a.EnterpriseCode = "110108123456789"
				return buildSealModel(dto.SealSpecCircular42, a)
			},
		},
		{
			name: "企业椭圆章/OVAL_45_30",
			model: func() *dto.SealModel {
				a := buildServiceAttr(dto.SealCategoryOval, dto.SealSpecOval4530)
				a.Content = "XX测试科技有限公司"
				a.Foot = "报关专用章"
				a.InnerContent = "检验检疫"
				return buildSealModel(dto.SealSpecOval4530, a)
			},
		},
		{
			name: "双行椭圆章/OVAL_45_30_32_18",
			model: func() *dto.SealModel {
				a := buildServiceAttr(dto.SealCategoryDoublelineOval, dto.SealSpecOval45303218)
				a.Content = "XX测试科技有限公司"
				a.Head = "合同专用章"
				a.Foot = "财务专用"
				a.InnerContent = "环形内圈文字"
				return buildSealModel(dto.SealSpecOval45303218, a)
			},
		},
		{
			name: "双行圆章/CIRCULAR_42",
			model: func() *dto.SealModel {
				a := buildServiceAttr(dto.SealCategoryDoublelineCir, dto.SealSpecCircular42)
				a.Content = "XX测试科技有限公司"
				a.Head = "人事专用章"
				a.Foot = "NO.0001"
				return buildSealModel(dto.SealSpecCircular42, a)
			},
		},
		{
			name: "双圆章/CIRCULAR_42",
			model: func() *dto.SealModel {
				a := buildServiceAttr(dto.SealCategoryDoubleCir, dto.SealSpecCircular42)
				a.Content = "XX测试科技有限公司"
				a.Foot = "财务专用章"
				a.InnerContent = "现金收讫"
				return buildSealModel(dto.SealSpecCircular42, a)
			},
		},
		{
			name: "党委章/CIRCULAR_42",
			model: func() *dto.SealModel {
				a := buildServiceAttr(dto.SealCategoryCommunist, dto.SealSpecCircular42)
				a.Content = "XX测试科技有限公司委员会"
				a.Foot = "机关委员会"
				return buildSealModel(dto.SealSpecCircular42, a)
			},
		},
		{
			name: "工会章/CIRCULAR_42",
			model: func() *dto.SealModel {
				a := buildServiceAttr(dto.SealCategoryLabour, dto.SealSpecCircular42)
				a.Content = "XX测试科技有限公司工会委员会"
				a.Foot = "工会委员会"
				return buildSealModel(dto.SealSpecCircular42, a)
			},
		},
		{
			name: "共青团章/CIRCULAR_42",
			model: func() *dto.SealModel {
				a := buildServiceAttr(dto.SealCategoryYouthLeague, dto.SealSpecCircular42)
				a.Content = "XX测试科技有限公司委员会"
				a.Head = "团委"
				a.Foot = "团支部"
				return buildSealModel(dto.SealSpecCircular42, a)
			},
		},
		{
			name: "多横排圆章/CIRCULAR_50",
			model: func() *dto.SealModel {
				a := buildServiceAttr(dto.SealCategoryMoreFootCircle, dto.SealSpecCircular50)
				a.Content = "XX测试科技有限公司"
				a.Head = "合同专用章"
				a.Foots = []dto.SealFoot{{Text: "第一分行"}, {Text: "第二分行"}}
				return buildSealModel(dto.SealSpecCircular50, a)
			},
		},
		{
			name: "多横排椭圆章/OVAL_45_30",
			model: func() *dto.SealModel {
				a := buildServiceAttr(dto.SealCategoryMoreFootOval, dto.SealSpecOval4530)
				a.Content = "XX测试科技有限公司"
				a.Head = "发票专用章"
				a.Foots = []dto.SealFoot{{Text: "第一税务分局"}}
				return buildSealModel(dto.SealSpecOval4530, a)
			},
		},
		{
			name: "多横排长方形章/RECTANGLE_50_20",
			model: func() *dto.SealModel {
				a := buildServiceAttr(dto.SealCategoryMoreFootRectangle, dto.SealSpecRectangle5020)
				a.Content = "XX测试科技有限公司"
				a.Foots = []dto.SealFoot{{Text: "人事部门"}, {Text: "财务部门"}}
				return buildSealModel(dto.SealSpecRectangle5020, a)
			},
		},
	}
}

// TestCompanySealService_PreviewAllSealCategories 全部章型走服务层 Preview 输出合法 PNG。
func TestCompanySealService_PreviewAllSealCategories(t *testing.T) {
	skipUnlessAllSealCategoriesLicensed(t)
	fontPath := serviceTestFontPath(t)
	svc := newRenderSealService(fontPath)
	for _, tc := range allServiceSealCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			data, err := svc.Preview(tc.model())
			if err != nil {
				t.Fatalf("Preview: %v", err)
			}
			assertServicePNG(t, data)
		})
	}
}

// TestCompanySealService_PreviewOpenImageBlur 开启模糊章面后 Preview 仍输出合法 PNG。
func TestCompanySealService_PreviewOpenImageBlur(t *testing.T) {
	skipIfNoLicense(t)
	fontPath := serviceTestFontPath(t)
	svc := newRenderSealService(fontPath)
	sm := buildSealModel(dto.SealSpecCircular42, func() *dto.CustomSealAttribute {
		a := buildServiceAttr(dto.SealCategoryEnterprise, dto.SealSpecCircular42)
		a.Content = "模糊章面校验专用公司名称"
		return a
	}())
	sm.ElectronicModel.OpenImageBlur = true
	data, err := svc.Preview(sm)
	if err != nil {
		t.Fatalf("Preview: %v", err)
	}
	assertServicePNG(t, data)
}

// TestCompanySealService_GetAutoSealImageAllSealCategories 全部章型生成清晰/模糊章面并上传。
func TestCompanySealService_GetAutoSealImageAllSealCategories(t *testing.T) {
	skipUnlessAllSealCategoriesLicensed(t)
	fontPath := serviceTestFontPath(t)
	svc := newRenderSealService(fontPath)
	for _, tc := range allServiceSealCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			sm := tc.model()
			img, err := svc.GetAutoSealImage(sm.ToSeal(), sm.ElectronicModel.AutoAttribute)
			if err != nil {
				t.Fatalf("GetAutoSealImage: %v", err)
			}
			if img.ImageKey == "" || img.BlurImageKey == "" {
				t.Fatalf("missing keys: %+v", img)
			}
			if img.Size <= 0 {
				t.Fatalf("size = %d", img.Size)
			}
		})
	}
}

// TestCompanySealService_PreviewNilAttribute 章面参数缺失时 Preview 报错。
func TestCompanySealService_PreviewNilAttribute(t *testing.T) {
	fontPath := serviceTestFontPath(t)
	svc := newRenderSealService(fontPath)
	sm := buildSealModel(dto.SealSpecCircular42, nil)
	if _, err := svc.Preview(sm); err != dto.ErrSealAttributeNotConfigured {
		t.Fatalf("err = %v, want ErrSealAttributeNotConfigured", err)
	}
}
