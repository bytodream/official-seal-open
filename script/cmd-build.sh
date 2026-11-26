#!/usr/bin/env bash
# 打包 cmd 可执行文件（纯 Go，无 cgo，字体与 Logo 已内嵌）。
#
# 用法:
#   ./script/cmd-build.sh                          # 当前平台 -> build-out/dist
#   ./script/cmd-build.sh --all                    # darwin/linux/windows × amd64/arm64
#   ./script/cmd-build.sh --os linux --arch amd64  # 指定单平台
#   ./script/cmd-build.sh --check                  # 打包前先跑 gofmt / vet / test
#   ./script/cmd-build.sh --version 1.2.0 --out dist
#
# 选项:
#   --all             构建全部平台矩阵
#   --os OS           目标系统（darwin/linux/windows），可配合 --arch
#   --arch ARCH       目标架构（amd64/arm64）
#   --out DIR         输出目录（默认 build-out/dist）
#   --version VER     写入二进制的版本号（默认取 git describe）
#   --check           打包前执行 gofmt、vet 与全量测试
#   --no-strip        保留符号表与调试信息（体积更大，便于排查）
#   -h, --help        显示帮助
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

# 清掉可能残留的 GOROOT：它与 PATH 中的 go 版本不一致时，标准库路径会指向旧版本导致编译失败。
# 优先使用 Homebrew 安装的 go。
unset GOROOT
if [[ -x /opt/homebrew/bin/go ]]; then
  export PATH="/opt/homebrew/bin:$PATH"
fi
if ! command -v go >/dev/null 2>&1; then
  echo "未找到 go 命令，请先安装 Go（见 go.mod 的 go 指令）" >&2
  exit 1
fi

# 校验 Go 版本不低于 go.mod 声明（只比较主次版本）。
ver_num() { echo "$1" | awk -F. '{printf "%d%03d", $1, $2}'; }
need_ver="$(awk '/^go /{print $2; exit}' go.mod)"
cur_ver="$(go env GOVERSION | sed 's/^go//')"
if [[ -n "$need_ver" && "$(ver_num "$cur_ver")" -lt "$(ver_num "$need_ver")" ]]; then
  echo "需要 Go >= ${need_ver}，当前为 ${cur_ver}" >&2
  exit 1
fi
echo "==> go ${cur_ver} -> ${need_ver}+"

APP="official_seal"
PKG="./seal/cmd"
VERSION_PKG="github.com/bytodream/official-seal-open/seal/utils.Version"
OUT_DIR="build-out/dist"
BUILD_ALL=0
TARGET_OS=""
TARGET_ARCH=""
VERSION=""
RUN_CHECK=0
STRIP=1

usage() {
  cat <<'EOF'
用法: ./script/cmd-build.sh [选项]

  --all             构建全部平台矩阵（darwin/linux/windows × amd64/arm64）
  --os OS           目标系统（darwin/linux/windows）
  --arch ARCH       目标架构（amd64/arm64）
  --out DIR         输出目录（默认 build-out/dist）
  --version VER     写入二进制的版本号（默认取 git describe）
  --check           打包前执行 gofmt、vet 与全量测试
  --no-strip        保留符号表与调试信息
  -h, --help        显示帮助
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --all) BUILD_ALL=1; shift ;;
    --os) TARGET_OS="${2:-}"; shift 2 ;;
    --arch) TARGET_ARCH="${2:-}"; shift 2 ;;
    --out) OUT_DIR="${2:-}"; shift 2 ;;
    --version) VERSION="${2:-}"; shift 2 ;;
    --check) RUN_CHECK=1; shift ;;
    --no-strip) STRIP=0; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "未知参数: $1" >&2; usage; exit 1 ;;
  esac
done

if [[ -z "$VERSION" ]]; then
  VERSION="$(git describe --tags --always --dirty 2>/dev/null || echo dev)"
fi

LDFLAGS="-X ${VERSION_PKG}=${VERSION}"
if [[ "$STRIP" -eq 1 ]]; then
  LDFLAGS="-s -w ${LDFLAGS}"
fi

if [[ "$RUN_CHECK" -eq 1 ]]; then
  echo "==> gofmt"
  unformatted="$(gofmt -l seal/ || true)"
  if [[ -n "$unformatted" ]]; then
    echo "以下文件未格式化：" >&2
    echo "$unformatted" >&2
    exit 1
  fi
  echo "==> go vet"
  go vet ./...
  echo "==> go test"
  go test -count=1 ./...
fi

mkdir -p "$OUT_DIR"

build_one() {
  local os="$1" arch="$2" ext=""
  [[ "$os" == "windows" ]] && ext=".exe"
  local out="${OUT_DIR}/${APP}_${os}_${arch}${ext}"
  echo "==> ${os}/${arch} -> ${out}"
  CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" \
    go build -trimpath -ldflags "$LDFLAGS" -o "$out" "$PKG"
}

if [[ "$BUILD_ALL" -eq 1 ]]; then
  if command -v go >/dev/null 2>&1; then :; else
    echo "未找到 go 命令" >&2; exit 1
  fi
  for os in darwin linux windows; do
    for arch in amd64 arm64; do
      build_one "$os" "$arch"
    done
  done
else
  os="${TARGET_OS:-$(go env GOOS)}"
  arch="${TARGET_ARCH:-$(go env GOARCH)}"
  build_one "$os" "$arch"
fi

echo
echo "==> 产物（版本 ${VERSION}）"
ls -lh "$OUT_DIR"