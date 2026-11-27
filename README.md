# official-seal

**official-seal** 是一个用 Go 编写的**企业电子公章章面生成库**。根据章型、规格与文案参数，在本地光栅化生成透明 PNG 章面，适用于电子签章、合同预览、印章制作等场景。集成方只需依赖 `seal` 包对外 API；绘制与渲染实现位于 `internal` 子包，不对外暴露。

- 模块路径：`github.com/bytodream/official-seal`
- 语言要求：Go 1.25+

---

## 核心能力

| 能力 | 说明 |
|------|------|
| **章面生成** | 纯 Go 矢量绘制 → PNG；支持多种国标/常用章型与毫米规格 |
| **企业服务** | `CompanySealCreateService`：预览、生成清晰/模糊章图、上传存储键 |
| **章面参数** | 环绕主文、脚注、企业代码弧排、内圈文字、多横排、中心图案、防伪等 |
| **字体** | 可配置 `OFFICIAL_SEAL_FONTS_PATH`；未配置时使用内嵌字体回退 |
| **随机种子** | `OFFICIAL_SEAL_RANDOM_SEED`：设为非 0 整数时，防伪纹与模糊加粗结果可复现（便于图像回归比对）；未设置则每次随机 |
| **许可证** | `OS_LICENSE` / `SetLicenseKey`；可用章型仅由许可证内的 `seal_categories` 决定（未配置即无章型权限） |
| **使用次数** | 非样章绘制成功后计入一次（`UsageLimit>0` 时），记录在 `OS_USAGE_FILE`（默认 `~/.official-seal/usage.json`）；该文件为本地明文记录，属于客户端诚实模型，服务端计费请以自己的账务为准 |
| **密钥工具** | `internal` 包内 RSA 混合封装/加密（签发工具用，非集成必选） |

---

## 对外 API（`seal` 包）

使用前加载许可证：

```go
import "github.com/bytodream/official-seal-open/seal"

_ = seal.SetLicenseKey(os.Getenv("OS_LICENSE"))
```

| API | 作用 |
|-----|------|
| `SetLicenseKey` / `GetLicenseKey` | 加载与查询许可证 |
| `CompanySealCreateService.Preview` | 按 `SealModel` 生成预览 PNG |
| `CompanySealCreateService.GetAutoSealImage` | 生成清晰/模糊图并上传 |
| `CompanySealCreateService.RenderAttributePNG` | 按 `CustomSealAttribute` 输出 PNG |
| `CompanySealCreateService.InitSealPermission` | 校验许可证是否已加载有效 |

章面生成前在**服务层**统一做许可证与章型校验；生成器固定为 `CommonSealGenerator`，不可注入自定义实现。

`GetAutoSealImage` 需要 `ImageService`（转码 + 上传）。库默认装配 `service.DefaultImageService`，它实现转码部分，行为对齐 Java `ImageServiceImpl.toImageBytes`（透明化 → 缩放 → 旋转 → PNG 编码）；上传需自行接入存储，内嵌即可复用转码：

```go
type myImageService struct{ service.DefaultImageService }

func (myImageService) UploadBytes(data []byte, attr dto.ImageAttr, companySeal bool) (string, error) {
    return myStorage.Upload(data, "image/png"), nil
}

// 内嵌类型的方法都是指针接收者，赋值时须取地址
svc := service.NewCompanySealCreateService().SetImageService(&myImageService{})
```

`ImageAttr` 的 `Transparent`、`Zoom`、`Rotate`、`LimitWidth`、`LimitHeight`、`Degree` 由默认实现消费；显式 `SetImageService(nil)` 后这些字段不再生效。

### 集成约定

- 请 **`import "github.com/bytodream/official-seal-open/seal"`**，类型使用本包导出的 `SealModel`、`CustomSealAttribute` 等；勿直接依赖 `seal/model`、`seal/internal`（`seal/utils` 为公开工具包）。
- 前端创建印章 JSON 可经 **`ParseSealModelFromCreateJSON`** 转为 `SealModel` 后调用 `Preview`。
- 服务层 `CustomSealAttribute` 经 JSON 转换进入内核，**不支持** `ToSealAttributeFunc` 回调；复杂映射请在业务侧组装好字段后再调用 API。
- 服务构造改为无参 + 链式 `Set`：由 `NewCompanySealCreateService(base)` 改为 `NewCompanySealCreateService().SetFontPath(p).SetImageService(s)`；公共依赖类型 `ElectronicAbstractSealCreateService` 不再导出。

---

## 架构分层

```
seal/                 对外门面（alias → service）
├── alias.go          门面别名：类型、错误、规格常量
├── doc.go            包说明
├── cmd/              辅助命令（如效果图生成）
├── service/          服务与许可证
│   └── dto/          集成方 DTO；ToDomain → 内核 model
├── model/            内核领域模型
├── utils/            通用工具（像素换算、日志、许可证、IoC）
└── internal/         内部实现
    ├── draw/         章面绘制（base 公共层 + module 各章型，每章自包含）
    ├── render/       PNG 渲染与图像效果
    └── private/      许可证密钥编解码
```

集成方类型请使用门面 `seal` 或 `seal/service/dto`。

### 绘制模块自包含

每个章型位于 `seal/internal/draw/module/<章型>/`，包内自带本章所需的全部绘制实现（`core_impl.go`），不引用其它章型包；`enterprise` 复用椭圆章逻辑的部分也已复制为 `oval_impl.go`。

- **按章交付**：单独授权某一章型时，只需该目录加上公共层（`base`、`model`、`utils`），无需暴露其它章型的绘制代码
- **零跨章依赖**：`draw/module` 内任意两个包之间不存在 import
- **同步成本**：公共实现的改动需同步到各章副本，改一处要改 11 处

公共层 `base` 只保留被跨章复用的基础设施：章面属性、绘制模型、字体解析、随机源与透明度换算。

---

## 章型效果图

下列样章由库内渲染管线生成（样例公司名：**示例科技有限公司**）。图片路径：`docs/images/seals/`。

| 章型 | 预览 |
|:-----|:----:|
| 企业圆章 42mm | <img class="no-margin" src="docs/images/seals/enterprise_circular_42.png" alt="企业圆章 42mm" width="50%" height="50%"> |
| 企业圆章 34mm | <img class="no-margin" src="docs/images/seals/enterprise_circular_34.png" alt="企业圆章 34mm" width="50%" height="50%"> |
| 企业椭圆章 | <img class="no-margin" src="docs/images/seals/oval_45_30.png" alt="企业椭圆章" width="50%" height="50%">                 |
| 双行椭圆章 | <img class="no-margin" src="docs/images/seals/doubleline_oval.png" alt="双行椭圆章" width="50%" height="50%">            |
| 双行圆章 | <img class="no-margin" src="docs/images/seals/doubleline_cir.png" alt="双行圆章" width="50%" height="50%">              |
| 双圆章 | <img class="no-margin" src="docs/images/seals/double_cir.png" alt="双圆章" width="50%" height="50%">                   |
| 党委章 | <img class="no-margin" src="docs/images/seals/communist.png" alt="党委章" width="50%" height="50%">                    |
| 工会章 | <img class="no-margin" src="docs/images/seals/labour.png" alt="工会章" width="50%" height="50%">                       |
| 团章 | <img class="no-margin" src="docs/images/seals/youth_league.png" alt="团章" width="50%" height="50%">                  |
| 多横排圆章 | <img class="no-margin" src="docs/images/seals/more_foot_circle.png" alt="多横排圆章" width="50%" height="50%">           |
| 多横排椭圆章 | <img class="no-margin" src="docs/images/seals/more_foot_oval.png" alt="多横排椭圆章" width="50%" height="50%">            |
| 多横排长方章 | <img class="no-margin" src="docs/images/seals/more_foot_rectangle.png" alt="多横排长方章" width="50%" height="50%">       |
| DIY 企业圆章 43mm | <img class="no-margin" src="docs/images/seals/diy_enterprise.png" alt="DIY 企业圆章 43mm" width="50%" height="50%">     |

---

## 支持的章型（SealCategory）

下表为代码中定义的章面分类。上节效果图覆盖主流通用章型；其余分类由各章型模块（`seal/internal/draw/module`）构建，可在商业版许可证下按模型启用。

| 分类值 | 名称 | 说明 |
|--------|----|------|
| `enterprise` | 企业圆章 | 标准圆形公章；主文弧排、可选脚注与企业信用代码、中心图案与防伪 |
| `OVAL` | 企业椭圆章 | 椭圆外廓；可含脚注、内椭圆文字 |
| `doubleline_cir` | 双行圆章 | 圆章 + 上排横字（Head）+ 下排脚注 |
| `doubleline_oval` | 双行椭圆章 | 椭圆 + 双行横排 + 可选内环/内椭圆文字 |
| `double_cir` | 双圆章 | 双圆圈层；外圈主文、内圈文字（InnerContent） |
| `communist` | 党委章 | 党徽体系圆章样式 |
| `labour` | 工会章 | 工会圆章样式 |
| `youth_league` | 团章 | 团章面（含团委/支部等排版） |
| `more_foot_circle` | 多横排圆章 | 圆章 + 多段横排脚注（Foots） |
| `more_foot_oval` | 多横排椭圆章 | 椭圆 + 多横排 |
| `more_foot_rectangle` | 多横排长方章 | 长方形版面 + 多横排文字 |
| `rectangle` | 普通长方形章 | 长方形印章 |
| `square` | 正方形章 | 正方形印章 |
| `CPPCC` | 政协章 | 人民政治协商会议章 |
| `rectangle_timestamp` | 长方形时间戳章 | 时间戳类长方章 |
| `csj_shusp` | 长三角一网通办 | 长三角一网通办专用样式 |
| `public_security` | 公检法司 | 公检法司类章型 |
| `double_oval` | 双椭圆章 | 双椭圆结构 |
| `archives` | 档案章 | 档案管理用章 |
| `rectangle_land` | 规土局方章 | 规划国土类方章 |
| `purpose_seal` | 用途章 | 用途类印章 |

### 常用规格示例（SealSpec）

| 规格 | 含义 |
|------|------|
| `CIRCULAR_42` | 直径 42 mm 圆章（默认） |
| `CIRCULAR_34` | 直径 34 mm 圆章 |
| `CIRCULAR_50` | 直径 50 mm 圆章（多横排圆章常用） |
| `OVAL_45_30` | 45×30 mm 椭圆 |
| `OVAL_45_30_32_18` | 双行椭圆复合规格 |
| `RECTANGLE_50_20` | 50×20 mm 长方形 |
| `DIY_SPEC` | 自定义毫米宽高、字号、边宽等 |

---

## 构建与打包

前置：Go ≥ 1.25（见 `go.mod`）。项目不使用 cgo，字体与 Logo 已通过 `go:embed` 内嵌，产物是**单文件、无外部依赖**。

### 一键脚本

`script/cmd-build.sh` 封装了交叉编译、版本注入与产物汇总：

```bash
./script/cmd-build.sh                      # 当前平台 -> build-out/dist
./script/cmd-build.sh --all                # darwin / linux / windows × amd64 / arm64
./script/cmd-build.sh --os linux --arch amd64
./script/cmd-build.sh --check              # 先跑 gofmt / vet / go test，再打包
./script/cmd-build.sh --version v1.2.0 --out dist
```

| 参数 | 说明 |
| --- | --- |
| `--all` | 构建全部平台矩阵 |
| `--os` / `--arch` | 指定目标系统与架构 |
| `--out DIR` | 输出目录，默认 `build-out/dist`（已在 `.gitignore` 中） |
| `--version VER` | 写入二进制的版本号，默认取 `git describe` |
| `--check` | 打包前执行 gofmt、vet 与全量测试 |
| `--no-strip` | 保留符号表与调试信息 |

### 手工等价命令

```bash
CGO_ENABLED=0 go build -trimpath \
  -ldflags "-s -w -X github.com/bytodream/official-seal-open/seal/utils.Version=$(git describe --tags --always)" \
  -o build-out/dist/official_seal ./seal/cmd
```

### 运行

```bash
./official_seal -version                     # 查看版本
./official_seal -out ./seals                 # 输出全部章型示例图
./official_seal -out ./seals -seed 42        # 固定随机种子，结果可复现
./official_seal -out ./seals -license-file ./license.txt
./official_seal -attr ./my-seal.json -name my-company -out ./seals
./official_seal -attr-text '{"sealCategory":"enterprise","content":"某某有限公司"}' -out ./seals
```

| 参数 | 说明 |
| --- | --- |
| `-out` | 输出目录；批量模式缺省 `<仓库根>/docs/images/seals`，单章模式缺省当前目录 |
| `-attr` | 自定义章面数据 JSON 文件；不传则渲染内置的 13 个样例 |
| `-attr-text` / `-attrText` | 自定义章面数据 JSON 文本，与 `-attr` 二选一（命令行直接传参时更方便） |
| `-name` | 单章模式的输出文件名（不含扩展名），缺省取 `-attr` 的文件名 |
| `-seed` | 随机种子，非 0 时防伪纹与模糊加粗结果可复现 |
| `-license` / `-license-file` | 授权密钥内容 / 文件；也可用 `OS_LICENSE` 环境变量 |
| `-font` | 字体路径，缺省回退 `fixture.fontPath`、`OFFICIAL_SEAL_FONTS_PATH` 或内嵌字体 |

### 自定义章面数据

`-attr`、`-attr-text` 传的是同一份数据，支持两种 JSON 写法。纯章面属性（字段同 `model.CustomSealAttribute`）：

```json
{
  "sealCategory": "enterprise",
  "sealSpec": "CIRCULAR_42",
  "content": "杭州某某科技有限公司",
  "foot": "合同专用章",
  "enterpriseCode": "330106123456789"
}
```

完整 fixture（可额外指定字体与透明度）：

```json
{
  "fontPath": "/path/to/font.ttf",
  "transparency": "20",
  "attribute": {
    "sealCategory": "more_foot_rectangle",
    "sealSpec": "RECTANGLE_50_20",
    "content": "上海某某贸易有限公司",
    "foots": [{ "text": "销售部门" }, { "text": "财务部门" }]
  }
}
```

缺省规则：`sealSpec` 缺失按 `CIRCULAR_42`、`sealCategory` 缺失按 `enterprise`、`transparency` 缺失取默认值。自定义数据与样例走同一条渲染链路，同样经过章型权限校验并使用次数统计。

> 二进制默认携带社区版授权（仅 `enterprise`）；生成其它章型需传入对应授权。macOS 对外分发需自行 `codesign` 与公证。

## 许可证与章型权限

- 环境变量 **`OS_LICENSE`**：HybridSeal 载荷（启动时自动加载）；兼容别名 `OFFICIAL_SEAL_LICENSE`。
- **Tier**：`community` / `individual` / `business` 对应不同默认可用章型；JSON 中 **`seal_categories`** 可显式限定列表。
- 未授权或章型不在许可范围内时，服务 API 返回 `ErrLicenseRequired` / `ErrSealCategoryNotLicensed` 等。

### 社区版许可证（开箱即用）

仓库已内置**社区版（community）**许可证载荷，开箱即可直接使用（源文件 `seal/internal/private/key/yunruan_license_community.txt` 属本地签发产物，`key/` 目录已被 `.gitignore` 忽略、不随仓库提交，故完整内容如下）：

```text
AgEAy7qZ7+eVbLwP0H7HR0MUd9IppoSqwvSgdr6IWTD8yzN5gesf5PgNMi4zkX58KTJB9Kk5F0YFtv7Pfr1TmoaDwv3ksIDMU63RNYEt8HAuM4j4r2b0kC52E96SgSmamG1aw3Gw2lDKDT+wdkzwl71ZuOx2b7prCcecTELdv9nJ+tyDJqYj79p1R1O0CzAHL1Ra1jmhzj7RVKwJNg5Ua7+cIzeyFp+pGfkTAuO0cHxoiT+ag5O+bnIgZdzSNMCAlLdrJt0s5qipt99dGzAjRTs4cDOBozMRE5YcFfG09IRjaMM4dpmSVTECHUonZ2RvyYpjNiLM2IsZo0Hzt1V26qUM6+97SwrB2tRnN9PveYzWZo5MKQRfyCWgYVQ+ORx8i48rg4iBbd0s0vjCwstLR83zeaE4/hnZL4uZzkdK3THT2lUSghfEYCdeGvXDkQF7rTsiNmUoJGEkU+j8wdo1c2f1I0zZ8fFFSYD0t6RbHlgagyRoHQF6145sOizJSgTXG68RP3EbabACcJ4FXglD7S/YFniDCg9tbwEchLUZyALALdERVyY7y8gbVYlfowJdeTvBA0Zvjp6SvrdjQK+5dPcncFTRN6vboghAcExbKDKYEWzn30Kt6lE7xMvxmjnnES5/LW/Py30mtGD8ULgFSmF0MWuTTTvTWuhW3awHfpA5Ku4pxmq2llTxjpd8EaAZyLfCTGbg
```

通过环境变量启用：

```bash
export OS_LICENSE='AgEAy7qZ7+eVbLwP0H7HR0MUd9IppoSqwvSgdr6IWTD8yzN5gesf5PgNMi4zkX58KTJB9Kk5F0YFtv7Pfr1TmoaDwv3ksIDMU63RNYEt8HAuM4j4r2b0kC52E96SgSmamG1aw3Gw2lDKDT+wdkzwl71ZuOx2b7prCcecTELdv9nJ+tyDJqYj79p1R1O0CzAHL1Ra1jmhzj7RVKwJNg5Ua7+cIzeyFp+pGfkTAuO0cHxoiT+ag5O+bnIgZdzSNMCAlLdrJt0s5qipt99dGzAjRTs4cDOBozMRE5YcFfG09IRjaMM4dpmSVTECHUonZ2RvyYpjNiLM2IsZo0Hzt1V26qUM6+97SwrB2tRnN9PveYzWZo5MKQRfyCWgYVQ+ORx8i48rg4iBbd0s0vjCwstLR83zeaE4/hnZL4uZzkdK3THT2lUSghfEYCdeGvXDkQF7rTsiNmUoJGEkU+j8wdo1c2f1I0zZ8fFFSYD0t6RbHlgagyRoHQF6145sOizJSgTXG68RP3EbabACcJ4FXglD7S/YFniDCg9tbwEchLUZyALALdERVyY7y8gbVYlfowJdeTvBA0Zvjp6SvrdjQK+5dPcncFTRN6vboghAcExbKDKYEWzn30Kt6lE7xMvxmjnnES5/LW/Py30mtGD8ULgFSmF0MWuTTTvTWuhW3awHfpA5Ku4pxmq2llTxjpd8EaAZyLfCTGbg'
```

> 社区版许可证仅覆盖 `enterprise`（企业章）章型；如需更多公章类型（党委、共青团、工会、双圆章、双线椭圆章、多横排章等），请联系作者获取商业版授权。

---

## 许可协议

仓库根目录 [LICENSE](LICENSE)（Apache 2.0）。商业功能与章型以正式 **OS_LICENSE** 为准。

---

## 💌 支持作者

如果觉得项目不错，或者已经在使用了，希望你可以去 <a target="_blank" href="https://github.com/bytodream/official-seal">Github</a> 帮我点个 ⭐ Star，这将是对我极大的鼓励与支持。

> 喝杯咖啡 ☕️ 或者来杯奶茶 🧋，让作者更有精神，写出更棒的代码！

<img class="no-margin" src="docs/images/zs.png" alt="微信打赏" width="200" height="200">
<img class="no-margin" src="docs/images/wx.png" alt="微信好友" width="200" height="200">

> **特别感谢：**
> **如果需要更多公章类型（党委、共青团、工会、双圆章、双线椭圆章、多横排章等），或者需要商业版授权，欢迎联系我！**
> **使用中如有问题，欢迎通过仓库 Issues 反馈！**
---

## 免责声明

本项目部分内容（包括但不限于文字、图片、代码示例等）可能来源于互联网或第三方，仅用于学习与研究目的。若相关内容涉及侵权，请联系本人进行删除：请通过仓库 Issues 或项目主页提供的联系方式告知，核实后将在第一时间删除相关内容。

---

英文说明见 [README.en.md](README.en.md)。
