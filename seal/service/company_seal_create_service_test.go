package service

import (
	"errors"
	"image"
	"os"
	"testing"

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

func TestCompanySealCreateService_IsSupport(t *testing.T) {
	svc := NewCompanySealCreateService()
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
	svc := NewCompanySealCreateService().SetFontPath(fontPath)
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
	svc := NewCompanySealCreateService().SetFontPath(fontPath).SetBlurGenerate(true)
	keys := []string{"image-key", "blur-key"}
	uploads := 0
	svc.SetUploadBytes(func(data []byte, attr dto.ImageAttr, companySeal bool) (string, error) {
		if uploads >= len(keys) {
			t.Fatalf("上传次数超出预期：%d", uploads)
		}
		key := keys[uploads]
		uploads++
		return key, nil
	})
	radius := 10
	attr := &dto.CustomSealAttribute{Content: "XX测试科技有限公司"}
	img, err := svc.GetAutoSealImage(&dto.Seal{Radius: &radius}, attr)
	if err != nil {
		t.Fatal(err)
	}
	if img.ImageKey != "image-key" || img.BlurImageKey != "blur-key" || img.Size <= 0 {
		t.Fatalf("unexpected seal image: %+v", img)
	}
	if uploads != 2 {
		t.Fatalf("上传次数 = %d, want 2（清晰图 + 模糊图）", uploads)
	}
}

// TestCompanySealCreateService_HasUploadBytes 覆盖上传能力探测与注入后的实际上传。
func TestCompanySealCreateService_HasUploadBytes(t *testing.T) {
	svc := NewCompanySealCreateService()
	if svc.HasUploadBytes() {
		t.Fatal("默认装配不应具备上传能力")
	}
	svc.SetUploadBytes(func(data []byte, attr dto.ImageAttr, companySeal bool) (string, error) {
		return "uploaded-key", nil
	})
	if !svc.HasUploadBytes() {
		t.Fatal("SetUploadBytes 后 HasUploadBytes 应为 true")
	}
	got, err := svc.ImageService.UploadBytes([]byte{1}, dto.ImageAttr{}, true)
	if err != nil || got != "uploaded-key" {
		t.Fatalf("UploadBytes = (%q, %v), want (uploaded-key, nil)", got, err)
	}
}

// TestCompanySealCreateService_IsBlurGenerate 覆盖模糊图开关的读写。
func TestCompanySealCreateService_IsBlurGenerate(t *testing.T) {
	svc := NewCompanySealCreateService()
	if svc.IsBlurGenerate() {
		t.Fatal("默认不应生成模糊图")
	}
	if !svc.SetBlurGenerate(true).IsBlurGenerate() {
		t.Fatal("SetBlurGenerate(true) 后 IsBlurGenerate 应为 true")
	}
	if svc.SetBlurGenerate(false).IsBlurGenerate() {
		t.Fatal("SetBlurGenerate(false) 后 IsBlurGenerate 应为 false")
	}
}

// TestCompanySealCreateService_TypedNilImageService 持有 nil 指针的接口值应等价于未注入，
// 而不是在调用方法时 panic。
func TestCompanySealCreateService_TypedNilImageService(t *testing.T) {
	svc := NewCompanySealCreateService().SetImageService((*DefaultImageService)(nil))
	if svc.HasUploadBytes() {
		t.Fatal("nil 图片服务的 HasUploadBytes 应为 false")
	}
	if svc.SetUploadBytes(nil).SetBlurGenerate(true); svc.IsBlurGenerate() {
		t.Fatal("nil 图片服务的配置应被忽略")
	}
	radius := 10
	if _, err := svc.GetAutoSealImage(&dto.Seal{Radius: &radius}, dto.NewCustomSealAttribute()); !errors.Is(err, dto.ErrImageServiceNotConfigured) {
		t.Fatalf("err = %v, want ErrImageServiceNotConfigured", err)
	}
}

// nonConfiguratorImageService 只实现 ImageService，不实现 ImageServiceConfigurator。
type nonConfiguratorImageService struct{}

func (nonConfiguratorImageService) ToImageBytes(image.Image, dto.ImageAttr) ([]byte, error) {
	return []byte("png"), nil
}

func (nonConfiguratorImageService) UploadBytes([]byte, dto.ImageAttr, bool) (string, error) {
	return "custom-key", nil
}

// TestCompanySealCreateService_NonConfiguratorImageService 未实现 ImageServiceConfigurator 时
// 配置入口被忽略、HasUploadBytes 恒为 false，GetAutoSealImage 在渲染前短路。
func TestCompanySealCreateService_NonConfiguratorImageService(t *testing.T) {
	svc := NewCompanySealCreateService().SetImageService(nonConfiguratorImageService{})
	if svc.HasUploadBytes() {
		t.Fatal("未实现 ImageServiceConfigurator 时 HasUploadBytes 应为 false")
	}
	svc.SetBlurGenerate(true)
	if svc.IsBlurGenerate() {
		t.Fatal("未实现 ImageServiceConfigurator 的模糊开关应被忽略")
	}
	radius := 10
	if _, err := svc.GetAutoSealImage(&dto.Seal{Radius: &radius}, dto.NewCustomSealAttribute()); !errors.Is(err, dto.ErrImageUploadNotConfigured) {
		t.Fatalf("err = %v, want ErrImageUploadNotConfigured", err)
	}
}
