// 生成章型示意图（PNG）。
//
// 不传 -attr / -attr-text 时按内置样例渲染全部 13 个章型（用于 README 配图）。
//
// 传入自己的数据（文件或直接给 JSON 文本）：
//
//	gen-readme-seals -attr ./my-seal.json -out ./seals
//	gen-readme-seals -attr-text '{"sealCategory":"enterprise","sealSpec":"CIRCULAR_42","content":"某某有限公司"}' -out ./seals
//
// 两种写法等价。JSON 支持纯章面属性：
//
//	{"sealCategory":"enterprise","sealSpec":"CIRCULAR_42","content":"某某有限公司","foot":"合同专用章"}
//
// 也支持完整 fixture（可额外指定字体与透明度）：
//
//	{"fontPath":"/path/to/font.ttf","transparency":"30","attribute":{...}}
//
// 仓库内直接运行（默认输出到 <仓库根>/docs/images/seals）：
//
//	go run ./seal/cmd/gen-readme-seals
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bytodream/official-seal-open/seal/internal/draw"
	"github.com/bytodream/official-seal-open/seal/internal/draw/base"
	"github.com/bytodream/official-seal-open/seal/internal/render"
	"github.com/bytodream/official-seal-open/seal/model"
	"github.com/bytodream/official-seal-open/seal/utils"
)

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

// loadAttribute 解析章面数据，兼容纯 CustomSealAttribute 与完整 fixture 两种格式。
func loadAttribute(data []byte) (*model.CustomSealAttribute, string, int, error) {
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, "", 0, fmt.Errorf("JSON 解析失败: %w", err)
	}
	if _, ok := probe["attribute"]; ok {
		return render.LoadSealRenderFixture(data)
	}
	// 纯属性：包一层 fixture 以复用默认值补齐逻辑（sealSpec / sealCategory / 透明度）。
	wrapped, err := json.Marshal(map[string]json.RawMessage{"attribute": data})
	if err != nil {
		return nil, "", 0, err
	}
	return render.LoadSealRenderFixture(wrapped)
}

// resolveOutDir 计算输出目录：显式指定优先；单章模式落到当前目录，批量模式落到仓库 docs 目录。
func resolveOutDir(flagOut string, single bool) string {
	if flagOut != "" {
		return flagOut
	}
	if single {
		return "."
	}
	root, err := findRepoRoot()
	if err != nil {
		// 独立二进制在任意目录运行：不再强制要求 go.mod。
		root = "."
	}
	return filepath.Join(root, "docs", "images", "seals")
}

// safeOutputName 只保留文件名最后一段，避免 -name 携带路径跳出输出目录。
func safeOutputName(name string) string {
	name = filepath.Base(strings.TrimSpace(name))
	if name == "." || name == ".." || name == string(filepath.Separator) {
		return ""
	}
	return name
}

// renderSingle 渲染单个章面数据。source 仅用于错误提示（文件路径或 -attr-text）。
func renderSingle(data []byte, source, name, flagOut, flagFont string) {
	attr, fixtureFont, transparency, err := loadAttribute(data)
	if err != nil {
		fatalf("%s: %v", source, err)
	}
	fontPath := flagFont
	if fontPath == "" {
		fontPath = fixtureFont
	}
	if draw.ResolveFont(fontPath) == "" {
		fatalf("无法解析字体，请通过 -font、fixture.fontPath 或 OFFICIAL_SEAL_FONTS_PATH 指定")
	}

	dir := resolveOutDir(flagOut, true)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fatalf("创建输出目录失败: %v", err)
	}

	png, err := render.PNGForShowcase(attr, fontPath, transparency)
	if err != nil {
		fatalf("%s: %v", source, err)
	}

	name = safeOutputName(name)
	if name == "" {
		name = "seal"
	}
	path := filepath.Join(dir, name+".png")
	if err := os.WriteFile(path, png, 0o644); err != nil {
		fatalf("写入 %s 失败: %v", path, err)
	}
	fmt.Println("wrote", path)
}

func main() {
	outDir := flag.String("out", "", "输出目录；单章模式缺省当前目录，批量模式缺省 <仓库根>/docs/images/seals")
	fontPath := flag.String("font", "", "字体文件路径；缺省回退 fixture.fontPath、OFFICIAL_SEAL_FONTS_PATH，再回退内嵌字体")
	attrFile := flag.String("attr", "", "自定义章面数据 JSON 文件（纯 CustomSealAttribute 或完整 fixture）；不传则渲染内置样例")
	var attrText string
	flag.StringVar(&attrText, "attr-text", "", "自定义章面数据 JSON 文本；与 -attr 二选一")
	flag.StringVar(&attrText, "attrText", "", "同 -attr-text（别名）")
	name := flag.String("name", "", "单章模式的输出文件名（不含扩展名），缺省取 -attr 的文件名，-attr-text 时为 seal")
	licenseText := flag.String("license", "", "授权密钥内容；也可用 OS_LICENSE 环境变量")
	licenseFile := flag.String("license-file", "", "授权密钥文件路径（优先于 -license）")
	seed := flag.Int64("seed", 0, "随机种子：非 0 时生成结果可复现，便于图像回归比对")
	showVersion := flag.Bool("version", false, "打印版本号后退出")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `用法: %s [选项]

不传 -attr / -attr-text：渲染内置的 13 个章型样例图。
传入其中之一：按自己的数据生成单个章面图。

示例:
  %s -out ./seals -seed 42 -license-file ./license.txt
  %s -attr ./my-seal.json -name my-company -out ./seals
  %s -attr-text '{"sealCategory":"enterprise","content":"某某有限公司"}' -out ./seals

选项:
`, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]), filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion {
		fmt.Println(utils.Version)
		return
	}

	// 归一化种子必须在首次绘制前设置，base.NewRand 会惰性读取该环境变量。
	if *seed != 0 {
		if err := os.Setenv(base.RandomSeedEnv, strconv.FormatInt(*seed, 10)); err != nil {
			fatalf("设置随机种子失败: %v", err)
		}
	}

	if *licenseFile != "" {
		data, err := os.ReadFile(*licenseFile)
		if err != nil {
			fatalf("读取授权文件失败: %v", err)
		}
		if err := utils.SetLicenseKey(string(data)); err != nil {
			fatalf("授权无效: %v", err)
		}
	} else if *licenseText != "" {
		if err := utils.SetLicenseKey(*licenseText); err != nil {
			fatalf("授权无效: %v", err)
		}
	}

	if *attrFile != "" && attrText != "" {
		fatalf("-attr 与 -attr-text 只能使用一个")
	}
	if *attrFile != "" {
		data, err := os.ReadFile(*attrFile)
		if err != nil {
			fatalf("读取章面数据失败: %v", err)
		}
		n := *name
		if n == "" {
			n = strings.TrimSuffix(filepath.Base(*attrFile), filepath.Ext(*attrFile))
		}
		renderSingle(data, *attrFile, n, *outDir, *fontPath)
		return
	}
	if attrText != "" {
		n := *name
		if n == "" {
			n = "seal"
		}
		renderSingle([]byte(attrText), "-attr-text", n, *outDir, *fontPath)
		return
	}

	dir := resolveOutDir(*outDir, false)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fatalf("创建输出目录失败: %v", err)
	}
	if draw.ResolveFont(*fontPath) == "" {
		fatalf("无法解析字体，请通过 -font 或 OFFICIAL_SEAL_FONTS_PATH 指定")
	}
	for _, item := range render.ShowcaseItems() {
		png, err := render.PNGForShowcase(item.Attr, *fontPath, model.DefaultSealTransparency)
		if err != nil {
			fatalf("%s: %v", item.Slug, err)
		}
		path := filepath.Join(dir, item.Slug+".png")
		if err := os.WriteFile(path, png, 0o644); err != nil {
			fatalf("写入 %s 失败: %v", path, err)
		}
		fmt.Println("wrote", path)
	}
}
