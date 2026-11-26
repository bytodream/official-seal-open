package render

import (
	"github.com/bytodream/official-seal-open/seal/internal/draw"
	"github.com/bytodream/official-seal-open/seal/model"

	"bytes"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// testdataDir 返回本包 testdata 目录绝对路径。
func testdataDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller")
	}
	return filepath.Join(filepath.Dir(file), "testdata")
}

func TestLoadSealRenderFixture_jiaxingCircular42(t *testing.T) {
	path := filepath.Join(testdataDir(t), "fixture.json")
	attr, _, trans, err := LoadSealRenderFixtureFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if attr.SealSpec != model.SealSpecCircular42 || attr.Content == "" {
		t.Fatalf("unexpected attr: spec=%s content=%q", attr.SealSpec, attr.Content)
	}
	if trans != draw.TransparencyFromSealString(ptr("6")) {
		t.Fatalf("transparency %d", trans)
	}
}

func TestRenderSealFixture_jiaxingCircular42(t *testing.T) {
	fontPath := os.Getenv("OFFICIAL_SEAL_FONTS_PATH")
	if draw.ResolveFont(fontPath) == "" {
		t.Skip("需要 OFFICIAL_SEAL_FONTS_PATH")
	}
	path := filepath.Join(testdataDir(t), "fixture.json")
	pngBytes, err := RenderSealRenderFixtureFile(path, fontPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(pngBytes) < 1000 {
		t.Fatalf("png too small %d", len(pngBytes))
	}
	_, _, trans, err := LoadSealRenderFixtureFile(path)
	if err != nil {
		t.Fatal(err)
	}
	attr, _, _, err := LoadSealRenderFixtureFile(path)
	if err != nil {
		t.Fatal(err)
	}
	againPNG, err := renderCustomSealPNG(attr, fontPath, trans)
	if err != nil {
		skipIfLicenseLimited(t, err)
		t.Fatal(err)
	}
	imgA, _ := png.Decode(bytes.NewReader(pngBytes))
	imgB, _ := png.Decode(bytes.NewReader(againPNG))
	diff := CompareImages(imgA, imgB, 0)
	if diff.TotalPixels <= 0 {
		t.Fatal("no stamp overlap")
	}
	if diff.MismatchedPixels > 0 {
		t.Fatalf("fixture 两次渲染不一致: mismatched=%d/%d", diff.MismatchedPixels, diff.TotalPixels)
	}
}

func ptr(s string) *string { return &s }
