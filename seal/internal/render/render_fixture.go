package render

import (
	"github.com/bytodream/official-seal-open/seal/internal/draw"
	"github.com/bytodream/official-seal-open/seal/model"

	"encoding/json"
	"fmt"
	"os"
)

// SealRenderFixture 章面渲染 fixture（JSON 输入，用于端到端回归与示例）。
type SealRenderFixture struct {
	// FontPath 字体文件路径；为空时原样返回空串，由调用方（base.ResolveFont）回退到
	// OFFICIAL_SEAL_FONTS_PATH 环境变量或内嵌字体。
	FontPath string `json:"fontPath"`
	// Transparency 透明度配置值（0–100 字符串，为空时用默认值）
	Transparency string `json:"transparency"`
	// Attribute 自定义章面属性（CustomSealAttribute JSON）
	Attribute json.RawMessage `json:"attribute"`
}

// LoadSealRenderFixture 解析 JSON 字节；缺失 sealSpec/sealCategory 时按内核默认值补齐
// （CIRCULAR_42 / enterprise），与 model.SealSpecEnum 的默认口径一致。
func LoadSealRenderFixture(data []byte) (attr *model.CustomSealAttribute, fontPath string, transparency int, err error) {
	var root SealRenderFixture
	if err = json.Unmarshal(data, &root); err != nil {
		return nil, "", 0, err
	}
	if len(root.Attribute) == 0 {
		return nil, "", 0, fmt.Errorf("attribute is required")
	}
	attr = model.NewCustomSealAttribute()
	if err = json.Unmarshal(root.Attribute, attr); err != nil {
		return nil, "", 0, err
	}
	if attr.SealSpec == "" {
		attr.SealSpec = model.SealSpecCircular42
	}
	if attr.SealCategory == nil && attr.SealCategoryType != "" {
		c := model.ConvertSealCategory(attr.SealCategoryType)
		attr.SealCategory = &c
	}
	if attr.SealCategory == nil {
		c := model.SealCategoryEnterprise
		attr.SealCategory = &c
	}
	fontPath = root.FontPath
	if root.Transparency != "" {
		s := root.Transparency
		transparency = draw.TransparencyFromSealString(&s)
	} else {
		transparency = model.DefaultSealTransparency
	}
	return attr, fontPath, transparency, nil
}

// LoadSealRenderFixtureFile 从 disk 读取 fixture。
func LoadSealRenderFixtureFile(path string) (*model.CustomSealAttribute, string, int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", 0, err
	}
	return LoadSealRenderFixture(data)
}

// RenderSealRenderFixtureFile 渲染 fixture 文件并输出 PNG 字节。
func RenderSealRenderFixtureFile(path, fontPath string) ([]byte, error) {
	attr, fixtureFont, transparency, err := LoadSealRenderFixtureFile(path)
	if err != nil {
		return nil, err
	}
	if fontPath == "" {
		fontPath = fixtureFont
	}
	return renderCustomSealPNG(attr, fontPath, transparency)
}
