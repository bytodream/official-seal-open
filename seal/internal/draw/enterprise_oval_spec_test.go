package draw_test

import (
	"errors"
	"os"
	"testing"

	"github.com/bytodream/official-seal-open/seal/internal/draw"
	"github.com/bytodream/official-seal-open/seal/internal/draw/module/enterprise"
	"github.com/bytodream/official-seal-open/seal/model"
	"github.com/bytodream/official-seal-open/seal/utils"
)

func TestGetEnterpriseOval4530(t *testing.T) {
	attr, err := enterprise.GetEnterprise(model.SealSpecOval4530, "财务专用章", "", "1234567890123", 1000, "", nil, true)
	if err != nil {
		t.Fatal(err)
	}
	if attr.Category != model.SealCategoryOval || attr.WidthPx == 0 || attr.HeightPx == 0 {
		t.Fatalf("%+v", attr)
	}
	if attr.EnterpriseCode == nil || attr.EnterpriseCode.Arc != 100 {
		t.Fatalf("code arc %+v", attr.EnterpriseCode)
	}
}

func TestGetEnterpriseOval4530WithInner(t *testing.T) {
	attr, err := enterprise.GetEnterprise(model.SealSpecOval4530, "", "", "", 1000, "内文", nil, true)
	if err != nil {
		t.Fatal(err)
	}
	if attr.InnerOval == nil || attr.InnerOval.Content != "内文" {
		t.Fatalf("%+v", attr)
	}
	if attr.SealLogo || attr.Logo != nil {
		t.Fatal("inner content should disable logo")
	}
}

func TestDrawOval4530WithFont(t *testing.T) {
	fontPath := os.Getenv("OFFICIAL_SEAL_FONTS_PATH")
	fontFile := draw.ResolveFont(fontPath)
	if fontFile == "" {
		t.Skip("OFFICIAL_SEAL_FONTS_PATH 未设置，跳过绘制冒烟")
	}
	if _, err := os.Stat(fontFile); err != nil {
		t.Skip("字体不存在")
	}
	cat := model.SealCategoryEnterprise
	a := &model.CustomSealAttribute{
		SealCategory:   &cat,
		SealSpec:       model.SealSpecOval4530,
		Content:        "XX测试科技有限公司",
		EnterpriseCode: "1234567890123",
		Foot:           "财务专用章",
	}
	bundle, err := draw.BuildSealAttributeBundle(a)
	if err != nil {
		t.Fatal(err)
	}
	m := draw.BuildSealPaintModel(bundle.Paint, model.DefaultSealTransparency, fontFile)
	img, err := draw.DrawSeal(m)
	if err != nil {
		if errors.Is(err, utils.ErrLicenseRequired) || errors.Is(err, utils.ErrLicenseExpired) ||
			errors.Is(err, utils.ErrSealCategoryNotLicensed) || errors.Is(err, utils.ErrUsageLimitExceeded) {
			t.Skipf("许可证限制，跳过绘制冒烟：%v", err)
		}
		t.Fatal(err)
	}
	b := img.Bounds()
	if b.Dx() < 100 || b.Dy() < 50 {
		t.Fatalf("bounds %v", b)
	}
}

func TestToSealAttributeOval45303218Spec(t *testing.T) {
	cat := model.SealCategoryEnterprise
	a := &model.CustomSealAttribute{
		SealCategory:   &cat,
		SealSpec:       model.SealSpecOval45303218,
		InnerContent:   "测试环绕",
		Head:           "专用章",
		Foot:           "财务",
		EnterpriseCode: "1234567890123",
	}
	b, err := draw.BuildSealAttributeBundle(a)
	if err != nil {
		t.Fatal(err)
	}
	c := b.Paint
	if c.Category != model.SealCategoryDoublelineOval || c.Content != "测试环绕" {
		t.Fatalf("%+v", c)
	}
}

func TestToSealAttributeFullMoreFootCircle(t *testing.T) {
	cat := model.SealCategoryMoreFootCircle
	a := &model.CustomSealAttribute{
		SealCategory:   &cat,
		SealSpec:       model.SealSpecCircular50,
		Content:        "测试公司",
		Head:           "上",
		Foots:          []model.SealFoot{{Text: "下"}},
		EnterpriseCode: "1234567890123",
	}
	bundle, err := draw.BuildSealAttributeBundle(a)
	if err != nil {
		t.Fatal(err)
	}
	c := bundle.Paint
	if c.MoreFootHead == nil || len(c.MoreFoots) != 1 {
		t.Fatalf("more foot head/foots %+v %+v", c.MoreFootHead, c.MoreFoots)
	}
	if c.Content != "测试公司" {
		t.Fatalf("content %q", c.Content)
	}
}
