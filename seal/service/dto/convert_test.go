package dto

import (
	"testing"

	model "github.com/bytodream/official-seal-open/seal/model"
)

func TestCustomSealAttribute_ToDomain_roundTrip(t *testing.T) {
	cat := SealCategoryEnterprise
	a := NewCustomSealAttribute()
	a.Content = "示例公司"
	a.SealSpec = SealSpecCircular42
	a.SealCategory = &cat
	a.EdgeWidth = f32ptr(1.2)

	domain, err := a.ToDomain()
	if err != nil {
		t.Fatal(err)
	}
	if domain == nil || domain.Content != a.Content || domain.SealSpec != a.SealSpec {
		t.Fatalf("domain mismatch: %+v", domain)
	}
	if domain.SealCategory == nil || *domain.SealCategory != model.SealCategory(cat) {
		t.Fatal("seal category not converted")
	}
}

func TestSealModel_ToDomain_nilElectronic(t *testing.T) {
	m := &SealModel{Category: BusinessTypeElectronic}
	domain, err := m.ToDomain()
	if err != nil {
		t.Fatal(err)
	}
	if domain.Category != model.BusinessTypeElectronic {
		t.Fatalf("category=%s", domain.Category)
	}
}

func TestSealModel_ToDomain_FromDomain_roundTrip(t *testing.T) {
	cat := SealCategoryEnterprise
	radius := 3
	transparency := "80"
	auto := NewCustomSealAttribute()
	auto.Content = "往返测试"
	auto.SealCategory = &cat
	dm := &model.SealModel{
		Category: model.BusinessTypeElectronic,
		ElectronicModel: &model.ElectronicModel{
			SealSpec:      model.SealSpecCircular42,
			OpenImageBlur: true,
			ShowWatermark: true,
			Radius:        &radius,
			Transparency:  &transparency,
			AutoAttribute: &model.CustomSealAttribute{
				Content:      auto.Content,
				SealSpec:     model.SealSpecCircular42,
				SealCategory: (*model.SealCategory)(&cat),
			},
		},
	}
	view, err := FromDomainSealModel(dm)
	if err != nil {
		t.Fatal(err)
	}
	if view.ElectronicModel == nil || !view.ElectronicModel.OpenImageBlur || !view.ElectronicModel.ShowWatermark {
		t.Fatalf("electronic model fields lost: %+v", view.ElectronicModel)
	}
	if view.ElectronicModel.AutoAttribute == nil || view.ElectronicModel.AutoAttribute.Content != auto.Content {
		t.Fatal("auto attribute not round-tripped via FromDomain")
	}
	back, err := view.ToDomain()
	if err != nil {
		t.Fatal(err)
	}
	if back.ElectronicModel == nil || back.ElectronicModel.SealSpec != model.SealSpecCircular42 {
		t.Fatalf("ToDomain after FromDomain: %+v", back.ElectronicModel)
	}
	if back.ElectronicModel.Radius == nil || *back.ElectronicModel.Radius != radius {
		t.Fatal("radius mismatch after round trip")
	}
}

func f32ptr(v float32) *float32 { return &v }
