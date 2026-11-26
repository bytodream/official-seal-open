package service

import (
	"os"
	"testing"

	"github.com/bytodream/official-seal-open/seal/internal/render"
	"github.com/bytodream/official-seal-open/seal/service/dto"
)

func skipIfNoLicense(t *testing.T) {
	t.Helper()
	lk := GetLicenseKey()
	if lk == nil || !lk.IsLicensed() {
		t.Skip("需要有效 OS_LICENSE（或 SetLicenseKey）")
	}
}

// allSealCategories 覆盖全部章型的授权清单，用于判断"全章型"用例是否可跑。
func allSealCategories() []string {
	return []string{
		"enterprise",
		"OVAL",
		"double_cir",
		"doubleline_cir",
		"doubleline_oval",
		"more_foot_circle",
		"more_foot_oval",
		"more_foot_rectangle",
		"labour",
		"communist",
		"youth_league",
	}
}

// skipUnlessAllSealCategoriesLicensed 全章型用例的守卫：
// 内置社区版授权仅放开 enterprise，需显式提供覆盖全部章型的授权才执行，否则跳过。
func skipUnlessAllSealCategoriesLicensed(t *testing.T) {
	t.Helper()
	lk := GetLicenseKey()
	if lk == nil || !lk.IsLicensed() {
		t.Skip("需要有效 OS_LICENSE（或 SetLicenseKey）")
	}
	have := make(map[string]struct{}, len(lk.SealCategories))
	for _, c := range lk.SealCategories {
		have[c] = struct{}{}
	}
	var missing []string
	for _, c := range allSealCategories() {
		if _, ok := have[c]; !ok {
			missing = append(missing, c)
		}
	}
	if len(missing) > 0 {
		t.Skipf("当前授权未覆盖章型 %v，跳过全章型用例", missing)
	}
}

func testFontPathOrSkip(t *testing.T) string {
	t.Helper()
	if os.Getenv("OFFICIAL_SEAL_FONTS_PATH") == "" {
		t.Skip("需要 OFFICIAL_SEAL_FONTS_PATH")
	}
	return os.Getenv("OFFICIAL_SEAL_FONTS_PATH")
}

type stubImageService struct{}

func (stubImageService) ToImageBytes(img dto.RasterImage, _ dto.ImageAttr) ([]byte, error) {
	return render.PNGBytes(img)
}

func (stubImageService) UploadBytes(_ []byte, _ dto.ImageAttr, _ bool) (string, error) {
	return "image-key", nil
}

func (stubImageService) UploadRasterImage(_ dto.RasterImage, _ dto.ImageAttr, _ bool) (string, error) {
	return "blur-key", nil
}

func TestCompanySealCreateService_IsSupport(t *testing.T) {
	svc := NewCompanySealCreateService(ElectronicAbstractSealCreateService{})
	if !svc.IsSupport(dto.SealTypeCompany, dto.BusinessTypeElectronic) {
		t.Fatal("expected support company electronic seal")
	}
	if svc.IsSupport(dto.SealTypeCompany, dto.BusinessType("PHYSICS")) {
		t.Fatal("expected no support for physics")
	}
}

func TestCompanySealCreateService_GenerateImageDefaults(t *testing.T) {
	skipIfNoLicense(t)
	fontPath := testFontPathOrSkip(t)
	svc := NewCompanySealCreateService(ElectronicAbstractSealCreateService{FontPath: fontPath})
	attr := &dto.CustomSealAttribute{Content: "XX测试科技有限公司"}
	domainSeal, err := (&dto.Seal{}).ToDomain()
	if err != nil {
		t.Fatal(err)
	}
	domainAttr, err := domainAttributeFromService(attr)
	if err != nil {
		t.Fatal(err)
	}
	img, err := svc.generateImage(domainSeal, domainAttr)
	if err != nil {
		t.Fatalf("generateImage: %v", err)
	}
	if img == nil {
		t.Fatal("nil image")
	}
	if attr.SealCategory == nil || *attr.SealCategory != dto.SealCategoryEnterprise {
		t.Fatal("seal category not defaulted on service view")
	}
	if attr.SealSpec != dto.SealSpecCircular42 {
		t.Fatalf("seal spec = %q, want %q", attr.SealSpec, dto.SealSpecCircular42)
	}
	if attr.SealCategory == nil || *attr.SealCategory != dto.SealCategoryEnterprise {
		t.Fatalf("seal category not defaulted to enterprise")
	}
}

func TestCompanySealCreateService_GetAutoSealImage(t *testing.T) {
	skipIfNoLicense(t)
	fontPath := testFontPathOrSkip(t)
	svc := NewCompanySealCreateService(ElectronicAbstractSealCreateService{
		FontPath:     fontPath,
		ImageService: stubImageService{},
	})
	attr := &dto.CustomSealAttribute{Content: "XX测试科技有限公司"}
	img, err := svc.GetAutoSealImage(&dto.Seal{}, attr)
	if err != nil {
		t.Fatal(err)
	}
	if img.ImageKey != "image-key" || img.BlurImageKey != "blur-key" || img.Size <= 0 {
		t.Fatalf("unexpected seal image: %+v", img)
	}
}
