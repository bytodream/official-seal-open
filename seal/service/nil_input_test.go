package service

import (
	"errors"
	"testing"

	"github.com/bytodream/official-seal-open/seal/service/dto"
)

// TestPreviewNilModel 空入参返回错误而非 panic。
func TestPreviewNilModel(t *testing.T) {
	svc := NewCompanySealCreateService(ElectronicAbstractSealCreateService{})
	if _, err := svc.Preview(nil); !errors.Is(err, dto.ErrSealAttributeNotConfigured) {
		t.Fatalf("err = %v, want ErrSealAttributeNotConfigured", err)
	}
}

// TestGetAutoSealImageNilSeal 空 Seal 入参返回错误而非 panic。
func TestGetAutoSealImageNilSeal(t *testing.T) {
	svc := NewCompanySealCreateService(ElectronicAbstractSealCreateService{ImageService: stubImageService{}})
	if _, err := svc.GetAutoSealImage(nil, dto.NewCustomSealAttribute()); !errors.Is(err, dto.ErrSealAttributeNotConfigured) {
		t.Fatalf("err = %v, want ErrSealAttributeNotConfigured", err)
	}
}

// TestGetAutoSealImageMissingImageService 未注入图片服务时返回明确错误。
func TestGetAutoSealImageMissingImageService(t *testing.T) {
	svc := NewCompanySealCreateService(ElectronicAbstractSealCreateService{})
	radius := 10
	seal := &dto.Seal{Radius: &radius}
	if _, err := svc.GetAutoSealImage(seal, dto.NewCustomSealAttribute()); !errors.Is(err, dto.ErrImageServiceNotConfigured) {
		t.Fatalf("err = %v, want ErrImageServiceNotConfigured", err)
	}
}
