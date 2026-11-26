#!/usr/bin/env bash
# 打 tag 并推送到源码库
# 用法: ./tag.sh [tag] [说明]
# 示例: ./tag.sh v1.0.0 "第一个正式版本"
#       ./tag.sh              # 在最新 vX.Y.Z 上 patch +1（如 v1.0.1 -> v1.0.2）
#       ./tag.sh "修复某问题"   # 自动递增版本，说明为参数

set -euo pipefail

semver_tag_re='^v[0-9]+\.[0-9]+\.[0-9]+$'

# 取最新 vX.Y.Z，patch 加 1；无 tag 时从 v1.0.0 开始。
next_semver_tag() {
  local latest major minor patch ver
  latest="$(
    git tag -l 'v[0-9]*.[0-9]*.[0-9]*' --sort=-v:refname 2>/dev/null \
      | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' \
      | head -1
  )"
  if [[ -z "$latest" ]]; then
    echo "v1.0.0"
    return
  fi
  ver="${latest#v}"
  IFS=. read -r major minor patch <<<"$ver"
  echo "v${major}.${minor}.$((patch + 1))"
}

if [[ $# -eq 0 ]]; then
  TAG="$(next_semver_tag)"
  MESSAGE="release ${TAG}"
elif [[ $# -eq 1 && ! "$1" =~ $semver_tag_re ]]; then
  TAG="$(next_semver_tag)"
  MESSAGE="$1"
elif [[ "$1" =~ $semver_tag_re ]]; then
  TAG="$1"
  MESSAGE="${2:-release ${TAG}}"
else
  echo "用法: $0 [tag] [说明]"
  echo "示例: $0 v1.0.0 \"第一个正式版本\""
  echo "      $0                    # 自动 patch +1"
  echo "      $0 \"修复描述\"           # 自动 patch +1 并使用说明"
  exit 1
fi

cd "$(git rev-parse --show-toplevel)"

# 校验 tag 格式（可选，要求 vX.Y.Z）
if [[ ! "$TAG" =~ $semver_tag_re ]]; then
  echo "警告: tag \"${TAG}\" 不符合 vX.Y.Z 格式"
  read -rp "仍要继续吗? [y/N] " confirm
  [[ "$confirm" =~ ^[yY]$ ]] || exit 1
fi

# 不允许重复打 tag
if git rev-parse -q --verify "refs/tags/${TAG}" >/dev/null; then
  echo "错误: tag ${TAG} 已存在"
  exit 1
fi

# 工作区必须干净
if [[ -n "$(git status --porcelain)" ]]; then
  echo "错误: 工作区有未提交的改动，请先提交"
  git status --short
  exit 1
fi

echo "将创建 tag: ${TAG}（${MESSAGE}）"

git tag -a "$TAG" -m "$MESSAGE"
echo "已创建 tag: ${TAG}"

git push origin "$TAG"
echo "已推送 tag ${TAG} 到 origin"